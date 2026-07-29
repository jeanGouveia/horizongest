package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware para permitir requisições de origens diferentes
// Sprint 4A: Implementa whitelist de origens por ambiente, nunca reflete Origin
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Obter whitelist de origens permitidas do ambiente
		allowedOrigins := getAllowedOrigins()

		// Verificar se a origem está na whitelist
		allowed := false
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
				// Suporte a wildcard no final (ex: https://*.example.com)
				if strings.HasSuffix(allowedOrigin, "*") {
					prefix := strings.TrimSuffix(allowedOrigin, "*")
					if strings.HasPrefix(origin, prefix) {
						allowed = true
						break
					}
				}
			}
		}

		// Se a origem for permitida, definir o header
		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getAllowedOrigins retorna a whitelist de origens permitidas baseada no ambiente
func getAllowedOrigins() []string {
	env := os.Getenv("ENVIRONMENT")

	// Configuração por ambiente
	switch env {
	case "production":
		// Produção: apenas origens explicitamente configuradas
		origins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if origins == "" {
			// Fallback seguro: nenhuma origem permitida se não configurado
			return []string{}
		}
		return strings.Split(origins, ",")
	case "staging":
		// Staging: origens configuradas ou localhost
		origins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if origins == "" {
			return []string{
				"http://localhost:5173",
				"http://localhost:3000",
				"http://127.0.0.1:5173",
				"http://127.0.0.1:3000",
			}
		}
		return strings.Split(origins, ",")
	default:
		// Desenvolvimento: permite localhost
		return []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:3000",
		}
	}
}
