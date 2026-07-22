package service

import (
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// Simplified test for AuthService focusing on JWT generation with dynamic issuer
// This is critical for branding support (Sprint 3.7)

func TestAuthService_generateJWT_DynamicIssuer(t *testing.T) {
	tests := []struct {
		name    string
		issuer  string
		user    *domain.User
		wantErr bool
	}{
		{
			name:   "success with custom issuer",
			issuer: "TestPlatform",
			user: &domain.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name:   "success with empty issuer (fallback)",
			issuer: "",
			user: &domain.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal service instance just for testing JWT generation
			svc := &AuthService{
				secret: []byte("test-secret"),
				expiry: 24 * 3600, // 24 hours in seconds
				issuer: tt.issuer,
			}

			token, err := svc.generateJWT(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.generateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Error("AuthService.generateJWT() returned empty token")
			}
		})
	}
}
