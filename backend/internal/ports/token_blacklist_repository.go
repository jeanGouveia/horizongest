package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type TokenBlacklistRepository interface {
	Add(ctx context.Context, entry *domain.TokenBlacklist) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanExpired(ctx context.Context) error
}
