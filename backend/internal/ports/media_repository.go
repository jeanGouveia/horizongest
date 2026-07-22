package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type MediaRepository interface {
	CreateMedia(ctx context.Context, m *domain.Media) error
	FindMediaByID(ctx context.Context, id uint) (*domain.Media, error)
	FindMediaByEntity(ctx context.Context, entityType string, entityID uint) ([]domain.Media, error)
	DeleteMedia(ctx context.Context, id uint) error
	DeleteMediaByEntity(ctx context.Context, entityType string, entityID uint) error
}
