package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
	"github.com/jeanGouveia/pratoOnline/backend/internal/util"
)

type ProductHandler struct {
	svc       *service.ProductService
	sanitizer *util.Sanitizer
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{
		svc:       svc,
		sanitizer: util.NewSanitizer(),
	}
}

// ── Produto ──────────────────────────────────────────────────────────────────

// POST /api/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var in service.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}

	// Sanitize inputs (Sprint 3.4 - Security Hardening)
	sanitizedName, err := h.sanitizer.SanitizeName(in.Name)
	if err != nil {
		jsonError(w, "nome inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	in.Name = sanitizedName

	if in.Description != "" {
		sanitizedDesc, err := h.sanitizer.SanitizeDescription(in.Description)
		if err != nil {
			jsonError(w, "descrição inválida: "+err.Error(), http.StatusBadRequest)
			return
		}
		in.Description = sanitizedDesc
	}

	if in.Slug != "" {
		sanitizedSlug, err := h.sanitizer.SanitizeSlug(in.Slug)
		if err != nil {
			jsonError(w, "slug inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
		in.Slug = sanitizedSlug
	}

	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	p, err := h.svc.CreateProduct(r.Context(), in)
	if err != nil {
		jsonError(w, "não foi possível criar o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, p)
}

// GET /api/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar os produtos. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, products)
}

// GET /api/products/active
func (h *ProductHandler) ListActiveProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListActiveProducts(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar os produtos ativos. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, products)
}

// GET /api/products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	p, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, p)
}

// DELETE /api/products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível remover o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "produto removido com sucesso"})
}

// PUT /api/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	p, err := h.svc.UpdateProduct(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível atualizar o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, p)
}

// ── Ingrediente ──────────────────────────────────────────────────────────────

// POST /api/ingredients
func (h *ProductHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	var in service.CreateIngredientInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	i, err := h.svc.CreateIngredient(r.Context(), in)
	if err != nil {
		jsonError(w, "não foi possível criar o ingrediente. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, i)
}

// GET /api/ingredients
func (h *ProductHandler) ListIngredients(w http.ResponseWriter, r *http.Request) {
	ingredients, err := h.svc.ListIngredients(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar os ingredientes. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, ingredients)
}

// GET /api/ingredients/{id}
func (h *ProductHandler) GetIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do ingrediente inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	i, err := h.svc.GetIngredient(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrIngredientNotFound) {
			jsonError(w, "ingrediente não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar o ingrediente. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, i)
}

// PATCH /api/ingredients/{id}/stock
func (h *ProductHandler) UpdateIngredientStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do ingrediente inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateStockInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	i, err := h.svc.UpdateIngredientStock(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrIngredientNotFound) {
			jsonError(w, "ingrediente não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível ajustar o estoque. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, i)
}

// PUT /api/ingredients/{id}
func (h *ProductHandler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do ingrediente inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateIngredientInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	i, err := h.svc.UpdateIngredient(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrIngredientNotFound) {
			jsonError(w, "ingrediente não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível atualizar o ingrediente. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, i)
}

// DELETE /api/ingredients/{id}
func (h *ProductHandler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do ingrediente inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteIngredient(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrIngredientNotFound) {
			jsonError(w, "ingrediente não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível remover o ingrediente. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "ingrediente removido com sucesso"})
}

// ── Ficha técnica ────────────────────────────────────────────────────────────

// PUT /api/products/{id}/ingredients
func (h *ProductHandler) SetProductIngredients(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.SetProductIngredientsInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	if err := h.svc.SetProductIngredients(r.Context(), id, in); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível atualizar a ficha técnica. Verifique os ingredientes informados.", http.StatusBadRequest)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "ficha técnica atualizada com sucesso"})
}

// GET /api/products/{id}/ingredients
func (h *ProductHandler) GetProductIngredients(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	p, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar a ficha técnica. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, p.Ingredients)
}

// POST /api/products/{id}/duplicate
func (h *ProductHandler) DuplicateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	product, err := h.svc.DuplicateProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível duplicar o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, product)
}

// POST /api/products/{id}/archive
func (h *ProductHandler) ArchiveProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do produto inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	if err := h.svc.ArchiveProduct(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			jsonError(w, "produto não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível arquivar o produto. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "produto arquivado com sucesso"})
}

// ── helper ───────────────────────────────────────────────────────────────────

func parseID(r *http.Request, param string) (uint, error) {
	raw := chi.URLParam(r, param)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
