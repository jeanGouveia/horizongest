package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type PlatformBrandHandler struct {
	brandService *service.PlatformBrandService
	validate     *validator.Validate
}

func NewPlatformBrandHandler(brandService *service.PlatformBrandService) *PlatformBrandHandler {
	return &PlatformBrandHandler{
		brandService: brandService,
		validate:     validator.New(),
	}
}

// GetPlatformBrand returns the current platform brand configuration
func (h *PlatformBrandHandler) GetPlatformBrand(w http.ResponseWriter, r *http.Request) {
	brand, err := h.brandService.Get(r.Context())
	if err != nil {
		jsonError(w, "não foi possível obter configuração de marca da plataforma", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, brand)
}

// GetPublicPlatformBrand returns the public platform branding information
// This endpoint does not require authentication and is used by the frontend
// to display branding dynamically (Sprint 3.6 - White Label)
func (h *PlatformBrandHandler) GetPublicPlatformBrand(w http.ResponseWriter, r *http.Request) {
	brand, err := h.brandService.Get(r.Context())
	if err != nil {
		jsonError(w, "não foi possível obter configuração de marca da plataforma", http.StatusInternalServerError)
		return
	}

	// Return only public-safe information
	publicBrand := map[string]interface{}{
		"platformName":      brand.PlatformName,
		"platformShortName": brand.PlatformShortName,
		"website":           brand.Website,
		"logoPath":          brand.LogoPath,
		"faviconPath":       brand.FaviconPath,
		"logoLight":         brand.LogoLight,
		"logoDark":          brand.LogoDark,
		"icon":              brand.Icon,
		"loginBackground":   brand.LoginBackground,
		"loginIllustration": brand.LoginIllustration,
		"copyright":         brand.Copyright,
		"primaryColor":      brand.PrimaryColor,
		"secondaryColor":    brand.SecondaryColor,
	}

	jsonResponse(w, http.StatusOK, publicBrand)
}

// UpdatePlatformBrandInput represents the input for updating platform brand configuration
type UpdatePlatformBrandInput struct {
	PlatformName       string `json:"platformName" validate:"required"`
	PlatformShortName  string `json:"platformShortName" validate:"required"`
	OwnerCompanyName   string `json:"ownerCompanyName" validate:"required"`
	OwnerDocument      string `json:"ownerDocument"`
	Website            string `json:"website" validate:"required,url"`
	SupportEmail       string `json:"supportEmail" validate:"required,email"`
	SupportURL         string `json:"supportUrl" validate:"required,url"`
	LogoPath           string `json:"logoPath"`
	FaviconPath        string `json:"faviconPath"`
	LogoLight          string `json:"logoLight"`
	LogoDark           string `json:"logoDark"`
	Icon               string `json:"icon"`
	LoginBackground    string `json:"loginBackground"`
	LoginIllustration  string `json:"loginIllustration"`
	Copyright          string `json:"copyright" validate:"required"`
	PrivacyPolicyURL   string `json:"privacyPolicyUrl"`
	TermsURL           string `json:"termsUrl"`
	InstagramURL       string `json:"instagramUrl"`
	FacebookURL        string `json:"facebookUrl"`
	LinkedInURL        string `json:"linkedinUrl"`
	YoutubeURL         string `json:"youtubeUrl"`
	DefaultLanguage    string `json:"defaultLanguage"`
	DefaultTimezone    string `json:"defaultTimezone"`
	MaintenanceMode    bool   `json:"maintenanceMode"`
	MaintenanceMessage string `json:"maintenanceMessage"`
	PrimaryColor       string `json:"primaryColor" validate:"required"`
	SecondaryColor     string `json:"secondaryColor" validate:"required"`
}

// UpdatePlatformBrand updates the platform brand configuration
// Only platform admins should be allowed to call this
func (h *PlatformBrandHandler) UpdatePlatformBrand(w http.ResponseWriter, r *http.Request) {
	var in UpdatePlatformBrandInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}

	// Get user ID from context (should be set by platform auth middleware)
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		jsonError(w, "não autorizado", http.StatusUnauthorized)
		return
	}

	brand := &domain.PlatformBrandConfig{
		PlatformName:       in.PlatformName,
		PlatformShortName:  in.PlatformShortName,
		OwnerCompanyName:   in.OwnerCompanyName,
		OwnerDocument:      in.OwnerDocument,
		Website:            in.Website,
		SupportEmail:       in.SupportEmail,
		SupportURL:         in.SupportURL,
		LogoPath:           in.LogoPath,
		FaviconPath:        in.FaviconPath,
		LogoLight:          in.LogoLight,
		LogoDark:           in.LogoDark,
		Icon:               in.Icon,
		LoginBackground:    in.LoginBackground,
		LoginIllustration:  in.LoginIllustration,
		Copyright:          in.Copyright,
		PrivacyPolicyURL:   in.PrivacyPolicyURL,
		TermsURL:           in.TermsURL,
		InstagramURL:       in.InstagramURL,
		FacebookURL:        in.FacebookURL,
		LinkedInURL:        in.LinkedInURL,
		YoutubeURL:         in.YoutubeURL,
		DefaultLanguage:    in.DefaultLanguage,
		DefaultTimezone:    in.DefaultTimezone,
		MaintenanceMode:    in.MaintenanceMode,
		MaintenanceMessage: in.MaintenanceMessage,
		PrimaryColor:       in.PrimaryColor,
		SecondaryColor:     in.SecondaryColor,
	}

	err := h.brandService.Update(r.Context(), brand, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPlatformBrand) {
			jsonError(w, "configuração de marca inválida", http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível atualizar configuração de marca", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "configuração de marca atualizada com sucesso"})
}
