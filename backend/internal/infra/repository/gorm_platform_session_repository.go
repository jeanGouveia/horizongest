package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type GormPlatformSession struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	PlatformUserID uint      `gorm:"not null;index"`
	Token          string    `gorm:"not null;uniqueIndex"`
	ExpiresAt      time.Time `gorm:"not null;index"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (GormPlatformSession) TableName() string {
	return "platform_sessions"
}

type GormPlatformSessionRepository struct {
	db *gorm.DB
}

func NewGormPlatformSessionRepository(db *gorm.DB) *GormPlatformSessionRepository {
	return &GormPlatformSessionRepository{db: db}
}

func (r *GormPlatformSessionRepository) Create(ctx context.Context, session *domain.PlatformSession) error {
	gormSession := r.toGorm(session)
	if err := r.db.WithContext(ctx).Create(gormSession).Error; err != nil {
		return fmt.Errorf("PlatformSessionRepository.Create: %w", err)
	}
	session.ID = gormSession.ID
	return nil
}

func (r *GormPlatformSessionRepository) FindByToken(ctx context.Context, token string) (*domain.PlatformSession, error) {
	var gormSession GormPlatformSession
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&gormSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("PlatformSessionRepository.FindByToken: %w", err)
	}
	return r.toDomain(&gormSession), nil
}

func (r *GormPlatformSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return fmt.Errorf("PlatformSessionRepository.DeleteByToken: %w", r.db.WithContext(ctx).Where("token = ?", token).Delete(&GormPlatformSession{}).Error)
}

func (r *GormPlatformSessionRepository) DeleteByPlatformUserID(ctx context.Context, platformUserID uint) error {
	return fmt.Errorf("PlatformSessionRepository.DeleteByPlatformUserID: %w", r.db.WithContext(ctx).Where("platform_user_id = ?", platformUserID).Delete(&GormPlatformSession{}).Error)
}

func (r *GormPlatformSessionRepository) DeleteExpired(ctx context.Context) error {
	return fmt.Errorf("PlatformSessionRepository.DeleteExpired: %w", r.db.WithContext(ctx).Where("expires_at < ?", time.Now().Unix()).Delete(&GormPlatformSession{}).Error)
}

func (r *GormPlatformSessionRepository) toGorm(session *domain.PlatformSession) *GormPlatformSession {
	return &GormPlatformSession{
		ID:             session.ID,
		PlatformUserID: session.PlatformUserID,
		Token:          session.Token,
		ExpiresAt:      session.ExpiresAt,
		CreatedAt:      session.CreatedAt,
	}
}

func (r *GormPlatformSessionRepository) toDomain(gormSession *GormPlatformSession) *domain.PlatformSession {
	return &domain.PlatformSession{
		ID:             gormSession.ID,
		PlatformUserID: gormSession.PlatformUserID,
		Token:          gormSession.Token,
		ExpiresAt:      gormSession.ExpiresAt,
		CreatedAt:      gormSession.CreatedAt,
	}
}
