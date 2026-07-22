package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type StockAdjustmentRepository interface {
	// CreateStockAdjustmentPending registra um ajuste de estoque pendente
	// Usado quando um pedido é cancelado para registrar quais ingredientes
	// poderiam ser devolvidos ao estoque após análise manual
	CreateStockAdjustmentPending(ctx context.Context, adjustment *domain.StockAdjustmentPending) error

	// FindPendingByOrderID busca todos os ajustes pendentes para um pedido específico
	FindPendingByOrderID(ctx context.Context, orderID uint) ([]domain.StockAdjustmentPending, error)

	// FindByOrderID busca todos os ajustes (todos os status) para um pedido específico
	FindByOrderID(ctx context.Context, orderID uint) ([]domain.StockAdjustmentPending, error)

	// FindPendingByIngredientID busca ajustes pendentes para um ingrediente específico
	FindPendingByIngredientID(ctx context.Context, ingredientID uint) ([]domain.StockAdjustmentPending, error)

	// ListPending busca todos os ajustes pendentes (para dashboard de aprovação)
	ListPending(ctx context.Context) ([]domain.StockAdjustmentPending, error)

	// UpdateStatus atualiza o status de um ajuste (pending → approved/rejected)
	UpdateStatus(ctx context.Context, id uint, status domain.StockAdjustmentStatus) error

	// Approve aprova um ajuste pendente, registrando quem aprovou e observações
	Approve(ctx context.Context, id uint, processedBy uint, notes string) error

	// ApproveAndRestoreStock aprova um ajuste pendente e repõe o estoque do ingrediente em transação atômica
	ApproveAndRestoreStock(ctx context.Context, id uint, processedBy uint, notes string) error

	// Reject rejeita um ajuste pendente, registrando quem rejeitou e observações
	Reject(ctx context.Context, id uint, processedBy uint, notes string) error

	// FindByID busca um ajuste por ID
	FindByID(ctx context.Context, id uint) (*domain.StockAdjustmentPending, error)
}
