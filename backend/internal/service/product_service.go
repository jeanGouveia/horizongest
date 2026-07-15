package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var (
	ErrProductNotFound    = errors.New("produto não encontrado")
	ErrIngredientNotFound = errors.New("ingrediente não encontrado")
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type CreateProductInput struct {
	Name        string  `json:"name"        validate:"required,min=2,max=120"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       validate:"required,gt=0"`
	IsComposto  bool    `json:"is_composto"`
}

type UpdateProductInput struct {
	Name        string  `json:"name"        validate:"required,min=2,max=120"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       validate:"required,gt=0"`
	IsComposto  bool    `json:"is_composto"`
	Active      bool    `json:"active"`
}

type CreateIngredientInput struct {
	Name          string  `json:"name"           validate:"required,min=2,max=120"`
	Unit          string  `json:"unit"           validate:"required,oneof=kg g L ml un"`
	StockQuantity float64 `json:"stock_quantity" validate:"gte=0"`
	MinStock      float64 `json:"min_stock"      validate:"gte=0"`
}

type UpdateIngredientInput struct {
	Name          string  `json:"name"           validate:"required,min=2,max=120"`
	Unit          string  `json:"unit"           validate:"required,oneof=kg g L ml un"`
	StockQuantity float64 `json:"stock_quantity" validate:"gte=0"`
	MinStock      float64 `json:"min_stock"      validate:"gte=0"`
}

type ProductIngredientInput struct {
	IngredientID uint    `json:"ingredient_id" validate:"required"`
	Quantity     float64 `json:"quantity"      validate:"required,gt=0"`
}

type SetProductIngredientsInput struct {
	Items []ProductIngredientInput `json:"items" validate:"required,dive"`
}

type UpdateStockInput struct {
	Quantity float64 `json:"quantity" validate:"required,gte=0"`
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (s *ProductService) CreateProduct(ctx context.Context, in CreateProductInput) (*domain.Product, error) {
	p := &domain.Product{
		Name: in.Name, Description: in.Description,
		Price: in.Price, IsComposto: in.IsComposto, Active: true,
	}
	if err := s.repo.CreateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("ProductService.CreateProduct: %w", err)
	}
	return p, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductService.ListProducts: %w", err)
	}
	return products, nil
}

func (s *ProductService) ListActiveProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.ListActiveProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductService.ListActiveProducts: %w", err)
	}
	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetProduct: %w", err)
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	// Enriquece com a ficha técnica
	ingredients, err := s.repo.GetProductIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetProduct ingredients: %w", err)
	}
	p.Ingredients = ingredients
	return p, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProductService.DeleteProduct: %w", err)
	}
	if p == nil {
		return ErrProductNotFound
	}
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, in UpdateProductInput) (*domain.Product, error) {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateProduct: %w", err)
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	p.Name = in.Name
	p.Description = in.Description
	p.Price = in.Price
	p.IsComposto = in.IsComposto
	p.Active = in.Active

	if err := s.repo.UpdateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateProduct: %w", err)
	}
	return p, nil
}

// ── Ingrediente ──────────────────────────────────────────────────────────────

func (s *ProductService) CreateIngredient(ctx context.Context, in CreateIngredientInput) (*domain.Ingredient, error) {
	i := &domain.Ingredient{
		Name: in.Name, Unit: in.Unit,
		StockQuantity: in.StockQuantity, MinStock: in.MinStock,
	}
	if err := s.repo.CreateIngredient(ctx, i); err != nil {
		return nil, fmt.Errorf("ProductService.CreateIngredient: %w", err)
	}
	return i, nil
}

func (s *ProductService) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	return s.repo.ListIngredients(ctx)
}

func (s *ProductService) GetIngredient(ctx context.Context, id uint) (*domain.Ingredient, error) {
	i, err := s.repo.FindIngredientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetIngredient: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}
	return i, nil
}

func (s *ProductService) UpdateIngredientStock(ctx context.Context, id uint, in UpdateStockInput) (*domain.Ingredient, error) {
	i, err := s.repo.FindIngredientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredientStock: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}
	i.StockQuantity = in.Quantity
	if err := s.repo.UpdateIngredient(ctx, i); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredientStock: %w", err)
	}
	return i, nil
}

func (s *ProductService) UpdateIngredient(ctx context.Context, id uint, in UpdateIngredientInput) (*domain.Ingredient, error) {
	i, err := s.repo.FindIngredientByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredient: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}

	i.Name = in.Name
	i.Unit = in.Unit
	i.StockQuantity = in.StockQuantity
	i.MinStock = in.MinStock

	if err := s.repo.UpdateIngredient(ctx, i); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredient: %w", err)
	}
	return i, nil
}

func (s *ProductService) DeleteIngredient(ctx context.Context, id uint) error {
	i, err := s.repo.FindIngredientByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProductService.DeleteIngredient: %w", err)
	}
	if i == nil {
		return ErrIngredientNotFound
	}
	return s.repo.DeleteIngredient(ctx, id)
}

// ── Ficha técnica ────────────────────────────────────────────────────────────

func (s *ProductService) SetProductIngredients(
	ctx context.Context, productID uint, in SetProductIngredientsInput,
) error {
	p, err := s.repo.FindProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("ProductService.SetProductIngredients: %w", err)
	}
	if p == nil {
		return ErrProductNotFound
	}

	items := make([]domain.ProductIngredient, len(in.Items))
	for i, item := range in.Items {
		// Valida que o ingrediente existe
		ing, err := s.repo.FindIngredientByID(ctx, item.IngredientID)
		if err != nil {
			return fmt.Errorf("ProductService.SetProductIngredients: %w", err)
		}
		if ing == nil {
			return fmt.Errorf("ingrediente id=%d não encontrado", item.IngredientID)
		}
		items[i] = domain.ProductIngredient{
			ProductID:    productID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
		}
	}
	return s.repo.SetProductIngredients(ctx, productID, items)
}
