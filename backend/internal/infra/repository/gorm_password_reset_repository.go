package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

type GormPasswordResetTokenModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (GormPasswordResetTokenModel) TableName() string { return "password_reset_tokens" }

var _ ports.PasswordResetRepository = (*GormPasswordResetRepository)(nil)

type GormPasswordResetRepository struct {
	db *gorm.DB
}

func NewGormPasswordResetRepository(db *gorm.DB) *GormPasswordResetRepository {
	return &GormPasswordResetRepository{db: db}
}

func (r *GormPasswordResetRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	model := GormPasswordResetTokenModel{
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("PasswordResetRepository.Create: %w", err)
	}
	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	token.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormPasswordResetRepository) FindByToken(ctx context.Context, tokenStr string) (*domain.PasswordResetToken, error) {
	var model GormPasswordResetTokenModel
	err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("PasswordResetRepository.FindByToken: %w", err)
	}
	return toDomainPasswordResetToken(&model), nil
}

func (r *GormPasswordResetRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.PasswordResetToken, error) {
	var models []GormPasswordResetTokenModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("PasswordResetRepository.FindByUserID: %w", err)
	}
	tokens := make([]*domain.PasswordResetToken, len(models))
	for i, model := range models {
		tokens[i] = toDomainPasswordResetToken(&model)
	}
	return tokens, nil
}

func (r *GormPasswordResetRepository) MarkAsUsed(ctx context.Context, token *domain.PasswordResetToken) error {
	model := GormPasswordResetTokenModel{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Used:      true,
	}
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("PasswordResetRepository.MarkAsUsed: %w", err)
	}
	token.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormPasswordResetRepository) Delete(ctx context.Context, token *domain.PasswordResetToken) error {
	if err := r.db.WithContext(ctx).Delete(&GormPasswordResetTokenModel{}, token.ID).Error; err != nil {
		return fmt.Errorf("PasswordResetRepository.Delete: %w", err)
	}
	return nil
}

func (r *GormPasswordResetRepository) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&GormPasswordResetTokenModel{}).Error; err != nil {
		return fmt.Errorf("PasswordResetRepository.DeleteExpired: %w", err)
	}
	return nil
}

func toDomainPasswordResetToken(m *GormPasswordResetTokenModel) *domain.PasswordResetToken {
	return &domain.PasswordResetToken{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		Used:      m.Used,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
