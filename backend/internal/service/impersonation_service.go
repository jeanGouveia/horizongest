package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrImpersonationNotAllowed = errors.New("impersonation not allowed: user must be platform admin")
	ErrCompanyOwnerNotFound    = errors.New("company owner not found")
	ErrAlreadyImpersonating    = errors.New("already impersonating a company")
	ErrNotImpersonating        = errors.New("not currently impersonating")
)

type ImpersonationService struct {
	authService            *AuthService
	companyRepo            ports.CompanyRepository
	userRepo               ports.UserRepository
	impersonationAuditRepo ports.ImpersonationAuditRepository
}

func NewImpersonationService(
	authService *AuthService,
	companyRepo ports.CompanyRepository,
	userRepo ports.UserRepository,
	impersonationAuditRepo ports.ImpersonationAuditRepository,
) *ImpersonationService {
	return &ImpersonationService{
		authService:            authService,
		companyRepo:            companyRepo,
		userRepo:               userRepo,
		impersonationAuditRepo: impersonationAuditRepo,
	}
}

type StartImpersonationInput struct {
	PlatformUserID uint   `json:"platformUserId"`
	CompanyID      uint   `json:"companyId"`
	IPAddress      string `json:"ipAddress"`
	UserAgent      string `json:"userAgent"`
}

type StartImpersonationOutput struct {
	Token       string `json:"token"`
	ExpiresAt   int64  `json:"expiresAt"`
	CompanyID   uint   `json:"companyId"`
	CompanyName string `json:"companyName"`
	OwnerEmail  string `json:"ownerEmail"`
}

type EndImpersonationInput struct {
	PlatformUserID uint   `json:"platformUserId"`
	IPAddress      string `json:"ipAddress"`
	UserAgent      string `json:"userAgent"`
}

// StartImpersonation begins a temporary impersonation session for a platform admin
// This method is now idempotent: if an active impersonation exists, it will be automatically
// ended before creating a new one. This allows platform admins to freely switch between companies.
func (s *ImpersonationService) StartImpersonation(ctx context.Context, input StartImpersonationInput) (*StartImpersonationOutput, error) {
	// Check if platform admin is already impersonating
	// If so, automatically end the previous impersonation to allow switching companies
	activeImpersonation, err := s.impersonationAuditRepo.FindActiveByPlatformUserID(ctx, input.PlatformUserID)
	if err == nil && activeImpersonation != nil {
		// Automatically end the active impersonation (regardless of age)
		// This makes the endpoint idempotent and allows free company switching
		activeImpersonation.End()
		if err := s.impersonationAuditRepo.Update(ctx, activeImpersonation); err != nil {
			return nil, fmt.Errorf("ImpersonationService.StartImpersonation: encerrar impersonation anterior: %w", err)
		}
	}

	// FORENSIC: Log before finding company
	log.Printf("[FORENSIC] StartImpersonation - ANTES FindByID company - CompanyID: %d", input.CompanyID)

	// Find the company
	company, err := s.companyRepo.FindByID(ctx, input.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("ImpersonationService.StartImpersonation: buscar empresa: %w", err)
	}
	if company == nil {
		return nil, fmt.Errorf("ImpersonationService.StartImpersonation: empresa não encontrada")
	}

	// FORENSIC: Log after finding company
	log.Printf("[FORENSIC] StartImpersonation - APÓS FindByID company - Company.ID: %d, Company.Name: %s", company.ID, company.Name)

	// Find the company owner (user with Role = Owner in this company)
	owner, err := s.findCompanyOwner(ctx, input.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("ImpersonationService.StartImpersonation: buscar owner: %w", err)
	}
	if owner == nil {
		return nil, ErrCompanyOwnerNotFound
	}

	// FORENSIC: Log owner found
	log.Printf("[FORENSIC] StartImpersonation - Owner encontrado - ID: %d, Name: %s, Email: %s, CompanyID: %d", owner.ID, owner.Name, owner.Email, owner.CompanyID)

	// Create impersonation token with impersonation claims
	token, err := s.authService.generateJWTWithImpersonation(owner, true, input.PlatformUserID)
	if err != nil {
		return nil, fmt.Errorf("ImpersonationService.StartImpersonation: gerar token: %w", err)
	}

	// FORENSIC: Log JWT generated
	log.Printf("[FORENSIC] StartImpersonation - JWT gerado: %s", token)

	// Calculate expiration (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	// Create audit record
	audit := &domain.ImpersonationAudit{
		PlatformUserID:     input.PlatformUserID,
		CompanyID:          input.CompanyID,
		CompanyOwnerUserID: owner.ID,
		StartedAt:          time.Now(),
		IPAddress:          input.IPAddress,
		UserAgent:          input.UserAgent,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.impersonationAuditRepo.Create(ctx, audit); err != nil {
		return nil, fmt.Errorf("ImpersonationService.StartImpersonation: criar registro de auditoria: %w", err)
	}

	return &StartImpersonationOutput{
		Token:       token,
		ExpiresAt:   expiresAt,
		CompanyID:   company.ID,
		CompanyName: company.Name,
		OwnerEmail:  owner.Email,
	}, nil
}

// EndImpersonation ends the current impersonation session
func (s *ImpersonationService) EndImpersonation(ctx context.Context, input EndImpersonationInput) error {
	// Find active impersonation session
	activeImpersonation, err := s.impersonationAuditRepo.FindActiveByPlatformUserID(ctx, input.PlatformUserID)
	if err != nil {
		return fmt.Errorf("ImpersonationService.EndImpersonation: buscar impersonation ativa: %w", err)
	}
	if activeImpersonation == nil {
		// Not currently impersonating, but this is not an error
		// Just return success to allow graceful handling
		return nil
	}

	// Mark as ended
	activeImpersonation.End()
	activeImpersonation.IPAddress = input.IPAddress
	activeImpersonation.UserAgent = input.UserAgent

	if err := s.impersonationAuditRepo.Update(ctx, activeImpersonation); err != nil {
		return fmt.Errorf("ImpersonationService.EndImpersonation: atualizar registro: %w", err)
	}

	return nil
}

// findCompanyOwner finds the owner user of a company
func (s *ImpersonationService) findCompanyOwner(ctx context.Context, companyID uint) (*domain.User, error) {
	// This is a simplified implementation. In a real scenario, you would have a method in UserRepository
	// to find users by company ID and role. For now, we'll use the List method and filter.
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// DEBUG: Log all users found
	log.Printf("[DEBUG] findCompanyOwner - Buscando owner para CompanyID: %d", companyID)
	log.Printf("[DEBUG] findCompanyOwner - Total de usuários encontrados: %d", len(users))
	for _, user := range users {
		log.Printf("[DEBUG] findCompanyOwner - Usuário: ID=%d, Nome=%s, CompanyID=%d, Role=%s", user.ID, user.Name, user.CompanyID, user.Role)
	}

	for _, user := range users {
		if user.CompanyID == companyID && user.Role == domain.RoleOwner {
			log.Printf("[DEBUG] findCompanyOwner - Owner escolhido: ID=%d, Nome=%s, CompanyID=%d", user.ID, user.Name, user.CompanyID)
			return user, nil
		}
	}

	return nil, nil
}

// GetActiveImpersonation returns the active impersonation session for a platform admin
func (s *ImpersonationService) GetActiveImpersonation(ctx context.Context, platformUserID uint) (*domain.ImpersonationAudit, error) {
	return s.impersonationAuditRepo.FindActiveByPlatformUserID(ctx, platformUserID)
}

// GetImpersonationHistory returns the impersonation history for a platform admin
func (s *ImpersonationService) GetImpersonationHistory(ctx context.Context, platformUserID uint) ([]*domain.ImpersonationAudit, error) {
	return s.impersonationAuditRepo.FindByPlatformUserID(ctx, platformUserID)
}
