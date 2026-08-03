package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type InvitationRepository interface {
	Create(ctx context.Context, invitation *domain.Invitation) error
	FindByID(ctx context.Context, id uint) (*domain.Invitation, error)
	FindByToken(ctx context.Context, token string) (*domain.Invitation, error)
	FindByTokenForUpdate(ctx context.Context, token string) (*domain.Invitation, error) // FASE A: For race condition fix
	FindByEmailAndCompanyID(ctx context.Context, email string, companyID uint) (*domain.Invitation, error)
	ListByCompanyID(ctx context.Context, companyID uint) ([]*domain.Invitation, error)
	Update(ctx context.Context, invitation *domain.Invitation) error
	Delete(ctx context.Context, id uint) error
}
