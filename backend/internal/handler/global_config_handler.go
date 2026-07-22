package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type GlobalConfigHandler struct {
	configService *service.GlobalConfigService
}

func NewGlobalConfigHandler(configService *service.GlobalConfigService) *GlobalConfigHandler {
	return &GlobalConfigHandler{
		configService: configService,
	}
}

// GetGlobalConfig retrieves the current global configuration
// Only platform admins should be allowed to call this
func (h *GlobalConfigHandler) GetGlobalConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.configService.Get(r.Context())
	if err != nil {
		jsonError(w, "não foi possível obter configuração global", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, config)
}

// UpdateGlobalConfigInput represents the input for updating global configuration
type UpdateGlobalConfigInput struct {
	DefaultTimezone    string `json:"defaultTimezone" validate:"required"`
	DefaultLocale      string `json:"defaultLocale" validate:"required"`
	MonetaryFormat     string `json:"monetaryFormat" validate:"required"`
	DateFormat         string `json:"dateFormat" validate:"required"`
	TimeFormat         string `json:"timeFormat" validate:"required"`
	MaxUploadSizeMB    int64  `json:"maxUploadSizeMb" validate:"required,min=1"`
	MaxImageSizeMB     int64  `json:"maxImageSizeMb" validate:"required,min=1"`
	AllowedImageTypes  string `json:"allowedImageTypes" validate:"required"`
	AllowedFileTypes   string `json:"allowedFileTypes" validate:"required"`
	MaintenanceMode    bool   `json:"maintenanceMode"`
	MaintenanceMessage string `json:"maintenanceMessage"`
	EnableFinance      bool   `json:"enableFinance"`
	EnablePurchasing   bool   `json:"enablePurchasing"`
	EnableInventory    bool   `json:"enableInventory"`
	EnableCRM          bool   `json:"enableCRM"`
	EnableCalendar     bool   `json:"enableCalendar"`
	EnablePOS          bool   `json:"enablePOS"`
	EnableAI           bool   `json:"enableAI"`
	EnableDelivery     bool   `json:"enableDelivery"`
	EnableMarketplace  bool   `json:"enableMarketplace"`
}

// UpdateGlobalConfig updates the global configuration
// Only platform admins should be allowed to call this
func (h *GlobalConfigHandler) UpdateGlobalConfig(w http.ResponseWriter, r *http.Request) {
	var in UpdateGlobalConfigInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Get user ID from context (should be set by platform auth middleware)
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	config := &domain.GlobalConfig{
		DefaultTimezone:    in.DefaultTimezone,
		DefaultLocale:      in.DefaultLocale,
		MonetaryFormat:     in.MonetaryFormat,
		DateFormat:         in.DateFormat,
		TimeFormat:         in.TimeFormat,
		MaxUploadSizeMB:    in.MaxUploadSizeMB,
		MaxImageSizeMB:     in.MaxImageSizeMB,
		AllowedImageTypes:  in.AllowedImageTypes,
		AllowedFileTypes:   in.AllowedFileTypes,
		MaintenanceMode:    in.MaintenanceMode,
		MaintenanceMessage: in.MaintenanceMessage,
		EnableFinance:      in.EnableFinance,
		EnablePurchasing:   in.EnablePurchasing,
		EnableInventory:    in.EnableInventory,
		EnableCRM:          in.EnableCRM,
		EnableCalendar:     in.EnableCalendar,
		EnablePOS:          in.EnablePOS,
		EnableAI:           in.EnableAI,
		EnableDelivery:     in.EnableDelivery,
		EnableMarketplace:  in.EnableMarketplace,
	}

	err := h.configService.Update(r.Context(), config, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidGlobalConfig) {
			jsonError(w, "configuração global inválida", http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível atualizar configuração global", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, config)
}

// GetModuleStatus checks if a specific module is enabled
// Only platform admins should be allowed to call this
func (h *GlobalConfigHandler) GetModuleStatus(w http.ResponseWriter, r *http.Request) {
	module := r.URL.Query().Get("module")
	if module == "" {
		jsonError(w, "parâmetro 'module' é obrigatório", http.StatusBadRequest)
		return
	}

	enabled, err := h.configService.IsModuleEnabled(r.Context(), module)
	if err != nil {
		jsonError(w, "não foi possível verificar status do módulo", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{
		"enabled": enabled,
	})
}
