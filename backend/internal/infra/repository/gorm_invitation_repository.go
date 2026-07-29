package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type GormInvitationModel struct {
	ID         uint       `gorm:"primaryKey;autoIncrement"`
	CompanyID  uint       `gorm:"not null;index"`
	Email      string     `gorm:"not null;index"`
	Role       string     `gorm:"not null"`
	Token      string     `gorm:"not null;uniqueIndex"`
	Status     string     `gorm:"not null;index;default:'pending'"`
	ExpiresAt  time.Time  `gorm:"not null;index"`
	AcceptedAt *time.Time `gorm:"index"`
	CreatedBy  uint       `gorm:"not null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}

func (GormInvitationModel) TableName() string { return "invitations" }

var _ ports.InvitationRepository = (*GormInvitationRepository)(nil)

type GormInvitationRepository struct {
	db *gorm.DB
}

func NewGormInvitationRepository(db *gorm.DB) *GormInvitationRepository {
	return &GormInvitationRepository{db: db}
}

func (r *GormInvitationRepository) Create(ctx context.Context, invitation *domain.Invitation) error {
	model := toGormInvitation(invitation)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("InvitationRepository.Create: %w", err)
	}
	invitation.ID = model.ID
	invitation.CreatedAt = model.CreatedAt
	invitation.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormInvitationRepository) FindByID(ctx context.Context, id uint) (*domain.Invitation, error) {
	var model GormInvitationModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("InvitationRepository.FindByID: %w", err)
	}
	return toDomainInvitation(&model), nil
}

func (r *GormInvitationRepository) FindByToken(ctx context.Context, token string) (*domain.Invitation, error) {
	var model GormInvitationModel
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("InvitationRepository.FindByToken: %w", err)
	}
	return toDomainInvitation(&model), nil
}

func (r *GormInvitationRepository) FindByEmailAndCompanyID(ctx context.Context, email string, companyID uint) (*domain.Invitation, error) {
	var model GormInvitationModel
	if err := r.db.WithContext(ctx).
		Where("email = ? AND company_id = ? AND status = ?", email, companyID, domain.InvitationStatusPending).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("InvitationRepository.FindByEmailAndCompanyID: %w", err)
	}
	return toDomainInvitation(&model), nil
}

func (r *GormInvitationRepository) ListByCompanyID(ctx context.Context, companyID uint) ([]*domain.Invitation, error) {
	var models []GormInvitationModel
	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("InvitationRepository.ListByCompanyID: %w", err)
	}

	invitations := make([]*domain.Invitation, len(models))
	for i, model := range models {
		invitations[i] = toDomainInvitation(&model)
	}
	return invitations, nil
}

func (r *GormInvitationRepository) Update(ctx context.Context, invitation *domain.Invitation) error {
	model := toGormInvitation(invitation)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("InvitationRepository.Update: %w", err)
	}
	invitation.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormInvitationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&GormInvitationModel{}, id).Error; err != nil {
		return fmt.Errorf("InvitationRepository.Delete: %w", err)
	}
	return nil
}

func toGormInvitation(invitation *domain.Invitation) GormInvitationModel {
	model := GormInvitationModel{
		ID:        invitation.ID,
		CompanyID: invitation.CompanyID,
		Email:     invitation.Email,
		Role:      invitation.Role.String(),
		Token:     invitation.Token,
		Status:    invitation.Status.String(),
		ExpiresAt: invitation.ExpiresAt,
		CreatedBy: invitation.CreatedBy,
		CreatedAt: invitation.CreatedAt,
		UpdatedAt: invitation.UpdatedAt,
	}

	if invitation.AcceptedAt != nil {
		model.AcceptedAt = invitation.AcceptedAt
	}

	return model
}

func toDomainInvitation(model *GormInvitationModel) *domain.Invitation {
	var acceptedAt *time.Time
	if model.AcceptedAt != nil {
		at := *model.AcceptedAt
		acceptedAt = &at
	}

	status, _ := domain.ParseInvitationStatus(model.Status)
	role, _ := domain.ParseRole(model.Role)

	return &domain.Invitation{
		ID:         model.ID,
		CompanyID:  model.CompanyID,
		Email:      model.Email,
		Role:       role,
		Token:      model.Token,
		Status:     status,
		ExpiresAt:  model.ExpiresAt,
		AcceptedAt: acceptedAt,
		CreatedBy:  model.CreatedBy,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}
