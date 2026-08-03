package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Context keys for correlation and request IDs (FASE A)
const (
	CorrelationIDKey contextKey = "correlation_id"
	RequestIDKey     contextKey = "request_id"
)

// generateUUID generates a random UUID using crypto/rand
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CorrelationMiddleware adds correlation and request IDs to the request context
// This enables tracing requests across the system
// FASE A.3: B21 - Enhanced correlation ID propagation
func CorrelationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or generate CorrelationID (can be passed from upstream services)
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateUUID()
		}

		// Generate RequestID (unique to this request)
		requestID := generateUUID()

		// Add to context
		ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		// Add to response headers
		w.Header().Set("X-Correlation-ID", correlationID)
		w.Header().Set("X-Request-ID", requestID)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCorrelationID retrieves the correlation ID from the context
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// WithCorrelationID adds a correlation ID to a context
// FASE A.3: B21 - Helper for adding correlation ID to context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// WithRequestID adds a request ID to a context
// FASE A.3: B21 - Helper for adding request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetCorrelationHeaders returns headers with correlation and request IDs
// FASE A.3: B21 - Helper for propagating correlation IDs to external services
func GetCorrelationHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)
	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		headers["X-Correlation-ID"] = correlationID
	}
	if requestID := GetRequestID(ctx); requestID != "" {
		headers["X-Request-ID"] = requestID
	}
	return headers
}
