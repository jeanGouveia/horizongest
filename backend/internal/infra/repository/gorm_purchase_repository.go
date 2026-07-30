package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
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
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateSupplier: %w", err)
	}
	supplier.CompanyID = companyID
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
	return suppliers, fmt.Errorf("ListSuppliers: %w", err)
}

func (r *GormPurchaseRepository) GetSupplierByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	var supplier domain.Supplier
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&supplier).Error
	return &supplier, fmt.Errorf("GetSupplierByID: %w", err)
}

func (r *GormPurchaseRepository) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	query := ApplyTenantFilterWithID(ctx, r.db, supplier.ID)
	return query.WithContext(ctx).Save(supplier).Error
}

func (r *GormPurchaseRepository) DeleteSupplier(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.Supplier{}).Error
}

// --- Pedidos de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreatePurchaseOrder: %w", err)
	}
	order.CompanyID = companyID
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
	return orders, fmt.Errorf("ListPurchaseOrders: %w", err)
}

func (r *GormPurchaseRepository) GetPurchaseOrderByID(ctx context.Context, id uint) (*domain.PurchaseOrder, error) {
	var order domain.PurchaseOrder
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Supplier").
		Preload("Items.Ingredient").
		First(&order).Error
	return &order, fmt.Errorf("GetPurchaseOrderByID: %w", err)
}

func (r *GormPurchaseRepository) UpdatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	query := ApplyTenantFilterWithID(ctx, r.db, order.ID)
	return query.WithContext(ctx).Save(order).Error
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderStatus(ctx context.Context, id uint, status domain.PurchaseOrderStatus) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.WithContext(ctx).Model(&domain.PurchaseOrder{}).
		Update("status", status).Error
}

func (r *GormPurchaseRepository) DeletePurchaseOrder(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.PurchaseOrder{}).Error
}

// --- Itens de Pedido de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormPurchaseRepository) ListPurchaseOrderItems(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseOrderItem, error) {
	var items []domain.PurchaseOrderItem
	query := ApplyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).
		Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Ingredient").
		Find(&items).Error
	return items, err
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	query := ApplyTenantFilterWithID(ctx, r.db, item.ID)
	return query.WithContext(ctx).Save(item).Error
}

func (r *GormPurchaseRepository) DeletePurchaseOrderItem(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.WithContext(ctx).Delete(&domain.PurchaseOrderItem{}).Error
}

// --- Recebimentos ---

func (r *GormPurchaseRepository) CreatePurchaseReceiving(ctx context.Context, receiving *domain.PurchaseReceiving) error {
	// PurchaseReceiving não tem CompanyID direto - o tenant é herdado através de PurchaseOrder
	// A validação de tenant deve ser feita no service layer verificando o PurchaseOrderID
	return r.db.WithContext(ctx).Create(receiving).Error
}

func (r *GormPurchaseRepository) GetPurchaseReceivingByID(ctx context.Context, id uint) (*domain.PurchaseReceiving, error) {
	var receiving domain.PurchaseReceiving
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Items.Ingredient").
		First(&receiving).Error
	return &receiving, fmt.Errorf("GetPurchaseReceivingByID: %w", err)
}

func (r *GormPurchaseRepository) ListPurchaseReceivings(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseReceiving, error) {
	var receivings []domain.PurchaseReceiving
	query := ApplyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).
		Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Items.Ingredient").
		Find(&receivings).Error
	return receivings, fmt.Errorf("ListPurchaseReceivings: %w", err)
}

func (r *GormPurchaseRepository) DeletePurchaseReceiving(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.PurchaseReceiving{}).Error
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
