package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type GormProduct struct {
	ID                     uint   `gorm:"primaryKey;autoIncrement"`
	Name                   string `gorm:"not null"`
	Description            string
	Price                  float64 `gorm:"not null;default:0"`
	IsComposto             bool    `gorm:"not null;default:false"`
	Active                 bool    `gorm:"not null;default:true"`
	PhotoURL               string
	CategoryID             *uint `gorm:"index"`
	CompanyID              uint  `gorm:"index;not null"` // Sprint 3: NOT NULL
	DisplayOrder           int   `gorm:"not null;default:0"`
	PreparationTimeMinutes int   `gorm:"not null;default:0"`
	Featured               bool  `gorm:"not null;default:false"`
	IsNew                  bool  `gorm:"not null;default:false"`
	PromotionPrice         *float64
	PromotionStart         *int64
	PromotionEnd           *int64
	AvailableFrom          string
	AvailableUntil         string
	SKU                    string
	InternalNotes          string
	DeletedAt              *int64 `gorm:"index"`
	CreatedAt              int64  `gorm:"autoCreateTime"`
	UpdatedAt              int64  `gorm:"autoUpdateTime"`

	// Sprint 4 - Ficha Técnica Avançada
	Cost           float64 `gorm:"default:0.0"`
	CMV            float64 `gorm:"default:0.0"`
	Margin         float64 `gorm:"default:0.0"`
	Profit         float64 `gorm:"default:0.0"`
	SuggestedPrice float64 `gorm:"default:0.0"`

	// SEO fields para Cardápio Digital
	Slug            string `gorm:"uniqueIndex"`
	MetaTitle       string
	MetaDescription string
	AltImage        string
	Canonical       string

	// iFood integration fields
	ExternalID    string
	MarketplaceID string
	SyncStatus    string
	LastSync      *int64
}

func (GormProduct) TableName() string { return "products" }

type GormIngredient struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"not null"`
	Unit          string  `gorm:"not null;default:'un'"`
	StockQuantity float64 `gorm:"not null;default:0"`
	MinStock      float64 `gorm:"not null;default:0"`
	Active        bool    `gorm:"not null;default:true"`
	CompanyID     uint    `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt     *int64  `gorm:"index"`
	CreatedAt     int64   `gorm:"autoCreateTime"`
	UpdatedAt     int64   `gorm:"autoUpdateTime"`
}

func (GormIngredient) TableName() string { return "ingredients" }

type GormProductIngredient struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	ProductID    uint    `gorm:"not null;index"`
	IngredientID uint    `gorm:"not null"`
	Quantity     float64 `gorm:"not null"`
	// Sprint 4 - Ficha Técnica Avançada
	Loss       float64        `gorm:"default:0.0"`
	Yield      float64        `gorm:"default:1.0"`
	UnitCost   float64        `gorm:"default:0.0"`
	TotalCost  float64        `gorm:"default:0.0"`
	DeletedAt  *int64         `gorm:"index"`
	Ingredient GormIngredient `gorm:"foreignKey:IngredientID"`
}

func (GormProductIngredient) TableName() string { return "product_ingredients" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.ProductRepository = (*GormProductRepository)(nil)

type GormProductRepository struct{ db *gorm.DB }

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateProduct: %w", err)
	}

	// Check for slug collision
	if p.Slug != "" {
		var existing GormProduct
		err := r.db.WithContext(ctx).Where("slug = ?", p.Slug).First(&existing).Error
		if err == nil {
			return fmt.Errorf("CreateProduct: slug '%s' já está em uso", p.Slug)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("CreateProduct: verificar slug: %w", err)
		}
	}

	m := GormProduct{
		Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
		CompanyID:              companyID, // Auto-filled from context
		DisplayOrder:           p.DisplayOrder,
		PreparationTimeMinutes: p.PreparationTimeMinutes,
		Featured:               p.Featured,
		IsNew:                  p.IsNew,
		PromotionPrice:         p.PromotionPrice,
		AvailableFrom:          p.AvailableFrom,
		AvailableUntil:         p.AvailableUntil,
		SKU:                    p.SKU,
		InternalNotes:          p.InternalNotes,
		Slug:                   p.Slug,
		MetaTitle:              p.MetaTitle,
		MetaDescription:        p.MetaDescription,
		AltImage:               p.AltImage,
		Canonical:              p.Canonical,
		ExternalID:             p.ExternalID,
		MarketplaceID:          p.MarketplaceID,
		SyncStatus:             p.SyncStatus,
	}
	if p.PromotionStart != nil {
		ps := p.PromotionStart.Unix()
		m.PromotionStart = &ps
	}
	if p.PromotionEnd != nil {
		pe := p.PromotionEnd.Unix()
		m.PromotionEnd = &pe
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateProduct: %w", err)
	}
	p.ID = m.ID
	p.CompanyID = m.CompanyID
	p.CreatedAt = time.Unix(m.CreatedAt, 0)
	p.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	var m GormProduct
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindProductByID: %w", err)
	}
	return productToDomain(&m), nil
}

func (r *GormProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var ms []GormProduct
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListProducts: %w", err)
	}
	out := make([]domain.Product, len(ms))
	for i, m := range ms {
		out[i] = *productToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) ListActiveProducts(ctx context.Context) ([]domain.Product, error) {
	var ms []GormProduct
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("active = ? AND deleted_at IS NULL", true).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListActiveProducts: %w", err)
	}
	out := make([]domain.Product, len(ms))
	for i, m := range ms {
		out[i] = *productToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	// First, verify the product belongs to the tenant
	var existing GormProduct
	query := ApplyTenantFilterWithID(ctx, r.db, p.ID)
	if err := query.Where("deleted_at IS NULL").First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("UpdateProduct: product not found or access denied")
		}
		return fmt.Errorf("UpdateProduct: %w", err)
	}

	// Check for slug collision (excluding current product)
	if p.Slug != "" && p.Slug != existing.Slug {
		var slugConflict GormProduct
		err := r.db.WithContext(ctx).Where("slug = ? AND id != ?", p.Slug, p.ID).First(&slugConflict).Error
		if err == nil {
			return fmt.Errorf("UpdateProduct: slug '%s' já está em uso", p.Slug)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("UpdateProduct: verificar slug: %w", err)
		}
	}

	// Update without changing CompanyID (immutable)
	m := GormProduct{
		ID: p.ID, Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
		CompanyID:              existing.CompanyID, // Preserve original CompanyID
		DisplayOrder:           p.DisplayOrder,
		PreparationTimeMinutes: p.PreparationTimeMinutes,
		Featured:               p.Featured,
		IsNew:                  p.IsNew,
		PromotionPrice:         p.PromotionPrice,
		AvailableFrom:          p.AvailableFrom,
		AvailableUntil:         p.AvailableUntil,
		SKU:                    p.SKU,
		InternalNotes:          p.InternalNotes,
		Slug:                   p.Slug,
		MetaTitle:              p.MetaTitle,
		MetaDescription:        p.MetaDescription,
		AltImage:               p.AltImage,
		Canonical:              p.Canonical,
		ExternalID:             p.ExternalID,
		MarketplaceID:          p.MarketplaceID,
		SyncStatus:             p.SyncStatus,
	}
	if p.PromotionStart != nil {
		ps := p.PromotionStart.Unix()
		m.PromotionStart = &ps
	}
	if p.PromotionEnd != nil {
		pe := p.PromotionEnd.Unix()
		m.PromotionEnd = &pe
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateProduct: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	// Soft delete: marca DeletedAt (princípio #8)
	now := time.Now().Unix()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormProduct{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteProduct: %w", err)
	}
	return nil
}

func (r *GormProductRepository) CanDeleteProduct(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	check := &domain.DependencyCheck{CanDelete: true, Reasons: []domain.DependencyReason{}}

	// Verificar pedidos que contêm este produto
	type OrderResult struct {
		ID     uint   `gorm:"column:id"`
		Status string `gorm:"column:status"`
		Date   int64  `gorm:"column:created_at"`
	}
	var orders []OrderResult
	if err := r.db.WithContext(ctx).Table("order_items").
		Select("orders.id, orders.status, orders.created_at").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("order_items.product_id = ? AND orders.deleted_at IS NULL", id).
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("CanDeleteProduct: verificar pedidos: %w", err)
	}

	for _, order := range orders {
		check.CanDelete = false
		check.Reasons = append(check.Reasons, domain.DependencyReason{
			Type:        "order",
			ID:          order.ID,
			Name:        fmt.Sprintf("Pedido #%d", order.ID),
			Description: fmt.Sprintf("Status: %s, Data: %s", order.Status, time.Unix(order.Date, 0).Format("02/01/2006")),
		})
	}

	// Verificar se produto composto é usado em fichas técnicas de outros produtos
	var recipeCount int64
	if err := r.db.WithContext(ctx).Model(&GormProductIngredient{}).
		Where("ingredient_id = ? AND deleted_at IS NULL", id).
		Count(&recipeCount).Error; err != nil {
		return nil, fmt.Errorf("CanDeleteProduct: verificar fichas técnicas: %w", err)
	}

	if recipeCount > 0 {
		check.CanDelete = false
		check.Reasons = append(check.Reasons, domain.DependencyReason{
			Type:        "recipe",
			ID:          id,
			Name:        "Fichas Técnicas",
			Description: fmt.Sprintf("Usado em %d fichas técnicas de produtos compostos", recipeCount),
		})
	}

	return check, nil
}

// ── Ingrediente ──────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateIngredient(ctx context.Context, i *domain.Ingredient) error {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateIngredient: %w", err)
	}

	m := GormIngredient{
		Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
		Active:    true,
		CompanyID: companyID, // Auto-filled from context
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateIngredient: %w", err)
	}
	i.ID = m.ID
	i.CompanyID = m.CompanyID
	i.Active = m.Active
	i.CreatedAt = time.Unix(m.CreatedAt, 0)
	i.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindIngredientByID(ctx context.Context, id uint) (*domain.Ingredient, error) {
	var m GormIngredient
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindIngredientByID: %w", err)
	}
	return ingredientToDomain(&m), nil
}

func (r *GormProductRepository) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ms []GormIngredient
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListIngredients: %w", err)
	}
	out := make([]domain.Ingredient, len(ms))
	for i, m := range ms {
		out[i] = *ingredientToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient) error {
	// First, verify the ingredient belongs to the tenant
	var existing GormIngredient
	query := ApplyTenantFilterWithID(ctx, r.db, i.ID)
	if err := query.Where("deleted_at IS NULL").First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("UpdateIngredient: ingredient not found or access denied")
		}
		return fmt.Errorf("UpdateIngredient: %w", err)
	}

	// Update without changing CompanyID (immutable)
	m := GormIngredient{
		ID: i.ID, Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
		Active:    i.Active,
		CompanyID: existing.CompanyID, // Preserve original CompanyID
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateIngredient: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteIngredient(ctx context.Context, id uint) error {
	// Soft delete: marca DeletedAt (princípio #8)
	now := time.Now().Unix()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormIngredient{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteIngredient: %w", err)
	}
	return nil
}

func (r *GormProductRepository) CanDeleteIngredient(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	check := &domain.DependencyCheck{CanDelete: true, Reasons: []domain.DependencyReason{}}

	// Verificar fichas técnicas que usam este ingrediente
	type ProductResult struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var products []ProductResult
	if err := r.db.WithContext(ctx).Table("product_ingredients").
		Select("products.id, products.name").
		Joins("JOIN products ON product_ingredients.product_id = products.id").
		Where("product_ingredients.ingredient_id = ? AND product_ingredients.deleted_at IS NULL AND products.deleted_at IS NULL", id).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("CanDeleteIngredient: verificar fichas técnicas: %w", err)
	}

	for _, product := range products {
		check.CanDelete = false
		check.Reasons = append(check.Reasons, domain.DependencyReason{
			Type:        "product",
			ID:          product.ID,
			Name:        product.Name,
			Description: "Usado na ficha técnica deste produto composto",
		})
	}

	return check, nil
}

// ── Ficha técnica ────────────────────────────────────────────────────────────

func (r *GormProductRepository) SetProductIngredients(
	ctx context.Context, productID uint, items []domain.ProductIngredient,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Apaga ficha anterior e recria (upsert simples)
		if err := tx.Where("product_id = ? AND deleted_at IS NULL", productID).
			Delete(&GormProductIngredient{}).Error; err != nil {
			return fmt.Errorf("SetProductIngredients delete: %w", err)
		}
		for _, item := range items {
			m := GormProductIngredient{
				ProductID:    productID,
				IngredientID: item.IngredientID,
				Quantity:     item.Quantity,
				Loss:         item.Loss,
				Yield:        item.Yield,
				UnitCost:     item.UnitCost,
				TotalCost:    item.TotalCost,
			}
			if err := tx.Create(&m).Error; err != nil {
				return fmt.Errorf("SetProductIngredients insert: %w", err)
			}
		}
		return nil
	})
}

func (r *GormProductRepository) GetProductIngredients(
	ctx context.Context, productID uint,
) ([]domain.ProductIngredient, error) {
	var ms []GormProductIngredient
	if err := r.db.WithContext(ctx).
		Preload("Ingredient").
		Where("product_id = ? AND deleted_at IS NULL", productID).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("GetProductIngredients: %w", err)
	}
	out := make([]domain.ProductIngredient, len(ms))
	for i, m := range ms {
		ing := ingredientToDomain(&m.Ingredient)
		out[i] = domain.ProductIngredient{
			ID:           m.ID,
			ProductID:    m.ProductID,
			IngredientID: m.IngredientID,
			Quantity:     m.Quantity,
			Loss:         m.Loss,
			Yield:        m.Yield,
			UnitCost:     m.UnitCost,
			TotalCost:    m.TotalCost,
			Ingredient:   ing,
		}
	}
	return out, nil
}

// ── Estoque ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) DecreaseIngredientStock(
	ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB,
	ingredientName string, currentStock float64,
) error {
	// Usa o DB da transação se fornecido, senão usa o DB padrão
	db := r.db
	if txDB != nil {
		db = txDB.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	// Apply tenant filter to stock operations
	query := ApplyTenantFilterWithID(ctx, db, ingredientID)

	// Usa UPDATE com CHECK inline para garantir que não vai negativo
	result := query.
		Model(&GormIngredient{}).
		Where("stock_quantity >= ?", qty).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))

	if result.Error != nil {
		return fmt.Errorf("DecreaseIngredientStock id=%d: %w", ingredientID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"estoque insuficiente para '%s': disponível=%.4f necessário=%.4f",
			ingredientName, currentStock, qty,
		)
	}
	return nil
}

func (r *GormProductRepository) IncreaseIngredientStock(
	ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB,
) error {
	db := r.db
	if txDB != nil {
		db = txDB.WithContext(ctx)
	} else {
		db = db.WithContext(ctx)
	}

	// Apply tenant filter to stock operations
	query := ApplyTenantFilterWithID(ctx, db, ingredientID)

	result := query.
		Model(&GormIngredient{}).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", qty))

	if result.Error != nil {
		return fmt.Errorf("IncreaseIngredientStock id=%d: %w", ingredientID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ingrediente id=%d não encontrado", ingredientID)
	}
	return nil
}

// ── Mappers ──────────────────────────────────────────────────────────────────

func productToDomain(m *GormProduct) *domain.Product {
	var deletedAt *time.Time
	if m.DeletedAt != nil {
		dt := time.Unix(*m.DeletedAt, 0)
		deletedAt = &dt
	}
	var promotionStart, promotionEnd, lastSync *time.Time
	if m.PromotionStart != nil {
		ps := time.Unix(*m.PromotionStart, 0)
		promotionStart = &ps
	}
	if m.PromotionEnd != nil {
		pe := time.Unix(*m.PromotionEnd, 0)
		promotionEnd = &pe
	}
	if m.LastSync != nil {
		ls := time.Unix(*m.LastSync, 0)
		lastSync = &ls
	}
	return &domain.Product{
		ID: m.ID, Name: m.Name, Description: m.Description,
		Price: m.Price, IsComposto: m.IsComposto, Active: m.Active,
		PhotoURL:               m.PhotoURL,
		CategoryID:             m.CategoryID,
		CompanyID:              m.CompanyID,
		DisplayOrder:           m.DisplayOrder,
		PreparationTimeMinutes: m.PreparationTimeMinutes,
		Featured:               m.Featured,
		IsNew:                  m.IsNew,
		PromotionPrice:         m.PromotionPrice,
		PromotionStart:         promotionStart,
		PromotionEnd:           promotionEnd,
		AvailableFrom:          m.AvailableFrom,
		AvailableUntil:         m.AvailableUntil,
		SKU:                    m.SKU,
		InternalNotes:          m.InternalNotes,
		DeletedAt:              deletedAt,
		CreatedAt:              time.Unix(m.CreatedAt, 0), UpdatedAt: time.Unix(m.UpdatedAt, 0),
		Slug:            m.Slug,
		MetaTitle:       m.MetaTitle,
		MetaDescription: m.MetaDescription,
		AltImage:        m.AltImage,
		Canonical:       m.Canonical,
		ExternalID:      m.ExternalID,
		MarketplaceID:   m.MarketplaceID,
		SyncStatus:      m.SyncStatus,
		LastSync:        lastSync,
	}
}

func ingredientToDomain(m *GormIngredient) *domain.Ingredient {
	var deletedAt *time.Time
	if m.DeletedAt != nil {
		dt := time.Unix(*m.DeletedAt, 0)
		deletedAt = &dt
	}
	return &domain.Ingredient{
		ID: m.ID, Name: m.Name, Unit: m.Unit,
		StockQuantity: m.StockQuantity, MinStock: m.MinStock,
		Active:    m.Active,
		CompanyID: m.CompanyID,
		DeletedAt: deletedAt,
		CreatedAt: time.Unix(m.CreatedAt, 0), UpdatedAt: time.Unix(m.UpdatedAt, 0),
	}
}
