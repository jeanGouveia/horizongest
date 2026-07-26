package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// TestJWT_GenerateToken tests JWT token generation
func TestJWT_GenerateToken(t *testing.T) {
	svc := &AuthService{
		secret: []byte("test-secret"),
		expiry: 24 * time.Hour,
		issuer: "TestPlatform",
	}

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CompanyID: 123,
		Role:      domain.RoleOwner,
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("generateJWT returned empty token")
	}

	// Verify token structure
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

	// Verify all expected claims
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

// TestJWT_TokenExpiration tests that tokens have correct expiration
func TestJWT_TokenExpiration(t *testing.T) {
	svc := &AuthService{
		secret: []byte("test-secret"),
		expiry: 1 * time.Hour, // Short expiry for testing
		issuer: "TestPlatform",
	}

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CompanyID: 123,
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

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

	// Verify expiration is set
	if claims.ExpiresAt == nil {
		t.Fatal("token has no expiration")
	}

	// Verify expiration is approximately correct (within 1 minute tolerance)
	expectedExpiry := time.Now().Add(svc.expiry)
	diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("token expiration not as expected: got %v, want approximately %v", claims.ExpiresAt.Time, expectedExpiry)
	}
}

// TestJWT_ImpersonationClaims tests that impersonation claims are included
func TestJWT_ImpersonationClaims(t *testing.T) {
	svc := &AuthService{
		secret: []byte("test-secret"),
		expiry: 24 * time.Hour,
		issuer: "TestPlatform",
	}

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CompanyID: 123,
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

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

	// Verify impersonation claims are present but false by default
	if claims.IsImpersonating {
		t.Error("IsImpersonating should be false by default")
	}
	if claims.OriginalPlatformUserID != 0 {
		t.Error("OriginalPlatformUserID should be 0 by default")
	}
}

// TestJWT_SecretValidation tests that tokens with wrong secret are rejected
func TestJWT_SecretValidation(t *testing.T) {
	svc := &AuthService{
		secret: []byte("test-secret"),
		expiry: 24 * time.Hour,
		issuer: "TestPlatform",
	}

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Test User",
		CompanyID: 123,
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

// TestJWT_HttpOnlyCookie tests that cookies are set with HttpOnly flag
func TestJWT_HttpOnlyCookie(t *testing.T) {
	// This test verifies that JWT cookies are set with HttpOnly flag
	// This would require testing the handler that sets cookies
	
	t.Skip("TODO: implement with handler test - verifies HttpOnly cookie flag")
}

// TestJWT_SecureCookie tests that cookies are set with Secure flag in production
func TestJWT_SecureCookie(t *testing.T) {
	// This test verifies that JWT cookies are set with Secure flag in production
	
	t.Skip("TODO: implement with handler test - verifies Secure cookie flag")
}

// TestJWT_SameSiteCookie tests that cookies are set with SameSite flag
func TestJWT_SameSiteCookie(t *testing.T) {
	// This test verifies that JWT cookies are set with SameSite flag for CSRF protection
	
	t.Skip("TODO: implement with handler test - verifies SameSite cookie flag")
}
