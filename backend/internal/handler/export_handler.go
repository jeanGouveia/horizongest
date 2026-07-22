package handler

import (
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

type ExportHandler struct {
	exportService *service.ExportService
}

func NewExportHandler(exportService *service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// POST /api/platform/export/companies
func (h *ExportHandler) ExportCompanies(w http.ResponseWriter, r *http.Request) {
	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		formatStr = "csv"
	}

	format := service.ExportFormat(formatStr)
	if format != service.ExportFormatCSV && format != service.ExportFormatJSON {
		jsonError(w, "formato inválido. Use 'csv' ou 'json'", http.StatusBadRequest)
		return
	}

	result, err := h.exportService.ExportCompanies(r.Context(), format)
	if err != nil {
		jsonError(w, "não foi possível exportar empresas", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, result)
}

// POST /api/platform/export/users
func (h *ExportHandler) ExportUsers(w http.ResponseWriter, r *http.Request) {
	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		formatStr = "csv"
	}

	format := service.ExportFormat(formatStr)
	if format != service.ExportFormatCSV && format != service.ExportFormatJSON {
		jsonError(w, "formato inválido. Use 'csv' ou 'json'", http.StatusBadRequest)
		return
	}

	result, err := h.exportService.ExportUsers(r.Context(), format)
	if err != nil {
		jsonError(w, "não foi possível exportar usuários", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, result)
}
