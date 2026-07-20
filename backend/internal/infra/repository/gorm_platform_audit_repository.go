package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type GormPlatformAudit struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	PlatformUserID *uint `gorm:"index"`
	Action        string `gorm:"not null;index"`
	EntityType    string `gorm:"not null;index"`
	EntityID      *uint  `gorm:"index"`
	Changes       string `gorm:"type:text"`
	IPAddress     string
	UserAgent     string
	CreatedAt     int64  `gorm:"autoCreateTime;index"`
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
		return err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	
	result := make([]*domain.PlatformAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormPlatformAuditRepository) toGorm(audit *domain.PlatformAudit) *GormPlatformAudit {
	return &GormPlatformAudit{
		ID:            audit.ID,
		PlatformUserID: audit.PlatformUserID,
		Action:        audit.Action,
		EntityType:    audit.EntityType,
		EntityID:      audit.EntityID,
		Changes:       audit.Changes,
		IPAddress:     audit.IPAddress,
		UserAgent:     audit.UserAgent,
		CreatedAt:     audit.CreatedAt.Unix(),
	}
}

func (r *GormPlatformAuditRepository) toDomain(gormAudit *GormPlatformAudit) *domain.PlatformAudit {
	return &domain.PlatformAudit{
		ID:            gormAudit.ID,
		PlatformUserID: gormAudit.PlatformUserID,
		Action:        gormAudit.Action,
		EntityType:    gormAudit.EntityType,
		EntityID:      gormAudit.EntityID,
		Changes:       gormAudit.Changes,
		IPAddress:     gormAudit.IPAddress,
		UserAgent:     gormAudit.UserAgent,
		CreatedAt:     time.Unix(gormAudit.CreatedAt, 0),
	}
}
