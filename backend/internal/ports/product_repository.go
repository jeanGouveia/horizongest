package ports

import (
	"context"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type ProductRepository interface {
	// Produto
	CreateProduct(ctx context.Context, p *domain.Product) error
	FindProductByID(ctx context.Context, id uint) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	ListActiveProducts(ctx context.Context) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, p *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	CanDeleteProduct(ctx context.Context, id uint) (*domain.DependencyCheck, error)

	// Ingrediente
	CreateIngredient(ctx context.Context, i *domain.Ingredient) error
	FindIngredientByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error)
	FindIngredientByIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error)
	ListIngredients(ctx context.Context) ([]domain.Ingredient, error)
	UpdateIngredient(ctx context.Context, i *domain.Ingredient, tx *gorm.DB) error
	DeleteIngredient(ctx context.Context, id uint) error
	CanDeleteIngredient(ctx context.Context, id uint) (*domain.DependencyCheck, error)

	// Ficha técnica (produto ↔ ingredientes)
	SetProductIngredients(ctx context.Context, productID uint, items []domain.ProductIngredient) error
	GetProductIngredients(ctx context.Context, productID uint) ([]domain.ProductIngredient, error)

	// Estoque — chamado dentro de transação pelo OrderRepository
	// Aceita um DB opcional para transações
	DecreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB, ingredientName string, currentStock float64) error
	IncreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB) error
}
