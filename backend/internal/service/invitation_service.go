package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrInvitationNotFound        = errors.New("convite não encontrado")
	ErrInvitationExpired         = errors.New("convite expirado")
	ErrInvitationAlreadyUsed     = errors.New("convite já utilizado")
	ErrInvitationRevoked         = errors.New("convite revogado")
	ErrDuplicateInvitation       = errors.New("já existe um convite pendente para este e-mail")
	ErrUserBelongsToOtherCompany = errors.New("usuário já pertence a outra empresa")
)

const (
	// DefaultInvitationExpiration is the default expiration time for邀请 (7 days)
	DefaultInvitationExpiration = 7 * 24 * time.Hour
)

type InvitationService struct {
	invitationRepo ports.InvitationRepository
	userRepo       ports.UserRepository
	companyRepo    ports.CompanyRepository
	rbacService    *RBACService
}

func NewInvitationService(
	invitationRepo ports.InvitationRepository,
	userRepo ports.UserRepository,
	companyRepo ports.CompanyRepository,
	rbacService *RBACService,
) *InvitationService {
	return &InvitationService{
		invitationRepo: invitationRepo,
		userRepo:       userRepo,
		companyRepo:    companyRepo,
		rbacService:    rbacService,
	}
}

// InvitationOutput represents invitation data for API responses
type InvitationOutput struct {
	ID         uint
	CompanyID  uint
	Email      string
	Role       string
	Token      string
	Status     string
	ExpiresAt  time.Time
	AcceptedAt *time.Time
	CreatedBy  uint
	CreatedAt  time.Time
}

// CreateInvitation creates a new invitation for a user to join a company
func (s *InvitationService) CreateInvitation(ctx context.Context, actorUserID uint, companyID uint, email string, role domain.Role) (*InvitationOutput, error) {
	// RBAC Validation: Only Owner and Admin can create invitations
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.CreateInvitation: verificar permissões: %w", err)
	}
	if !canManage {
		return nil, ErrPermissionDenied
	}

	// Check if user already exists and belongs to this company
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.CreateInvitation: verificar usuário: %w", err)
	}
	if user != nil {
		if user.CompanyID == companyID {
			return nil, ErrUserAlreadyInCompany
		}
		if user.CompanyID != 0 {
			return nil, ErrUserBelongsToOtherCompany
		}
	}

	// Check for duplicate pending invitation
	existing, err := s.invitationRepo.FindByEmailAndCompanyID(ctx, email, companyID)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.CreateInvitation: verificar convite existente: %w", err)
	}
	if existing != nil {
		return nil, ErrDuplicateInvitation
	}

	// Generate token
	token, err := domain.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("InvitationService.CreateInvitation: gerar token: %w", err)
	}

	// Create invitation
	invitation := &domain.Invitation{
		CompanyID: companyID,
		Email:     email,
		Role:      role,
		Token:     token,
		Status:    domain.InvitationStatusPending,
		ExpiresAt: time.Now().Add(DefaultInvitationExpiration),
		CreatedBy: actorUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, fmt.Errorf("InvitationService.CreateInvitation: criar convite: %w", err)
	}

	return s.toInvitationOutput(invitation), nil
}

// ListInvitations returns all invitations for a company
func (s *InvitationService) ListInvitations(ctx context.Context, companyID uint) ([]*InvitationOutput, error) {
	invitations, err := s.invitationRepo.ListByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.ListInvitations: listar convites: %w", err)
	}

	output := make([]*InvitationOutput, len(invitations))
	for i, inv := range invitations {
		output[i] = s.toInvitationOutput(inv)
	}
	return output, nil
}

// GetInvitation returns a specific invitation by ID
func (s *InvitationService) GetInvitation(ctx context.Context, companyID uint, invitationID uint) (*InvitationOutput, error) {
	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.GetInvitation: buscar convite: %w", err)
	}
	if invitation == nil {
		return nil, ErrInvitationNotFound
	}

	// Check if invitation belongs to the company
	if invitation.CompanyID != companyID {
		return nil, ErrInvitationNotFound
	}

	return s.toInvitationOutput(invitation), nil
}

// RevokeInvitation revokes an invitation
func (s *InvitationService) RevokeInvitation(ctx context.Context, actorUserID uint, companyID uint, invitationID uint) error {
	// RBAC Validation: Only Owner and Admin can revoke invitations
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return fmt.Errorf("InvitationService.RevokeInvitation: verificar permissões: %w", err)
	}
	if !canManage {
		return ErrPermissionDenied
	}

	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("InvitationService.RevokeInvitation: buscar convite: %w", err)
	}
	if invitation == nil {
		return ErrInvitationNotFound
	}

	// Check if invitation belongs to the company
	if invitation.CompanyID != companyID {
		return ErrInvitationNotFound
	}

	// Cannot revoke already accepted invitations
	if invitation.Status == domain.InvitationStatusAccepted {
		return errors.New("não é possível revogar convite já aceito")
	}

	// Revoke invitation
	invitation.Status = domain.InvitationStatusRevoked
	invitation.UpdatedAt = time.Now()

	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		return fmt.Errorf("InvitationService.RevokeInvitation: atualizar convite: %w", err)
	}

	return nil
}

// GetInvitationByToken returns an invitation by its token (public endpoint)
func (s *InvitationService) GetInvitationByToken(ctx context.Context, token string) (*InvitationOutput, error) {
	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("InvitationService.GetInvitationByToken: buscar convite: %w", err)
	}
	if invitation == nil {
		return nil, ErrInvitationNotFound
	}

	// Check if invitation is expired
	if invitation.IsExpired() {
		// Auto-expire the invitation
		invitation.Status = domain.InvitationStatusExpired
		invitation.UpdatedAt = time.Now()
		_ = s.invitationRepo.Update(ctx, invitation)
		return nil, ErrInvitationExpired
	}

	return s.toInvitationOutput(invitation), nil
}

// AcceptInvitation accepts an invitation and associates the user with the company
func (s *InvitationService) AcceptInvitation(ctx context.Context, token string, userEmail string) error {
	// FASE A: Use FindByTokenForUpdate with SELECT FOR UPDATE to prevent race condition
	invitation, err := s.invitationRepo.FindByTokenForUpdate(ctx, token)
	if err != nil {
		return fmt.Errorf("InvitationService.AcceptInvitation: buscar convite: %w", err)
	}
	if invitation == nil {
		return ErrInvitationNotFound
	}

	// Validate that the authenticated user's email matches the invitation email
	if invitation.Email != userEmail {
		return errors.New("o convite não pertence a este usuário")
	}

	// Check if invitation can be accepted (double-check after lock)
	if !invitation.CanBeAccepted() {
		if invitation.Status == domain.InvitationStatusAccepted {
			return ErrInvitationAlreadyUsed
		}
		if invitation.Status == domain.InvitationStatusRevoked {
			return ErrInvitationRevoked
		}
		if invitation.IsExpired() {
			return ErrInvitationExpired
		}
		return errors.New("convite não pode ser aceito")
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, invitation.Email)
	if err != nil {
		return fmt.Errorf("InvitationService.AcceptInvitation: buscar usuário: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Check if user already belongs to a company
	if user.CompanyID != 0 {
		return ErrUserBelongsToOtherCompany
	}

	// Associate user with company
	user.CompanyID = invitation.CompanyID
	user.Role = invitation.Role
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("InvitationService.AcceptInvitation: atualizar usuário: %w", err)
	}

	// Mark invitation as accepted
	now := time.Now()
	invitation.Status = domain.InvitationStatusAccepted
	invitation.AcceptedAt = &now
	invitation.UpdatedAt = now

	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		return fmt.Errorf("InvitationService.AcceptInvitation: atualizar convite: %w", err)
	}

	return nil
}

func (s *InvitationService) toInvitationOutput(invitation *domain.Invitation) *InvitationOutput {
	return &InvitationOutput{
		ID:         invitation.ID,
		CompanyID:  invitation.CompanyID,
		Email:      invitation.Email,
		Role:       invitation.Role.String(),
		Token:      invitation.Token,
		Status:     invitation.Status.String(),
		ExpiresAt:  invitation.ExpiresAt,
		AcceptedAt: invitation.AcceptedAt,
		CreatedBy:  invitation.CreatedBy,
		CreatedAt:  invitation.CreatedAt,
	}
}
