package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var (
	ErrUserNotFound         = errors.New("usuário não encontrado")
	ErrUserNotInCompany     = errors.New("usuário não pertence à empresa")
	ErrCannotAlterOwner     = errors.New("apenas Owner pode alterar papel de Owner")
	ErrCannotAlterAdmin     = errors.New("apenas Owner pode alterar papel de Admin")
	ErrCannotRemoveOwner    = errors.New("não é possível remover Owner da empresa")
	ErrUserAlreadyInCompany = errors.New("usuário já pertence a esta empresa")
	ErrPermissionDenied     = errors.New("permissão negada")
	ErrCannotDeactivateSelf = errors.New("não é possível desativar o próprio usuário")
)

type UserManagementService struct {
	userRepo    ports.UserRepository
	companyRepo ports.CompanyRepository
	rbacService *RBACService
}

func NewUserManagementService(userRepo ports.UserRepository, companyRepo ports.CompanyRepository, rbacService *RBACService) *UserManagementService {
	return &UserManagementService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		rbacService: rbacService,
	}
}

// UserOutput represents user data for API responses
type UserOutput struct {
	ID        uint
	Name      string
	Email     string
	Role      *string
	Active    bool
	CompanyID *uint
}

// ListUsers returns all users in the company
func (s *UserManagementService) ListUsers(ctx context.Context, companyID uint) ([]UserOutput, error) {
	// Get all users (repository should filter by companyID)
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: failed to list users: %w", err)
	}

	var result []UserOutput
	for _, user := range users {
		// Filter by companyID
		if user.CompanyID == companyID {
			role := user.Role.String()
			result = append(result, UserOutput{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Role:      &role,
				Active:    user.Active,
				CompanyID: &user.CompanyID,
			})
		}
	}

	return result, nil
}

// GetUser returns a specific user in the company
func (s *UserManagementService) GetUser(ctx context.Context, companyID uint, userID uint) (*UserOutput, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUser: failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if user belongs to the company
	if user.CompanyID != companyID {
		return nil, ErrUserNotInCompany
	}

	role := user.Role.String()

	return &UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      &role,
		Active:    user.Active,
		CompanyID: &user.CompanyID,
	}, nil
}

// ChangeRole changes a user's role within the company
func (s *UserManagementService) ChangeRole(ctx context.Context, actorUserID uint, targetUserID uint, newRole domain.Role) error {
	// Get actor user
	actor, err := s.userRepo.FindByID(ctx, actorUserID)
	if err != nil {
		return fmt.Errorf("ChangeRole: failed to get actor: %w", err)
	}
	if actor == nil {
		return ErrPermissionDenied
	}

	// Get target user
	target, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("ChangeRole: failed to get target: %w", err)
	}
	if target == nil {
		return ErrUserNotFound
	}

	// Check if target belongs to same company
	if target.CompanyID != actor.CompanyID {
		return ErrUserNotInCompany
	}

	// RBAC Validation: Only Owner can alter Owner role
	if target.Role == domain.RoleOwner {
		canAlter, err := s.rbacService.CanAlterOwnerRole(ctx, actorUserID)
		if err != nil {
			return err
		}
		if !canAlter {
			return ErrCannotAlterOwner
		}
	}

	// RBAC Validation: Only Owner can alter Admin role
	if target.Role == domain.RoleAdmin {
		canAlter, err := s.rbacService.CanAlterAdminRole(ctx, actorUserID)
		if err != nil {
			return err
		}
		if !canAlter {
			return ErrCannotAlterAdmin
		}
	}

	// RBAC Validation: Only Owner and Admin can change roles
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return ErrPermissionDenied
	}

	// Update target user's role
	target.Role = newRole
	if err := s.userRepo.Update(ctx, target); err != nil {
		return fmt.Errorf("ChangeRole: failed to update user: %w", err)
	}

	return nil
}

// RemoveFromCompany removes a user from the company (does not delete the user)
func (s *UserManagementService) RemoveFromCompany(ctx context.Context, actorUserID uint, targetUserID uint) error {
	// Get actor user
	actor, err := s.userRepo.FindByID(ctx, actorUserID)
	if err != nil {
		return fmt.Errorf("RemoveFromCompany: failed to get actor: %w", err)
	}
	if actor == nil {
		return ErrPermissionDenied
	}

	// Get target user
	target, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("RemoveFromCompany: failed to get target: %w", err)
	}
	if target == nil {
		return ErrUserNotFound
	}

	// Check if target belongs to same company
	if target.CompanyID != actor.CompanyID {
		return ErrUserNotInCompany
	}

	// RBAC Validation: Cannot remove Owner from company
	if target.Role == domain.RoleOwner {
		return ErrCannotRemoveOwner
	}

	// RBAC Validation: Only Owner and Admin can remove users
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return ErrPermissionDenied
	}

	// Remove user from company - NOT ALLOWED in Sprint 3
	// Users must always belong to a company. Use deletion instead.
	return errors.New("RemoveFromCompany não permitido - usuários devem pertencer a uma empresa")
}

// AddExistingUser adds an existing user to the company by email
func (s *UserManagementService) AddExistingUser(ctx context.Context, actorUserID uint, email string) (*UserOutput, error) {
	// Get actor user
	actor, err := s.userRepo.FindByID(ctx, actorUserID)
	if err != nil {
		return nil, fmt.Errorf("AddExistingUser: failed to get actor: %w", err)
	}
	if actor == nil {
		return nil, ErrPermissionDenied
	}

	// RBAC Validation: Only Owner and Admin can add users
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, ErrPermissionDenied
	}

	// Find user by email
	target, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("AddExistingUser: failed to find user: %w", err)
	}
	if target == nil {
		return nil, ErrUserNotFound
	}

	// Check if user already belongs to this company
	if target.CompanyID == actor.CompanyID {
		return nil, ErrUserAlreadyInCompany
	}

	// Check if user belongs to another company
	if target.CompanyID != 0 {
		return nil, errors.New("usuário já pertence a outra empresa")
	}

	// Add user to company with default role (Manager)
	target.CompanyID = actor.CompanyID
	defaultRole := domain.RoleManager
	target.Role = defaultRole

	if err := s.userRepo.Update(ctx, target); err != nil {
		return nil, fmt.Errorf("AddExistingUser: failed to update user: %w", err)
	}

	role := target.Role.String()

	return &UserOutput{
		ID:        target.ID,
		Name:      target.Name,
		Email:     target.Email,
		Role:      &role,
		Active:    target.Active,
		CompanyID: &target.CompanyID,
	}, nil
}

// SetUserActive sets the active status of a user
func (s *UserManagementService) SetUserActive(ctx context.Context, actorUserID uint, targetUserID uint, active bool) error {
	// Get actor user
	actor, err := s.userRepo.FindByID(ctx, actorUserID)
	if err != nil {
		return fmt.Errorf("SetUserActive: failed to get actor: %w", err)
	}
	if actor == nil {
		return ErrPermissionDenied
	}

	// Check if actor is trying to deactivate themselves
	if actorUserID == targetUserID && !active {
		return ErrCannotDeactivateSelf
	}

	// Get target user
	target, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("SetUserActive: failed to get target: %w", err)
	}
	if target == nil {
		return ErrUserNotFound
	}

	// Check if target belongs to same company
	if target.CompanyID != actor.CompanyID {
		return ErrUserNotInCompany
	}

	// RBAC Validation: Cannot deactivate Owner
	if !active && target.Role == domain.RoleOwner {
		return errors.New("não é possível desativar Owner da empresa")
	}

	// RBAC Validation: Only Owner and Admin can activate/deactivate users
	canManage, err := s.rbacService.CanManageUsers(ctx, actorUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return ErrPermissionDenied
	}

	// Update target user's active status
	target.Active = active
	if err := s.userRepo.Update(ctx, target); err != nil {
		return fmt.Errorf("SetUserActive: failed to update user: %w", err)
	}

	return nil
}
