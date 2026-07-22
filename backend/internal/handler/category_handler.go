package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// POST /api/categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var in service.CreateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.CreateCategory(r.Context(), in)
	if err != nil {
		jsonError(w, "não foi possível criar a categoria. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, c)
}

// GET /api/categories
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategories(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar as categorias. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, categories)
}

// GET /api/categories/{id}
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da categoria inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	c, err := h.svc.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			jsonError(w, "categoria não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar a categoria. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// PUT /api/categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da categoria inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	c, err := h.svc.UpdateCategory(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			jsonError(w, "categoria não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível atualizar a categoria. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, c)
}

// DELETE /api/categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID da categoria inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteCategory(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			jsonError(w, "categoria não encontrada. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível remover a categoria. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "categoria removida com sucesso"})
}
