package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
