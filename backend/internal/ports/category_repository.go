package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, c *domain.Category) error
	FindCategoryByID(ctx context.Context, id uint) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, c *domain.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	CanDeleteCategory(ctx context.Context, id uint) (*domain.DependencyCheck, error)
}
