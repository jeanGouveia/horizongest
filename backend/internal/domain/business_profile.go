package domain

import "time"

// BusinessProfile represents the complete business identity for a company
// It combines Company data with business-specific configuration
type BusinessProfile struct {
	// Company basic info
	CompanyID   uint
	CompanyName string
	CompanySlug string
	Active      bool
	
	// Business Engine fields
	BusinessType BusinessType
	Locale       string
	Currency     string
	Timezone     string
	
	// Theme fields (from White Label Engine)
	LogoURL        string
	PrimaryColor   string
	SecondaryColor string
	
	// Metadata
	LoadedAt time.Time
	IsDefault bool // true if using default profile (no company configured)
}

// DefaultBusinessProfile returns the default business profile
func DefaultBusinessProfile() *BusinessProfile {
	return &BusinessProfile{
		CompanyID:     0,
		CompanyName:   "PratoOnline",
		CompanySlug:   "pratoonline",
		Active:        true,
		BusinessType:  BusinessTypeGeneric,
		Locale:        "pt-BR",
		Currency:      "BRL",
		Timezone:      "America/Sao_Paulo",
		LogoURL:       "",
		PrimaryColor:  "#6366f1",
		SecondaryColor: "#4f46e5",
		LoadedAt:      time.Now(),
		IsDefault:     true,
	}
}

// BusinessProfileFromCompany creates a BusinessProfile from a Company entity
func BusinessProfileFromCompany(company *Company) *BusinessProfile {
	if company == nil {
		return DefaultBusinessProfile()
	}
	
	businessType := company.BusinessType
	if !businessType.IsValid() {
		businessType = BusinessTypeGeneric
	}
	
	locale := company.Locale
	if locale == "" {
		locale = "pt-BR"
	}
	
	currency := company.Currency
	if currency == "" {
		currency = "BRL"
	}
	
	timezone := company.Timezone
	if timezone == "" {
		timezone = "America/Sao_Paulo"
	}
	
	return &BusinessProfile{
		CompanyID:     company.ID,
		CompanyName:   company.Name,
		CompanySlug:   company.Slug,
		Active:        company.Active,
		BusinessType:  businessType,
		Locale:        locale,
		Currency:      currency,
		Timezone:      timezone,
		LogoURL:       company.LogoURL,
		PrimaryColor:  company.PrimaryColor,
		SecondaryColor: company.SecondaryColor,
		LoadedAt:      time.Now(),
		IsDefault:     false,
	}
}
