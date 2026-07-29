package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var _ ports.StockMovementRepository = (*GormStockMovementRepository)(nil)

type GormStockMovementRepository struct {
	db *gorm.DB
}

func NewGormStockMovementRepository(db *gorm.DB) *GormStockMovementRepository {
	return &GormStockMovementRepository{db: db}
}

// getDB retorna a transação se fornecida, senão retorna o DB padrão
// Sprint 4B.1 v2: Transaction propagation
func (r *GormStockMovementRepository) getDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

// --- Movimentações ---

func (r *GormStockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement, tx *gorm.DB) error {
	return r.getDB(ctx, tx).Create(movement).Error
}

func (r *GormStockMovementRepository) List(ctx context.Context, companyID uint, ingredientID *uint, limit, offset int) ([]domain.StockMovement, error) {
	var movements []domain.StockMovement
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if ingredientID != nil {
		query = query.Where("ingredient_id = ?", *ingredientID)
	}

	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	err := query.Preload("Ingredient").Preload("Performer").Find(&movements).Error
	return movements, err
}

func (r *GormStockMovementRepository) GetByID(ctx context.Context, id uint) (*domain.StockMovement, error) {
	var movement domain.StockMovement
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Ingredient").Preload("Performer").
		First(&movement).Error
	return &movement, err
}

func (r *GormStockMovementRepository) Delete(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.StockMovement{}).Error
}

// --- Inventários ---

func (r *GormStockMovementRepository) CreateInventory(ctx context.Context, inventory *domain.StockInventory, tx *gorm.DB) error {
	return r.getDB(ctx, tx).Create(inventory).Error
}

func (r *GormStockMovementRepository) GetInventoryByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error) {
	var inventory domain.StockInventory
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	err := query.Where("deleted_at IS NULL").
		Preload("Items.Ingredient").
		First(&inventory).Error
	return &inventory, err
}

// Sprint 4B.5: FindInventoryByIDForUpdate busca inventário com SELECT FOR UPDATE
// Isso previne double completion e modificações concorrentes durante CompleteInventory
func (r *GormStockMovementRepository) FindInventoryByIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error) {
	var inventory domain.StockInventory
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	err := query.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("deleted_at IS NULL").
		Preload("Items.Ingredient").
		First(&inventory).Error
	return &inventory, err
}

func (r *GormStockMovementRepository) ListInventories(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.StockInventory, error) {
	var inventories []domain.StockInventory
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	err := query.Find(&inventories).Error
	return inventories, err
}

func (r *GormStockMovementRepository) UpdateInventoryStatus(ctx context.Context, id uint, status string, tx *gorm.DB) error {
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	result := query.Model(&domain.StockInventory{}).
		Where("status = ?", "draft").
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("inventory already completed or not found")
	}
	return nil
}

func (r *GormStockMovementRepository) DeleteInventory(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.StockInventory{}).Error
}

// --- Itens de Inventário ---

func (r *GormStockMovementRepository) CreateInventoryItem(ctx context.Context, item *domain.StockInventoryItem, tx *gorm.DB) error {
	return r.getDB(ctx, tx).Create(item).Error
}

func (r *GormStockMovementRepository) ListInventoryItems(ctx context.Context, inventoryID uint, tx *gorm.DB) ([]domain.StockInventoryItem, error) {
	// Sprint 4A: NOTA - Query sem filtro de tenant explícito
	// NÃO é IDOR pois é chamado apenas após validação de inventory_id
	// O inventoryID deve ser validado pelo service caller antes de chamar este método
	// O service deve garantir que o inventory pertence ao tenant do usuário
	var items []domain.StockInventoryItem
	err := r.getDB(ctx, tx).Where("inventory_id = ? AND deleted_at IS NULL", inventoryID).
		Preload("Ingredient").
		Find(&items).Error
	return items, err
}

func (r *GormStockMovementRepository) DeleteInventoryItem(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.StockInventoryItem{}).Error
}

// --- Histórico ---

func (r *GormStockMovementRepository) GetMovementHistory(ctx context.Context, companyID uint, ingredientID uint, limit int) ([]domain.StockMovement, error) {
	var movements []domain.StockMovement
	err := r.db.WithContext(ctx).Where("company_id = ? AND ingredient_id = ? AND deleted_at IS NULL", companyID, ingredientID).
		Order("created_at DESC").
		Limit(limit).
		Preload("Performer").
		Find(&movements).Error
	return movements, err
}
