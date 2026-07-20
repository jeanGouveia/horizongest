package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var _ ports.PurchaseRepository = (*GormPurchaseRepository)(nil)

type GormPurchaseRepository struct {
	db *gorm.DB
}

func NewGormPurchaseRepository(db *gorm.DB) *GormPurchaseRepository {
	return &GormPurchaseRepository{db: db}
}

// --- Fornecedores ---

func (r *GormPurchaseRepository) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

func (r *GormPurchaseRepository) ListSuppliers(ctx context.Context, companyID uint, activeOnly bool, limit, offset int) ([]domain.Supplier, error) {
	var suppliers []domain.Supplier
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)
	
	if activeOnly {
		query = query.Where("active = ?", true)
	}
	
	query = query.Order("name ASC").Limit(limit).Offset(offset)
	
	err := query.Find(&suppliers).Error
	return suppliers, err
}

func (r *GormPurchaseRepository) GetSupplierByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	var supplier domain.Supplier
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&supplier).Error
	return &supplier, err
}

func (r *GormPurchaseRepository) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

func (r *GormPurchaseRepository) DeleteSupplier(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Supplier{}).Error
}

// --- Pedidos de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *GormPurchaseRepository) ListPurchaseOrders(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.PurchaseOrder, error) {
	var orders []domain.PurchaseOrder
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	query = query.Order("created_at DESC").Limit(limit).Offset(offset)
	
	err := query.Preload("Supplier").Preload("Items.Ingredient").Find(&orders).Error
	return orders, err
}

func (r *GormPurchaseRepository) GetPurchaseOrderByID(ctx context.Context, id uint) (*domain.PurchaseOrder, error) {
	var order domain.PurchaseOrder
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).
		Preload("Supplier").
		Preload("Items.Ingredient").
		First(&order).Error
	return &order, err
}

func (r *GormPurchaseRepository) UpdatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderStatus(ctx context.Context, id uint, status domain.PurchaseOrderStatus) error {
	return r.db.WithContext(ctx).Model(&domain.PurchaseOrder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *GormPurchaseRepository) DeletePurchaseOrder(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PurchaseOrder{}).Error
}

// --- Itens de Pedido de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormPurchaseRepository) ListPurchaseOrderItems(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseOrderItem, error) {
	var items []domain.PurchaseOrderItem
	err := r.db.WithContext(ctx).Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Ingredient").
		Find(&items).Error
	return items, err
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormPurchaseRepository) DeletePurchaseOrderItem(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PurchaseOrderItem{}).Error
}

// --- Recebimentos ---

func (r *GormPurchaseRepository) CreatePurchaseReceiving(ctx context.Context, receiving *domain.PurchaseReceiving) error {
	return r.db.WithContext(ctx).Create(receiving).Error
}

func (r *GormPurchaseRepository) GetPurchaseReceivingByID(ctx context.Context, id uint) (*domain.PurchaseReceiving, error) {
	var receiving domain.PurchaseReceiving
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).
		Preload("Items.Ingredient").
		First(&receiving).Error
	return &receiving, err
}

func (r *GormPurchaseRepository) ListPurchaseReceivings(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseReceiving, error) {
	var receivings []domain.PurchaseReceiving
	err := r.db.WithContext(ctx).Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Items.Ingredient").
		Find(&receivings).Error
	return receivings, err
}

func (r *GormPurchaseRepository) DeletePurchaseReceiving(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.PurchaseReceiving{}).Error
}

// --- Itens de Recebimento ---

func (r *GormPurchaseRepository) CreatePurchaseReceivingItem(ctx context.Context, item *domain.PurchaseReceivingItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormPurchaseRepository) ListPurchaseReceivingItems(ctx context.Context, receivingID uint) ([]domain.PurchaseReceivingItem, error) {
	var items []domain.PurchaseReceivingItem
	err := r.db.WithContext(ctx).Where("purchase_receiving_id = ? AND deleted_at IS NULL", receivingID).
		Preload("Ingredient").
		Find(&items).Error
	return items, err
}
