package repository

import (
	"context"

	"gorm.io/gorm"

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

// --- Movimentações ---

func (r *GormStockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	return r.db.WithContext(ctx).Create(movement).Error
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

func (r *GormStockMovementRepository) CreateInventory(ctx context.Context, inventory *domain.StockInventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *GormStockMovementRepository) GetInventoryByID(ctx context.Context, id uint) (*domain.StockInventory, error) {
	var inventory domain.StockInventory
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
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

func (r *GormStockMovementRepository) UpdateInventoryStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&domain.StockInventory{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *GormStockMovementRepository) DeleteInventory(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.StockInventory{}).Error
}

// --- Itens de Inventário ---

func (r *GormStockMovementRepository) CreateInventoryItem(ctx context.Context, item *domain.StockInventoryItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormStockMovementRepository) ListInventoryItems(ctx context.Context, inventoryID uint) ([]domain.StockInventoryItem, error) {
	var items []domain.StockInventoryItem
	err := r.db.WithContext(ctx).Where("inventory_id = ? AND deleted_at IS NULL", inventoryID).
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
