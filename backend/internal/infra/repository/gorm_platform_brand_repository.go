package repository

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type GormPlatformBrand struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	PlatformName       string `gorm:"column:platform_name;not null"`
	PlatformShortName  string `gorm:"column:platform_short_name;not null"`
	OwnerCompanyName   string `gorm:"column:owner_company_name;not null"`
	OwnerDocument      string `gorm:"column:owner_document"`
	Website            string `gorm:"not null"`
	SupportEmail       string `gorm:"not null"`
	SupportURL         string `gorm:"not null"`
	LogoPath           string
	FaviconPath        string
	LogoLight          string `gorm:"column:logo_light"`
	LogoDark           string `gorm:"column:logo_dark"`
	Icon               string
	LoginBackground    string `gorm:"column:login_background"`
	LoginIllustration  string `gorm:"column:login_illustration"`
	Copyright          string `gorm:"not null"`
	PrivacyPolicyURL   string `gorm:"column:privacy_policy_url"`
	TermsURL           string `gorm:"column:terms_url"`
	InstagramURL       string `gorm:"column:instagram_url"`
	FacebookURL        string `gorm:"column:facebook_url"`
	LinkedInURL        string `gorm:"column:linkedin_url"`
	YoutubeURL         string `gorm:"column:youtube_url"`
	DefaultLanguage    string `gorm:"column:default_language"`
	DefaultTimezone    string `gorm:"column:default_timezone"`
	MaintenanceMode    bool   `gorm:"column:maintenance_mode"`
	MaintenanceMessage string `gorm:"column:maintenance_message"`
	PrimaryColor       string `gorm:"not null"`
	SecondaryColor     string `gorm:"not null"`
	UpdatedAt          int64  `gorm:"autoUpdateTime"`
	UpdatedBy          uint
}

func (GormPlatformBrand) TableName() string {
	return "platform_brand_config"
}

// GormPlatformBrandRepository implements platform brand configuration persistence with in-memory caching
// TODO: Future White Label Support
// To support multiple platform brands (white label), modify this repository to:
// - Remove singleton pattern (ID=1) and support multiple brand configurations
// - Add methods like GetByBrandKey(brandKey) or GetByDomain(domain)
// - Update cache to be a map[string]*domain.PlatformBrandConfig keyed by brand identifier
// - Add brand_key field to GormPlatformBrand struct and database schema
type GormPlatformBrandRepository struct {
	db          *gorm.DB
	cache       *domain.PlatformBrandConfig
	cacheMutex  sync.RWMutex
	cacheLoaded bool
}

func NewGormPlatformBrandRepository(db *gorm.DB) *GormPlatformBrandRepository {
	return &GormPlatformBrandRepository{db: db}
}

// Get retrieves the platform brand configuration (singleton - always ID=1)
// Uses in-memory cache to avoid unnecessary database queries
func (r *GormPlatformBrandRepository) Get(ctx context.Context) (*domain.PlatformBrandConfig, error) {
	// Try to get from cache first (read lock)
	r.cacheMutex.RLock()
	if r.cacheLoaded && r.cache != nil {
		cached := r.cache
		r.cacheMutex.RUnlock()
		return cached, nil
	}
	r.cacheMutex.RUnlock()

	// Cache miss - load from database
	var gormBrand GormPlatformBrand
	err := r.db.WithContext(ctx).First(&gormBrand, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default brand if not found
			return domain.DefaultPlatformBrand(), nil
		}
		return nil, err
	}

	brand := r.toDomain(&gormBrand)

	// Update cache (write lock)
	r.cacheMutex.Lock()
	r.cache = brand
	r.cacheLoaded = true
	r.cacheMutex.Unlock()

	return brand, nil
}

// Update updates the platform brand configuration (singleton - always ID=1)
// Automatically invalidates cache after successful update
func (r *GormPlatformBrandRepository) Update(ctx context.Context, brand *domain.PlatformBrandConfig, updatedBy uint) error {
	gormBrand := r.toGorm(brand)
	gormBrand.ID = 1
	gormBrand.UpdatedBy = updatedBy
	err := r.db.WithContext(ctx).Save(gormBrand).Error
	if err != nil {
		return err
	}
	// Invalidate cache after successful update
	r.InvalidateCache()
	return nil
}

// Initialize creates the default platform brand configuration if it doesn't exist
func (r *GormPlatformBrandRepository) Initialize(ctx context.Context) error {
	var count int64
	err := r.db.WithContext(ctx).Table("platform_brand_config").Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		defaultBrand := r.toGorm(domain.DefaultPlatformBrand())
		defaultBrand.ID = 1
		return r.db.WithContext(ctx).Create(defaultBrand).Error
	}
	return nil
}

func (r *GormPlatformBrandRepository) toGorm(brand *domain.PlatformBrandConfig) *GormPlatformBrand {
	return &GormPlatformBrand{
		PlatformName:       brand.PlatformName,
		PlatformShortName:  brand.PlatformShortName,
		OwnerCompanyName:   brand.OwnerCompanyName,
		OwnerDocument:      brand.OwnerDocument,
		Website:            brand.Website,
		SupportEmail:       brand.SupportEmail,
		SupportURL:         brand.SupportURL,
		LogoPath:           brand.LogoPath,
		FaviconPath:        brand.FaviconPath,
		LogoLight:          brand.LogoLight,
		LogoDark:           brand.LogoDark,
		Icon:               brand.Icon,
		LoginBackground:    brand.LoginBackground,
		LoginIllustration:  brand.LoginIllustration,
		Copyright:          brand.Copyright,
		PrivacyPolicyURL:   brand.PrivacyPolicyURL,
		TermsURL:           brand.TermsURL,
		InstagramURL:       brand.InstagramURL,
		FacebookURL:        brand.FacebookURL,
		LinkedInURL:        brand.LinkedInURL,
		YoutubeURL:         brand.YoutubeURL,
		DefaultLanguage:    brand.DefaultLanguage,
		DefaultTimezone:    brand.DefaultTimezone,
		MaintenanceMode:    brand.MaintenanceMode,
		MaintenanceMessage: brand.MaintenanceMessage,
		PrimaryColor:       brand.PrimaryColor,
		SecondaryColor:     brand.SecondaryColor,
		UpdatedAt:          brand.UpdatedAt.Unix(),
		UpdatedBy:          brand.UpdatedBy,
	}
}

func (r *GormPlatformBrandRepository) toDomain(gormBrand *GormPlatformBrand) *domain.PlatformBrandConfig {
	return &domain.PlatformBrandConfig{
		PlatformName:       gormBrand.PlatformName,
		PlatformShortName:  gormBrand.PlatformShortName,
		OwnerCompanyName:   gormBrand.OwnerCompanyName,
		OwnerDocument:      gormBrand.OwnerDocument,
		Website:            gormBrand.Website,
		SupportEmail:       gormBrand.SupportEmail,
		SupportURL:         gormBrand.SupportURL,
		LogoPath:           gormBrand.LogoPath,
		FaviconPath:        gormBrand.FaviconPath,
		LogoLight:          gormBrand.LogoLight,
		LogoDark:           gormBrand.LogoDark,
		Icon:               gormBrand.Icon,
		LoginBackground:    gormBrand.LoginBackground,
		LoginIllustration:  gormBrand.LoginIllustration,
		Copyright:          gormBrand.Copyright,
		PrivacyPolicyURL:   gormBrand.PrivacyPolicyURL,
		TermsURL:           gormBrand.TermsURL,
		InstagramURL:       gormBrand.InstagramURL,
		FacebookURL:        gormBrand.FacebookURL,
		LinkedInURL:        gormBrand.LinkedInURL,
		YoutubeURL:         gormBrand.YoutubeURL,
		DefaultLanguage:    gormBrand.DefaultLanguage,
		DefaultTimezone:    gormBrand.DefaultTimezone,
		MaintenanceMode:    gormBrand.MaintenanceMode,
		MaintenanceMessage: gormBrand.MaintenanceMessage,
		PrimaryColor:       gormBrand.PrimaryColor,
		SecondaryColor:     gormBrand.SecondaryColor,
		UpdatedAt:          time.Unix(gormBrand.UpdatedAt, 0),
		UpdatedBy:          gormBrand.UpdatedBy,
	}
}

// InvalidateCache clears the in-memory cache
// Should be called after any update to ensure consistency
func (r *GormPlatformBrandRepository) InvalidateCache() {
	r.cacheMutex.Lock()
	r.cache = nil
	r.cacheLoaded = false
	r.cacheMutex.Unlock()
}

// ReloadCache forces a cache reload from the database
// Useful when you need to ensure fresh data
func (r *GormPlatformBrandRepository) ReloadCache(ctx context.Context) error {
	r.InvalidateCache()
	_, err := r.Get(ctx)
	return err
}
