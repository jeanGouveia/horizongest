package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	FindByID(ctx context.Context, id uint) (*domain.Plan, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]*domain.Plan, error)
	ListActive(ctx context.Context) ([]*domain.Plan, error)
}
