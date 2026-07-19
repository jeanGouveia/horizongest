package service

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

type ThemeService struct {
	companyRepo ports.CompanyRepository
	userRepo    ports.UserRepository
}

func NewThemeService(companyRepo ports.CompanyRepository, userRepo ports.UserRepository) *ThemeService {
	return &ThemeService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

// GetThemeForUser retrieves the theme for a specific user based on their company
// If the user has no company or company has no custom theme, returns default theme
func (s *ThemeService) GetThemeForUser(ctx context.Context, userID uint) (*domain.Theme, error) {
	// Get user to find their company
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ThemeService.GetThemeForUser: failed to get user: %w", err)
	}
	
	if user == nil {
		// User not found, return default theme
		return domain.DefaultTheme(), nil
	}
	
	// If user has no company, return default theme
	if user.CompanyID == nil {
		return domain.DefaultTheme(), nil
	}
	
	// Get company to retrieve theme configuration
	company, err := s.companyRepo.FindByID(ctx, *user.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("ThemeService.GetThemeForUser: failed to get company: %w", err)
	}
	
	if company == nil {
		// Company not found, return default theme
		return domain.DefaultTheme(), nil
	}
	
	// Create theme from company
	return domain.ThemeFromCompany(company), nil
}

// GetDefaultTheme returns the default PratoOnline theme
func (s *ThemeService) GetDefaultTheme() *domain.Theme {
	return domain.DefaultTheme()
}
