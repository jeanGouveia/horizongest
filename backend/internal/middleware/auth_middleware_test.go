package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/mocks"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

// TestAuthMiddleware_MissingToken tests that requests without token are rejected
func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenError = service.ErrInvalidCredentials
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// TestAuthMiddleware_InvalidToken tests that requests with invalid token are rejected
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenError = service.ErrInvalidCredentials
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

// TestAuthMiddleware_ValidToken tests that requests with valid token are accepted
func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenResult = &service.JWTClaims{
		UserID:    1,
		CompanyID: 1,
		Email:     "test@example.com",
		Name:      "Test User",
	}
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	userID, ok := GetUserIDFromContext(capturedContext)
	if !ok || userID != 1 {
		t.Errorf("Expected userID 1 in context, got %d", userID)
	}
}

// TestAuthMiddleware_ClaimsExtraction tests that JWT claims are properly extracted
func TestAuthMiddleware_ClaimsExtraction(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenResult = &service.JWTClaims{
		UserID:    1,
		CompanyID: 1,
		Email:     "test@example.com",
		Name:      "Test User",
	}
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	claims, ok := GetClaimsFromContext(capturedContext)
	if !ok {
		t.Fatal("Expected claims in context")
	}
	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}
	if claims.CompanyID != 1 {
		t.Errorf("Expected CompanyID 1, got %d", claims.CompanyID)
	}
}

// TestAuthMiddleware_ImpersonationClaims tests that impersonation claims are handled correctly
func TestAuthMiddleware_ImpersonationClaims(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenResult = &service.JWTClaims{
		UserID:                 1,
		CompanyID:              1,
		Email:                  "test@example.com",
		Name:                   "Test User",
		IsImpersonating:        true,
		OriginalPlatformUserID: 100,
	}
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	var capturedContext context.Context
	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContext = r.Context()
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	isImpersonating, ok := GetIsImpersonating(capturedContext)
	if !ok || !isImpersonating {
		t.Error("Expected impersonation flag in context")
	}

	originalUserID, ok := GetOriginalPlatformUserID(capturedContext)
	if !ok || originalUserID != 100 {
		t.Errorf("Expected original platform user ID 100, got %d", originalUserID)
	}
}

// TestGetUserIDFromContext tests that UserID can be extracted from context
func TestGetUserIDFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyUserID, uint(123))
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		t.Fatal("Expected to extract UserID from context")
	}
	if userID != 123 {
		t.Errorf("Expected UserID 123, got %d", userID)
	}
}

// TestGetClaimsFromContext tests that JWT claims can be extracted from context
func TestGetClaimsFromContext(t *testing.T) {
	claims := &service.JWTClaims{
		UserID:    1,
		CompanyID: 1,
		Email:     "test@example.com",
		Name:      "Test User",
	}
	ctx := context.WithValue(context.Background(), ContextKeyClaims, claims)
	extractedClaims, ok := GetClaimsFromContext(ctx)
	if !ok {
		t.Fatal("Expected to extract claims from context")
	}
	if extractedClaims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", extractedClaims.UserID)
	}
}

// TestAuthMiddleware_CookieBasedAuth tests that cookie-based authentication works
func TestAuthMiddleware_CookieBasedAuth(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenResult = &service.JWTClaims{
		UserID:    1,
		CompanyID: 1,
		Email:     "test@example.com",
		Name:      "Test User",
	}
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "valid_token"})
	w := httptest.NewRecorder()

	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestAuthMiddleware_HeaderBasedAuth tests that header-based authentication works
func TestAuthMiddleware_HeaderBasedAuth(t *testing.T) {
	mockAuth := mocks.NewMockAuthService()
	mockAuth.ValidateTokenResult = &service.JWTClaims{
		UserID:    1,
		CompanyID: 1,
		Email:     "test@example.com",
		Name:      "Test User",
	}
	authMw := NewAuthMiddleware(mockAuth)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	handler := authMw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
