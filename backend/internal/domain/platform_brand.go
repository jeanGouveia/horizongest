package domain

import "time"

// PlatformBrandConfig represents the institutional branding configuration for the platform.
// This is completely separate from Tenant Branding (company-specific branding).
// Only Platform Admin can modify this configuration.
//
// The platform name (e.g., "HorizonGest") is separate from the owner company information.
// This allows the software to maintain its brand identity even if the owning company changes.
//
// TODO: Future White Label Support
// This architecture is prepared for future multi-brand support where each installation
// can have its own PlatformBrandConfig. To implement:
// - Add a unique identifier (e.g., installation_id or brand_id) to distinguish different brands
// - Modify repository to support multiple brand configurations instead of singleton (ID=1)
// - Add context-based brand selection based on domain, subdomain, or header
// - Consider adding a "brand_key" field for routing to specific brand configs
type PlatformBrandConfig struct {
	// Platform Information (software brand)
	PlatformName      string // Full platform name (e.g., "HorizonGest")
	PlatformShortName string // Short platform name (e.g., "Horizon")

	// Owner Company Information (business entity that owns the platform)
	OwnerCompanyName string // Owner company name (e.g., "HorizonGest Inc.")
	OwnerDocument    string // Owner company legal document (CNPJ, EIN, etc.)

	// Platform Website
	Website string // Platform website (e.g., "https://horizongest.com")

	// Support Information
	SupportEmail string // Support email address
	SupportURL   string // Support/help center URL

	// Branding Assets
	LogoPath          string // Path to platform logo
	FaviconPath       string // Path to platform favicon
	LogoLight         string // Path to light mode logo (optional)
	LogoDark          string // Path to dark mode logo (optional)
	Icon              string // Path to platform icon (optional)
	LoginBackground   string // Path to login page background (optional)
	LoginIllustration string // Path to login page illustration (optional)

	// Legal Information
	Copyright string // Copyright notice

	// Legal URLs
	PrivacyPolicyURL string // URL to privacy policy (optional)
	TermsURL         string // URL to terms of service (optional)

	// Social Media URLs
	InstagramURL string // Instagram profile URL (optional)
	FacebookURL  string // Facebook page URL (optional)
	LinkedInURL  string // LinkedIn company page URL (optional)
	YoutubeURL   string // YouTube channel URL (optional)

	// Localization
	DefaultLanguage string // Default language code (e.g., "pt-BR", optional)
	DefaultTimezone string // Default timezone (e.g., "America/Sao_Paulo", optional)

	// Maintenance Mode
	MaintenanceMode    bool   // Whether platform is in maintenance mode (optional)
	MaintenanceMessage string // Message shown during maintenance (optional)

	// Branding Colors
	PrimaryColor   string // Platform primary color
	SecondaryColor string // Platform secondary color

	// Metadata
	UpdatedAt time.Time // Last update timestamp
	UpdatedBy uint      // ID of platform admin who last updated
}

// DefaultPlatformBrand returns an empty platform branding configuration.
// All institutional data should come from the database via migration.
// This is only used as a fallback when the database is empty.
func DefaultPlatformBrand() *PlatformBrandConfig {
	return &PlatformBrandConfig{
		PlatformName:       "",
		PlatformShortName:  "",
		OwnerCompanyName:   "",
		OwnerDocument:      "",
		Website:            "",
		SupportEmail:       "",
		SupportURL:         "",
		LogoPath:           "",
		FaviconPath:        "",
		LogoLight:          "",
		LogoDark:           "",
		Icon:               "",
		LoginBackground:    "",
		LoginIllustration:  "",
		Copyright:          "",
		PrivacyPolicyURL:   "",
		TermsURL:           "",
		InstagramURL:       "",
		FacebookURL:        "",
		LinkedInURL:        "",
		YoutubeURL:         "",
		DefaultLanguage:    "",
		DefaultTimezone:    "",
		MaintenanceMode:    false,
		MaintenanceMessage: "",
		PrimaryColor:       "",
		SecondaryColor:     "",
		UpdatedAt:          time.Now(),
		UpdatedBy:          0,
	}
}
