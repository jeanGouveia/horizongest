package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockOrderRepository is a mock implementation of ports.OrderRepository
type MockOrderRepository struct {
	Orders     map[uint]*domain.Order
	CreateError error
	FindError  error
	UpdateError error
	DeleteError error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[uint]*domain.Order),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if order.ID == 0 {
		order.ID = uint(len(m.Orders) + 1)
	}
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) FindOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Orders[id], nil
}

func (m *MockOrderRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	for _, order := range m.Orders {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (m *MockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Orders, id)
	return nil
}

func (m *MockOrderRepository) CreateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	return nil
}

func (m *MockOrderRepository) FindOrderItemsByOrderID(ctx context.Context, orderID uint) ([]domain.OrderItem, error) {
	return []domain.OrderItem{}, nil
}
