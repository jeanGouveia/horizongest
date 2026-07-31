package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

var (
	ErrInvalidGlobalConfig = errors.New("invalid global configuration")
)

// GlobalConfigRepository interface to avoid import cycle
// The repository implementation handles caching internally
type GlobalConfigRepository interface {
	Get(ctx context.Context) (*domain.GlobalConfig, error)
	Update(ctx context.Context, config *domain.GlobalConfig, updatedBy uint) error
	Initialize(ctx context.Context) error
}

type GlobalConfigService struct {
	configRepo GlobalConfigRepository
}

func NewGlobalConfigService(configRepo GlobalConfigRepository) *GlobalConfigService {
	return &GlobalConfigService{
		configRepo: configRepo,
	}
}

// Get retrieves the current global configuration
func (s *GlobalConfigService) Get(ctx context.Context) (*domain.GlobalConfig, error) {
	return s.configRepo.Get(ctx)
}

// Update updates the global configuration
// Only platform admins should be allowed to call this
func (s *GlobalConfigService) Update(ctx context.Context, config *domain.GlobalConfig, updatedBy uint) error {
	// Validate configuration
	if err := s.validateConfig(config); err != nil {
		return fmt.Errorf("GlobalConfigService.Update: validar configuração: %w", err)
	}

	return s.configRepo.Update(ctx, config, updatedBy)
}

// Initialize ensures the default global configuration exists
func (s *GlobalConfigService) Initialize(ctx context.Context) error {
	return s.configRepo.Initialize(ctx)
}

// IsModuleEnabled checks if a specific module is enabled via feature flags
func (s *GlobalConfigService) IsModuleEnabled(ctx context.Context, module string) (bool, error) {
	config, err := s.configRepo.Get(ctx)
	if err != nil {
		return false, err
	}

	switch module {
	case "finance":
		return config.EnableFinance, nil
	case "purchasing":
		return config.EnablePurchasing, nil
	case "inventory":
		return config.EnableInventory, nil
	case "crm":
		return config.EnableCRM, nil
	case "calendar":
		return config.EnableCalendar, nil
	case "pos":
		return config.EnablePOS, nil
	case "ai":
		return config.EnableAI, nil
	case "delivery":
		return config.EnableDelivery, nil
	case "marketplace":
		return config.EnableMarketplace, nil
	default:
		return false, nil
	}
}

func (s *GlobalConfigService) validateConfig(config *domain.GlobalConfig) error {
	if config.DefaultTimezone == "" {
		return errors.New("default timezone is required")
	}
	if config.DefaultLocale == "" {
		return errors.New("default locale is required")
	}
	if config.MonetaryFormat == "" {
		return errors.New("monetary format is required")
	}
	if config.DateFormat == "" {
		return errors.New("date format is required")
	}
	if config.TimeFormat == "" {
		return errors.New("time format is required")
	}
	if config.MaxUploadSizeMB <= 0 {
		return errors.New("max upload size must be greater than 0")
	}
	if config.MaxImageSizeMB <= 0 {
		return errors.New("max image size must be greater than 0")
	}
	if config.AllowedImageTypes == "" {
		return errors.New("allowed image types is required")
	}
	if config.AllowedFileTypes == "" {
		return errors.New("allowed file types is required")
	}
	return nil
}
