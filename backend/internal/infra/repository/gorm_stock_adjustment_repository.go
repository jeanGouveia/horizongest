package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type GormStockAdjustmentPending struct {
	ID              uint    `gorm:"primaryKey;autoIncrement"`
	OrderID         uint    `gorm:"not null;index"`
	IngredientID    uint    `gorm:"not null;index"`
	Quantity        float64 `gorm:"not null"`
	OrderStatus     string  `gorm:"not null"`
	Status          string  `gorm:"not null;default:'pending';index"`
	CreatedAt       int64   `gorm:"autoCreateTime"`
	ProcessedAt     *int64  `gorm:"index"`
	ProcessedBy     *uint   `gorm:"index"`
	ProcessingNotes string  `gorm:"type:text"`
	IngredientName  string  `gorm:"not null"`              // snapshot do nome
	IngredientUnit  string  `gorm:"not null;default:'un'"` // snapshot da unidade
	DeletedAt       *int64  `gorm:"index"`
}

func (GormStockAdjustmentPending) TableName() string { return "stock_adjustments_pending" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.StockAdjustmentRepository = (*GormStockAdjustmentRepository)(nil)

type GormStockAdjustmentRepository struct {
	db          *gorm.DB
	productRepo ports.ProductRepository
}

func NewGormStockAdjustmentRepository(db *gorm.DB, productRepo ports.ProductRepository) *GormStockAdjustmentRepository {
	return &GormStockAdjustmentRepository{db: db, productRepo: productRepo}
}

func (r *GormStockAdjustmentRepository) CreateStockAdjustmentPending(
	ctx context.Context, adjustment *domain.StockAdjustmentPending,
) error {
	return r.CreateStockAdjustmentPendingWithTx(ctx, adjustment, nil)
}

// CreateStockAdjustmentPendingWithTx cria um ajuste de estoque pendente, opcionalmente dentro de uma transação
// Se txDB for fornecido, usa o DB da transação; senão usa o DB padrão
func (r *GormStockAdjustmentRepository) CreateStockAdjustmentPendingWithTx(
	ctx context.Context, adjustment *domain.StockAdjustmentPending, txDB *gorm.DB,
) error {
	log.Printf("[STOCK_REPO] CreateStockAdjustmentPendingWithTx chamado: order_id=%d, ingredient_id=%d, quantity=%.4f", adjustment.OrderID, adjustment.IngredientID, adjustment.Quantity)
	db := r.db
	if txDB != nil {
		db = txDB.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	gAdjustment := GormStockAdjustmentPending{
		OrderID:        adjustment.OrderID,
		IngredientID:   adjustment.IngredientID,
		Quantity:       adjustment.Quantity,
		OrderStatus:    adjustment.OrderStatus,
		Status:         string(adjustment.Status),
		IngredientName: adjustment.IngredientName, // snapshot
		IngredientUnit: adjustment.IngredientUnit, // snapshot
	}
	if err := db.Create(&gAdjustment).Error; err != nil {
		log.Printf("[STOCK_REPO] Erro ao criar ajuste: %v", err)
		return fmt.Errorf("CreateStockAdjustmentPendingWithTx: %w", err)
	}
	log.Printf("[STOCK_REPO] Ajuste criado com sucesso: id=%d", gAdjustment.ID)
	adjustment.ID = gAdjustment.ID
	adjustment.CreatedAt = time.Unix(gAdjustment.CreatedAt, 0)
	return nil
}

func (r *GormStockAdjustmentRepository) FindPendingByOrderID(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []GormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("order_id = ? AND status = ? AND deleted_at IS NULL", orderID, domain.StockAdjustmentStatusPending).
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindPendingByOrderID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) FindByOrderID(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []GormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindByOrderID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) FindPendingByIngredientID(
	ctx context.Context, ingredientID uint,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []GormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("ingredient_id = ? AND status = ? AND deleted_at IS NULL", ingredientID, domain.StockAdjustmentStatusPending).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("FindPendingByIngredientID: %w", err)
	}
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) ListPending(
	ctx context.Context,
) ([]domain.StockAdjustmentPending, error) {
	var gAdjustments []GormStockAdjustmentPending
	if err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", domain.StockAdjustmentStatusPending).
		Order("created_at desc").
		Find(&gAdjustments).Error; err != nil {
		return nil, fmt.Errorf("ListPending: %w", err)
	}
	log.Printf("[STOCK_REPO] ListPending retornou %d registros", len(gAdjustments))
	return r.mapToDomainSlice(gAdjustments), nil
}

func (r *GormStockAdjustmentRepository) UpdateStatus(
	ctx context.Context, id uint, status domain.StockAdjustmentStatus,
) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status": string(status),
	}
	if status != domain.StockAdjustmentStatusPending {
		updates["processed_at"] = now
	}
	if err := r.db.WithContext(ctx).
		Model(&GormStockAdjustmentPending{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}
	return nil
}

func (r *GormStockAdjustmentRepository) Approve(
	ctx context.Context, id uint, processedBy uint, notes string,
) error {
	return r.approveWithTx(ctx, id, processedBy, notes, nil)
}

func (r *GormStockAdjustmentRepository) approveWithTx(
	ctx context.Context, id uint, processedBy uint, notes string, txDB *gorm.DB,
) error {
	db := r.db
	if txDB != nil {
		db = txDB.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status":           string(domain.StockAdjustmentStatusApproved),
		"processed_at":     now,
		"processed_by":     processedBy,
		"processing_notes": notes,
	}
	result := db.
		Model(&GormStockAdjustmentPending{}).
		Where("id = ? AND status = ?", id, domain.StockAdjustmentStatusPending).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("approveWithTx: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ajuste id=%d não encontrado ou já processado", id)
	}
	return nil
}

// ApproveAndRestoreStock aprova o ajuste e repõe estoque do ingrediente em transação atômica.
func (r *GormStockAdjustmentRepository) ApproveAndRestoreStock(
	ctx context.Context, id uint, processedBy uint, notes string,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var gAdjustment GormStockAdjustmentPending
		if err := tx.Where("deleted_at IS NULL").First(&gAdjustment, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("ajuste id=%d não encontrado", id)
			}
			return fmt.Errorf("ApproveAndRestoreStock: buscar ajuste: %w", err)
		}
		if gAdjustment.Status != string(domain.StockAdjustmentStatusPending) {
			return fmt.Errorf("ajuste já processado (status: %s)", gAdjustment.Status)
		}

		if err := r.approveWithTx(ctx, id, processedBy, notes, tx); err != nil {
			return err
		}

		if err := r.productRepo.IncreaseIngredientStock(ctx, gAdjustment.IngredientID, gAdjustment.Quantity, tx); err != nil {
			return fmt.Errorf("ApproveAndRestoreStock: repor estoque: %w", err)
		}

		return nil
	})
}

func (r *GormStockAdjustmentRepository) Reject(
	ctx context.Context, id uint, processedBy uint, notes string,
) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"status":           string(domain.StockAdjustmentStatusRejected),
		"processed_at":     now,
		"processed_by":     processedBy,
		"processing_notes": notes,
	}
	if err := r.db.WithContext(ctx).
		Model(&GormStockAdjustmentPending{}).
		Where("id = ? AND status = ?", id, domain.StockAdjustmentStatusPending).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("Reject: %w", err)
	}
	return nil
}

func (r *GormStockAdjustmentRepository) FindByID(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	var gAdjustment GormStockAdjustmentPending
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&gAdjustment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return r.mapToDomain(&gAdjustment), nil
}

// ─── Mappers ──────────────────────────────────────────────────────────────────

func (r *GormStockAdjustmentRepository) mapToDomain(g *GormStockAdjustmentPending) *domain.StockAdjustmentPending {
	var deletedAt *time.Time
	if g.DeletedAt != nil {
		dt := time.Unix(*g.DeletedAt, 0)
		deletedAt = &dt
	}
	adjustment := &domain.StockAdjustmentPending{
		ID:              g.ID,
		OrderID:         g.OrderID,
		IngredientID:    g.IngredientID,
		Quantity:        g.Quantity,
		OrderStatus:     g.OrderStatus,
		Status:          domain.StockAdjustmentStatus(g.Status),
		CreatedAt:       time.Unix(g.CreatedAt, 0),
		ProcessedBy:     g.ProcessedBy,
		ProcessingNotes: g.ProcessingNotes,
		IngredientName:  g.IngredientName, // snapshot
		IngredientUnit:  g.IngredientUnit, // snapshot
		DeletedAt:       deletedAt,
	}
	if g.ProcessedAt != nil {
		processedAt := time.Unix(*g.ProcessedAt, 0)
		adjustment.ProcessedAt = &processedAt
	}
	return adjustment
}

func (r *GormStockAdjustmentRepository) mapToDomainSlice(gAdjustments []GormStockAdjustmentPending) []domain.StockAdjustmentPending {
	out := make([]domain.StockAdjustmentPending, len(gAdjustments))
	for i, g := range gAdjustments {
		out[i] = *r.mapToDomain(&g)
	}
	return out
}
