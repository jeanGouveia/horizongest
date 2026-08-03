package repository

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

// AuditRepository implements the audit log repository
// FASE A.3: B17 - Audit log repository implementation
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *gorm.DB) ports.AuditRepository {
	return &AuditRepository{db: db}
}

// CreateAuditLog creates a new audit log entry
func (r *AuditRepository) CreateAuditLog(ctx context.Context, audit *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(audit).Error
}

// FindAuditLogsByCompany retrieves audit logs for a company
func (r *AuditRepository) FindAuditLogsByCompany(ctx context.Context, companyID uint, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	query := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("timestamp DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&logs).Error
	return logs, err
}

// FindAuditLogsByUser retrieves audit logs for a specific user
func (r *AuditRepository) FindAuditLogsByUser(ctx context.Context, userID uint, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("timestamp DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&logs).Error
	return logs, err
}
