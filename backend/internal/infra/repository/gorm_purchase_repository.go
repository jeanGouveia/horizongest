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
		return fmt.Errorf("PurchaseRepository.CreateSupplier: %w", err)
	}
	supplier.CompanyID = companyID
	if err := r.db.WithContext(ctx).Create(supplier).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.CreateSupplier: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) ListSuppliers(ctx context.Context, companyID uint, activeOnly bool, limit, offset int) ([]domain.Supplier, error) {
	var suppliers []domain.Supplier
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if activeOnly {
		query = query.Where("active = ?", true)
	}

	query = query.Order("name ASC").Limit(limit).Offset(offset)

	err := query.Find(&suppliers).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.ListSuppliers: %w", err)
	}
	return suppliers, nil
}

func (r *GormPurchaseRepository) GetSupplierByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	var supplier domain.Supplier
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&supplier).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.GetSupplierByID: %w", err)
	}
	return &supplier, nil
}

func (r *GormPurchaseRepository) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	query := ApplyTenantFilterWithID(ctx, r.db, supplier.ID)
	if err := query.WithContext(ctx).Save(supplier).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.UpdateSupplier: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) DeleteSupplier(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.Delete(&domain.Supplier{}).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.DeleteSupplier: %w", err)
	}
	return nil
}

// --- Pedidos de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("PurchaseRepository.CreatePurchaseOrder: %w", err)
	}
	order.CompanyID = companyID
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.CreatePurchaseOrder: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) ListPurchaseOrders(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.PurchaseOrder, error) {
	var orders []domain.PurchaseOrder
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	err := query.Preload("Supplier").Preload("Items.Ingredient").Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.ListPurchaseOrders: %w", err)
	}
	return orders, nil
}

func (r *GormPurchaseRepository) GetPurchaseOrderByID(ctx context.Context, id uint) (*domain.PurchaseOrder, error) {
	var order domain.PurchaseOrder
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Supplier").
		Preload("Items.Ingredient").
		First(&order).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.GetPurchaseOrderByID: %w", err)
	}
	return &order, nil
}

func (r *GormPurchaseRepository) UpdatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	query := ApplyTenantFilterWithID(ctx, r.db, order.ID)
	if err := query.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.UpdatePurchaseOrder: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderStatus(ctx context.Context, id uint, status domain.PurchaseOrderStatus) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&domain.PurchaseOrder{}).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.UpdatePurchaseOrderStatus: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) DeletePurchaseOrder(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.Delete(&domain.PurchaseOrder{}).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.DeletePurchaseOrder: %w", err)
	}
	return nil
}

// --- Itens de Pedido de Compra ---

func (r *GormPurchaseRepository) CreatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.CreatePurchaseOrderItem: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) ListPurchaseOrderItems(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseOrderItem, error) {
	var items []domain.PurchaseOrderItem
	query := ApplyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).
		Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Ingredient").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.ListPurchaseOrderItems: %w", err)
	}
	return items, nil
}

func (r *GormPurchaseRepository) UpdatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error {
	query := ApplyTenantFilterWithID(ctx, r.db, item.ID)
	if err := query.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.UpdatePurchaseOrderItem: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) DeletePurchaseOrderItem(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Delete(&domain.PurchaseOrderItem{}).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.DeletePurchaseOrderItem: %w", err)
	}
	return nil
}

// --- Recebimentos ---

func (r *GormPurchaseRepository) CreatePurchaseReceiving(ctx context.Context, receiving *domain.PurchaseReceiving) error {
	// PurchaseReceiving não tem CompanyID direto - o tenant é herdado através de PurchaseOrder
	// A validação de tenant deve ser feita no service layer verificando o PurchaseOrderID
	if err := r.db.WithContext(ctx).Create(receiving).Error; err != nil {
		return fmt.Errorf("PurchaseRepository.CreatePurchaseReceiving: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) GetPurchaseReceivingByID(ctx context.Context, id uint) (*domain.PurchaseReceiving, error) {
	var receiving domain.PurchaseReceiving
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Items.Ingredient").
		First(&receiving).Error
	if err != nil {
		return nil, fmt.Errorf("PurchaseRepository.GetPurchaseReceivingByID: %w", err)
	}
	return &receiving, nil
}

func (r *GormPurchaseRepository) ListPurchaseReceivings(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseReceiving, error) {
	var receivings []domain.PurchaseReceiving
	query := ApplyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).
		Where("purchase_order_id = ? AND deleted_at IS NULL", purchaseOrderID).
		Preload("Items.Ingredient").
		Find(&receivings).Error
	if err != nil {
		return nil, fmt.Errorf("ListPurchaseReceivings: %w", err)
	}
	return receivings, nil
}

func (r *GormPurchaseRepository) DeletePurchaseReceiving(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.Delete(&domain.PurchaseReceiving{}).Error; err != nil {
		return fmt.Errorf("DeletePurchaseReceiving: %w", err)
	}
	return nil
}

// --- Itens de Recebimento ---

func (r *GormPurchaseRepository) CreatePurchaseReceivingItem(ctx context.Context, item *domain.PurchaseReceivingItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("CreatePurchaseReceivingItem: %w", err)
	}
	return nil
}

func (r *GormPurchaseRepository) ListPurchaseReceivingItems(ctx context.Context, receivingID uint) ([]domain.PurchaseReceivingItem, error) {
	var items []domain.PurchaseReceivingItem
	err := r.db.WithContext(ctx).Where("purchase_receiving_id = ? AND deleted_at IS NULL", receivingID).
		Preload("Ingredient").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("ListPurchaseReceivingItems: %w", err)
	}
	return items, nil
}
