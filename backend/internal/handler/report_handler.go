package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		r.Get("/sales", h.GetSalesReport)
		r.Get("/products", h.GetProductsReport)
		r.Get("/cmv", h.GetCMVReport)
		r.Get("/profit", h.GetProfitReport)
		r.Get("/stock", h.GetStockReport)
		r.Get("/purchases", h.GetPurchasesReport)
		r.Get("/financial", h.GetFinancialReport)
	})
}

// GetSalesReport retorna relatório de vendas
func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
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

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetSalesReport(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetProductsReport retorna relatório de produtos
func (h *ReportHandler) GetProductsReport(w http.ResponseWriter, r *http.Request) {
	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetProductsReport(r.Context(), companyID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetCMVReport retorna relatório de CMV
func (h *ReportHandler) GetCMVReport(w http.ResponseWriter, r *http.Request) {
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

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetCMVReport(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetProfitReport retorna relatório de lucro
func (h *ReportHandler) GetProfitReport(w http.ResponseWriter, r *http.Request) {
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

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetProfitReport(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetStockReport retorna relatório de estoque
func (h *ReportHandler) GetStockReport(w http.ResponseWriter, r *http.Request) {
	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetStockReport(r.Context(), companyID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetPurchasesReport retorna relatório de compras
func (h *ReportHandler) GetPurchasesReport(w http.ResponseWriter, r *http.Request) {
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

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetPurchasesReport(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}

// GetFinancialReport retorna relatório financeiro
func (h *ReportHandler) GetFinancialReport(w http.ResponseWriter, r *http.Request) {
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

	tenantCtx, ok := domain.GetTenantContextFromContext(r.Context())
	if !ok {
		jsonError(w, "contexto tenant não encontrado", http.StatusUnauthorized)
		return
	}
	companyID := tenantCtx.CompanyID

	report, err := h.reportService.GetFinancialReport(r.Context(), companyID, startDate, endDate)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, report)
}
