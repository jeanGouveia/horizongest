package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound    = errors.New("produto não encontrado")
	ErrIngredientNotFound = errors.New("ingrediente não encontrado")
)

type ProductService struct {
	repo ports.ProductRepository
	db   *gorm.DB
}

func NewProductService(repo ports.ProductRepository, db *gorm.DB) *ProductService {
	return &ProductService{repo: repo, db: db}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func generateSlug(name string) string {
	// Converter para minúsculas
	slug := strings.ToLower(name)

	// Remover acentos (simplificado)
	replacements := map[rune]string{
		'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a",
		'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
		'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
		'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o",
		'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
		'ç': "c",
		'ñ': "n",
	}
	for char, replacement := range replacements {
		slug = strings.ReplaceAll(slug, string(char), replacement)
	}

	// Remover caracteres especiais (exceto hífen)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, slug)

	// Substituir espaços e underscores por hífen
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remover múltiplos hífens consecutivos
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remover hífens no início e fim
	slug = strings.Trim(slug, "-")

	// Limitar tamanho
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type CreateProductInput struct {
	Name                   string        `json:"name"                    validate:"required,min=2,max=120"`
	Description            string        `json:"description"`
	Price                  domain.Money  `json:"price"                   validate:"required,gt=0"`
	IsComposto             bool          `json:"is_composto"`
	PhotoURL               string        `json:"photo_url"`
	CategoryID             *uint         `json:"category_id"`
	DisplayOrder           int           `json:"display_order"           validate:"gte=0"`
	PreparationTimeMinutes int           `json:"preparation_time_minutes" validate:"gte=0"`
	Featured               bool          `json:"featured"`
	IsNew                  bool          `json:"is_new"`
	PromotionPrice         *domain.Money `json:"promotion_price"          validate:"omitempty,gt=0"`
	PromotionStart         *time.Time    `json:"promotion_start"`
	PromotionEnd           *time.Time    `json:"promotion_end"`
	AvailableFrom          string        `json:"available_from"`
	AvailableUntil         string        `json:"available_until"`
	SKU                    string        `json:"sku"`
	InternalNotes          string        `json:"internal_notes"`
	Slug                   string        `json:"slug"`
	MetaTitle              string        `json:"meta_title"`
	MetaDescription        string        `json:"meta_description"`
	AltImage               string        `json:"alt_image"`
	Canonical              string        `json:"canonical"`
	ExternalID             string        `json:"external_id"`
	MarketplaceID          string        `json:"marketplace_id"`
	SyncStatus             string        `json:"sync_status"`
}

type UpdateProductInput struct {
	Name                   string        `json:"name"                    validate:"required,min=2,max=120"`
	Description            string        `json:"description"`
	Price                  domain.Money  `json:"price"                   validate:"required,gt=0"`
	IsComposto             *bool         `json:"is_composto"`
	Active                 *bool         `json:"active"`
	PhotoURL               string        `json:"photo_url"`
	CategoryID             *uint         `json:"category_id"`
	DisplayOrder           int           `json:"display_order"           validate:"gte=0"`
	PreparationTimeMinutes int           `json:"preparation_time_minutes" validate:"gte=0"`
	Featured               *bool         `json:"featured"`
	IsNew                  *bool         `json:"is_new"`
	PromotionPrice         *domain.Money `json:"promotion_price"          validate:"omitempty,gt=0"`
	PromotionStart         *time.Time    `json:"promotion_start"`
	PromotionEnd           *time.Time    `json:"promotion_end"`
	AvailableFrom          string        `json:"available_from"`
	AvailableUntil         string        `json:"available_until"`
	SKU                    string        `json:"sku"`
	InternalNotes          string        `json:"internal_notes"`
	Slug                   string        `json:"slug"`
	MetaTitle              string        `json:"meta_title"`
	MetaDescription        string        `json:"meta_description"`
	AltImage               string        `json:"alt_image"`
	Canonical              string        `json:"canonical"`
	ExternalID             string        `json:"external_id"`
	MarketplaceID          string        `json:"marketplace_id"`
	SyncStatus             string        `json:"sync_status"`
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
	Quantity      float64 `json:"quantity" validate:"omitempty,gte=0"`
	StockQuantity float64 `json:"stock_quantity" validate:"omitempty,gte=0"`
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (s *ProductService) CreateProduct(ctx context.Context, in CreateProductInput) (*domain.Product, error) {
	// Gerar slug automaticamente se não fornecido
	slug := in.Slug
	if slug == "" {
		slug = generateSlug(in.Name)
	}

	p := &domain.Product{
		Name: in.Name, Description: in.Description,
		Price: in.Price, IsComposto: in.IsComposto, Active: true,
		PhotoURL:               in.PhotoURL,
		CategoryID:             in.CategoryID,
		DisplayOrder:           in.DisplayOrder,
		PreparationTimeMinutes: in.PreparationTimeMinutes,
		Featured:               in.Featured,
		IsNew:                  in.IsNew,
		PromotionPrice:         in.PromotionPrice,
		PromotionStart:         in.PromotionStart,
		PromotionEnd:           in.PromotionEnd,
		AvailableFrom:          in.AvailableFrom,
		AvailableUntil:         in.AvailableUntil,
		SKU:                    in.SKU,
		InternalNotes:          in.InternalNotes,
		Slug:                   slug,
		MetaTitle:              in.MetaTitle,
		MetaDescription:        in.MetaDescription,
		AltImage:               in.AltImage,
		Canonical:              in.Canonical,
		ExternalID:             in.ExternalID,
		MarketplaceID:          in.MarketplaceID,
		SyncStatus:             in.SyncStatus,
	}
	if err := s.repo.CreateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("ProductService.CreateProduct: criar produto: %w", err)
	}
	return p, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductService.ListProducts: listar produtos: %w", err)
	}
	return products, nil
}

func (s *ProductService) ListActiveProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.ListActiveProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProductService.ListActiveProducts: listar produtos ativos: %w", err)
	}
	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetProduct: buscar produto: %w", err)
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	// Enriquece com a ficha técnica
	ingredients, err := s.repo.GetProductIngredients(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetProduct: buscar ingredientes: %w", err)
	}
	p.Ingredients = ingredients
	return p, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProductService.DeleteProduct: buscar produto: %w", err)
	}
	if p == nil {
		return ErrProductNotFound
	}
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, in UpdateProductInput) (*domain.Product, error) {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateProduct: buscar produto: %w", err)
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	p.Name = in.Name
	p.Description = in.Description
	p.Price = in.Price
	if in.IsComposto != nil {
		p.IsComposto = *in.IsComposto
	}
	if in.Active != nil {
		p.Active = *in.Active
	}
	p.PhotoURL = in.PhotoURL
	p.CategoryID = in.CategoryID
	p.DisplayOrder = in.DisplayOrder
	p.PreparationTimeMinutes = in.PreparationTimeMinutes
	if in.Featured != nil {
		p.Featured = *in.Featured
	}
	if in.IsNew != nil {
		p.IsNew = *in.IsNew
	}
	p.PromotionPrice = in.PromotionPrice
	p.PromotionStart = in.PromotionStart
	p.PromotionEnd = in.PromotionEnd
	p.AvailableFrom = in.AvailableFrom
	p.AvailableUntil = in.AvailableUntil
	p.SKU = in.SKU
	p.InternalNotes = in.InternalNotes

	// Gerar slug automaticamente se não fornecido
	if in.Slug == "" {
		p.Slug = generateSlug(in.Name)
	} else {
		p.Slug = in.Slug
	}
	p.MetaTitle = in.MetaTitle
	p.MetaDescription = in.MetaDescription
	p.AltImage = in.AltImage
	p.Canonical = in.Canonical
	p.ExternalID = in.ExternalID
	p.MarketplaceID = in.MarketplaceID
	p.SyncStatus = in.SyncStatus

	if err := s.repo.UpdateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateProduct: atualizar produto: %w", err)
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
		return nil, fmt.Errorf("ProductService.CreateIngredient: criar ingrediente: %w", err)
	}
	return i, nil
}

func (s *ProductService) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	return s.repo.ListIngredients(ctx)
}

func (s *ProductService) GetIngredient(ctx context.Context, id uint) (*domain.Ingredient, error) {
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	i, err := s.repo.FindIngredientByID(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetIngredient: buscar ingrediente: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}
	return i, nil
}

func (s *ProductService) UpdateIngredientStock(ctx context.Context, id uint, in UpdateStockInput) (*domain.Ingredient, error) {
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	i, err := s.repo.FindIngredientByID(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredientStock: buscar ingrediente: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}
	// Support both 'quantity' and 'stock_quantity' field names
	// At least one must be provided
	if in.StockQuantity == 0 && in.Quantity == 0 {
		return nil, errors.New("quantity or stock_quantity is required")
	}
	if in.StockQuantity != 0 {
		i.StockQuantity = in.StockQuantity
	} else {
		i.StockQuantity = in.Quantity
	}
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	if err := s.repo.UpdateIngredient(ctx, i, nil); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredientStock: atualizar ingrediente: %w", err)
	}
	return i, nil
}

func (s *ProductService) UpdateIngredient(ctx context.Context, id uint, in UpdateIngredientInput) (*domain.Ingredient, error) {
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	i, err := s.repo.FindIngredientByID(ctx, id, nil)
	if err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredient: buscar ingrediente: %w", err)
	}
	if i == nil {
		return nil, ErrIngredientNotFound
	}

	i.Name = in.Name
	i.Unit = in.Unit
	i.StockQuantity = in.StockQuantity
	i.MinStock = in.MinStock

	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	if err := s.repo.UpdateIngredient(ctx, i, nil); err != nil {
		return nil, fmt.Errorf("ProductService.UpdateIngredient: atualizar ingrediente: %w", err)
	}
	return i, nil
}

func (s *ProductService) DeleteIngredient(ctx context.Context, id uint) error {
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	i, err := s.repo.FindIngredientByID(ctx, id, nil)
	if err != nil {
		return fmt.Errorf("ProductService.DeleteIngredient: buscar ingrediente: %w", err)
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
		return fmt.Errorf("ProductService.SetProductIngredients: buscar produto: %w", err)
	}
	if p == nil {
		return ErrProductNotFound
	}

	items := make([]domain.ProductIngredient, len(in.Items))
	for i, item := range in.Items {
		// Valida que o ingrediente existe
		// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
		ing, err := s.repo.FindIngredientByID(ctx, item.IngredientID, nil)
		if err != nil {
			return fmt.Errorf("ProductService.SetProductIngredients: buscar ingrediente: %w", err)
		}
		if ing == nil {
			return fmt.Errorf("ProductService.SetProductIngredients: ingrediente id=%d não encontrado", item.IngredientID)
		}
		items[i] = domain.ProductIngredient{
			ProductID:    productID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
		}
	}
	return s.repo.SetProductIngredients(ctx, productID, items)
}

// DuplicateProduct creates a copy of an existing product in transação atômica
func (s *ProductService) DuplicateProduct(ctx context.Context, id uint) (*domain.Product, error) {
	var duplicate *domain.Product

	// Executar toda a operação em transação atômica (se db disponível)
	executeInTx := func(fn func() error) error {
		if s.db != nil {
			return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return fn()
			})
		}
		return fn()
	}

	err := executeInTx(func() error {
		// Get the original product
		original, err := s.repo.FindProductByID(ctx, id)
		if err != nil {
			return fmt.Errorf("ProductService.DuplicateProduct: buscar produto original: %w", err)
		}
		if original == nil {
			return ErrProductNotFound
		}

		// Get original ingredients
		ingredients, err := s.repo.GetProductIngredients(ctx, id)
		if err != nil {
			return fmt.Errorf("ProductService.DuplicateProduct: buscar ingredientes: %w", err)
		}

		// Create duplicate with modified name
		duplicate = &domain.Product{
			Name:                   original.Name + " (Cópia)",
			Description:            original.Description,
			Price:                  original.Price,
			IsComposto:             original.IsComposto,
			Active:                 true, // Always activate duplicates
			PhotoURL:               original.PhotoURL,
			CategoryID:             original.CategoryID,
			DisplayOrder:           original.DisplayOrder,
			PreparationTimeMinutes: original.PreparationTimeMinutes,
			Featured:               false, // Reset featured flag
			IsNew:                  true,  // Mark as new
			PromotionPrice:         nil,   // Reset promotions
			PromotionStart:         nil,
			PromotionEnd:           nil,
			AvailableFrom:          original.AvailableFrom,
			AvailableUntil:         original.AvailableUntil,
			SKU:                    "", // Reset SKU
			InternalNotes:          original.InternalNotes,
			Slug:                   generateSlug(original.Name + "-copia"),
			MetaTitle:              original.MetaTitle,
			MetaDescription:        original.MetaDescription,
			AltImage:               original.AltImage,
			Canonical:              original.Canonical,
			ExternalID:             original.ExternalID,
			MarketplaceID:          original.MarketplaceID,
			SyncStatus:             original.SyncStatus,
		}

		if err := s.repo.CreateProduct(ctx, duplicate); err != nil {
			return fmt.Errorf("ProductService.DuplicateProduct: criar duplicata: %w", err)
		}

		// Copy ingredients if it's a composite product
		if len(ingredients) > 0 {
			duplicateIngredients := make([]domain.ProductIngredient, len(ingredients))
			for i, ing := range ingredients {
				duplicateIngredients[i] = domain.ProductIngredient{
					ProductID:    duplicate.ID,
					IngredientID: ing.IngredientID,
					Quantity:     ing.Quantity,
				}
			}
			if err := s.repo.SetProductIngredients(ctx, duplicate.ID, duplicateIngredients); err != nil {
				return fmt.Errorf("ProductService.DuplicateProduct: copiar ingredientes: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return duplicate, nil
}

// ArchiveProduct sets a product's active status to false
func (s *ProductService) ArchiveProduct(ctx context.Context, id uint) error {
	p, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ProductService.ArchiveProduct: buscar produto: %w", err)
	}
	if p == nil {
		return ErrProductNotFound
	}

	p.Active = false
	if err := s.repo.UpdateProduct(ctx, p); err != nil {
		return fmt.Errorf("ProductService.ArchiveProduct: atualizar produto: %w", err)
	}
	return nil
}
