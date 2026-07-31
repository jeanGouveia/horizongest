package service

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type BusinessService struct {
	companyRepo ports.CompanyRepository
	userRepo    ports.UserRepository
}

func NewBusinessService(companyRepo ports.CompanyRepository, userRepo ports.UserRepository) *BusinessService {
	return &BusinessService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

// GetBusinessProfile retrieves the business profile for a specific user based on their company
// If the user has no company or company has no configuration, returns default profile
func (s *BusinessService) GetBusinessProfile(ctx context.Context, userID uint) (*domain.BusinessProfile, error) {
	// Get user to find their company
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BusinessService.GetBusinessProfile: buscar usuário: %w", err)
	}

	if user == nil {
		// User not found, return default profile
		return domain.DefaultBusinessProfile(), nil
	}

	// If user has no company (CompanyID = 0), return default profile
	if user.CompanyID == 0 {
		return domain.DefaultBusinessProfile(), nil
	}

	// Get company to retrieve business profile
	company, err := s.companyRepo.FindByID(ctx, user.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("BusinessService.GetBusinessProfile: buscar empresa: %w", err)
	}

	if company == nil {
		// Company not found, return default profile
		return domain.DefaultBusinessProfile(), nil
	}

	// Create business profile from company
	return domain.BusinessProfileFromCompany(company), nil
}

// GetDefaultBusinessProfile returns the default tenant business profile
func (s *BusinessService) GetDefaultBusinessProfile() *domain.BusinessProfile {
	return domain.DefaultBusinessProfile()
}
