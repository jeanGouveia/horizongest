package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type DashboardRepository interface {
	GetDashboard(ctx context.Context) (*domain.Dashboard, error)
}
