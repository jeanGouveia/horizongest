package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var in service.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	log.Printf("Handler - Payload recebido: %+v", in)
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	order, err := h.svc.CreateOrder(r.Context(), in)
	if err != nil {
		log.Printf("Handler - Erro ao criar pedido: %v", err)
		// Erros de estoque insuficiente chegam aqui com mensagem amigável
		jsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	log.Printf("Handler - Pedido criado com sucesso: ID=%d", order.ID)
	jsonResponse(w, http.StatusCreated, order)
}

// GET /api/orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ListOrders(r.Context())
	if err != nil {
		jsonError(w, "não foi possível carregar os pedidos. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, orders)
}

// GET /api/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do pedido inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	order, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			jsonError(w, "pedido não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível carregar o pedido. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, order)
}

// PATCH /api/orders/{id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do pedido inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateOrderStatusInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	order, err := h.svc.UpdateOrderStatus(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			jsonError(w, "pedido não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInvalidOrderStatus) {
			jsonError(w, "transição de status inválida. Verifique o fluxo do pedido.", http.StatusBadRequest)
			return
		}
		jsonError(w, "não foi possível atualizar o status do pedido. Tente novamente.", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, order)
}

// PUT /api/orders/{id}
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID do pedido inválido. Verifique o valor informado.", http.StatusBadRequest)
		return
	}
	var in service.UpdateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, "formato dos dados inválido. Verifique o JSON enviado.", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(in); err != nil {
		jsonValidationError(w, err)
		return
	}
	order, err := h.svc.UpdateOrder(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			jsonError(w, "pedido não encontrado. Verifique o ID informado.", http.StatusNotFound)
			return
		}
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse(w, http.StatusOK, order)
}
