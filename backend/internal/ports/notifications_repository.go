package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type NotificationsRepository interface {
	GetNotifications(ctx context.Context) (*domain.Notifications, error)
}
