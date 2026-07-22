package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type StockAdjustmentHandler struct {
	svc *service.StockAdjustmentService
}

func NewStockAdjustmentHandler(svc *service.StockAdjustmentService) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{svc: svc}
}

// ListPendingAdjustments lista ajustes de estoque pendentes com filtros opcionais
// GET /api/stock-adjustments/pending?status=pending&order_id=1&ingredient_id=2
func (h *StockAdjustmentHandler) ListPendingAdjustments(w http.ResponseWriter, r *http.Request) {
	// Extrair filtros da query string
	query := r.URL.Query()

	var statusFilter string
	if status := query.Get("status"); status != "" {
		statusFilter = status
	}

	var orderIDFilter *uint
	if orderIDStr := query.Get("order_id"); orderIDStr != "" {
		orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err != nil {
			jsonError(w, "order_id inválido", http.StatusBadRequest)
			return
		}
		oid := uint(orderID)
		orderIDFilter = &oid
	}

	var ingredientIDFilter *uint
	if ingredientIDStr := query.Get("ingredient_id"); ingredientIDStr != "" {
		ingredientID, err := strconv.ParseUint(ingredientIDStr, 10, 32)
		if err != nil {
			jsonError(w, "ingredient_id inválido", http.StatusBadRequest)
			return
		}
		iid := uint(ingredientID)
		ingredientIDFilter = &iid
	}

	// Buscar ajustes com filtros aplicados
	adjustments, err := h.svc.ListPendingAdjustmentsWithFilters(
		r.Context(),
		statusFilter,
		orderIDFilter,
		ingredientIDFilter,
	)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[HANDLER] Retornando %d ajustes pendentes", len(adjustments))
	jsonResponse(w, http.StatusOK, adjustments)
}

// ApproveAdjustmentInput representa o corpo da requisição de aprovação
type ApproveAdjustmentInput struct {
	Notes string `json:"notes"`
}

// ApproveAdjustment aprova um ajuste de estoque pendente
// POST /api/stock-adjustments/{id}/approve
func (h *StockAdjustmentHandler) ApproveAdjustment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Extrair usuário do contexto (injetado pelo middleware de autenticação)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	log.Printf("[HANDLER] UserID recuperado do contexto: %d, ok: %v", userID, ok)
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	// Parsear corpo da requisição
	var input ApproveAdjustmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Aprovar ajuste
	if err := h.svc.ApproveAdjustment(r.Context(), uint(id), userID, input.Notes); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "ajuste aprovado com sucesso"})
}

// RejectAdjustmentInput representa o corpo da requisição de rejeição
type RejectAdjustmentInput struct {
	Notes string `json:"notes"`
}

// RejectAdjustment rejeita um ajuste de estoque pendente
// POST /api/stock-adjustments/{id}/reject
func (h *StockAdjustmentHandler) RejectAdjustment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Extrair usuário do contexto (injetado pelo middleware de autenticação)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	log.Printf("[HANDLER] UserID recuperado do contexto: %d, ok: %v", userID, ok)
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	// Parsear corpo da requisição
	var input RejectAdjustmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Rejeitar ajuste
	if err := h.svc.RejectAdjustment(r.Context(), uint(id), userID, input.Notes); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "ajuste rejeitado com sucesso"})
}
