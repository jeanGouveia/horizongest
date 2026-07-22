package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	FindByID(ctx context.Context, id uint) (*domain.Company, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Company, error)
	List(ctx context.Context) ([]domain.Company, error)
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id uint) error
}
