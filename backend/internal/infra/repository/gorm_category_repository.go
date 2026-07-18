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

var _ ports.CategoryRepository = (*GormCategoryRepository)(nil)

type GormCategoryRepository struct{ db *gorm.DB }

func NewGormCategoryRepository(db *gorm.DB) *GormCategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (r *GormCategoryRepository) CreateCategory(ctx context.Context, c *domain.Category) error {
	m := GormCategory{
		Name:         c.Name,
		Description:  c.Description,
		DisplayOrder: c.DisplayOrder,
		Active:       c.Active,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("CreateCategory: %w", err)
	}
	c.ID = m.ID
	c.CreatedAt = time.Unix(m.CreatedAt, 0)
	c.UpdatedAt = time.Unix(m.UpdatedAt, 0)
	return nil
}

func (r *GormCategoryRepository) FindCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	var m GormCategory
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindCategoryByID: %w", err)
	}
	return categoryToDomain(&m), nil
}

func (r *GormCategoryRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	var ms []GormCategory
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("display_order ASC, name ASC").Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("ListCategories: %w", err)
	}
	out := make([]domain.Category, len(ms))
	for i, m := range ms {
		out[i] = *categoryToDomain(&m)
	}
	return out, nil
}

func (r *GormCategoryRepository) UpdateCategory(ctx context.Context, c *domain.Category) error {
	m := GormCategory{
		ID:           c.ID,
		Name:         c.Name,
		Description:  c.Description,
		DisplayOrder: c.DisplayOrder,
		Active:       c.Active,
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("UpdateCategory: %w", err)
	}
	return nil
}

func (r *GormCategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	now := time.Now().Unix()
	if err := r.db.WithContext(ctx).Model(&GormCategory{}).
		Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteCategory: %w", err)
	}
	return nil
}

func (r *GormCategoryRepository) CanDeleteCategory(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	check := &domain.DependencyCheck{CanDelete: true, Reasons: []domain.DependencyReason{}}

	// Verificar produtos que usam esta categoria
	type ProductResult struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var products []ProductResult
	if err := r.db.WithContext(ctx).Table("products").
		Select("id, name").
		Where("category_id = ? AND deleted_at IS NULL", id).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("CanDeleteCategory: verificar produtos: %w", err)
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
		dt := time.Unix(*m.DeletedAt, 0)
		deletedAt = &dt
	}
	return &domain.Category{
		ID:           m.ID,
		Name:         m.Name,
		Description:  m.Description,
		DisplayOrder: m.DisplayOrder,
		Active:       m.Active,
		DeletedAt:    deletedAt,
		CreatedAt:    time.Unix(m.CreatedAt, 0),
		UpdatedAt:    time.Unix(m.UpdatedAt, 0),
	}
}
