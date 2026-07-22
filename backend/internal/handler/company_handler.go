package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type CompanyHandler struct {
	svc *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// POST /api/companies
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var in service.CreateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.CreateCompany(r.Context(), in)
	if err != nil {
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			jsonError(w, "slug já está em uso. Escolha outro identificador.", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível criar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, c)
}

// GET /api/companies
func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.svc.ListCompanies(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar as empresas. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, companies)
}

// GET /api/companies/{id}
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	c, err := h.svc.GetCompany(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// PUT /api/companies/{id}
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.UpdateCompany(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			jsonError(w, "slug já está em uso. Escolha outro identificador.", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível atualizar a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// DELETE /api/companies/{id}
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da empresa inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteCompany(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrCompanyNotFound) {
			jsonError(w, "empresa não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível remover a empresa. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "empresa removida com sucesso"})
}
