package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCompanyAlreadyExists = errors.New("company with this slug already exists")
	// ErrPermissionDenied and ErrCompanyNotFound are defined in user_management_service.go and company_service.go
)

type PlatformService struct {
	companyRepo       ports.CompanyRepository
	userRepo          ports.UserRepository
	platformUserRepo  PlatformUserRepository
	platformAuditRepo PlatformAuditRepository
	emailService      *EmailService
}

type PlatformAuditRepository interface {
	Create(ctx context.Context, audit *domain.PlatformAudit) error
}

func NewPlatformService(
	companyRepo ports.CompanyRepository,
	userRepo ports.UserRepository,
	platformUserRepo PlatformUserRepository,
	platformAuditRepo PlatformAuditRepository,
	emailService *EmailService,
) *PlatformService {
	return &PlatformService{
		companyRepo:       companyRepo,
		userRepo:          userRepo,
		platformUserRepo:  platformUserRepo,
		platformAuditRepo: platformAuditRepo,
		emailService:      emailService,
	}
}

type PlatformCreateCompanyInput struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Slug         string `json:"slug" validate:"required,min=2,max=50"`
	Description  string `json:"description"`
	BusinessType string `json:"business_type" validate:"omitempty,oneof=restaurant bakery confectionery coffee_shop pizzeria burger ice_cream acai food_truck dark_kitchen generic"`
	Locale       string `json:"locale" validate:"omitempty"`
	Currency     string `json:"currency" validate:"omitempty"`
	Timezone     string `json:"timezone" validate:"omitempty"`
}

type CreateCompanyOutput struct {
	CompanyID uint `json:"company_id"`
	UserID    uint `json:"user_id"`
}

// CreateCompany creates a new company and assigns an initial owner
func (s *PlatformService) CreateCompany(ctx context.Context, platformUserID uint, input PlatformCreateCompanyInput, ownerEmail, ownerPassword, ownerName string) (*CreateCompanyOutput, error) {
	// Verify platform user has permission (PlatformAdmin only)
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.CreateCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil {
		return nil, ErrPermissionDenied
	}
	if platformUser.Role != domain.PlatformRoleAdmin {
		return nil, ErrPermissionDenied
	}

	// Check if company slug already exists
	existing, err := s.companyRepo.FindBySlug(ctx, input.Slug)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.CreateCompany: verificar slug: %w", err)
	}
	if existing != nil {
		return nil, ErrCompanyAlreadyExists
	}

	// Create company
	company := &domain.Company{
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		BusinessType: domain.BusinessType(input.BusinessType),
		Active:       true,
		Locale:       input.Locale,
		Currency:     input.Currency,
		Timezone:     input.Timezone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("PlatformService.CreateCompany: criar empresa: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(ownerPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.CreateCompany: gerar hash senha: %w", err)
	}

	// Create owner user
	owner := &domain.User{
		Name:         ownerName,
		Email:        ownerEmail,
		PasswordHash: string(hashedPassword),
		CompanyID:    company.ID,
		Role:         domain.RoleOwner,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, owner); err != nil {
		return nil, fmt.Errorf("PlatformService.CreateCompany: criar owner: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "create_company", "company", company.ID, map[string]interface{}{
		"name": company.Name,
		"slug": company.Slug,
	})

	// Send welcome email with temporary password
	if s.emailService != nil {
		_ = s.emailService.SendWelcomeEmail(ownerEmail, ownerName, company.Name, ownerPassword)
	}

	return &CreateCompanyOutput{
		CompanyID: company.ID,
		UserID:    owner.ID,
	}, nil
}

// ListCompanies returns all companies (platform admin only)
func (s *PlatformService) ListCompanies(ctx context.Context, platformUserID uint) ([]*domain.Company, error) {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.ListCompanies: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return nil, ErrPermissionDenied
	}

	companies, err := s.companyRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.ListCompanies: listar empresas: %w", err)
	}

	// Convert []domain.Company to []*domain.Company
	result := make([]*domain.Company, len(companies))
	for i := range companies {
		result[i] = &companies[i]
	}

	return result, nil
}

// GetCompany returns a specific company by ID (platform admin only)
func (s *PlatformService) GetCompany(ctx context.Context, platformUserID uint, companyID uint) (*domain.Company, error) {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return nil, ErrPermissionDenied
	}

	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	return company, nil
}

// UpdateCompany updates a company (platform admin only)
func (s *PlatformService) UpdateCompany(ctx context.Context, platformUserID uint, companyID uint, input PlatformCreateCompanyInput) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.UpdateCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get existing company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.UpdateCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Check if new slug conflicts with another company
	if input.Slug != company.Slug {
		existing, err := s.companyRepo.FindBySlug(ctx, input.Slug)
		if err != nil {
			return fmt.Errorf("PlatformService.UpdateCompany: verificar slug: %w", err)
		}
		if existing != nil && existing.ID != companyID {
			return ErrCompanyAlreadyExists
		}
	}

	// Update company
	company.Name = input.Name
	company.Slug = input.Slug
	company.Description = input.Description
	company.BusinessType = domain.BusinessType(input.BusinessType)
	company.Locale = input.Locale
	company.Currency = input.Currency
	company.Timezone = input.Timezone
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.UpdateCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "update_company", "company", companyID, map[string]interface{}{
		"name": company.Name,
		"slug": company.Slug,
	})

	return nil
}

// DeactivateCompany deactivates a company (platform admin only)
func (s *PlatformService) DeactivateCompany(ctx context.Context, platformUserID uint, companyID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.DeactivateCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get existing company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.DeactivateCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Deactivate company
	company.Active = false
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.DeactivateCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "deactivate_company", "company", companyID, nil)

	return nil
}

// ActivateCompany activates a company (platform admin only)
func (s *PlatformService) ActivateCompany(ctx context.Context, platformUserID uint, companyID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.ActivateCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get existing company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.ActivateCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Activate company
	company.Active = true
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.ActivateCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "activate_company", "company", companyID, nil)

	return nil
}

// logAudit logs an audit event
func (s *PlatformService) logAudit(ctx context.Context, platformUserID uint, action string, entityType string, entityID uint, changes map[string]interface{}) {
	if s.platformAuditRepo == nil {
		return
	}

	// Convert changes to JSON string (simplified)
	// In production, use json.Marshal
	changesStr := ""
	for k, v := range changes {
		if changesStr != "" {
			changesStr += ", "
		}
		changesStr += fmt.Sprintf("%s=%v", k, v)
	}

	audit := &domain.PlatformAudit{
		PlatformUserID: &platformUserID,
		Action:         action,
		EntityType:     entityType,
		EntityID:       &entityID,
		Changes:        changesStr,
		CreatedAt:      time.Now(),
	}

	// Log audit asynchronously (ignore errors)
	_ = s.platformAuditRepo.Create(ctx, audit)
}

type DashboardStats struct {
	TotalCompanies   uint `json:"total_companies"`
	TotalOwners      uint `json:"total_owners"`
	TotalUsers       uint `json:"total_users"`
	BlockedCompanies uint `json:"blocked_companies"`
	TrialCompanies   uint `json:"trial_companies"`
	PaidCompanies    uint `json:"paid_companies"`
}

// GetDashboardStats returns platform dashboard statistics
func (s *PlatformService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	companies, err := s.companyRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetDashboardStats: listar empresas: %w", err)
	}

	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetDashboardStats: listar usuários: %w", err)
	}

	stats := &DashboardStats{
		TotalCompanies: uint(len(companies)),
		TotalUsers:     uint(len(users)),
	}

	// Count owners, blocked, trial, and paid companies
	for _, company := range companies {
		if !company.Active {
			stats.BlockedCompanies++
		}
		// Note: Trial and Paid status would require additional fields in Company entity
		// For now, we'll set them to 0
	}

	// Count owners (users with RoleOwner)
	for _, user := range users {
		if user.Role == domain.RoleOwner {
			stats.TotalOwners++
		}
	}

	return stats, nil
}

// GetCompanyOwner returns the owner of a company
func (s *PlatformService) GetCompanyOwner(ctx context.Context, platformUserID uint, companyID uint) (*domain.User, error) {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetCompanyOwner: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return nil, ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetCompanyOwner: buscar empresa: %w", err)
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	// Find owner user for this company
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("PlatformService.GetCompanyOwner: listar usuários: %w", err)
	}

	for i := range users {
		if users[i].CompanyID == companyID && users[i].Role == domain.RoleOwner {
			return users[i], nil
		}
	}

	return nil, fmt.Errorf("PlatformService.GetCompanyOwner: owner não encontrado")
}

// ResetOwnerPassword resets the password of a company owner
func (s *PlatformService) ResetOwnerPassword(ctx context.Context, platformUserID uint, companyID uint, newPassword string) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Find owner user
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: listar usuários: %w", err)
	}

	var owner *domain.User
	for i := range users {
		if users[i].CompanyID == companyID && users[i].Role == domain.RoleOwner {
			owner = users[i]
			break
		}
	}

	if owner == nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: owner não encontrado")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: gerar hash senha: %w", err)
	}

	// Update password
	owner.PasswordHash = string(hashedPassword)
	owner.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, owner); err != nil {
		return fmt.Errorf("PlatformService.ResetOwnerPassword: atualizar owner: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "reset_owner_password", "user", owner.ID, map[string]interface{}{
		"company_id": companyID,
	})

	return nil
}

// BlockUser blocks a user
func (s *PlatformService) BlockUser(ctx context.Context, platformUserID uint, userID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.BlockUser: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("PlatformService.BlockUser: buscar usuário: %w", err)
	}
	if user == nil {
		return fmt.Errorf("PlatformService.BlockUser: usuário não encontrado")
	}

	// Block user
	user.Active = false
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("PlatformService.BlockUser: atualizar usuário: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "block_user", "user", userID, map[string]interface{}{
		"company_id": user.CompanyID,
	})

	return nil
}

// UnblockUser unblocks a user
func (s *PlatformService) UnblockUser(ctx context.Context, platformUserID uint, userID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.UnblockUser: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("PlatformService.UnblockUser: buscar usuário: %w", err)
	}
	if user == nil {
		return fmt.Errorf("PlatformService.UnblockUser: usuário não encontrado")
	}

	// Unblock user
	user.Active = true
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("PlatformService.UnblockUser: atualizar usuário: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "unblock_user", "user", userID, map[string]interface{}{
		"company_id": user.CompanyID,
	})

	return nil
}

// LoginAsCompany creates a temporary session for a platform user to login as a company user
func (s *PlatformService) LoginAsCompany(ctx context.Context, platformUserID uint, companyID uint) (string, error) {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return "", fmt.Errorf("PlatformService.LoginAsCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin && platformUser.Role != domain.PlatformRoleSupport {
		return "", ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return "", fmt.Errorf("PlatformService.LoginAsCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return "", ErrCompanyNotFound
	}

	// Find owner user for this company
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return "", fmt.Errorf("PlatformService.LoginAsCompany: listar usuários: %w", err)
	}

	var owner *domain.User
	for i := range users {
		if users[i].CompanyID == companyID && users[i].Role == domain.RoleOwner {
			owner = users[i]
			break
		}
	}

	if owner == nil {
		return "", fmt.Errorf("PlatformService.LoginAsCompany: owner não encontrado")
	}

	// Generate a temporary JWT token for the company user
	// This token should have a short expiration (e.g., 1 hour)
	// For now, we'll return the owner ID and let the auth service handle token generation
	// In a real implementation, this would use the auth service to generate a temporary token

	// Log audit
	s.logAudit(ctx, platformUserID, "login_as_company", "company", companyID, map[string]interface{}{
		"owner_id":    owner.ID,
		"owner_email": owner.Email,
	})

	return owner.Email, nil
}

// SetCompanyTrial sets a company to trial status with an end date
func (s *PlatformService) SetCompanyTrial(ctx context.Context, platformUserID uint, companyID uint, trialEndsAt time.Time) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.SetCompanyTrial: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.SetCompanyTrial: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Set trial status
	company.Status = "trial"
	company.TrialEndsAt = &trialEndsAt
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.SetCompanyTrial: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "set_company_trial", "company", companyID, map[string]interface{}{
		"trial_ends_at": trialEndsAt,
	})

	return nil
}

// SuspendCompany suspends a company
func (s *PlatformService) SuspendCompany(ctx context.Context, platformUserID uint, companyID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.SuspendCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.SuspendCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Suspend company
	company.Status = "suspended"
	company.Active = false
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.SuspendCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "suspend_company", "company", companyID, nil)

	return nil
}

// CancelCompany cancels a company
func (s *PlatformService) CancelCompany(ctx context.Context, platformUserID uint, companyID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.CancelCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.CancelCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Cancel company
	company.Status = "cancelled"
	company.Active = false
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.CancelCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "cancel_company", "company", companyID, nil)

	return nil
}

// ReactivateCompany reactivates a suspended or cancelled company
func (s *PlatformService) ReactivateCompany(ctx context.Context, platformUserID uint, companyID uint) error {
	// Verify platform user has permission
	platformUser, err := s.platformUserRepo.FindByID(ctx, platformUserID)
	if err != nil {
		return fmt.Errorf("PlatformService.ReactivateCompany: buscar usuário plataforma: %w", err)
	}
	if platformUser == nil || platformUser.Role != domain.PlatformRoleAdmin {
		return ErrPermissionDenied
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("PlatformService.ReactivateCompany: buscar empresa: %w", err)
	}
	if company == nil {
		return ErrCompanyNotFound
	}

	// Reactivate company
	company.Status = "active"
	company.Active = true
	company.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("PlatformService.ReactivateCompany: atualizar empresa: %w", err)
	}

	// Log audit
	s.logAudit(ctx, platformUserID, "reactivate_company", "company", companyID, nil)

	return nil
}
