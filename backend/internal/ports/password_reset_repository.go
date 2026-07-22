package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, token *domain.PasswordResetToken) error
	FindByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, token *domain.PasswordResetToken) error
	Delete(ctx context.Context, token *domain.PasswordResetToken) error
	DeleteExpired(ctx context.Context) error
}
