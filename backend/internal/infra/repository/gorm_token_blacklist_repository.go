package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type GormTokenBlacklistRepository struct {
	db *gorm.DB
}

func NewGormTokenBlacklistRepository(db *gorm.DB) ports.TokenBlacklistRepository {
	return &GormTokenBlacklistRepository{db: db}
}

// GormTokenBlacklist é o modelo GORM para persistência
type GormTokenBlacklist struct {
	Token     string    `gorm:"primaryKey;size:500"`
	RevokedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}

func (r *GormTokenBlacklistRepository) Add(ctx context.Context, entry *domain.TokenBlacklist) error {
	gormEntry := GormTokenBlacklist{
		Token:     entry.Token,
		RevokedAt: entry.RevokedAt,
		ExpiresAt: entry.ExpiresAt,
	}
	if err := r.db.WithContext(ctx).Create(&gormEntry).Error; err != nil {
		return fmt.Errorf("TokenBlacklistRepository.Add: %w", err)
	}
	return nil
}

func (r *GormTokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	var entry GormTokenBlacklist
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("TokenBlacklistRepository.IsBlacklisted: %w", err)
	}

	// Se o token expirou, não está mais na blacklist efetivamente
	if time.Now().After(entry.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

func (r *GormTokenBlacklistRepository) CleanExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&GormTokenBlacklist{})
	if result.Error != nil {
		return fmt.Errorf("TokenBlacklistRepository.CleanExpired: %w", result.Error)
	}
	return nil
}
