package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCSRFMiddleware_GenerateToken(t *testing.T) {
	middleware := NewCSRFMiddleware()

	token := middleware.GenerateToken()
	if token == "" {
		t.Error("expected token to be non-empty")
	}

	if len(token) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("expected token length 64, got %d", len(token))
	}
}

func TestCSRFMiddleware_ValidateToken(t *testing.T) {
	middleware := NewCSRFMiddleware()

	token := middleware.GenerateToken()
	if !middleware.ValidateToken(token) {
		t.Error("expected valid token to pass validation")
	}

	if middleware.ValidateToken("invalid-token") {
		t.Error("expected invalid token to fail validation")
	}
}

func TestCSRFMiddleware_ValidateToken_Expired(t *testing.T) {
	middleware := NewCSRFMiddleware()
	middleware.tokenExpiry = 1 * time.Millisecond

	token := middleware.GenerateToken()
	time.Sleep(10 * time.Millisecond)

	if middleware.ValidateToken(token) {
		t.Error("expected expired token to fail validation")
	}
}

func TestCSRFMiddleware_Protect_SafeMethods(t *testing.T) {
	middleware := NewCSRFMiddleware()

	handler := middleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	safeMethods := []string{"GET", "HEAD", "OPTIONS"}

	for _, method := range safeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/test", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status 200 for %s, got %d", method, rr.Code)
			}
		})
	}
}

func TestCSRFMiddleware_Protect_APIEndpoints(t *testing.T) {
	middleware := NewCSRFMiddleware()

	handler := middleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// API endpoints should skip CSRF validation
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 for API endpoint, got %d", rr.Code)
	}
}

func TestCSRFMiddleware_Protect_NonAPIEndpoints(t *testing.T) {
	middleware := NewCSRFMiddleware()

	handler := middleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Non-API endpoints should require CSRF token
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403 for non-API endpoint without CSRF token, got %d", rr.Code)
	}
}

func TestCSRFMiddleware_Protect_WithValidToken(t *testing.T) {
	middleware := NewCSRFMiddleware()

	handler := middleware.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := middleware.GenerateToken()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("X-CSRF-Token", token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 with valid CSRF token, got %d", rr.Code)
	}
}

func TestCSRFMiddleware_CleanupExpiredTokens(t *testing.T) {
	middleware := NewCSRFMiddleware()
	middleware.tokenExpiry = 1 * time.Millisecond

	token1 := middleware.GenerateToken()
	time.Sleep(10 * time.Millisecond)
	token2 := middleware.GenerateToken()

	middleware.CleanupExpiredTokens()

	if middleware.ValidateToken(token1) {
		t.Error("expected expired token to be cleaned up")
	}

	if !middleware.ValidateToken(token2) {
		t.Error("expected valid token to still exist after cleanup")
	}
}
