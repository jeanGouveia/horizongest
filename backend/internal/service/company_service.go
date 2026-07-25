package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrCompanyNotFound   = errors.New("empresa não encontrada")
	ErrSlugAlreadyExists = errors.New("slug já está em uso")
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

type CreateCompanyInput struct {
	Name           string `json:"name" validate:"required,min=2,max=120"`
	Slug           string `json:"slug" validate:"required,min=2,max=100"`
	Description    string `json:"description"`
	LogoURL        string `json:"logo_url"`
	PrimaryColor   string `json:"primary_color" validate:"omitempty,hexcolor"`
	SecondaryColor string `json:"secondary_color" validate:"omitempty,hexcolor"`
	// Business Engine fields (Sprint 3)
	BusinessType string `json:"business_type" validate:"omitempty"`
	Locale       string `json:"locale" validate:"omitempty"`
	Currency     string `json:"currency" validate:"omitempty"`
	Timezone     string `json:"timezone" validate:"omitempty"`
}

type UpdateCompanyInput struct {
	Name           string `json:"name" validate:"required,min=2,max=120"`
	Slug           string `json:"slug" validate:"required,min=2,max=100"`
	Description    string `json:"description"`
	Active         *bool  `json:"active"`
	LogoURL        string `json:"logo_url"`
	PrimaryColor   string `json:"primary_color" validate:"omitempty,hexcolor"`
	SecondaryColor string `json:"secondary_color" validate:"omitempty,hexcolor"`
	// Business Engine fields (Sprint 3)
	BusinessType string `json:"business_type" validate:"omitempty"`
	Locale       string `json:"locale" validate:"omitempty"`
	Currency     string `json:"currency" validate:"omitempty"`
	Timezone     string `json:"timezone" validate:"omitempty"`
}

func (s *CompanyService) CreateCompany(ctx context.Context, in CreateCompanyInput) (*domain.Company, error) {
	// Normalize slug to lowercase
	slug := strings.ToLower(in.Slug)

	// Check if slug already exists
	existing, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.CreateCompany: %w", err)
	}
	if existing != nil {
		return nil, ErrSlugAlreadyExists
	}

	// Set default colors if not provided
	primaryColor := in.PrimaryColor
	if primaryColor == "" {
		primaryColor = "#3b82f6"
	}
	secondaryColor := in.SecondaryColor
	if secondaryColor == "" {
		secondaryColor = "#1e40af"
	}

	// Set default Business Engine fields with fallback
	businessType := domain.BusinessType(in.BusinessType)
	if !businessType.IsValid() {
		businessType = domain.BusinessTypeGeneric
	}
	locale := in.Locale
	if locale == "" {
		locale = "pt-BR"
	}
	currency := in.Currency
	if currency == "" {
		currency = "BRL"
	}
	timezone := in.Timezone
	if timezone == "" {
		timezone = "America/Sao_Paulo"
	}

	c := &domain.Company{
		Name:           in.Name,
		Slug:           slug,
		Description:    in.Description,
		Active:         true,
		LogoURL:        in.LogoURL,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
		BusinessType:   businessType,
		Locale:         locale,
		Currency:       currency,
		Timezone:       timezone,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("CompanyService.CreateCompany: %w", err)
	}
	return c, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	companies, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.ListCompanies: %w", err)
	}
	return companies, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, id uint) (*domain.Company, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.GetCompany: %w", err)
	}
	if c == nil {
		return nil, ErrCompanyNotFound
	}
	return c, nil
}

// GetCurrentCompany retrieves the company by ID from the tenant context
// This is the secure method that uses the CompanyID from the JWT/tenant context
// The handler should extract the CompanyID from the tenant context and pass it here
func (s *CompanyService) GetCurrentCompany(ctx context.Context, companyID uint) (*domain.Company, error) {
	if companyID == 0 {
		return nil, errors.New("company ID cannot be zero")
	}

	c, err := s.repo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.GetCurrentCompany: %w", err)
	}
	if c == nil {
		return nil, ErrCompanyNotFound
	}
	return c, nil
}

func (s *CompanyService) GetCompanyBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	c, err := s.repo.FindBySlug(ctx, strings.ToLower(slug))
	if err != nil {
		return nil, fmt.Errorf("CompanyService.GetCompanyBySlug: %w", err)
	}
	if c == nil {
		return nil, ErrCompanyNotFound
	}
	return c, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id uint, in UpdateCompanyInput) (*domain.Company, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.UpdateCompany: %w", err)
	}
	if c == nil {
		return nil, ErrCompanyNotFound
	}

	// Check if new slug conflicts with another company
	slug := strings.ToLower(in.Slug)
	if slug != c.Slug {
		existing, err := s.repo.FindBySlug(ctx, slug)
		if err != nil {
			return nil, fmt.Errorf("CompanyService.UpdateCompany: %w", err)
		}
		if existing != nil {
			return nil, ErrSlugAlreadyExists
		}
	}

	c.Name = in.Name
	c.Slug = slug
	c.Description = in.Description
	c.LogoURL = in.LogoURL

	// Update colors if provided
	if in.PrimaryColor != "" {
		c.PrimaryColor = in.PrimaryColor
	}
	if in.SecondaryColor != "" {
		c.SecondaryColor = in.SecondaryColor
	}

	// Update Business Engine fields if provided
	if in.BusinessType != "" {
		bt := domain.BusinessType(in.BusinessType)
		if bt.IsValid() {
			c.BusinessType = bt
		}
	}
	if in.Locale != "" {
		c.Locale = in.Locale
	}
	if in.Currency != "" {
		c.Currency = in.Currency
	}
	if in.Timezone != "" {
		c.Timezone = in.Timezone
	}

	// Update active status if provided
	if in.Active != nil {
		c.Active = *in.Active
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("CompanyService.UpdateCompany: %w", err)
	}
	return c, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id uint) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CompanyService.DeleteCompany: %w", err)
	}
	if c == nil {
		return ErrCompanyNotFound
	}
	return s.repo.Delete(ctx, id)
}
