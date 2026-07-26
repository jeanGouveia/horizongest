package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type PurchaseHandler struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseHandler(purchaseService *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: purchaseService}
}

func (h *PurchaseHandler) RegisterRoutes(r chi.Router) {
	r.Route("/suppliers", func(r chi.Router) {
		r.Post("/", h.CreateSupplier)
		r.Get("/", h.ListSuppliers)
		r.Get("/{id}", h.GetSupplierByID)
		r.Put("/{id}", h.UpdateSupplier)
		r.Delete("/{id}", h.DeleteSupplier)
	})

	r.Route("/purchase-orders", func(r chi.Router) {
		r.Post("/", h.CreatePurchaseOrder)
		r.Get("/", h.ListPurchaseOrders)
		r.Get("/{id}", h.GetPurchaseOrderByID)
		r.Put("/{id}", h.UpdatePurchaseOrder)
		r.Patch("/{id}/status", h.UpdatePurchaseOrderStatus)
		r.Delete("/{id}", h.DeletePurchaseOrder)
		r.Post("/{id}/receivings", h.CreatePurchaseReceiving)
		r.Get("/{id}/receivings", h.ListPurchaseReceivings)
	})

	r.Route("/purchase-receivings", func(r chi.Router) {
		r.Get("/{id}", h.GetPurchaseReceivingByID)
		r.Delete("/{id}", h.DeletePurchaseReceiving)
	})
}

// --- Fornecedores ---

func (h *PurchaseHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var input service.CreateSupplierInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	supplier, err := h.purchaseService.CreateSupplier(r.Context(), companyID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, supplier)
}

func (h *PurchaseHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"
	limit := 50
	offset := 0

	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	suppliers, err := h.purchaseService.ListSuppliers(r.Context(), companyID, activeOnly, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, suppliers)
}

func (h *PurchaseHandler) GetSupplierByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	supplier, err := h.purchaseService.GetSupplierByID(r.Context(), id)
	if err != nil {
		jsonError(w, "fornecedor não encontrado", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, supplier)
}

func (h *PurchaseHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	_, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var supplier service.CreateSupplierInput
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get supplier from repository and update
	// Por enquanto, retornar erro
	jsonError(w, "funcionalidade não implementada", http.StatusNotImplemented)
}

func (h *PurchaseHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.purchaseService.DeleteSupplier(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "fornecedor removido"})
}

// --- Pedidos de Compra ---

func (h *PurchaseHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var input service.CreatePurchaseOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
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

	order, err := h.purchaseService.CreatePurchaseOrder(r.Context(), companyID, userID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, order)
}

func (h *PurchaseHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit := 50
	offset := 0

	tenantCtx, ok := middleware.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	orders, err := h.purchaseService.ListPurchaseOrders(r.Context(), companyID, status, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, orders)
}

func (h *PurchaseHandler) GetPurchaseOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	order, err := h.purchaseService.GetPurchaseOrderByID(r.Context(), id)
	if err != nil {
		jsonError(w, "pedido de compra não encontrado", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, order)
}

func (h *PurchaseHandler) UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	_, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input service.CreatePurchaseOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get order from repository and update
	// Por enquanto, retornar erro
	jsonError(w, "funcionalidade não implementada", http.StatusNotImplemented)
}

func (h *PurchaseHandler) UpdatePurchaseOrderStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input struct {
		Status string `json:"status" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	status := domain.PurchaseOrderStatus(input.Status)
	if err := h.purchaseService.UpdatePurchaseOrderStatus(r.Context(), id, status); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "status atualizado"})
}

func (h *PurchaseHandler) DeletePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.purchaseService.DeletePurchaseOrder(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "pedido de compra removido"})
}

// --- Recebimentos ---

func (h *PurchaseHandler) CreatePurchaseReceiving(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input service.CreatePurchaseReceivingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get userID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		jsonError(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	receiving, err := h.purchaseService.CreatePurchaseReceiving(r.Context(), id, userID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, receiving)
}

func (h *PurchaseHandler) ListPurchaseReceivings(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	receivings, err := h.purchaseService.ListPurchaseReceivings(r.Context(), id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, receivings)
}

func (h *PurchaseHandler) GetPurchaseReceivingByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	receiving, err := h.purchaseService.GetPurchaseReceivingByID(r.Context(), id)
	if err != nil {
		jsonError(w, "recebimento não encontrado", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, receiving)
}

func (h *PurchaseHandler) DeletePurchaseReceiving(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.purchaseService.DeletePurchaseReceiving(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "recebimento removido"})
}
