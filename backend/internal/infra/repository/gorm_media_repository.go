package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type GormMedia struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	FileName      string `gorm:"not null"`
	OriginalName  string `gorm:"not null"`
	FilePath      string `gorm:"not null"`
	ThumbnailPath string
	FileSize      int64  `gorm:"not null"`
	MimeType      string `gorm:"not null"`
	Width         *int
	Height        *int
	AltText       string
	EntityType    string `gorm:"not null;index"`
	EntityID      *uint  `gorm:"index"`
	CompanyID     uint   `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt     *time.Time `gorm:"index"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (GormMedia) TableName() string { return "media" }

// ─── Repository ─────────────────────────────────────────────────────────────

type GormMediaRepository struct{ db *gorm.DB }

func NewGormMediaRepository(db *gorm.DB) *GormMediaRepository {
	return &GormMediaRepository{db: db}
}

func (r *GormMediaRepository) CreateMedia(ctx context.Context, m *domain.Media) error {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateMedia: %w", err)
	}

	gm := GormMedia{
		FileName:      m.FileName,
		OriginalName:  m.OriginalName,
		FilePath:      m.FilePath,
		ThumbnailPath: m.ThumbnailPath,
		FileSize:      m.FileSize,
		MimeType:      m.MimeType,
		Width:         m.Width,
		Height:        m.Height,
		AltText:       m.AltText,
		EntityType:    m.EntityType,
		EntityID:      m.EntityID,
		CompanyID:     companyID, // Auto-filled from context
	}
	if err := r.db.WithContext(ctx).Create(&gm).Error; err != nil {
		return fmt.Errorf("CreateMedia: %w", err)
	}
	m.ID = gm.ID
	m.CompanyID = gm.CompanyID
	m.CreatedAt = gm.CreatedAt
	m.UpdatedAt = gm.UpdatedAt
	return nil
}

func (r *GormMediaRepository) FindMediaByID(ctx context.Context, id uint) (*domain.Media, error) {
	var gm GormMedia
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&gm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindMediaByID: %w", err)
	}
	return mediaToDomain(&gm), nil
}

func (r *GormMediaRepository) FindMediaByEntity(ctx context.Context, entityType string, entityID uint) ([]domain.Media, error) {
	var gms []GormMedia
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND deleted_at IS NULL", entityType, entityID).
		Find(&gms).Error; err != nil {
		return nil, fmt.Errorf("FindMediaByEntity: %w", err)
	}
	out := make([]domain.Media, len(gms))
	for i, gm := range gms {
		out[i] = *mediaToDomain(&gm)
	}
	return out, nil
}

func (r *GormMediaRepository) DeleteMedia(ctx context.Context, id uint) error {
	now := time.Now().Unix()
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormMedia{}).
		Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteMedia: %w", err)
	}
	return nil
}

func (r *GormMediaRepository) DeleteMediaByEntity(ctx context.Context, entityType string, entityID uint) error {
	now := time.Now().Unix()
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Model(&GormMedia{}).
		Where("entity_type = ? AND entity_id = ? AND deleted_at IS NULL", entityType, entityID).
		Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("DeleteMediaByEntity: %w", err)
	}
	return nil
}

// ── Mapper ─────────────────────────────────────────────────────────────────

func mediaToDomain(m *GormMedia) *domain.Media {
	var deletedAt *time.Time
	if m.DeletedAt != nil {
		dt := *m.DeletedAt
		deletedAt = &dt
	}
	return &domain.Media{
		ID:            m.ID,
		FileName:      m.FileName,
		OriginalName:  m.OriginalName,
		FilePath:      m.FilePath,
		ThumbnailPath: m.ThumbnailPath,
		FileSize:      m.FileSize,
		MimeType:      m.MimeType,
		Width:         m.Width,
		Height:        m.Height,
		AltText:       m.AltText,
		EntityType:    m.EntityType,
		EntityID:      m.EntityID,
		CompanyID:     m.CompanyID,
		DeletedAt:     deletedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
