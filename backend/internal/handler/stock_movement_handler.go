package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type StockMovementHandler struct {
	stockMovementService *service.StockMovementService
	roleMw               *middleware.RoleMiddleware
}

func NewStockMovementHandler(stockMovementService *service.StockMovementService, roleMw *middleware.RoleMiddleware) *StockMovementHandler {
	return &StockMovementHandler{stockMovementService: stockMovementService, roleMw: roleMw}
}

func (h *StockMovementHandler) RegisterRoutes(r chi.Router) {
	r.Route("/stock-movements", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/", h.ListStockMovements)
			r.Get("/{id}", h.GetStockMovementByID)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Post("/", h.CreateStockMovement)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Delete("/{id}", h.DeleteStockMovement)
		})
	})

	r.Route("/stock-inventories", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/", h.ListInventories)
			r.Get("/{id}", h.GetInventoryByID)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Post("/", h.CreateInventory)
			r.Delete("/{id}", h.DeleteInventory)
			r.Post("/{id}/complete", h.CompleteInventory)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Post("/{id}/items", h.AddInventoryItem)
		})
	})
}

// --- Movimentações ---

func (h *StockMovementHandler) CreateStockMovement(w http.ResponseWriter, r *http.Request) {
	var input service.CreateStockMovementInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	movement, err := h.stockMovementService.CreateStockMovement(r.Context(), companyID, userID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, movement)
}

func (h *StockMovementHandler) ListStockMovements(w http.ResponseWriter, r *http.Request) {
	ingredientIDStr := r.URL.Query().Get("ingredient_id")
	var ingredientID *uint
	if ingredientIDStr != "" {
		id, err := strconv.ParseUint(ingredientIDStr, 10, 32)
		if err != nil {
			jsonError(w, "ingredient_id inválido", http.StatusBadRequest)
			return
		}
		uid := uint(id)
		ingredientID = &uid
	}

	limit := 50
	offset := 0

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	movements, err := h.stockMovementService.ListStockMovements(r.Context(), companyID, ingredientID, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, movements)
}

func (h *StockMovementHandler) GetStockMovementByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	movement, err := h.stockMovementService.GetStockMovementByID(r.Context(), id)
	if err != nil {
		jsonError(w, "movimentação não encontrada", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, movement)
}

func (h *StockMovementHandler) DeleteStockMovement(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.stockMovementService.DeleteStockMovement(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "movimentação removida"})
}

// --- Inventários ---

func (h *StockMovementHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var input service.CreateInventoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	inventory, err := h.stockMovementService.CreateInventory(r.Context(), companyID, userID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, inventory)
}

func (h *StockMovementHandler) ListInventories(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit := 50
	offset := 0

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	inventories, err := h.stockMovementService.ListInventories(r.Context(), companyID, status, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, inventories)
}

func (h *StockMovementHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	inventory, err := h.stockMovementService.GetInventoryByID(r.Context(), id)
	if err != nil {
		jsonError(w, "inventário não encontrado", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, inventory)
}

func (h *StockMovementHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.stockMovementService.DeleteInventory(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "inventário removido"})
}

func (h *StockMovementHandler) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input struct {
		IngredientID  uint    `json:"ingredientId"`
		ExpectedStock float64 `json:"expectedStock"`
		ActualStock   float64 `json:"actualStock"`
		Reason        string  `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	item, err := h.stockMovementService.AddInventoryItem(r.Context(), id, input.IngredientID, input.ExpectedStock, input.ActualStock, input.Reason)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, item)
}

func (h *StockMovementHandler) CompleteInventory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	if err := h.stockMovementService.CompleteInventory(r.Context(), id, userID); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "inventário concluído"})
}
