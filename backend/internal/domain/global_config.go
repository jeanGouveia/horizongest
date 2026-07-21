package domain

import "time"

// GlobalConfig represents technical configuration for the platform.
// This is separate from PlatformBrandConfig (branding/institutional) and contains
// only technical settings like locale, timezone, upload limits, etc.
type GlobalConfig struct {
	// Localization
	DefaultTimezone string // Default timezone (e.g., "America/Sao_Paulo")
	DefaultLocale   string // Default locale (e.g., "pt-BR")

	// Formats
	MonetaryFormat string // Monetary format (e.g., "BRL R$ 1.000,00")
	DateFormat     string // Date format (e.g., "DD/MM/YYYY")
	TimeFormat     string // Time format (e.g., "HH:mm")

	// Upload Limits
	MaxUploadSizeMB    int64 // Maximum upload size in MB
	MaxImageSizeMB     int64 // Maximum image size in MB
	AllowedImageTypes  string // Comma-separated allowed image types (e.g., "jpg,png,webp")
	AllowedFileTypes   string // Comma-separated allowed file types (e.g., "pdf,doc,xlsx")

	// Maintenance
	MaintenanceMode   bool   // Whether platform is in maintenance mode
	MaintenanceMessage string // Message shown during maintenance

	// Feature Flags (stored here for convenience, managed via FeatureFlagService)
	EnableFinance     bool // Enable Finance module
	EnablePurchasing   bool // Enable Purchasing module
	EnableInventory    bool // Enable Inventory module
	EnableCRM          bool // Enable CRM module
	EnableCalendar     bool // Enable Calendar module
	EnablePOS          bool // Enable POS module
	EnableAI           bool // Enable AI module
	EnableDelivery     bool // Enable Delivery module
	EnableMarketplace  bool // Enable Marketplace module

	// Metadata
	UpdatedAt time.Time
	UpdatedBy uint
}

// DefaultGlobalConfig returns the default technical configuration.
func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		DefaultTimezone:   "America/Sao_Paulo",
		DefaultLocale:     "pt-BR",
		MonetaryFormat:   "BRL R$ 1.000,00",
		DateFormat:       "DD/MM/YYYY",
		TimeFormat:       "HH:mm",
		MaxUploadSizeMB:  10,
		MaxImageSizeMB:   5,
		AllowedImageTypes: "jpg,png,webp,gif",
		AllowedFileTypes:  "pdf,doc,docx,xlsx,xls,txt",
		MaintenanceMode:   false,
		MaintenanceMessage: "",
		EnableFinance:    true,
		EnablePurchasing:  true,
		EnableInventory:   true,
		EnableCRM:         false,
		EnableCalendar:    false,
		EnablePOS:         false,
		EnableAI:          false,
		EnableDelivery:    false,
		EnableMarketplace: false,
		UpdatedAt:         time.Now(),
		UpdatedBy:         0,
	}
}
