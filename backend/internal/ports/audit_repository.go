package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// AuditRepository defines the interface for audit log persistence
// FASE A.3: B17 - Audit repository interface
type AuditRepository interface {
	CreateAuditLog(ctx context.Context, audit *domain.AuditLog) error
	FindAuditLogsByCompany(ctx context.Context, companyID uint, limit int) ([]domain.AuditLog, error)
	FindAuditLogsByUser(ctx context.Context, userID uint, limit int) ([]domain.AuditLog, error)
}
