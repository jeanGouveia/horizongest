package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// TestAuthService_Login_Success tests successful login flow
func TestAuthService_Login_Success(t *testing.T) {
	// Create a minimal service instance
	svc := &AuthService{
		secret:     []byte("test-secret"),
		expiry:     24 * time.Hour,
		bcryptCost: bcrypt.DefaultCost,
		issuer:     "TestPlatform",
	}

	// Mock user would be set up via repository in real tests
	// This test verifies JWT generation logic
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // bcrypt hash of "password"
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("generateJWT returned empty token")
	}

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return svc.secret, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, claims.Email)
	}
	if claims.CompanyID != user.CompanyID {
		t.Errorf("expected CompanyID %d, got %d", user.CompanyID, claims.CompanyID)
	}
	if claims.Issuer != svc.issuer {
		t.Errorf("expected Issuer %s, got %s", svc.issuer, claims.Issuer)
	}
}

// TestAuthService_Login_InactiveUser tests that inactive users cannot login
func TestAuthService_Login_InactiveUser(t *testing.T) {
	// This test would verify that inactive users are rejected
	// In a real test, we would set up a mock repository returning an inactive user
	// For now, we document the expected behavior
	t.Skip("TODO: implement with mock repository")
}

// TestAuthService_Logout tests logout functionality
func TestAuthService_Logout(t *testing.T) {
	// This test would verify that logout adds token to blacklist
	// In a real test, we would set up a mock token blacklist repository
	t.Skip("TODO: implement with mock repository")
}

// TestAuthService_ValidateToken tests token validation
func TestAuthService_ValidateToken(t *testing.T) {
	svc := &AuthService{
		secret:     []byte("test-secret"),
		expiry:     24 * time.Hour,
		bcryptCost: bcrypt.DefaultCost,
		issuer:     "TestPlatform",
	}

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CompanyID: 1,
		Role:      domain.RoleOwner,
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	// Verify token is valid
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return svc.secret, nil
	})
	if err != nil {
		t.Fatalf("failed to parse valid token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("token is not valid")
	}
}

// TestAuthService_ValidateToken_InvalidSecret tests that tokens with wrong secret are rejected
func TestAuthService_ValidateToken_InvalidSecret(t *testing.T) {
	svc := &AuthService{
		secret:     []byte("test-secret"),
		expiry:     24 * time.Hour,
		bcryptCost: bcrypt.DefaultCost,
		issuer:     "TestPlatform",
	}

	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	// Try to validate with wrong secret
	wrongSecret := []byte("wrong-secret")
	_, err = jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return wrongSecret, nil
	})

	if err == nil {
		t.Fatal("expected error when validating token with wrong secret")
	}
}

// TestAuthService_generateJWT_DynamicIssuer tests JWT generation with dynamic issuer
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
