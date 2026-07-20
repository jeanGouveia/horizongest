package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var (
	ErrMediaNotFound      = errors.New("mídia não encontrada")
	ErrInvalidFileType    = errors.New("tipo de arquivo inválido")
	ErrFileTooLarge       = errors.New("arquivo muito grande")
	ErrInvalidImageFormat = errors.New("formato de imagem inválido")
)

type MediaService struct {
	repo ports.MediaRepository
}

func NewMediaService(repo ports.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type UploadMediaInput struct {
	FileName     string
	OriginalName string
	FileSize     int64
	MimeType     string
	AltText      string
	EntityType   string
	EntityID     *uint
}

type ResizeConfig struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
}

// ── Service ──────────────────────────────────────────────────────────────────

func (s *MediaService) UploadMedia(
	ctx context.Context,
	fileData []byte,
	in UploadMediaInput,
) (*domain.Media, error) {
	// Validar tipo de arquivo
	if !isValidMimeType(in.MimeType) {
		return nil, ErrInvalidFileType
	}

	// Validar tamanho (máximo 5MB)
	if in.FileSize > 5*1024*1024 {
		return nil, ErrFileTooLarge
	}

	// Gerar nome único
	fileName := generateUniqueFileName(in.FileName)

	// Criar diretórios se não existiren
	uploadDir := "uploads/products"
	thumbDir := "uploads/products/thumbs"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("UploadMedia: criar diretório: %w", err)
	}
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, fmt.Errorf("UploadMedia: criar diretório thumbs: %w", err)
	}

	// Salvar arquivo original
	filePath := filepath.Join(uploadDir, fileName)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("UploadMedia: salvar arquivo: %w", err)
	}

	// Thumbnail generation to be implemented in future sprint
	thumbnailPath := ""

	// Criar registro no banco
	media := &domain.Media{
		FileName:      fileName,
		OriginalName:  in.OriginalName,
		FilePath:      filePath,
		ThumbnailPath: thumbnailPath,
		FileSize:      in.FileSize,
		MimeType:      in.MimeType,
		AltText:       in.AltText,
		EntityType:    in.EntityType,
		EntityID:      in.EntityID,
	}

	if err := s.repo.CreateMedia(ctx, media); err != nil {
		// Rollback: deletar arquivo se falhar no banco
		os.Remove(filePath)
		return nil, fmt.Errorf("UploadMedia: %w", err)
	}

	return media, nil
}

func (s *MediaService) GetMedia(ctx context.Context, id uint) (*domain.Media, error) {
	m, err := s.repo.FindMediaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("MediaService.GetMedia: %w", err)
	}
	if m == nil {
		return nil, ErrMediaNotFound
	}
	return m, nil
}

func (s *MediaService) GetMediaByEntity(ctx context.Context, entityType string, entityID uint) ([]domain.Media, error) {
	media, err := s.repo.FindMediaByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("MediaService.GetMediaByEntity: %w", err)
	}
	return media, nil
}

func (s *MediaService) DeleteMedia(ctx context.Context, id uint) error {
	m, err := s.repo.FindMediaByID(ctx, id)
	if err != nil {
		return fmt.Errorf("MediaService.DeleteMedia: %w", err)
	}
	if m == nil {
		return ErrMediaNotFound
	}

	// Deletar arquivo físico
	if m.FilePath != "" {
		os.Remove(m.FilePath)
	}
	if m.ThumbnailPath != "" {
		os.Remove(m.ThumbnailPath)
	}

	// Deletar registro do banco
	return s.repo.DeleteMedia(ctx, id)
}

func (s *MediaService) DeleteMediaByEntity(ctx context.Context, entityType string, entityID uint) error {
	media, err := s.repo.FindMediaByEntity(ctx, entityType, entityID)
	if err != nil {
		return fmt.Errorf("MediaService.DeleteMediaByEntity: %w", err)
	}

	// Deletar arquivos físicos
	for _, m := range media {
		if m.FilePath != "" {
			os.Remove(m.FilePath)
		}
		if m.ThumbnailPath != "" {
			os.Remove(m.ThumbnailPath)
		}
	}

	// Deletar registros do banco
	return s.repo.DeleteMediaByEntity(ctx, entityType, entityID)
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func isValidMimeType(mimeType string) bool {
	validTypes := []string{
		"image/png",
		"image/jpeg",
		"image/webp",
	}
	for _, t := range validTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d%s", name, timestamp, ext)
}
