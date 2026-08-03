package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrelationMiddleware(t *testing.T) {
	middleware := CorrelationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that correlation ID is in context
		correlationID := GetCorrelationID(r.Context())
		if correlationID == "" {
			t.Error("expected correlation ID to be set in context")
		}

		// Check that request ID is in context
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("expected request ID to be set in context")
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	// Check response headers
	correlationIDHeader := rr.Header().Get("X-Correlation-ID")
	if correlationIDHeader == "" {
		t.Error("expected X-Correlation-ID header in response")
	}

	requestIDHeader := rr.Header().Get("X-Request-ID")
	if requestIDHeader == "" {
		t.Error("expected X-Request-ID header in response")
	}
}

func TestCorrelationMiddleware_WithExistingCorrelationID(t *testing.T) {
	existingCorrelationID := "existing-correlation-123"

	middleware := CorrelationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the existing correlation ID is preserved
		correlationID := GetCorrelationID(r.Context())
		if correlationID != existingCorrelationID {
			t.Errorf("expected correlation ID to be %s, got %s", existingCorrelationID, correlationID)
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Correlation-ID", existingCorrelationID)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	// Check that the existing correlation ID is returned
	correlationIDHeader := rr.Header().Get("X-Correlation-ID")
	if correlationIDHeader != existingCorrelationID {
		t.Errorf("expected X-Correlation-ID header to be %s, got %s", existingCorrelationID, correlationIDHeader)
	}
}

func TestGetCorrelationID_NotSet(t *testing.T) {
	// Test with a context that doesn't have the correlation ID set
	correlationID := GetCorrelationID(context.Background())
	if correlationID != "" {
		t.Error("expected empty correlation ID when not set")
	}
}

func TestGetRequestID_NotSet(t *testing.T) {
	// Test with a context that doesn't have the request ID set
	requestID := GetRequestID(context.Background())
	if requestID != "" {
		t.Error("expected empty request ID when not set")
	}
}
