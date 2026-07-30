package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type GormPlatformAudit struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	PlatformUserID *uint  `gorm:"index"`
	Action         string `gorm:"not null;index"`
	EntityType     string `gorm:"not null;index"`
	EntityID       *uint  `gorm:"index"`
	Changes        string `gorm:"type:text"`
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time `gorm:"autoCreateTime;index"`
}

func (GormPlatformAudit) TableName() string {
	return "platform_audit"
}

type GormPlatformAuditRepository struct {
	db *gorm.DB
}

func NewGormPlatformAuditRepository(db *gorm.DB) *GormPlatformAuditRepository {
	return &GormPlatformAuditRepository{db: db}
}

func (r *GormPlatformAuditRepository) Create(ctx context.Context, audit *domain.PlatformAudit) error {
	gormAudit := r.toGorm(audit)
	if err := r.db.WithContext(ctx).Create(gormAudit).Error; err != nil {
		return fmt.Errorf("PlatformAuditRepository.Create: %w", err)
	}
	audit.ID = gormAudit.ID
	return nil
}

func (r *GormPlatformAuditRepository) ListByPlatformUserID(ctx context.Context, platformUserID uint, limit int) ([]*domain.PlatformAudit, error) {
	var gormAudits []GormPlatformAudit
	err := r.db.WithContext(ctx).
		Where("platform_user_id = ?", platformUserID).
		Order("created_at DESC").
		Limit(limit).
		Find(&gormAudits).Error
	if err != nil {
		return nil, fmt.Errorf("PlatformAuditRepository.ListByPlatformUserID: %w", err)
	}

	result := make([]*domain.PlatformAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormPlatformAuditRepository) ListByEntity(ctx context.Context, entityType string, entityID uint, limit int) ([]*domain.PlatformAudit, error) {
	var gormAudits []GormPlatformAudit
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Limit(limit).
		Find(&gormAudits).Error
	if err != nil {
		return nil, fmt.Errorf("PlatformAuditRepository.ListByEntity: %w", err)
	}

	result := make([]*domain.PlatformAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormPlatformAuditRepository) ListRecent(ctx context.Context, limit int) ([]*domain.PlatformAudit, error) {
	var gormAudits []GormPlatformAudit
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&gormAudits).Error
	if err != nil {
		return nil, fmt.Errorf("PlatformAuditRepository.ListRecent: %w", err)
	}

	result := make([]*domain.PlatformAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormPlatformAuditRepository) toGorm(audit *domain.PlatformAudit) *GormPlatformAudit {
	return &GormPlatformAudit{
		ID:             audit.ID,
		PlatformUserID: audit.PlatformUserID,
		Action:         audit.Action,
		EntityType:     audit.EntityType,
		EntityID:       audit.EntityID,
		Changes:        audit.Changes,
		IPAddress:      audit.IPAddress,
		UserAgent:      audit.UserAgent,
		CreatedAt:      audit.CreatedAt,
	}
}

func (r *GormPlatformAuditRepository) toDomain(gormAudit *GormPlatformAudit) *domain.PlatformAudit {
	return &domain.PlatformAudit{
		ID:             gormAudit.ID,
		PlatformUserID: gormAudit.PlatformUserID,
		Action:         gormAudit.Action,
		EntityType:     gormAudit.EntityType,
		EntityID:       gormAudit.EntityID,
		Changes:        gormAudit.Changes,
		IPAddress:      gormAudit.IPAddress,
		UserAgent:      gormAudit.UserAgent,
		CreatedAt:      gormAudit.CreatedAt,
	}
}
