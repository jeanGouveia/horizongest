package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

// StockMovementRepository define a interface para repositório de movimentações de estoque
type StockMovementRepository interface {
	// Movimentações
	Create(ctx context.Context, movement *domain.StockMovement) error
	List(ctx context.Context, companyID uint, ingredientID *uint, limit, offset int) ([]domain.StockMovement, error)
	GetByID(ctx context.Context, id uint) (*domain.StockMovement, error)
	Delete(ctx context.Context, id uint) error
	
	// Inventários
	CreateInventory(ctx context.Context, inventory *domain.StockInventory) error
	GetInventoryByID(ctx context.Context, id uint) (*domain.StockInventory, error)
	ListInventories(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.StockInventory, error)
	UpdateInventoryStatus(ctx context.Context, id uint, status string) error
	DeleteInventory(ctx context.Context, id uint) error
	
	// Itens de Inventário
	CreateInventoryItem(ctx context.Context, item *domain.StockInventoryItem) error
	ListInventoryItems(ctx context.Context, inventoryID uint) ([]domain.StockInventoryItem, error)
	DeleteInventoryItem(ctx context.Context, id uint) error
	
	// Histórico
	GetMovementHistory(ctx context.Context, companyID uint, ingredientID uint, limit int) ([]domain.StockMovement, error)
}
