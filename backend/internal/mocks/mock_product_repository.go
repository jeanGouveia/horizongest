package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/gorm"
)

// MockProductRepository is a mock implementation of ports.ProductRepository
type MockProductRepository struct {
	Products           map[uint]*domain.Product
	Ingredients        map[uint]*domain.Ingredient
	ProductIngredients map[uint][]domain.ProductIngredient
	CreateError        error
	FindError          error
	UpdateError        error
	DeleteError        error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products:           make(map[uint]*domain.Product),
		Ingredients:        make(map[uint]*domain.Ingredient),
		ProductIngredients: make(map[uint][]domain.ProductIngredient),
	}
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if p.ID == 0 {
		p.ID = uint(len(m.Products) + 1)
	}
	m.Products[p.ID] = p
	return nil
}

func (m *MockProductRepository) FindProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Products[id], nil
}

func (m *MockProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	for _, p := range m.Products {
		products = append(products, *p)
	}
	return products, nil
}

func (m *MockProductRepository) ListActiveProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	for _, p := range m.Products {
		if p.Active {
			products = append(products, *p)
		}
	}
	return products, nil
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Products[p.ID] = p
	return nil
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Products, id)
	return nil
}

func (m *MockProductRepository) CanDeleteProduct(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	return &domain.DependencyCheck{CanDelete: true}, nil
}

func (m *MockProductRepository) CreateIngredient(ctx context.Context, i *domain.Ingredient) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if i.ID == 0 {
		i.ID = uint(len(m.Ingredients) + 1)
	}
	m.Ingredients[i.ID] = i
	return nil
}

func (m *MockProductRepository) FindIngredientByID(ctx context.Context, id uint) (*domain.Ingredient, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Ingredients[id], nil
}

func (m *MockProductRepository) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	for _, i := range m.Ingredients {
		ingredients = append(ingredients, *i)
	}
	return ingredients, nil
}

func (m *MockProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Ingredients[i.ID] = i
	return nil
}

func (m *MockProductRepository) DeleteIngredient(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Ingredients, id)
	return nil
}

func (m *MockProductRepository) CanDeleteIngredient(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	return &domain.DependencyCheck{CanDelete: true}, nil
}

func (m *MockProductRepository) SetProductIngredients(ctx context.Context, productID uint, items []domain.ProductIngredient) error {
	m.ProductIngredients[productID] = items
	return nil
}

func (m *MockProductRepository) GetProductIngredients(ctx context.Context, productID uint) ([]domain.ProductIngredient, error) {
	return m.ProductIngredients[productID], nil
}

func (m *MockProductRepository) DecreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB, ingredientName string, currentStock float64) error {
	if ing, ok := m.Ingredients[ingredientID]; ok {
		ing.StockQuantity -= qty
	}
	return nil
}

func (m *MockProductRepository) IncreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB) error {
	if ing, ok := m.Ingredients[ingredientID]; ok {
		ing.StockQuantity += qty
	}
	return nil
}
