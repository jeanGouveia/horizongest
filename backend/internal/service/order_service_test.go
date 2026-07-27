package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockOrderRepository is a mock implementation of ports.OrderRepository
type MockOrderRepository struct {
	orders      map[uint]*domain.Order
	createError error
	findError   error
	updateError error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders: make(map[uint]*domain.Order),
	}
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) error {
	if m.createError != nil {
		return m.createError
	}
	order.ID = uint(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) FindOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	return m.orders[id], nil
}

func (m *MockOrderRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	for _, order := range m.orders {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error {
	if order, exists := m.orders[id]; exists {
		order.Status = status
	}
	return nil
}

func (m *MockOrderRepository) UpdateOrderStatusWithAdjustments(ctx context.Context, orderID uint, status domain.OrderStatus, productIngredients map[uint][]domain.ProductIngredient, items []domain.OrderItem) error {
	if order, exists := m.orders[orderID]; exists {
		order.Status = status
	}
	return nil
}

func (m *MockOrderRepository) UpdateOrder(ctx context.Context, id uint, items []domain.OrderItem, total float64, notes string, productIngredients map[uint][]domain.ProductIngredient) error {
	if order, exists := m.orders[id]; exists {
		order.Items = items
		order.TotalPrice = total
		order.Notes = notes
	}
	return nil
}

func (m *MockOrderRepository) ValidateStock(ctx context.Context, items []domain.OrderItem, productIngredients map[uint][]domain.ProductIngredient) (*domain.StockValidationResponse, error) {
	return &domain.StockValidationResponse{Valid: true}, nil
}

func NewMockProductRepositoryForOrder() *MockProductRepository {
	return NewMockProductRepository()
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup product
	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: true,
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 2},
		},
		Notes: "Test order",
	}

	order, err := svc.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order.Status != domain.OrderStatusPending {
		t.Errorf("expected status pending, got %s", order.Status)
	}
	if order.TotalPrice != 20.0 {
		t.Errorf("expected total 20.0, got %f", order.TotalPrice)
	}
	if len(order.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(order.Items))
	}
}

func TestOrderService_CreateOrder_ProductNotFound(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 999, Quantity: 1},
		},
	}

	_, err := svc.CreateOrder(context.Background(), input)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
}

func TestOrderService_CreateOrder_InactiveProduct(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup inactive product
	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: false,
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1},
		},
	}

	_, err := svc.CreateOrder(context.Background(), input)
	if err == nil {
		t.Error("expected error for inactive product")
	}
}

func TestOrderService_CreateOrder_InsufficientStock(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup product with ingredients
	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: true,
	}

	mockProductRepo.ingredients[2] = &domain.Ingredient{
		ID:            2,
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 0.5, // Only 0.5 kg available
	}

	mockProductRepo.productIngredients[1] = []domain.ProductIngredient{
		{IngredientID: 2, Quantity: 1.0}, // Requires 1.0 kg per cake
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1}, // Needs 1.0 kg, only have 0.5 kg
		},
	}

	_, err := svc.CreateOrder(context.Background(), input)
	if err == nil {
		t.Error("expected error for insufficient stock")
	}
	var stockErr *InsufficientStockError
	if !errors.As(err, &stockErr) {
		t.Errorf("expected InsufficientStockError, got %T", err)
	}
}

func TestOrderService_CreateOrder_SimpleProduct(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup simple product without ingredients
	mockProductRepo.products[1] = &domain.Product{
		ID:         1,
		Name:       "Water",
		Price:      2.0,
		Active:     true,
		IsComposto: false,
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 5},
		},
	}

	order, err := svc.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order.TotalPrice != 10.0 {
		t.Errorf("expected total 10.0, got %f", order.TotalPrice)
	}
}

func TestOrderService_ListOrders(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup orders
	mockOrderRepo.orders[1] = &domain.Order{
		ID:         1,
		Status:     domain.OrderStatusPending,
		TotalPrice: 10.0,
	}

	orders, err := svc.ListOrders(context.Background())
	if err != nil {
		t.Fatalf("ListOrders failed: %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(orders))
	}
}

func TestOrderService_GetOrder(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:         1,
		Status:     domain.OrderStatusPending,
		TotalPrice: 10.0,
	}

	order, err := svc.GetOrder(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}
	if order.ID != 1 {
		t.Errorf("expected ID 1, got %d", order.ID)
	}
}

func TestOrderService_GetOrder_NotFound(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	_, err := svc.GetOrder(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent order")
	}
	if !errors.Is(err, ErrOrderNotFound) && err.Error() != ErrOrderNotFound.Error() {
		t.Errorf("expected ErrOrderNotFound, got: %v", err)
	}
}

func TestOrderService_UpdateOrderStatus_ValidTransition(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusPending,
	}

	input := UpdateOrderStatusInput{
		Status: "confirmed",
	}

	order, err := svc.UpdateOrderStatus(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateOrderStatus failed: %v", err)
	}
	if order.Status != domain.OrderStatusConfirmed {
		t.Errorf("expected status confirmed, got %s", order.Status)
	}
}

func TestOrderService_UpdateOrderStatus_InvalidTransition(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusDelivered,
	}

	input := UpdateOrderStatusInput{
		Status: "preparing",
	}

	_, err := svc.UpdateOrderStatus(context.Background(), 1, input)
	if err == nil {
		t.Error("expected error for invalid status transition")
	}
}

func TestOrderService_UpdateOrderStatus_CancelOrder(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusConfirmed,
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 1},
		},
	}

	mockProductRepo.productIngredients[1] = []domain.ProductIngredient{
		{IngredientID: 2, Quantity: 1.0},
	}

	input := UpdateOrderStatusInput{
		Status: "cancelled",
	}

	order, err := svc.UpdateOrderStatus(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateOrderStatus failed: %v", err)
	}
	if order.Status != domain.OrderStatusCancelled {
		t.Errorf("expected status cancelled, got %s", order.Status)
	}
}

func TestOrderService_UpdateOrderStatus_NotFound(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	input := UpdateOrderStatusInput{
		Status: "confirmed",
	}

	_, err := svc.UpdateOrderStatus(context.Background(), 999, input)
	if err == nil {
		t.Error("expected error for non-existent order")
	}
	if !errors.Is(err, ErrOrderNotFound) && err.Error() != ErrOrderNotFound.Error() {
		t.Errorf("expected ErrOrderNotFound, got: %v", err)
	}
}

func TestOrderService_UpdateOrder_Success(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusPending,
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 1, UnitPrice: 10.0},
		},
	}

	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: true,
	}

	mockProductRepo.products[2] = &domain.Product{
		ID:     2,
		Name:   "Coffee",
		Price:  5.0,
		Active: true,
	}

	input := UpdateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 2, Quantity: 2},
		},
		Notes: "Updated order",
	}

	order, err := svc.UpdateOrder(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateOrder failed: %v", err)
	}
	if len(order.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(order.Items))
	}
	if order.TotalPrice != 10.0 {
		t.Errorf("expected total 10.0, got %f", order.TotalPrice)
	}
}

func TestOrderService_UpdateOrder_InvalidStatus(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusDelivered,
	}

	input := UpdateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1},
		},
	}

	_, err := svc.UpdateOrder(context.Background(), 1, input)
	if err == nil {
		t.Error("expected error for invalid order status")
	}
}

func TestOrderService_UpdateOrder_InsufficientStock(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	mockOrderRepo.orders[1] = &domain.Order{
		ID:     1,
		Status: domain.OrderStatusPending,
	}

	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: true,
	}

	mockProductRepo.ingredients[2] = &domain.Ingredient{
		ID:            2,
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 0.5,
	}

	mockProductRepo.productIngredients[1] = []domain.ProductIngredient{
		{IngredientID: 2, Quantity: 1.0},
	}

	input := UpdateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1},
		},
	}

	_, err := svc.UpdateOrder(context.Background(), 1, input)
	if err == nil {
		t.Error("expected error for insufficient stock")
	}
	var stockErr *InsufficientStockError
	if !errors.As(err, &stockErr) {
		t.Errorf("expected InsufficientStockError, got %T", err)
	}
}

func TestOrderService_UpdateOrder_NotFound(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	input := UpdateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1},
		},
	}

	_, err := svc.UpdateOrder(context.Background(), 999, input)
	if err == nil {
		t.Error("expected error for non-existent order")
	}
	if !errors.Is(err, ErrOrderNotFound) && err.Error() != ErrOrderNotFound.Error() {
		t.Errorf("expected ErrOrderNotFound, got: %v", err)
	}
}

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		name     string
		current  domain.OrderStatus
		new      domain.OrderStatus
		expected bool
	}{
		{"Pending to Confirmed", domain.OrderStatusPending, domain.OrderStatusConfirmed, true},
		{"Pending to Cancelled", domain.OrderStatusPending, domain.OrderStatusCancelled, true},
		{"Confirmed to Preparing", domain.OrderStatusConfirmed, domain.OrderStatusPreparing, true},
		{"Confirmed to Cancelled", domain.OrderStatusConfirmed, domain.OrderStatusCancelled, true},
		{"Preparing to Ready", domain.OrderStatusPreparing, domain.OrderStatusReady, true},
		{"Preparing to Cancelled", domain.OrderStatusPreparing, domain.OrderStatusCancelled, true},
		{"Ready to Delivered", domain.OrderStatusReady, domain.OrderStatusDelivered, true},
		{"Ready to Cancelled", domain.OrderStatusReady, domain.OrderStatusCancelled, true},
		{"Delivered to any", domain.OrderStatusDelivered, domain.OrderStatusPending, false},
		{"Cancelled to any", domain.OrderStatusCancelled, domain.OrderStatusPending, false},
		{"Invalid transition", domain.OrderStatusPending, domain.OrderStatusDelivered, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTransition(tt.current, tt.new)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestOrderService_CreateOrder_MultipleItems(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup products
	mockProductRepo.products[1] = &domain.Product{
		ID:     1,
		Name:   "Cake",
		Price:  10.0,
		Active: true,
	}

	mockProductRepo.products[2] = &domain.Product{
		ID:     2,
		Name:   "Coffee",
		Price:  5.0,
		Active: true,
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}

	order, err := svc.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order.TotalPrice != 35.0 {
		t.Errorf("expected total 35.0, got %f", order.TotalPrice)
	}
	if len(order.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(order.Items))
	}
}

func TestOrderService_CreateOrder_WithProductSnapshots(t *testing.T) {
	mockOrderRepo := NewMockOrderRepository()
	mockProductRepo := NewMockProductRepositoryForOrder()
	svc := NewOrderService(mockOrderRepo, mockProductRepo)

	// Setup product with all fields
	mockProductRepo.products[1] = &domain.Product{
		ID:             1,
		Name:           "Cake",
		Description:    "Delicious cake",
		Price:          10.0,
		Active:         true,
		PhotoURL:       "http://example.com/cake.jpg",
		IsComposto:     true,
		PromotionPrice: &[]float64{8.0}[0],
		Featured:       true,
		IsNew:          true,
	}

	input := CreateOrderInput{
		Items: []OrderItemInput{
			{ProductID: 1, Quantity: 1},
		},
	}

	order, err := svc.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	item := order.Items[0]
	if item.ProductName != "Cake" {
		t.Errorf("expected product name 'Cake', got '%s'", item.ProductName)
	}
	if item.ProductDescription != "Delicious cake" {
		t.Errorf("expected product description 'Delicious cake', got '%s'", item.ProductDescription)
	}
	if item.UnitPrice != 10.0 {
		t.Errorf("expected unit price 10.0, got %f", item.UnitPrice)
	}
	if !item.ProductIsComposto {
		t.Error("expected product to be marked as composite")
	}
	if item.ProductPromotionPrice == nil || *item.ProductPromotionPrice != 8.0 {
		t.Error("expected promotion price to be set")
	}
	if !item.ProductFeatured {
		t.Error("expected product to be marked as featured")
	}
	if !item.ProductIsNew {
		t.Error("expected product to be marked as new")
	}
}
