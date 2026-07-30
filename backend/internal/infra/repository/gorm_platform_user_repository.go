package repository

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type GormPlatformUser struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Name         string     `gorm:"not null"`
	Email        string     `gorm:"not null;uniqueIndex"`
	PasswordHash string     `gorm:"not null"`
	Role         string     `gorm:"not null;default:'PlatformSupport'"`
	Active       bool       `gorm:"not null;default:true"`
	DeletedAt    *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

func (GormPlatformUser) TableName() string {
	return "platform_users"
}

type GormPlatformUserRepository struct {
	db *gorm.DB
}

func NewGormPlatformUserRepository(db *gorm.DB) *GormPlatformUserRepository {
	return &GormPlatformUserRepository{db: db}
}

func (r *GormPlatformUserRepository) Create(ctx context.Context, user *domain.PlatformUser) error {
	gormUser := r.toGorm(user)
	if err := r.db.WithContext(ctx).Create(gormUser).Error; err != nil {
		return fmt.Errorf("PlatformUserRepository.Create: %w", err)
	}
	user.ID = gormUser.ID
	return nil
}

func (r *GormPlatformUserRepository) FindByEmail(ctx context.Context, email string) (*domain.PlatformUser, error) {
	var gormUser GormPlatformUser
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&gormUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("PlatformUserRepository.FindByEmail: %w", err)
	}
	return r.toDomain(&gormUser), nil
}

func (r *GormPlatformUserRepository) FindByID(ctx context.Context, id uint) (*domain.PlatformUser, error) {
	var gormUser GormPlatformUser
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&gormUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("PlatformUserRepository.FindByID: %w", err)
	}
	return r.toDomain(&gormUser), nil
}

func (r *GormPlatformUserRepository) Update(ctx context.Context, user *domain.PlatformUser) error {
	gormUser := r.toGorm(user)
	return fmt.Errorf("PlatformUserRepository.Update: %w", r.db.WithContext(ctx).Save(gormUser).Error)
}

func (r *GormPlatformUserRepository) toGorm(user *domain.PlatformUser) *GormPlatformUser {
	return &GormPlatformUser{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role.String(),
		Active:       user.Active,
		DeletedAt:    user.DeletedAt,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func (r *GormPlatformUserRepository) toDomain(gormUser *GormPlatformUser) *domain.PlatformUser {
	var deletedAt *time.Time
	if gormUser.DeletedAt != nil {
		t := *gormUser.DeletedAt
		deletedAt = &t
	}

	role, _ := domain.ParsePlatformRole(gormUser.Role)

	return &domain.PlatformUser{
		ID:           gormUser.ID,
		Name:         gormUser.Name,
		Email:        gormUser.Email,
		PasswordHash: gormUser.PasswordHash,
		Role:         role,
		Active:       gormUser.Active,
		DeletedAt:    deletedAt,
		CreatedAt:    gormUser.CreatedAt,
		UpdatedAt:    gormUser.UpdatedAt,
	}
}

// CreateInitialAdmin creates the initial platform admin if it doesn't exist
func (r *GormPlatformUserRepository) CreateInitialAdmin(ctx context.Context, email, password string) error {
	// Check if admin already exists
	existing, err := r.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("PlatformUserRepository.CreateInitialAdmin: %w", err)
	}
	if existing != nil {
		return nil // Already exists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("PlatformUserRepository.CreateInitialAdmin: %w", err)
	}

	// Create admin
	admin := &domain.PlatformUser{
		Name:         "Platform Admin",
		Email:        email,
		PasswordHash: string(hash),
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return r.Create(ctx, admin)
}
