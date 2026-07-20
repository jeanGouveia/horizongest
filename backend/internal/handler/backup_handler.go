package handler

import (
	"net/http"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type BackupHandler struct {
	backupService *service.BackupService
}

func NewBackupHandler(backupService *service.BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

// POST /api/platform/backup
func (h *BackupHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	result, err := h.backupService.CreateBackup(r.Context())
	if err != nil {
		jsonError(w, "não foi possível criar backup", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, result)
}

// GET /api/platform/backup
func (h *BackupHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	backups, err := h.backupService.ListBackups(r.Context())
	if err != nil {
		jsonError(w, "não foi possível listar backups", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"backups": backups,
	})
}

// DELETE /api/platform/backup/:filename
func (h *BackupHandler) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		jsonError(w, "nome do arquivo não fornecido", http.StatusBadRequest)
		return
	}

	if err := h.backupService.DeleteBackup(r.Context(), fileName); err != nil {
		jsonError(w, "não foi possível deletar backup", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "backup deletado com sucesso",
	})
}
