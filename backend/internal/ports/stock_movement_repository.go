package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/gorm"
)

// StockMovementRepository define a interface para repositório de movimentações de estoque
type StockMovementRepository interface {
	// Movimentações
	Create(ctx context.Context, movement *domain.StockMovement, tx *gorm.DB) error
	List(ctx context.Context, companyID uint, ingredientID *uint, limit, offset int) ([]domain.StockMovement, error)
	GetByID(ctx context.Context, id uint) (*domain.StockMovement, error)
	Delete(ctx context.Context, id uint) error

	// Inventários
	CreateInventory(ctx context.Context, inventory *domain.StockInventory, tx *gorm.DB) error
	GetInventoryByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error)
	FindInventoryByIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error) // Sprint 4B.5: SELECT FOR UPDATE
	ListInventories(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.StockInventory, error)
	UpdateInventoryStatus(ctx context.Context, id uint, status string, tx *gorm.DB) error
	DeleteInventory(ctx context.Context, id uint) error

	// Itens de Inventário
	CreateInventoryItem(ctx context.Context, item *domain.StockInventoryItem, tx *gorm.DB) error
	ListInventoryItems(ctx context.Context, inventoryID uint, tx *gorm.DB) ([]domain.StockInventoryItem, error)
	DeleteInventoryItem(ctx context.Context, id uint) error

	// Histórico
	GetMovementHistory(ctx context.Context, companyID uint, ingredientID uint, limit int) ([]domain.StockMovement, error)
}
