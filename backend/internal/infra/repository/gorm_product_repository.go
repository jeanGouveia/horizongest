package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

// ─── Helper functions for Money conversion ───────────────────────────────────

func convertMoneyPtrToInt64Ptr(m *domain.Money) *int64 {
	if m == nil {
		return nil
	}
	v := int64(*m)
	return &v
}

func convertInt64PtrToMoneyPtr(i *int64) *domain.Money {
	if i == nil {
		return nil
	}
	v := domain.Money(*i)
	return &v
}

// ─── GORM models ────────────────────────────────────────────────────────────

type GormProduct struct {
	ID                     uint   `gorm:"primaryKey;autoIncrement"`
	Name                   string `gorm:"not null"`
	Description            string
	Price                  int64 `gorm:"not null;default:0"`
	IsComposto             bool  `gorm:"not null;default:false"`
	Active                 bool  `gorm:"not null;default:true"`
	PhotoURL               string
	CategoryID             *uint `gorm:"index"`
	CompanyID              uint  `gorm:"index;not null"` // Sprint 3: NOT NULL
	DisplayOrder           int   `gorm:"not null;default:0"`
	PreparationTimeMinutes int   `gorm:"not null;default:0"`
	Featured               bool  `gorm:"not null;default:false"`
	IsNew                  bool  `gorm:"not null;default:false"`
	PromotionPrice         *int64
	PromotionStart         *int64
	PromotionEnd           *int64
	AvailableFrom          string
	AvailableUntil         string
	SKU                    string
	InternalNotes          string
	DeletedAt              *time.Time `gorm:"index"`
	CreatedAt              time.Time  `gorm:"autoCreateTime"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime"`

	// Sprint 4 - Ficha Técnica Avançada
	Cost           int64   `gorm:"default:0"`
	CMV            int64   `gorm:"default:0"`
	Margin         float64 `gorm:"default:0.0"`
	Profit         int64   `gorm:"default:0"`
	SuggestedPrice int64   `gorm:"default:0"`

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
	ID            uint       `gorm:"primaryKey;autoIncrement"`
	Name          string     `gorm:"not null"`
	Unit          string     `gorm:"not null;default:'un'"`
	StockQuantity float64    `gorm:"not null;default:0"`
	MinStock      float64    `gorm:"not null;default:0"`
	Active        bool       `gorm:"not null;default:true"`
	CompanyID     uint       `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt     *time.Time `gorm:"index"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
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
	UnitCost   int64          `gorm:"default:0"`
	TotalCost  int64          `gorm:"default:0"`
	DeletedAt  *time.Time     `gorm:"index"`
	Ingredient GormIngredient `gorm:"foreignKey:IngredientID"`
}

func (GormProductIngredient) TableName() string { return "product_ingredients" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.ProductRepository = (*GormProductRepository)(nil)

type GormProductRepository struct{ db *gorm.DB }

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// getDB retorna a transação se fornecida, senão retorna o DB padrão
// Sprint 4B.1 v2: Transaction propagation
func (r *GormProductRepository) getDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

// ── Produto ──────────────────────────────────────────────────────────────────

func (r *GormProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("ProductRepository.CreateProduct: %w", err)
	}

	// Check for slug collision
	if p.Slug != "" {
		var existing GormProduct
		err := r.db.WithContext(ctx).Where("slug = ?", p.Slug).First(&existing).Error
		if err == nil {
			return fmt.Errorf("ProductRepository.CreateProduct: slug '%s' já está em uso", p.Slug)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ProductRepository.CreateProduct: verificar slug: %w", err)
		}
	}

	m := GormProduct{
		Name: p.Name, Description: p.Description,
		Price: int64(p.Price), IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
		CompanyID:              companyID, // Auto-filled from context
		DisplayOrder:           p.DisplayOrder,
		PreparationTimeMinutes: p.PreparationTimeMinutes,
		Featured:               p.Featured,
		IsNew:                  p.IsNew,
		PromotionPrice:         convertMoneyPtrToInt64Ptr(p.PromotionPrice),
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
		return fmt.Errorf("ProductRepository.CreateProduct: %w", err)
	}
	p.ID = m.ID
	p.CompanyID = m.CompanyID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
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
		return nil, fmt.Errorf("ProductRepository.FindProductByID: %w", err)
	}
	return productToDomain(&m), nil
}

func (r *GormProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var ms []GormProduct
	query := ApplyTenantFilter(ctx, r.db)
	// Sprint 5D.3 - Performance Hardening: Add default LIMIT to prevent timeouts
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Limit(100).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ProductRepository.ListProducts: %w", err)
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
	// Sprint 5D.3 - Performance Hardening: Add default LIMIT to prevent timeouts
	if err := query.WithContext(ctx).Where("active = ? AND deleted_at IS NULL", true).Limit(100).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ProductRepository.ListActiveProducts: %w", err)
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
			return fmt.Errorf("ProductRepository.UpdateProduct: produto não encontrado ou acesso negado")
		}
		return fmt.Errorf("ProductRepository.UpdateProduct: %w", err)
	}

	// Check for slug collision (excluding current product)
	if p.Slug != "" && p.Slug != existing.Slug {
		var slugConflict GormProduct
		err := r.db.WithContext(ctx).Where("slug = ? AND id != ?", p.Slug, p.ID).First(&slugConflict).Error
		if err == nil {
			return fmt.Errorf("ProductRepository.UpdateProduct: slug '%s' já está em uso", p.Slug)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ProductRepository.UpdateProduct: verificar slug: %w", err)
		}
	}

	// Update without changing CompanyID (immutable)
	m := GormProduct{
		ID: p.ID, Name: p.Name, Description: p.Description,
		Price: int64(p.Price), IsComposto: p.IsComposto, Active: p.Active,
		PhotoURL:               p.PhotoURL,
		CategoryID:             p.CategoryID,
		CompanyID:              existing.CompanyID, // Preserve original CompanyID
		DisplayOrder:           p.DisplayOrder,
		PreparationTimeMinutes: p.PreparationTimeMinutes,
		Featured:               p.Featured,
		IsNew:                  p.IsNew,
		PromotionPrice:         convertMoneyPtrToInt64Ptr(p.PromotionPrice),
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
		return fmt.Errorf("ProductRepository.UpdateProduct: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	// Soft delete: marca DeletedAt (princípio #8)
	now := time.Now()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormProduct{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("ProductRepository.DeleteProduct: %w", err)
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
		return nil, fmt.Errorf("ProductRepository.CanDeleteProduct: verificar pedidos: %w", err)
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
		return nil, fmt.Errorf("ProductRepository.CanDeleteProduct: verificar fichas técnicas: %w", err)
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
		return fmt.Errorf("ProductRepository.CreateIngredient: %w", err)
	}

	m := GormIngredient{
		Name: i.Name, Unit: i.Unit,
		StockQuantity: i.StockQuantity, MinStock: i.MinStock,
		Active:    true,
		CompanyID: companyID, // Auto-filled from context
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("ProductRepository.CreateIngredient: %w", err)
	}
	i.ID = m.ID
	i.CompanyID = m.CompanyID
	i.Active = m.Active
	i.CreatedAt = m.CreatedAt
	i.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormProductRepository) FindIngredientByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error) {
	var m GormIngredient
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	err := query.Where("deleted_at IS NULL").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("ProductRepository.FindIngredientByID: ingrediente não encontrado ou acesso negado")
	}
	if err != nil {
		return nil, fmt.Errorf("ProductRepository.FindIngredientByID: %w", err)
	}
	return ingredientToDomain(&m), nil
}

// FindIngredientByIDForUpdate busca ingrediente com SELECT FOR UPDATE
// Sprint 4B.1 v2: Lock pessimista real
func (r *GormProductRepository) FindIngredientByIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error) {
	var m GormIngredient
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	err := query.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("deleted_at IS NULL").
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ProductRepository.FindIngredientByIDForUpdate: %w", err)
	}
	return ingredientToDomain(&m), nil
}

func (r *GormProductRepository) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ms []GormIngredient
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ProductRepository.ListIngredients: %w", err)
	}
	out := make([]domain.Ingredient, len(ms))
	for i, m := range ms {
		out[i] = *ingredientToDomain(&m)
	}
	return out, nil
}

// Sprint 4B.4: Adicionado SELECT FOR UPDATE para prevenir lost update
func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient, tx *gorm.DB) error {
	// Sprint 4B.4: SELECT FOR UPDATE antes de qualquer leitura de campos
	// Isso previne lost update bloqueando o registro durante a transação
	var existing GormIngredient
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), i.ID)
	if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("deleted_at IS NULL").
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ProductRepository.UpdateIngredient: ingrediente não encontrado ou acesso negado")
		}
		return fmt.Errorf("ProductRepository.UpdateIngredient: erro ao atualizar ingrediente: %w", err)
	}

	// Sprint 4B.4: Usar Updates() em vez de Save() para evitar atualizar timestamps desnecessários
	// Update sem changing CompanyID (immutable)
	if err := r.getDB(ctx, tx).Model(&GormIngredient{}).
		Where("id = ? AND deleted_at IS NULL", i.ID).
		Updates(map[string]interface{}{
			"name":           i.Name,
			"unit":           i.Unit,
			"stock_quantity": i.StockQuantity,
			"min_stock":      i.MinStock,
			"active":         i.Active,
		}).Error; err != nil {
		return fmt.Errorf("ProductRepository.UpdateIngredient: %w", err)
	}
	return nil
}

func (r *GormProductRepository) DeleteIngredient(ctx context.Context, id uint) error {
	// Soft delete: marca DeletedAt (princípio #8)
	now := time.Now()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormIngredient{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("ProductRepository.DeleteIngredient: %w", err)
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
		return nil, fmt.Errorf("ProductRepository.CanDeleteIngredient: verificar fichas técnicas: %w", err)
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
			return fmt.Errorf("ProductRepository.SetProductIngredients: deletar ingredientes anteriores: %w", err)
		}
		for _, item := range items {
			m := GormProductIngredient{
				ProductID:    productID,
				IngredientID: item.IngredientID,
				Quantity:     item.Quantity,
				Loss:         item.Loss,
				Yield:        item.Yield,
				UnitCost:     int64(item.UnitCost),
				TotalCost:    int64(item.TotalCost),
			}
			if err := tx.Create(&m).Error; err != nil {
				return fmt.Errorf("ProductRepository.SetProductIngredients: inserir ingrediente: %w", err)
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
		return nil, fmt.Errorf("ProductRepository.GetProductIngredients: %w", err)
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
			UnitCost:     domain.Money(m.UnitCost),
			TotalCost:    domain.Money(m.TotalCost),
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
	db := r.getDB(ctx, txDB)

	// Build WHERE clause manually
	whereClause := "id = ? AND deleted_at IS NULL"
	whereArgs := []interface{}{ingredientID}

	// Apply tenant filter if context has tenant
	if tenantCtx, ok := domain.GetTenantContextFromContext(ctx); ok {
		whereClause += " AND company_id = ?"
		whereArgs = append(whereArgs, tenantCtx.GetCompanyID())
	}

	// Sprint 4B.1 v2: SELECT FOR UPDATE real com Clauses(clause.Locking{Strength: "UPDATE"})
	// Isso impede que outras transações leiam/modifiquem o ingrediente durante esta operação
	var ingredient GormIngredient
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(whereClause, whereArgs...).
		First(&ingredient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ProductRepository.DecreaseIngredientStock: ingrediente não encontrado")
		}
		return fmt.Errorf("ProductRepository.DecreaseIngredientStock: lock ingrediente: %w", err)
	}

	// Verificar estoque suficiente antes do UPDATE (validação adicional)
	if ingredient.StockQuantity < qty {
		return fmt.Errorf("ProductRepository.DecreaseIngredientStock: estoque insuficiente para '%s': disponível=%.4f necessário=%.4f",
			ingredientName, ingredient.StockQuantity, qty,
		)
	}

	// Usa UPDATE com CHECK inline como garantia adicional (defesa em profundidade)
	updateWhere := whereClause + " AND stock_quantity >= ?"
	updateArgs := append(whereArgs, qty)
	result := db.Model(&GormIngredient{}).
		Where(updateWhere, updateArgs...).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", qty))

	if result.Error != nil {
		return fmt.Errorf("ProductRepository.DecreaseIngredientStock: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ProductRepository.DecreaseIngredientStock: estoque insuficiente (concorrência)")
	}
	return nil
}

func (r *GormProductRepository) IncreaseIngredientStock(
	ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB,
) error {
	db := r.getDB(ctx, txDB)

	// Build WHERE clause manually
	whereClause := "id = ? AND deleted_at IS NULL"
	whereArgs := []interface{}{ingredientID}

	// Apply tenant filter if context has tenant
	if tenantCtx, ok := domain.GetTenantContextFromContext(ctx); ok {
		whereClause += " AND company_id = ?"
		whereArgs = append(whereArgs, tenantCtx.GetCompanyID())
	}

	// Sprint 4B.1 v2: SELECT FOR UPDATE real com Clauses(clause.Locking{Strength: "UPDATE"})
	// Consistente com DecreaseIngredientStock para evitar race condition
	var ingredient GormIngredient
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(whereClause, whereArgs...).
		First(&ingredient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ProductRepository.IncreaseIngredientStock: ingrediente não encontrado")
		}
		return fmt.Errorf("ProductRepository.IncreaseIngredientStock: lock ingrediente: %w", err)
	}

	result := db.Model(&GormIngredient{}).
		Where(whereClause, whereArgs...).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", qty))

	if result.Error != nil {
		return fmt.Errorf("ProductRepository.IncreaseIngredientStock: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ProductRepository.IncreaseIngredientStock: ingrediente id=%d não encontrado", ingredientID)
	}
	return nil
}

// ── Mappers ──────────────────────────────────────────────────────────────────

func productToDomain(m *GormProduct) *domain.Product {
	var deletedAt *time.Time
	if m.DeletedAt != nil {
		dt := *m.DeletedAt
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
		Price: domain.Money(m.Price), IsComposto: m.IsComposto, Active: m.Active,
		PhotoURL:               m.PhotoURL,
		CategoryID:             m.CategoryID,
		CompanyID:              m.CompanyID,
		DisplayOrder:           m.DisplayOrder,
		PreparationTimeMinutes: m.PreparationTimeMinutes,
		Featured:               m.Featured,
		IsNew:                  m.IsNew,
		PromotionPrice:         convertInt64PtrToMoneyPtr(m.PromotionPrice),
		PromotionStart:         promotionStart,
		PromotionEnd:           promotionEnd,
		AvailableFrom:          m.AvailableFrom,
		AvailableUntil:         m.AvailableUntil,
		SKU:                    m.SKU,
		InternalNotes:          m.InternalNotes,
		DeletedAt:              deletedAt,
		CreatedAt:              m.CreatedAt, UpdatedAt: m.UpdatedAt,
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
	return &domain.Ingredient{
		ID: m.ID, Name: m.Name, Unit: m.Unit,
		StockQuantity: m.StockQuantity, MinStock: m.MinStock,
		Active:    m.Active,
		CompanyID: m.CompanyID,
		DeletedAt: m.DeletedAt,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
