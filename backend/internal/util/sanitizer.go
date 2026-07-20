package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides input sanitization functions (Sprint 3.4 - Security Hardening)
type Sanitizer struct {
	maxNameLength        int
	maxDescriptionLength int
	maxSlugLength        int
	maxEmailLength       int
	maxPhoneLength       int
	maxNotesLength       int
}

// NewSanitizer creates a new sanitizer with default limits
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		maxNameLength:        255,
		maxDescriptionLength: 5000,
		maxSlugLength:        255,
		maxEmailLength:       255,
		maxPhoneLength:       20,
		maxNotesLength:       2000,
	}
}

// SanitizeString trims whitespace and validates length
func (s *Sanitizer) SanitizeString(input string, maxLength int, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > maxLength {
		return "", fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, maxLength)
	}
	return trimmed, nil
}

// SanitizeName sanitizes a name field
func (s *Sanitizer) SanitizeName(input string) (string, error) {
	return s.SanitizeString(input, s.maxNameLength, "name")
}

// SanitizeDescription sanitizes a description field
func (s *Sanitizer) SanitizeDescription(input string) (string, error) {
	return s.SanitizeString(input, s.maxDescriptionLength, "description")
}

// SanitizeSlug sanitizes a slug field (alphanumeric, hyphens, underscores only)
func (s *Sanitizer) SanitizeSlug(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > s.maxSlugLength {
		return "", fmt.Errorf("slug exceeds maximum length of %d characters", s.maxSlugLength)
	}
	
	// Validate slug format: lowercase alphanumeric, hyphens, underscores
	slugRegex := regexp.MustCompile(`^[a-z0-9-_]+$`)
	if !slugRegex.MatchString(trimmed) {
		return "", fmt.Errorf("slug must contain only lowercase letters, numbers, hyphens, and underscores")
	}
	
	return trimmed, nil
}

// SanitizeEmail validates email format and length
func (s *Sanitizer) SanitizeEmail(input string) (string, error) {
	trimmed := strings.TrimSpace(strings.ToLower(input))
	if len(trimmed) > s.maxEmailLength {
		return "", fmt.Errorf("email exceeds maximum length of %d characters", s.maxEmailLength)
	}
	
	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(trimmed) {
		return "", fmt.Errorf("invalid email format")
	}
	
	return trimmed, nil
}

// SanitizePhone sanitizes a phone number (digits, spaces, hyphens, parentheses, plus)
func (s *Sanitizer) SanitizePhone(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > s.maxPhoneLength {
		return "", fmt.Errorf("phone exceeds maximum length of %d characters", s.maxPhoneLength)
	}
	
	// Remove all non-allowed characters
	phoneRegex := regexp.MustCompile(`[^0-9+\-\(\) ]`)
	cleaned := phoneRegex.ReplaceAllString(trimmed, "")
	
	if len(cleaned) < 10 {
		return "", fmt.Errorf("phone number must have at least 10 digits")
	}
	
	return cleaned, nil
}

// SanitizeNotes sanitizes notes/observations field
func (s *Sanitizer) SanitizeNotes(input string) (string, error) {
	return s.SanitizeString(input, s.maxNotesLength, "notes")
}

// SanitizeCompanyName sanitizes company name
func (s *Sanitizer) SanitizeCompanyName(input string) (string, error) {
	return s.SanitizeString(input, s.maxNameLength, "company name")
}

// SanitizeText removes potentially dangerous characters from free text
func (s *Sanitizer) SanitizeText(input string) string {
	// Remove null bytes and other control characters except newline and tab
	var sb strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// ValidateMaxLength checks if a string exceeds maximum length
func (s *Sanitizer) ValidateMaxLength(input string, maxLength int, fieldName string) error {
	if len(input) > maxLength {
		return fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, maxLength)
	}
	return nil
}

// SanitizeURL validates URL format
func (s *Sanitizer) SanitizeURL(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil // Empty URL is allowed
	}
	
	if len(trimmed) > 500 {
		return "", fmt.Errorf("URL exceeds maximum length of 500 characters")
	}
	
	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}(/.*)?$`)
	if !urlRegex.MatchString(trimmed) {
		return "", fmt.Errorf("invalid URL format")
	}
	
	return trimmed, nil
}

// SanitizeColor validates hex color format
func (s *Sanitizer) SanitizeColor(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", nil // Empty color is allowed
	}
	
	colorRegex := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	if !colorRegex.MatchString(trimmed) {
		return "", fmt.Errorf("color must be in hex format (e.g., #3b82f6)")
	}
	
	return trimmed, nil
}
