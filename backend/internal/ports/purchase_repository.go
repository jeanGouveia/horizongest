package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

// PurchaseRepository define a interface para repositório de compras
type PurchaseRepository interface {
	// Fornecedores
	CreateSupplier(ctx context.Context, supplier *domain.Supplier) error
	ListSuppliers(ctx context.Context, companyID uint, activeOnly bool, limit, offset int) ([]domain.Supplier, error)
	GetSupplierByID(ctx context.Context, id uint) (*domain.Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error
	DeleteSupplier(ctx context.Context, id uint) error
	
	// Pedidos de Compra
	CreatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error
	ListPurchaseOrders(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.PurchaseOrder, error)
	GetPurchaseOrderByID(ctx context.Context, id uint) (*domain.PurchaseOrder, error)
	UpdatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error
	UpdatePurchaseOrderStatus(ctx context.Context, id uint, status domain.PurchaseOrderStatus) error
	DeletePurchaseOrder(ctx context.Context, id uint) error
	
	// Itens de Pedido de Compra
	CreatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error
	ListPurchaseOrderItems(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseOrderItem, error)
	UpdatePurchaseOrderItem(ctx context.Context, item *domain.PurchaseOrderItem) error
	DeletePurchaseOrderItem(ctx context.Context, id uint) error
	
	// Recebimentos
	CreatePurchaseReceiving(ctx context.Context, receiving *domain.PurchaseReceiving) error
	GetPurchaseReceivingByID(ctx context.Context, id uint) (*domain.PurchaseReceiving, error)
	ListPurchaseReceivings(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseReceiving, error)
	DeletePurchaseReceiving(ctx context.Context, id uint) error
	
	// Itens de Recebimento
	CreatePurchaseReceivingItem(ctx context.Context, item *domain.PurchaseReceivingItem) error
	ListPurchaseReceivingItems(ctx context.Context, receivingID uint) ([]domain.PurchaseReceivingItem, error)
}
