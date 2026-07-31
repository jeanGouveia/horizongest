package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

var (
	ErrInvalidPlatformBrand = errors.New("invalid platform brand configuration")
	ErrUnauthorizedUpdate   = errors.New("unauthorized: only platform admin can update platform brand")
)

// PlatformBrandRepository interface to avoid import cycle
// The repository implementation handles caching internally
// The service layer is unaware of cache details
type PlatformBrandRepository interface {
	Get(ctx context.Context) (*domain.PlatformBrandConfig, error)
	Update(ctx context.Context, brand *domain.PlatformBrandConfig, updatedBy uint) error
	Initialize(ctx context.Context) error
}

type PlatformBrandService struct {
	brandRepo PlatformBrandRepository
}

func NewPlatformBrandService(brandRepo PlatformBrandRepository) *PlatformBrandService {
	return &PlatformBrandService{
		brandRepo: brandRepo,
	}
}

// Get retrieves the current platform brand configuration
func (s *PlatformBrandService) Get(ctx context.Context) (*domain.PlatformBrandConfig, error) {
	return s.brandRepo.Get(ctx)
}

// Update updates the platform brand configuration
// Only platform admins should be allowed to call this
func (s *PlatformBrandService) Update(ctx context.Context, brand *domain.PlatformBrandConfig, updatedBy uint) error {
	// Validate configuration
	if err := s.validateBrand(brand); err != nil {
		return fmt.Errorf("PlatformBrandService.Update: validar configuração: %w", err)
	}

	return s.brandRepo.Update(ctx, brand, updatedBy)
}

// Initialize ensures the default platform brand configuration exists
func (s *PlatformBrandService) Initialize(ctx context.Context) error {
	return s.brandRepo.Initialize(ctx)
}

func (s *PlatformBrandService) validateBrand(brand *domain.PlatformBrandConfig) error {
	if brand.PlatformName == "" {
		return errors.New("platform name is required")
	}
	if brand.PlatformShortName == "" {
		return errors.New("platform short name is required")
	}
	if brand.OwnerCompanyName == "" {
		return errors.New("owner company name is required")
	}
	if brand.Website == "" {
		return errors.New("website is required")
	}
	if brand.SupportEmail == "" {
		return errors.New("support email is required")
	}
	if brand.Copyright == "" {
		return errors.New("copyright is required")
	}
	if brand.PrimaryColor == "" {
		return errors.New("primary color is required")
	}
	if brand.SecondaryColor == "" {
		return errors.New("secondary color is required")
	}
	return nil
}
