package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
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
	DeletedAt     *int64  `gorm:"index"`
	CreatedAt     int64   `gorm:"autoCreateTime"`
	UpdatedAt     int64   `gorm:"autoUpdateTime"`
}

func (GormIngredient) TableName() string { return "ingredients" }

type GormProductIngredient struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	ProductID    uint           `gorm:"not null;index"`
	IngredientID uint           `gorm:"not null"`
	Quantity     float64        `gorm:"not null"`
	DeletedAt    *int64         `gorm:"index"`
	Ingredient   GormIngredient `gorm:"foreignKey:IngredientID"`
}

func (GormProductIngredient) TableName() string { return "product_ingredients" }

type GormCategory struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"not null"`
	Description  string
	DisplayOrder int    `gorm:"not null;default:0"`
	Active       bool   `gorm:"not null;default:true"`
	DeletedAt    *int64 `gorm:"index"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
}

func (GormCategory) TableName() string { return "categories" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.ProductRepository = (*GormProductRepository)(nil)

type GormProductRepository struct{ db *gorm.DB }

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	m := GormProduct{
		Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
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
	p.CreatedAt = time.Unix(m.CreatedAt, 0)
	p.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	var m GormProduct
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error
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
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&ms).Error; err != nil {
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
	if err := r.db.WithContext(ctx).Where("active = ? AND deleted_at IS NULL", true).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListActiveProducts: %w", err)
	}
	out := make([]domain.Product, len(ms))
	for i, m := range ms {
		out[i] = *productToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	m := GormProduct{
		ID: p.ID, Name: p.Name, Description: p.Description,
		Price: p.Price, IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
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
	if err := r.db.WithContext(ctx).Model(&GormProduct{}).
		Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteProduct: %w", err)
	}
	return nil
}

// ── Ingrediente ──────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateIngredient(ctx context.Context, i *domain.Ingredient) error {
	m := GormIngredient{
		Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
		Active: true,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateIngredient: %w", err)
	}
	i.ID = m.ID
	i.Active = m.Active
	i.CreatedAt = time.Unix(m.CreatedAt, 0)
	i.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormProductRepository) FindIngredientByID(ctx context.Context, id uint) (*domain.Ingredient, error) {
	var m GormIngredient
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error
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
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListIngredients: %w", err)
	}
	out := make([]domain.Ingredient, len(ms))
	for i, m := range ms {
		out[i] = *ingredientToDomain(&m)
	}
	return out, nil
}

func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient) error {
	m := GormIngredient{
		ID: i.ID, Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
		Active: i.Active,
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateIngredient: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteIngredient(ctx context.Context, id uint) error {
	// Soft delete: marca DeletedAt (princípio #8)
	now := time.Now().Unix()
	if err := r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteIngredient: %w", err)
	}
	return nil
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
			ID: m.ID, ProductID: m.ProductID,
			IngredientID: m.IngredientID, Quantity: m.Quantity,
			Ingredient: ing,
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

	// Usa UPDATE com CHECK inline para garantir que não vai negativo
	result := db.
		Model(&GormIngredient{}).
		Where("id = ? AND stock_quantity >= ?", ingredientID, qty).
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

	result := db.
		Model(&GormIngredient{}).
		Where("id = ?", ingredientID).
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
		DeletedAt: deletedAt,
		CreatedAt: time.Unix(m.CreatedAt, 0), UpdatedAt: time.Unix(m.UpdatedAt, 0),
	}
}
