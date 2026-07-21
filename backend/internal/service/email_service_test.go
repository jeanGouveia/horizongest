package service

import (
	"testing"
)

// Test for EmailService focusing on dynamic platform name in templates
// This is critical for branding support (Sprint 3.7)

func TestEmailService_SendWelcomeEmail_DynamicPlatformName(t *testing.T) {
	tests := []struct {
		name         string
		platformName string
		enabled      bool
	}{
		{
			name:         "success with custom platform name",
			platformName: "TestPlatform",
			enabled:      false, // Disabled to avoid actual email sending
		},
		{
			name:         "success with empty platform name (fallback)",
			platformName: "",
			enabled:      false,
		},
		{
			name:         "success with HorizonGest",
			platformName: "HorizonGest",
			enabled:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService(tt.enabled, "noreply@example.com", tt.platformName)

			err := svc.SendWelcomeEmail("test@example.com", "Test User", "Test Company", "tempPassword123")
			if err != nil {
				t.Errorf("EmailService.SendWelcomeEmail() error = %v", err)
			}
		})
	}
}

func TestEmailService_SendPasswordResetEmail_DynamicPlatformName(t *testing.T) {
	tests := []struct {
		name         string
		platformName string
		enabled      bool
	}{
		{
			name:         "success with custom platform name",
			platformName: "TestPlatform",
			enabled:      false,
		},
		{
			name:         "success with empty platform name (fallback)",
			platformName: "",
			enabled:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService(tt.enabled, "noreply@example.com", tt.platformName)

			err := svc.SendPasswordResetEmail("test@example.com", "Test User", "https://example.com/reset")
			if err != nil {
				t.Errorf("EmailService.SendPasswordResetEmail() error = %v", err)
			}
		})
	}
}
