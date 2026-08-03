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

type GormUserModel struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	Name         string     `gorm:"not null"`
	Email        string     `gorm:"uniqueIndex;not null"`
	PasswordHash string     `gorm:"not null"`
	Active       bool       `gorm:"not null;default:true"`
	CompanyID    uint       `gorm:"index;not null"` // Sprint 3: NOT NULL
	Role         string     `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt    *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
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
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	role := user.Role.String()
	model := GormUserModel{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Active:       user.Active,
		CompanyID:    companyID, // Preenchido automaticamente do contexto
		Role:         role,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

// CreateBootstrapOwner creates an owner user without requiring tenant context
// Used for platform bootstrap when creating the first company
func (r *GormUserRepository) CreateBootstrapOwner(ctx context.Context, user *domain.User, companyID uint) error {
	return r.CreateBootstrapOwnerWithTx(ctx, user, companyID, nil)
}

// CreateBootstrapOwnerWithTx creates an owner user with optional transaction
func (r *GormUserRepository) CreateBootstrapOwnerWithTx(ctx context.Context, user *domain.User, companyID uint, tx *gorm.DB) error {
	role := user.Role.String()
	model := GormUserModel{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Active:       user.Active,
		CompanyID:    companyID, // Fornecido explicitamente como parâmetro
		Role:         role,
	}

	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}

	if err := db.Create(&model).Error; err != nil {
		return fmt.Errorf("UserRepository.CreateBootstrapOwner: %w", err)
	}
	user.ID = model.ID
	user.CompanyID = companyID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model GormUserModel
	query := ApplyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&model).Error
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
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindByID: %w", err)
	}
	return toDomainUser(&model), nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
	role := user.Role.String()
	model := GormUserModel{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Active:       user.Active,
		CompanyID:    user.CompanyID,
		Role:         role,
	}

	query := ApplyTenantFilterWithID(ctx, r.db, user.ID)
	if err := query.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("UserRepository.Update: %w", err)
	}
	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	// Usar ApplyTenantFilter para isolamento de tenant e performance
	query := ApplyTenantFilter(ctx, r.db)
	var models []GormUserModel
	if err := query.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.List: %w", err)
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = toDomainUser(&model)
	}
	return users, nil
}

func toDomainUser(m *GormUserModel) *domain.User {
	role, _ := domain.ParseRole(m.Role)

	return &domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Active:       m.Active,
		CompanyID:    m.CompanyID,
		Role:         role,
		DeletedAt:    m.DeletedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
