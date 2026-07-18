package ports

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

type DashboardRepository interface {
	GetDashboard(ctx context.Context) (*domain.Dashboard, error)
}
