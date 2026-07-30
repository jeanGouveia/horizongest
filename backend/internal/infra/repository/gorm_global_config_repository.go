package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type GormGlobalConfig struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	DefaultTimezone    string    `gorm:"column:default_timezone"`
	DefaultLocale      string    `gorm:"column:default_locale"`
	MonetaryFormat     string    `gorm:"column:monetary_format"`
	DateFormat         string    `gorm:"column:date_format"`
	TimeFormat         string    `gorm:"column:time_format"`
	MaxUploadSizeMB    int64     `gorm:"column:max_upload_size_mb"`
	MaxImageSizeMB     int64     `gorm:"column:max_image_size_mb"`
	AllowedImageTypes  string    `gorm:"column:allowed_image_types"`
	AllowedFileTypes   string    `gorm:"column:allowed_file_types"`
	MaintenanceMode    bool      `gorm:"column:maintenance_mode"`
	MaintenanceMessage string    `gorm:"column:maintenance_message"`
	EnableFinance      bool      `gorm:"column:enable_finance"`
	EnablePurchasing   bool      `gorm:"column:enable_purchasing"`
	EnableInventory    bool      `gorm:"column:enable_inventory"`
	EnableCRM          bool      `gorm:"column:enable_crm"`
	EnableCalendar     bool      `gorm:"column:enable_calendar"`
	EnablePOS          bool      `gorm:"column:enable_pos"`
	EnableAI           bool      `gorm:"column:enable_ai"`
	EnableDelivery     bool      `gorm:"column:enable_delivery"`
	EnableMarketplace  bool      `gorm:"column:enable_marketplace"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	UpdatedBy          uint
}

func (GormGlobalConfig) TableName() string {
	return "global_config"
}

// GormGlobalConfigRepository implements global configuration persistence with in-memory caching
type GormGlobalConfigRepository struct {
	db          *gorm.DB
	cache       *domain.GlobalConfig
	cacheMutex  sync.RWMutex
	cacheLoaded bool
}

func NewGormGlobalConfigRepository(db *gorm.DB) *GormGlobalConfigRepository {
	return &GormGlobalConfigRepository{db: db}
}

// Get retrieves the global configuration (singleton - always ID=1)
// Uses in-memory cache to avoid unnecessary database queries
func (r *GormGlobalConfigRepository) Get(ctx context.Context) (*domain.GlobalConfig, error) {
	// Try to get from cache first (read lock)
	r.cacheMutex.RLock()
	if r.cacheLoaded && r.cache != nil {
		cached := r.cache
		r.cacheMutex.RUnlock()
		return cached, nil
	}
	r.cacheMutex.RUnlock()

	// Cache miss - load from database
	var gormConfig GormGlobalConfig
	err := r.db.WithContext(ctx).First(&gormConfig, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default config if not found
			return domain.DefaultGlobalConfig(), nil
		}
		return nil, fmt.Errorf("GlobalConfigRepository.Get: %w", err)
	}

	config := r.toDomain(&gormConfig)

	// Update cache (write lock)
	r.cacheMutex.Lock()
	r.cache = config
	r.cacheLoaded = true
	r.cacheMutex.Unlock()

	return config, nil
}

// Update updates the global configuration (singleton - always ID=1)
// Automatically invalidates cache after successful update
func (r *GormGlobalConfigRepository) Update(ctx context.Context, config *domain.GlobalConfig, updatedBy uint) error {
	gormConfig := r.toGorm(config)
	gormConfig.ID = 1
	gormConfig.UpdatedBy = updatedBy
	err := r.db.WithContext(ctx).Save(gormConfig).Error
	if err != nil {
		return fmt.Errorf("GlobalConfigRepository.Update: %w", err)
	}
	// Invalidate cache after successful update
	r.InvalidateCache()
	return nil
}

// Initialize creates the default global configuration if it doesn't exist
func (r *GormGlobalConfigRepository) Initialize(ctx context.Context) error {
	var count int64
	err := r.db.WithContext(ctx).Table("global_config").Count(&count).Error
	if err != nil {
		return fmt.Errorf("GlobalConfigRepository.Initialize: %w", err)
	}

	if count == 0 {
		defaultConfig := r.toGorm(domain.DefaultGlobalConfig())
		defaultConfig.ID = 1
		if err := r.db.WithContext(ctx).Create(defaultConfig).Error; err != nil {
			return fmt.Errorf("GlobalConfigRepository.Initialize: %w", err)
		}
	}
	return nil
}

func (r *GormGlobalConfigRepository) toGorm(config *domain.GlobalConfig) *GormGlobalConfig {
	return &GormGlobalConfig{
		DefaultTimezone:    config.DefaultTimezone,
		DefaultLocale:      config.DefaultLocale,
		MonetaryFormat:     config.MonetaryFormat,
		DateFormat:         config.DateFormat,
		TimeFormat:         config.TimeFormat,
		MaxUploadSizeMB:    config.MaxUploadSizeMB,
		MaxImageSizeMB:     config.MaxImageSizeMB,
		AllowedImageTypes:  config.AllowedImageTypes,
		AllowedFileTypes:   config.AllowedFileTypes,
		MaintenanceMode:    config.MaintenanceMode,
		MaintenanceMessage: config.MaintenanceMessage,
		EnableFinance:      config.EnableFinance,
		EnablePurchasing:   config.EnablePurchasing,
		EnableInventory:    config.EnableInventory,
		EnableCRM:          config.EnableCRM,
		EnableCalendar:     config.EnableCalendar,
		EnablePOS:          config.EnablePOS,
		EnableAI:           config.EnableAI,
		EnableDelivery:     config.EnableDelivery,
		EnableMarketplace:  config.EnableMarketplace,
		UpdatedAt:          config.UpdatedAt,
		UpdatedBy:          config.UpdatedBy,
	}
}

func (r *GormGlobalConfigRepository) toDomain(gormConfig *GormGlobalConfig) *domain.GlobalConfig {
	return &domain.GlobalConfig{
		DefaultTimezone:    gormConfig.DefaultTimezone,
		DefaultLocale:      gormConfig.DefaultLocale,
		MonetaryFormat:     gormConfig.MonetaryFormat,
		DateFormat:         gormConfig.DateFormat,
		TimeFormat:         gormConfig.TimeFormat,
		MaxUploadSizeMB:    gormConfig.MaxUploadSizeMB,
		MaxImageSizeMB:     gormConfig.MaxImageSizeMB,
		AllowedImageTypes:  gormConfig.AllowedImageTypes,
		AllowedFileTypes:   gormConfig.AllowedFileTypes,
		MaintenanceMode:    gormConfig.MaintenanceMode,
		MaintenanceMessage: gormConfig.MaintenanceMessage,
		EnableFinance:      gormConfig.EnableFinance,
		EnablePurchasing:   gormConfig.EnablePurchasing,
		EnableInventory:    gormConfig.EnableInventory,
		EnableCRM:          gormConfig.EnableCRM,
		EnableCalendar:     gormConfig.EnableCalendar,
		EnablePOS:          gormConfig.EnablePOS,
		EnableAI:           gormConfig.EnableAI,
		EnableDelivery:     gormConfig.EnableDelivery,
		EnableMarketplace:  gormConfig.EnableMarketplace,
		UpdatedAt:          gormConfig.UpdatedAt,
		UpdatedBy:          gormConfig.UpdatedBy,
	}
}

// InvalidateCache clears the in-memory cache
func (r *GormGlobalConfigRepository) InvalidateCache() {
	r.cacheMutex.Lock()
	r.cache = nil
	r.cacheLoaded = false
	r.cacheMutex.Unlock()
}

// ReloadCache forces a cache reload from the database
func (r *GormGlobalConfigRepository) ReloadCache(ctx context.Context) error {
	r.InvalidateCache()
	_, err := r.Get(ctx)
	return fmt.Errorf("GlobalConfigRepository.ReloadCache: %w", err)
}
