package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type ImpersonationHandler struct {
	impersonationService *service.ImpersonationService
}

func NewImpersonationHandler(impersonationService *service.ImpersonationService) *ImpersonationHandler {
	return &ImpersonationHandler{
		impersonationService: impersonationService,
	}
}

func (h *ImpersonationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/impersonation", func(r chi.Router) {
		r.Post("/start", h.StartImpersonation)
		r.Post("/end", h.EndImpersonation)
		r.Get("/active", h.GetActiveImpersonation)
		r.Get("/history", h.GetImpersonationHistory)
	})
}

// StartImpersonation begins a temporary impersonation session
func (h *ImpersonationHandler) StartImpersonation(w http.ResponseWriter, r *http.Request) {
	// FORENSIC: Log request complete
	log.Printf("[FORENSIC impersonation/start] REQUEST_COMPLETE - Method: %s, URL: %s", r.Method, r.URL.Path)
	log.Printf("[FORENSIC impersonation/start] REQUEST_COMPLETE - Headers: %v", r.Header)
	log.Printf("[FORENSIC impersonation/start] REQUEST_COMPLETE - Cookies: %v", r.Cookies())

	var input service.StartImpersonationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Debug: check what's in context
	log.Printf("Context values: platformUserID=%v, platformRole=%v", r.Context().Value("platformUserID"), r.Context().Value("platformRole"))

	// Get platform user ID from context
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		log.Printf("Failed to extract platformUserID as uint, got: %T", r.Context().Value("platformUserID"))
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Get platform role from context
	platformRole, ok := r.Context().Value("platformRole").(domain.PlatformRole)
	if !ok {
		log.Printf("Failed to extract platformRole as domain.PlatformRole, got: %T", r.Context().Value("platformRole"))
		jsonError(w, "não autenticado", http.StatusUnauthorized)
		return
	}

	// Only platform admins can start impersonation
	if platformRole != domain.PlatformRoleAdmin {
		jsonError(w, "apenas administradores da plataforma podem iniciar impersonation", http.StatusForbidden)
		return
	}

	// Set platform user ID from context
	input.PlatformUserID = platformUserID

	// Get IP address and user agent
	input.IPAddress = r.RemoteAddr
	input.UserAgent = r.UserAgent()

	// FORENSIC: Log JWT received
	authHeader := r.Header.Get("Authorization")
	log.Printf("[FORENSIC impersonation/start] JWT_RECEIVED - Authorization: %s", authHeader)
	log.Printf("[FORENSIC impersonation/start] JWT_VALIDATED - PlatformUserID: %d, PlatformRole: %s", platformUserID, platformRole)
	log.Printf("[FORENSIC impersonation/start] COMPANY_ID - CompanyID recebido: %d", input.CompanyID)
	log.Printf("[FORENSIC impersonation/start] USER - PlatformUserID: %d", platformUserID)

	result, err := h.impersonationService.StartImpersonation(r.Context(), input)
	if err != nil {
		if err == service.ErrCompanyOwnerNotFound {
			jsonError(w, "empresa não possui owner definido", http.StatusBadRequest)
			return
		}
		jsonError(w, "erro ao iniciar impersonation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the JWT as HttpOnly cookie - same pattern as regular login
	secureCookie := os.Getenv("ENVIRONMENT") == "production"
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// FORENSIC: Log cookie set
	log.Printf("[FORENSIC impersonation/start] COOKIE_SET - auth_token set as HttpOnly cookie")
	log.Printf("[FORENSIC impersonation/start] COOKIE_ATTRIBUTES - HttpOnly: true, Secure: %v, SameSite: Lax, Path: /, Expires: 24h", secureCookie)

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}

// EndImpersonation ends the current impersonation session
func (h *ImpersonationHandler) EndImpersonation(w http.ResponseWriter, r *http.Request) {
	var input service.EndImpersonationInput

	// Get platform user ID from context
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	// Set platform user ID from context
	input.PlatformUserID = platformUserID

	// Get IP address and user agent
	input.IPAddress = r.RemoteAddr
	input.UserAgent = r.UserAgent()

	err := h.impersonationService.EndImpersonation(r.Context(), input)
	if err != nil {
		if err == service.ErrNotImpersonating {
			jsonError(w, "não existe sessão de impersonation ativa", http.StatusBadRequest)
			return
		}
		jsonError(w, "erro ao encerrar impersonation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "impersonation encerrada com sucesso"})
}

// GetActiveImpersonation returns the active impersonation session
func (h *ImpersonationHandler) GetActiveImpersonation(w http.ResponseWriter, r *http.Request) {
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	active, err := h.impersonationService.GetActiveImpersonation(r.Context(), platformUserID)
	if err != nil {
		jsonError(w, "erro ao buscar impersonation ativa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if active == nil {
		jsonResponse(w, http.StatusOK, nil)
		return
	}

	jsonResponse(w, http.StatusOK, active)
}

// GetImpersonationHistory returns the impersonation history for the platform admin
func (h *ImpersonationHandler) GetImpersonationHistory(w http.ResponseWriter, r *http.Request) {
	platformUserID, ok := r.Context().Value("platformUserID").(uint)
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	history, err := h.impersonationService.GetImpersonationHistory(r.Context(), platformUserID)
	if err != nil {
		jsonError(w, "erro ao buscar histórico de impersonation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, history)
}
