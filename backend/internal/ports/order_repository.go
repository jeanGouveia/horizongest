package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type OrderRepository interface {
	// CreateOrder persiste o pedido, seus itens e executa a baixa de estoque
	// de todos os ingredientes em uma única transação atômica.
	// Se qualquer ingrediente ficar com estoque negativo → rollback total.
	// productIngredients é um mapa de product_id -> ingredientes pré-carregados
	// Sprint 4C: Retorna o pedido criado ou o pedido existente em caso de colisão de idempotency_key
	CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) (*domain.Order, error)

	// FindByIdempotencyKey busca um pedido pela chave de idempotência
	// Retorna nil se não encontrado (não é erro)
	FindByIdempotencyKey(ctx context.Context, companyID uint, idempotencyKey string) (*domain.Order, error)

	FindOrderByID(ctx context.Context, id uint) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error

	// ValidateStock verifica se há estoque suficiente para os itens do pedido
	ValidateStock(ctx context.Context, items []domain.OrderItem, productIngredients map[uint][]domain.ProductIngredient) (*domain.StockValidationResponse, error)

	// UpdateOrderStatusWithAdjustments atualiza o status do pedido e registra ajustes de estoque
	// em uma única transação atômica. Se o status for 'cancelled', registra ajustes pendentes.
	// Se qualquer etapa falhar, rollback completo. Garante consistência entre status e auditoria.
	UpdateOrderStatusWithAdjustments(
		ctx context.Context,
		id uint,
		status domain.OrderStatus,
		productIngredients map[uint][]domain.ProductIngredient,
		orderItems []domain.OrderItem,
	) error

	// UpdateOrder atualiza os itens e notas de um pedido existente
	// Ajusta o estoque automaticamente para refletir as mudanças nos itens
	// productIngredients é um mapa de product_id -> ingredientes pré-carregados
	UpdateOrder(
		ctx context.Context,
		id uint,
		items []domain.OrderItem,
		total domain.Money,
		notes string,
		productIngredients map[uint][]domain.ProductIngredient,
	) error
}
