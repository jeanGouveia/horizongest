package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type GormImpersonationAudit struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	PlatformUserID     uint       `gorm:"not null;index"`
	CompanyID          uint       `gorm:"not null;index"`
	CompanyOwnerUserID uint       `gorm:"not null;index"`
	StartedAt          time.Time  `gorm:"not null"`
	EndedAt            *time.Time `gorm:"index"`
	IPAddress          string     `gorm:"size:45"`
	UserAgent          string     `gorm:"type:text"`
	CreatedAt          time.Time  `gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime"`
	DeletedAt          *time.Time `gorm:"index"`
}

func (GormImpersonationAudit) TableName() string {
	return "impersonation_audit"
}

type GormImpersonationAuditRepository struct {
	db *gorm.DB
}

func NewGormImpersonationAuditRepository(db *gorm.DB) *GormImpersonationAuditRepository {
	return &GormImpersonationAuditRepository{db: db}
}

func (r *GormImpersonationAuditRepository) Create(ctx context.Context, audit *domain.ImpersonationAudit) error {
	gormAudit := r.toGorm(audit)
	if err := r.db.WithContext(ctx).Create(gormAudit).Error; err != nil {
		return fmt.Errorf("ImpersonationAuditRepository.Create: %w", err)
	}
	audit.ID = gormAudit.ID
	return nil
}

func (r *GormImpersonationAuditRepository) FindByID(ctx context.Context, id uint) (*domain.ImpersonationAudit, error) {
	var gormAudit GormImpersonationAudit
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormAudit).Error
	if err != nil {
		return nil, fmt.Errorf("ImpersonationAuditRepository.FindByID: %w", err)
	}
	return r.toDomain(&gormAudit), nil
}

func (r *GormImpersonationAuditRepository) FindByPlatformUserID(ctx context.Context, platformUserID uint) ([]*domain.ImpersonationAudit, error) {
	var gormAudits []GormImpersonationAudit
	err := r.db.WithContext(ctx).
		Where("platform_user_id = ?", platformUserID).
		Order("started_at DESC").
		Find(&gormAudits).Error
	if err != nil {
		return nil, fmt.Errorf("ImpersonationAuditRepository.FindByPlatformUserID: %w", err)
	}

	result := make([]*domain.ImpersonationAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormImpersonationAuditRepository) FindByCompanyID(ctx context.Context, companyID uint) ([]*domain.ImpersonationAudit, error) {
	var gormAudits []GormImpersonationAudit
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("started_at DESC").
		Find(&gormAudits).Error
	if err != nil {
		return nil, fmt.Errorf("ImpersonationAuditRepository.FindByCompanyID: %w", err)
	}

	result := make([]*domain.ImpersonationAudit, len(gormAudits))
	for i := range gormAudits {
		result[i] = r.toDomain(&gormAudits[i])
	}
	return result, nil
}

func (r *GormImpersonationAuditRepository) FindActiveByPlatformUserID(ctx context.Context, platformUserID uint) (*domain.ImpersonationAudit, error) {
	var gormAudit GormImpersonationAudit
	err := r.db.WithContext(ctx).
		Where("platform_user_id = ? AND ended_at IS NULL", platformUserID).
		Order("started_at DESC").
		First(&gormAudit).Error
	if err != nil {
		return nil, fmt.Errorf("ImpersonationAuditRepository.FindActiveByPlatformUserID: %w", err)
	}
	return r.toDomain(&gormAudit), nil
}

func (r *GormImpersonationAuditRepository) Update(ctx context.Context, audit *domain.ImpersonationAudit) error {
	gormAudit := r.toGorm(audit)
	if err := r.db.WithContext(ctx).Save(gormAudit).Error; err != nil {
		return fmt.Errorf("ImpersonationAuditRepository.Update: %w", err)
	}
	return nil
}

func (r *GormImpersonationAuditRepository) toGorm(audit *domain.ImpersonationAudit) *GormImpersonationAudit {
	return &GormImpersonationAudit{
		ID:                 audit.ID,
		PlatformUserID:     audit.PlatformUserID,
		CompanyID:          audit.CompanyID,
		CompanyOwnerUserID: audit.CompanyOwnerUserID,
		StartedAt:          audit.StartedAt,
		EndedAt:            audit.EndedAt,
		IPAddress:          audit.IPAddress,
		UserAgent:          audit.UserAgent,
		CreatedAt:          audit.CreatedAt,
		UpdatedAt:          audit.UpdatedAt,
		DeletedAt:          audit.DeletedAt,
	}
}

func (r *GormImpersonationAuditRepository) toDomain(gormAudit *GormImpersonationAudit) *domain.ImpersonationAudit {
	return &domain.ImpersonationAudit{
		ID:                 gormAudit.ID,
		PlatformUserID:     gormAudit.PlatformUserID,
		CompanyID:          gormAudit.CompanyID,
		CompanyOwnerUserID: gormAudit.CompanyOwnerUserID,
		StartedAt:          gormAudit.StartedAt,
		EndedAt:            gormAudit.EndedAt,
		IPAddress:          gormAudit.IPAddress,
		UserAgent:          gormAudit.UserAgent,
		CreatedAt:          gormAudit.CreatedAt,
		UpdatedAt:          gormAudit.UpdatedAt,
		DeletedAt:          gormAudit.DeletedAt,
	}
}
