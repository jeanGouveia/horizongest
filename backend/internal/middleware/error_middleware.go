package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// ErrorResponseMiddleware padroniza todas as respostas de erro
func ErrorResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gerar request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)

		// Criar response writer customizado para capturar status code
		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		// Se houve erro e ainda não foi escrito, padronizar
		if rw.status >= 400 && !rw.written {
			sendErrorResponse(w, rw.status, "error", http.StatusText(rw.status), requestID)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}

func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendErrorResponse(w http.ResponseWriter, status int, code, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResp := domain.NewErrorResponse(code, message, requestID)
	json.NewEncoder(w).Encode(errorResp)
}
