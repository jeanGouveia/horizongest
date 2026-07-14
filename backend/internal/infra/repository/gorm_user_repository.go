package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

type GormUserModel struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
}

func (GormUserModel) TableName() string { return "users" }

var _ ports.UserRepository = (*GormUserRepository)(nil)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	model := GormUserModel{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	user.ID = model.ID
	user.CreatedAt = time.Unix(model.CreatedAt, 0)
	user.UpdatedAt = time.Unix(model.UpdatedAt, 0)
	return nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model GormUserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindByEmail: %w", err)
	}
	return toDomainUser(&model), nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var model GormUserModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindByID: %w", err)
	}
	return toDomainUser(&model), nil
}

func toDomainUser(m *GormUserModel) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    time.Unix(m.CreatedAt, 0),
		UpdatedAt:    time.Unix(m.UpdatedAt, 0),
	}
}
