package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	// CreateBootstrapOwner creates an owner user without requiring tenant context
	// Used for platform bootstrap when creating the first company
	CreateBootstrapOwner(ctx context.Context, user *domain.User, companyID uint) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	List(ctx context.Context) ([]*domain.User, error)
}
