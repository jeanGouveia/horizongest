package handler

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type MediaHandler struct {
	svc *service.MediaService
}

func NewMediaHandler(svc *service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

// POST /api/media/upload
func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// Limitar tamanho do upload (5MB)
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	// Parse multipart form
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		jsonError(w, "arquivo muito grande (máximo 5MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonError(w, "arquivo não encontrado no formulário", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ler arquivo
	fileData, err := io.ReadAll(file)
	if err != nil {
		jsonError(w, "erro ao ler arquivo", http.StatusInternalServerError)
		return
	}

	// Determinar MIME type
	mimeType := http.DetectContentType(fileData)
	if !strings.HasPrefix(mimeType, "image/") {
		jsonError(w, "apenas imagens são permitidas (PNG, JPG, WEBP)", http.StatusBadRequest)
		return
	}

	// Extrair parâmetros
	entityType := r.FormValue("entity_type")
	entityIDStr := r.FormValue("entity_id")
	altText := r.FormValue("alt_text")

	var entityID *uint
	if entityIDStr != "" {
		id, err := strconv.ParseUint(entityIDStr, 10, 64)
		if err != nil {
			jsonError(w, "entity_id inválido", http.StatusBadRequest)
			return
		}
		eid := uint(id)
		entityID = &eid
	}

	// Validar entity_type
	if entityType == "" {
		entityType = "product" // padrão
	}

	// Criar input
	in := service.UploadMediaInput{
		FileName:     header.Filename,
		OriginalName: header.Filename,
		FileSize:     int64(len(fileData)),
		MimeType:     mimeType,
		AltText:      altText,
		EntityType:   entityType,
		EntityID:     entityID,
	}

	media, err := h.svc.UploadMedia(r.Context(), fileData, in)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileType) {
			jsonError(w, "tipo de arquivo inválido", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrFileTooLarge) {
			jsonError(w, "arquivo muito grande", http.StatusBadRequest)
			return
		}
		jsonError(w, "erro ao fazer upload da mídia", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, media)
}

// GET /api/media/{id}
func (h *MediaHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	media, err := h.svc.GetMedia(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMediaNotFound) {
			jsonError(w, "mídia não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "erro ao buscar mídia", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, media)
}

// DELETE /api/media/{id}
func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteMedia(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrMediaNotFound) {
			jsonError(w, "mídia não encontrada", http.StatusNotFound)
			return
		}
		jsonError(w, "erro ao deletar mídia", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "mídia removida com sucesso"})
}

// GET /api/media/entity/{entity_type}/{entity_id}
func (h *MediaHandler) GetMediaByEntity(w http.ResponseWriter, r *http.Request) {
	entityType := chi.URLParam(r, "entity_type")
	entityIDStr := chi.URLParam(r, "entity_id")

	entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
	if err != nil {
		jsonError(w, "entity_id inválido", http.StatusBadRequest)
		return
	}

	media, err := h.svc.GetMediaByEntity(r.Context(), entityType, uint(entityID))
	if err != nil {
		jsonError(w, "erro ao buscar mídias", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, media)
}

// GET /uploads/{path}
func (h *MediaHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	filePath := chi.URLParam(r, "*")

	// Validar caminho para evitar directory traversal
	if strings.Contains(filePath, "..") {
		http.Error(w, "caminho inválido", http.StatusBadRequest)
		return
	}

	// Construir caminho completo
	fullPath := filepath.Join("uploads", filePath)

	// Verificar se arquivo existe
	if _, err := http.Dir("uploads").Open(filePath); err != nil {
		http.Error(w, "arquivo não encontrado", http.StatusNotFound)
		return
	}

	// Servir arquivo
	http.ServeFile(w, r, fullPath)
}
