package domain

import "time"

// Theme represents the visual branding configuration for a company
type Theme struct {
	// Branding colors from Company
	PrimaryColor   string
	SecondaryColor string
	LogoURL        string

	// Extended theme configuration (future expansion)
	FontFamily   string
	BorderRadius string

	// Metadata
	LoadedAt  time.Time
	IsDefault bool // true if using default theme (no company configured)
}

// DefaultTheme returns the default tenant theme
func DefaultTheme() *Theme {
	return &Theme{
		PrimaryColor:   "#6366f1", // Indigo-500
		SecondaryColor: "#4f46e5", // Indigo-600
		LogoURL:        "",
		FontFamily:     "Inter",
		BorderRadius:   "8px",
		LoadedAt:       time.Now(),
		IsDefault:      true,
	}
}

// ThemeFromCompany creates a Theme from a Company entity
func ThemeFromCompany(company *Company) *Theme {
	if company == nil {
		return DefaultTheme()
	}

	primaryColor := company.PrimaryColor
	if primaryColor == "" {
		primaryColor = "#6366f1" // Default
	}

	secondaryColor := company.SecondaryColor
	if secondaryColor == "" {
		secondaryColor = "#4f46e5" // Default
	}

	return &Theme{
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
		LogoURL:        company.LogoURL,
		FontFamily:     "Inter",
		BorderRadius:   "8px",
		LoadedAt:       time.Now(),
		IsDefault:      false,
	}
}
