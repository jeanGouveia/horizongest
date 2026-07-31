package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrUserNoCompany = errors.New("user does not have a company")
)

type CompanySettingsService struct {
	companyRepo ports.CompanyRepository
	userRepo    ports.UserRepository
}

func NewCompanySettingsService(companyRepo ports.CompanyRepository, userRepo ports.UserRepository) *CompanySettingsService {
	return &CompanySettingsService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

// GetSettingsInput represents the input for getting company settings
type GetSettingsInput struct {
	UserID uint
}

// GetSettingsOutput represents the output for company settings
type GetSettingsOutput struct {
	Name           string
	Slug           string
	Description    string
	LogoURL        string
	PrimaryColor   string
	SecondaryColor string
	BusinessType   string
	Locale         string
	Currency       string
	Timezone       string
}

// UpdateSettingsInput represents the input for updating company settings
type UpdateSettingsInput struct {
	Name           *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description    *string `json:"description" validate:"omitempty,max=500"`
	LogoURL        *string `json:"logo_url" validate:"omitempty,url"`
	PrimaryColor   *string `json:"primary_color" validate:"omitempty,hexcolor|rgb|rgba"`
	SecondaryColor *string `json:"secondary_color" validate:"omitempty,hexcolor|rgb|rgba"`
	BusinessType   *string `json:"business_type" validate:"omitempty,oneof=restaurant bakery cafe bar food_truck catering other"`
	Locale         *string `json:"locale" validate:"omitempty,min=2,max=10"`
	Currency       *string `json:"currency" validate:"omitempty,min=3,max=3,uppercase"`
	Timezone       *string `json:"timezone" validate:"omitempty,max=50"`
}

// GetSettings retrieves the settings for the user's company
func (s *CompanySettingsService) GetSettings(ctx context.Context, userID uint) (*GetSettingsOutput, error) {
	log.Printf("[FORENSIC] CompanySettingsService - GetSettings - UserID recebido: %d", userID)

	// FORENSIC: ANTES do FindByID user
	log.Printf("[FORENSIC] CompanySettingsService - ANTES FindByID user - UserID: %d", userID)

	// Get user to find their company
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("CompanySettingsService.GetSettings: buscar usuário: %w", err)
	}

	if user == nil {
		return nil, ErrUserNoCompany
	}

	if user.CompanyID == 0 {
		return nil, ErrUserNoCompany
	}

	// FORENSIC: APÓS do FindByID user
	log.Printf("[FORENSIC] CompanySettingsService - APÓS FindByID user - user.ID=%d, user.CompanyID=%d, user.Name=%s", user.ID, user.CompanyID, user.Name)

	// FORENSIC: ANTES do CompanyRepo
	log.Printf("[FORENSIC] CompanySettingsService - ANTES CompanyRepo.FindByID - CompanyID usado: %d", user.CompanyID)

	// Get company
	company, err := s.companyRepo.FindByID(ctx, user.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("CompanySettingsService.GetSettings: buscar empresa: %w", err)
	}

	if company == nil {
		return nil, errors.New("company not found")
	}

	// FORENSIC: APÓS do CompanyRepo
	log.Printf("[FORENSIC] CompanySettingsService - APÓS CompanyRepo.FindByID - empresa carregada: ID=%d, Name=%s", company.ID, company.Name)

	return &GetSettingsOutput{
		Name:           company.Name,
		Slug:           company.Slug,
		Description:    company.Description,
		LogoURL:        company.LogoURL,
		PrimaryColor:   company.PrimaryColor,
		SecondaryColor: company.SecondaryColor,
		BusinessType:   string(company.BusinessType),
		Locale:         company.Locale,
		Currency:       company.Currency,
		Timezone:       company.Timezone,
	}, nil
}

// UpdateSettings updates the settings for the user's company
func (s *CompanySettingsService) UpdateSettings(ctx context.Context, userID uint, input UpdateSettingsInput) error {
	// Get user to find their company
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("CompanySettingsService.UpdateSettings: buscar usuário: %w", err)
	}

	if user == nil {
		return ErrUserNoCompany
	}

	if user.CompanyID == 0 {
		return ErrUserNoCompany
	}

	// Get current company
	company, err := s.companyRepo.FindByID(ctx, user.CompanyID)
	if err != nil {
		return fmt.Errorf("CompanySettingsService.UpdateSettings: buscar empresa: %w", err)
	}

	if company == nil {
		return errors.New("company not found")
	}

	// Update fields if provided
	if input.Name != nil {
		company.Name = *input.Name
	}
	if input.Description != nil {
		company.Description = *input.Description
	}
	if input.LogoURL != nil {
		company.LogoURL = *input.LogoURL
	}
	if input.PrimaryColor != nil {
		company.PrimaryColor = *input.PrimaryColor
	}
	if input.SecondaryColor != nil {
		company.SecondaryColor = *input.SecondaryColor
	}
	if input.BusinessType != nil {
		company.BusinessType = domain.BusinessType(*input.BusinessType)
	}
	if input.Locale != nil {
		company.Locale = *input.Locale
	}
	if input.Currency != nil {
		company.Currency = *input.Currency
	}
	if input.Timezone != nil {
		company.Timezone = *input.Timezone
	}

	// Update company
	if err := s.companyRepo.Update(ctx, company); err != nil {
		return fmt.Errorf("CompanySettingsService.UpdateSettings: atualizar empresa: %w", err)
	}

	return nil
}
