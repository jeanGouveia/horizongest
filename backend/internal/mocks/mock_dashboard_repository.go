package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockDashboardRepository is a mock implementation of ports.DashboardRepository
type MockDashboardRepository struct {
	Dashboard  *domain.Dashboard
	FindError  error
}

func NewMockDashboardRepository() *MockDashboardRepository {
	return &MockDashboardRepository{
		Dashboard: &domain.Dashboard{},
	}
}

func (m *MockDashboardRepository) GetDashboard(ctx context.Context) (*domain.Dashboard, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Dashboard, nil
}
