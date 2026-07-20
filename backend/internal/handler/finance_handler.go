package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type FinanceHandler struct {
	financeService *service.FinanceService
}

func NewFinanceHandler(financeService *service.FinanceService) *FinanceHandler {
	return &FinanceHandler{financeService: financeService}
}

func (h *FinanceHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transaction-categories", func(r chi.Router) {
		r.Post("/", h.CreateTransactionCategory)
		r.Get("/", h.ListTransactionCategories)
		r.Get("/{id}", h.GetTransactionCategoryByID)
		r.Put("/{id}", h.UpdateTransactionCategory)
		r.Delete("/{id}", h.DeleteTransactionCategory)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.CreateTransaction)
		r.Get("/", h.ListTransactions)
		r.Get("/{id}", h.GetTransactionByID)
		r.Put("/{id}", h.UpdateTransaction)
		r.Delete("/{id}", h.DeleteTransaction)
		r.Get("/cash-flow", h.GetCashFlow)
		r.Get("/summary", h.GetFinancialSummary)
	})
}

// --- Categorias ---

func (h *FinanceHandler) CreateTransactionCategory(w http.ResponseWriter, r *http.Request) {
	var input service.CreateTransactionCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get companyID from context
	companyID := uint(1) // Placeholder

	category, err := h.financeService.CreateTransactionCategory(r.Context(), companyID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, category)
}

func (h *FinanceHandler) ListTransactionCategories(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")
	var transactionType *domain.TransactionType
	if typeStr != "" {
		t := domain.TransactionType(typeStr)
		transactionType = &t
	}

	limit := 50
	offset := 0

	// TODO: Get companyID from context
	companyID := uint(1) // Placeholder

	categories, err := h.financeService.ListTransactionCategories(r.Context(), companyID, transactionType, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, categories)
}

func (h *FinanceHandler) GetTransactionCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	category, err := h.financeService.GetTransactionCategoryByID(r.Context(), id)
	if err != nil {
		jsonError(w, "categoria não encontrada", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, category)
}

func (h *FinanceHandler) UpdateTransactionCategory(w http.ResponseWriter, r *http.Request) {
	_, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input service.CreateTransactionCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get category from repository and update
	// Por enquanto, retornar erro
	jsonError(w, "funcionalidade não implementada", http.StatusNotImplemented)
}

func (h *FinanceHandler) DeleteTransactionCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.financeService.DeleteTransactionCategory(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "categoria removida"})
}

// --- Transações ---

func (h *FinanceHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input service.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get companyID and userID from context
	companyID := uint(1) // Placeholder
	userID := uint(1)    // Placeholder

	transaction, err := h.financeService.CreateTransaction(r.Context(), companyID, userID, input)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse(w, http.StatusCreated, transaction)
}

func (h *FinanceHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")
	var transactionType *domain.TransactionType
	if typeStr != "" {
		t := domain.TransactionType(typeStr)
		transactionType = &t
	}

	startDateStr := r.URL.Query().Get("start_date")
	var startDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}

	endDateStr := r.URL.Query().Get("end_date")
	var endDate *time.Time
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	limit := 50
	offset := 0

	// TODO: Get companyID from context
	companyID := uint(1) // Placeholder

	transactions, err := h.financeService.ListTransactions(r.Context(), companyID, transactionType, startDate, endDate, limit, offset)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, transactions)
}

func (h *FinanceHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	transaction, err := h.financeService.GetTransactionByID(r.Context(), id)
	if err != nil {
		jsonError(w, "transação não encontrada", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, transaction)
}

func (h *FinanceHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	_, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input service.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}

	// TODO: Get transaction from repository and update
	// Por enquanto, retornar erro
	jsonError(w, "funcionalidade não implementada", http.StatusNotImplemented)
}

func (h *FinanceHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.financeService.DeleteTransaction(r.Context(), id); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "transação removida"})
}

// --- Resumos ---

func (h *FinanceHandler) GetCashFlow(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		jsonError(w, "data inicial inválida", http.StatusBadRequest)
		return
	}

	endDateStr := r.URL.Query().Get("end_date")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		jsonError(w, "data final inválida", http.StatusBadRequest)
		return
	}

	// TODO: Get companyID from context
	companyID := uint(1) // Placeholder

	cashFlow, err := h.financeService.GetCashFlow(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, cashFlow)
}

func (h *FinanceHandler) GetFinancialSummary(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		jsonError(w, "data inicial inválida", http.StatusBadRequest)
		return
	}

	endDateStr := r.URL.Query().Get("end_date")
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		jsonError(w, "data final inválida", http.StatusBadRequest)
		return
	}

	// TODO: Get companyID from context
	companyID := uint(1) // Placeholder

	summary, err := h.financeService.GetFinancialSummary(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, summary)
}
