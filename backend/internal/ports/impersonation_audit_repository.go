package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type ImpersonationAuditRepository interface {
	Create(ctx context.Context, audit *domain.ImpersonationAudit) error
	FindByID(ctx context.Context, id uint) (*domain.ImpersonationAudit, error)
	FindByPlatformUserID(ctx context.Context, platformUserID uint) ([]*domain.ImpersonationAudit, error)
	FindByCompanyID(ctx context.Context, companyID uint) ([]*domain.ImpersonationAudit, error)
	FindActiveByPlatformUserID(ctx context.Context, platformUserID uint) (*domain.ImpersonationAudit, error)
	Update(ctx context.Context, audit *domain.ImpersonationAudit) error
}
