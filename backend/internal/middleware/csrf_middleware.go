package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// CSRFMiddleware provides CSRF protection for state-changing operations
// FASE A.2: B10 - CSRF Protection
type CSRFMiddleware struct {
	tokens      map[string]csrfToken
	mu          sync.RWMutex
	tokenExpiry time.Duration
	secure      bool // HTTPS only in production
	sameSite    http.SameSite
}

type csrfToken struct {
	value     string
	expiresAt time.Time
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware() *CSRFMiddleware {
	env := os.Getenv("ENVIRONMENT")
	secure := env == "production"
	
	return &CSRFMiddleware{
		tokens:      make(map[string]csrfToken),
		tokenExpiry: 24 * time.Hour,
		secure:      secure,
		sameSite:    http.SameSiteStrictMode,
	}
}

// GenerateToken generates a new CSRF token
func (m *CSRFMiddleware) GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.tokens[token] = csrfToken{
		value:     token,
		expiresAt: time.Now().Add(m.tokenExpiry),
	}
	
	return token
}

// ValidateToken validates a CSRF token
func (m *CSRFMiddleware) ValidateToken(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	csrfToken, exists := m.tokens[token]
	if !exists {
		return false
	}
	
	// Check if token is expired
	if time.Now().After(csrfToken.expiresAt) {
		return false
	}
	
	return true
}

// CleanupExpiredTokens removes expired tokens
func (m *CSRFMiddleware) CleanupExpiredTokens() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := time.Now()
	for token, csrfToken := range m.tokens {
		if now.After(csrfToken.expiresAt) {
			delete(m.tokens, token)
		}
	}
}

// Protect applies CSRF protection to the handler
func (m *CSRFMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF validation for safe methods
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		
		// Skip CSRF validation for API endpoints that use JWT authentication
		// JWT-based APIs are less vulnerable to CSRF since tokens are stored in localStorage
		// and sent via Authorization header, not cookies
		if strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		
		// For non-API endpoints, validate CSRF token
		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			csrfToken = r.FormValue("csrf_token")
		}
		
		if !m.ValidateToken(csrfToken) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// GetCSRFToken returns a CSRF token for the current session
func (m *CSRFMiddleware) GetCSRFToken(w http.ResponseWriter, r *http.Request) string {
	token := m.GenerateToken()
	
	// Set CSRF token in cookie
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		Secure:   m.secure,
		HttpOnly: true,
		SameSite: m.sameSite,
		MaxAge:   int(m.tokenExpiry.Seconds()),
	}
	
	http.SetCookie(w, cookie)
	
	return token
}
