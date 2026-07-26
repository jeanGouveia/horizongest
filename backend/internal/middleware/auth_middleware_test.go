package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

// MockAuthService is a mock implementation of service.AuthServiceInterface
type MockAuthService struct {
	ValidateTokenResult *service.JWTClaims
	ValidateTokenError  error
	LogoutError         error
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) Login(ctx context.Context, input service.LoginInput) (*service.LoginResult, error) {
	return &service.LoginResult{}, nil
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return m.LogoutError
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*service.JWTClaims, error) {
	return m.ValidateTokenResult, m.ValidateTokenError
}

func (m *MockAuthService) UpdateProfile(ctx context.Context, userID uint, input service.UpdateProfileInput) (*domain.User, error) {
	return &domain.User{}, nil
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID uint, input service.ChangePasswordInput) error {
	return nil
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, input service.RequestPasswordResetInput) error {
	return nil
}

func (m *MockAuthService) ResetPassword(ctx context.Context, input service.ResetPasswordInput) error {
	return nil
}

// TestAuthMiddleware_MissingToken tests that requests without token are rejected
func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
	mockAuth := NewMockAuthService()
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
