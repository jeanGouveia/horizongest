package repository

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type GormPlatformUser struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"not null;uniqueIndex"`
	PasswordHash  string `gorm:"not null"`
	Role          string `gorm:"not null;default:'PlatformSupport'"`
	Active        bool   `gorm:"not null;default:true"`
	DeletedAt     *int64 `gorm:"index"`
	CreatedAt     int64  `gorm:"autoCreateTime"`
	UpdatedAt     int64  `gorm:"autoUpdateTime"`
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
		return err
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
		return nil, err
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
		return nil, err
	}
	return r.toDomain(&gormUser), nil
}

func (r *GormPlatformUserRepository) Update(ctx context.Context, user *domain.PlatformUser) error {
	gormUser := r.toGorm(user)
	return r.db.WithContext(ctx).Save(gormUser).Error
}

func (r *GormPlatformUserRepository) toGorm(user *domain.PlatformUser) *GormPlatformUser {
	var deletedAt *int64
	if user.DeletedAt != nil {
		t := user.DeletedAt.Unix()
		deletedAt = &t
	}
	return &GormPlatformUser{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role.String(),
		Active:       user.Active,
		DeletedAt:    deletedAt,
		CreatedAt:    user.CreatedAt.Unix(),
		UpdatedAt:    user.UpdatedAt.Unix(),
	}
}

func (r *GormPlatformUserRepository) toDomain(gormUser *GormPlatformUser) *domain.PlatformUser {
	var deletedAt *time.Time
	if gormUser.DeletedAt != nil {
		t := time.Unix(*gormUser.DeletedAt, 0)
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
		CreatedAt:    time.Unix(gormUser.CreatedAt, 0),
		UpdatedAt:    time.Unix(gormUser.UpdatedAt, 0),
	}
}

// CreateInitialAdmin creates the initial platform admin if it doesn't exist
func (r *GormPlatformUserRepository) CreateInitialAdmin(ctx context.Context, email, password string) error {
	// Check if admin already exists
	existing, err := r.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil // Already exists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
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
