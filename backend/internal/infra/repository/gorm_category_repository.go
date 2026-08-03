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

var _ ports.CategoryRepository = (*GormCategoryRepository)(nil)

type GormCategory struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Name         string     `gorm:"not null"`
	Description  string     `gorm:"type:text"`
	DisplayOrder int        `gorm:"not null;default:0"`
	Active       bool       `gorm:"not null;default:true"`
	CompanyID    uint       `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt    *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

func (GormCategory) TableName() string { return "categories" }

type GormCategoryRepository struct{ db *gorm.DB }

func NewGormCategoryRepository(db *gorm.DB) *GormCategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (r *GormCategoryRepository) CreateCategory(ctx context.Context, c *domain.Category) error {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CategoryRepository.CreateCategory: %w", err)
	}

	m := GormCategory{
		Name:         c.Name,
		Description:  c.Description,
		DisplayOrder: c.DisplayOrder,
		Active:       c.Active,
		CompanyID:    companyID, // Auto-filled from context
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CategoryRepository.CreateCategory: %w", err)
	}
	c.ID = m.ID
	c.CompanyID = m.CompanyID
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *GormCategoryRepository) FindCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	var m GormCategory
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("CategoryRepository.FindCategoryByID: %w", err)
	}
	return categoryToDomain(&m), nil
}

func (r *GormCategoryRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	var ms []GormCategory
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Order("display_order ASC, name ASC").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("CategoryRepository.ListCategories: %w", err)
	}
	out := make([]domain.Category, len(ms))
	for i, m := range ms {
		out[i] = *categoryToDomain(&m)
	}
	return out, nil
}

func (r *GormCategoryRepository) UpdateCategory(ctx context.Context, c *domain.Category) error {
	// First, verify the category belongs to the tenant
	var existing GormCategory
	query := ApplyTenantFilterWithID(ctx, r.db, c.ID)
	if err := query.Where("deleted_at IS NULL").First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("CategoryRepository.UpdateCategory: categoria não encontrada ou acesso negado")
		}
		return fmt.Errorf("CategoryRepository.UpdateCategory: %w", err)
	}

	// Update without changing CompanyID (immutable)
	m := GormCategory{
		ID:           c.ID,
		Name:         c.Name,
		Description:  c.Description,
		DisplayOrder: c.DisplayOrder,
		Active:       c.Active,
		CompanyID:    existing.CompanyID, // Preserve original CompanyID
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("CategoryRepository.UpdateCategory: %w", err)
	}
	return nil
}

func (r *GormCategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	now := time.Now()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormCategory{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("CategoryRepository.DeleteCategory: %w", err)
	}
	return nil
}

func (r *GormCategoryRepository) CanDeleteCategory(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	check := &domain.DependencyCheck{CanDelete: true, Reasons: []domain.DependencyReason{}}

	// Verificar produtos que usam esta categoria (respecting tenant isolation)
	type ProductResult struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var products []ProductResult
	query := ApplyTenantFilter(ctx, r.db.Table("products"))
	if err := query.WithContext(ctx).
		Select("id, name").
		Where("category_id = ? AND deleted_at IS NULL", id).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("CategoryRepository.CanDeleteCategory: verificar produtos: %w", err)
	}

	for _, product := range products {
		check.CanDelete = false
		check.Reasons = append(check.Reasons, domain.DependencyReason{
			Type:        "product",
			ID:          product.ID,
			Name:        product.Name,
			Description: "Produto vinculado a esta categoria",
		})
	}

	return check, nil
}

func categoryToDomain(m *GormCategory) *domain.Category {
	var deletedAt *time.Time
	if m.DeletedAt != nil {
		dt := *m.DeletedAt
		deletedAt = &dt
	}
	return &domain.Category{
		ID:           m.ID,
		Name:         m.Name,
		Description:  m.Description,
		DisplayOrder: m.DisplayOrder,
		Active:       m.Active,
		CompanyID:    m.CompanyID,
		DeletedAt:    deletedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
