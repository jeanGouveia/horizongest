package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/gorm"
)

// MockProductRepository is a mock implementation of ports.ProductRepository
type MockProductRepository struct {
	products           map[uint]*domain.Product
	ingredients        map[uint]*domain.Ingredient
	productIngredients map[uint][]domain.ProductIngredient
	createError        error
	findError          error
	updateError        error
	deleteError        error
	ingredientsError   error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products:           make(map[uint]*domain.Product),
		ingredients:        make(map[uint]*domain.Ingredient),
		productIngredients: make(map[uint][]domain.ProductIngredient),
	}
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	if m.createError != nil {
		return m.createError
	}
	product.ID = uint(len(m.products) + 1)
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) FindProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	return m.products[id], nil
}

func (m *MockProductRepository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	for _, product := range m.products {
		products = append(products, *product)
	}
	return products, nil
}

func (m *MockProductRepository) ListActiveProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	for _, product := range m.products {
		if product.Active {
			products = append(products, *product)
		}
	}
	return products, nil
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id uint) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) CreateIngredient(ctx context.Context, ingredient *domain.Ingredient) error {
	ingredient.ID = uint(len(m.ingredients) + 1)
	m.ingredients[ingredient.ID] = ingredient
	return nil
}

func (m *MockProductRepository) FindIngredientByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error) {
	if m.ingredientsError != nil {
		return nil, m.ingredientsError
	}
	return m.ingredients[id], nil
}

// FindIngredientByIDForUpdate mock implementation
func (m *MockProductRepository) FindIngredientByIDForUpdate(ctx context.Context, id uint, tx *gorm.DB) (*domain.Ingredient, error) {
	if m.ingredientsError != nil {
		return nil, m.ingredientsError
	}
	return m.ingredients[id], nil
}

func (m *MockProductRepository) ListIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	for _, ingredient := range m.ingredients {
		ingredients = append(ingredients, *ingredient)
	}
	return ingredients, nil
}

func (m *MockProductRepository) UpdateIngredient(ctx context.Context, ingredient *domain.Ingredient, tx *gorm.DB) error {
	m.ingredients[ingredient.ID] = ingredient
	return nil
}

func (m *MockProductRepository) DeleteIngredient(ctx context.Context, id uint) error {
	delete(m.ingredients, id)
	return nil
}

func (m *MockProductRepository) GetProductIngredients(ctx context.Context, productID uint) ([]domain.ProductIngredient, error) {
	return m.productIngredients[productID], nil
}

func (m *MockProductRepository) SetProductIngredients(ctx context.Context, productID uint, ingredients []domain.ProductIngredient) error {
	m.productIngredients[productID] = ingredients
	return nil
}

func (m *MockProductRepository) CanDeleteProduct(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	return &domain.DependencyCheck{CanDelete: true}, nil
}

func (m *MockProductRepository) CanDeleteIngredient(ctx context.Context, id uint) (*domain.DependencyCheck, error) {
	return &domain.DependencyCheck{CanDelete: true}, nil
}

func (m *MockProductRepository) DecreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB, ingredientName string, currentStock float64) error {
	if ingredient, exists := m.ingredients[ingredientID]; exists {
		ingredient.StockQuantity -= qty
	}
	return nil
}

func (m *MockProductRepository) IncreaseIngredientStock(ctx context.Context, ingredientID uint, qty float64, txDB *gorm.DB) error {
	if ingredient, exists := m.ingredients[ingredientID]; exists {
		ingredient.StockQuantity += qty
	}
	return nil
}

func TestProductService_CreateProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	input := CreateProductInput{
		Name:  "Test Product",
		Price: domain.FromFloat64(10.99),
	}

	product, err := svc.CreateProduct(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}
	if product.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", product.Name)
	}
	if product.Price != domain.FromFloat64(10.99) {
		t.Errorf("expected price 10.99, got %f", product.Price.ToFloat64())
	}
	if !product.Active {
		t.Error("expected product to be active by default")
	}
	if product.Slug == "" {
		t.Error("expected slug to be generated")
	}
}

func TestProductService_CreateProduct_WithSlug(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	input := CreateProductInput{
		Name:  "Test Product",
		Slug:  "custom-slug",
		Price: domain.FromFloat64(10.99),
	}

	product, err := svc.CreateProduct(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}
	if product.Slug != "custom-slug" {
		t.Errorf("expected slug 'custom-slug', got '%s'", product.Slug)
	}
}

func TestProductService_ListProducts(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	// Create some products
	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Product 1", Active: true}
	mockRepo.products[2] = &domain.Product{ID: 2, Name: "Product 2", Active: true}

	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatalf("ListProducts failed: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestProductService_ListActiveProducts(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	// Create active and inactive products
	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Active Product", Active: true}
	mockRepo.products[2] = &domain.Product{ID: 2, Name: "Inactive Product", Active: false}

	products, err := svc.ListActiveProducts(context.Background())
	if err != nil {
		t.Fatalf("ListActiveProducts failed: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 active product, got %d", len(products))
	}
	if products[0].Name != "Active Product" {
		t.Errorf("expected 'Active Product', got '%s'", products[0].Name)
	}
}

func TestProductService_GetProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Test Product", Active: true}

	product, err := svc.GetProduct(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetProduct failed: %v", err)
	}
	if product.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", product.Name)
	}
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	_, err := svc.GetProduct(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) && err.Error() != ErrProductNotFound.Error() {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Test Product", Active: true}

	err := svc.DeleteProduct(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}

	if _, exists := mockRepo.products[1]; exists {
		t.Error("expected product to be deleted")
	}
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	err := svc.DeleteProduct(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) && err.Error() != ErrProductNotFound.Error() {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Old Name", Price: domain.FromFloat64(9.99), Active: true}

	input := UpdateProductInput{
		Name:  "New Name",
		Price: domain.FromFloat64(19.99),
	}

	product, err := svc.UpdateProduct(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}
	if product.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", product.Name)
	}
	if product.Price != domain.FromFloat64(19.99) {
		t.Errorf("expected price 19.99, got %f", product.Price.ToFloat64())
	}
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	input := UpdateProductInput{
		Name:  "New Name",
		Price: domain.FromFloat64(19.99),
	}

	_, err := svc.UpdateProduct(context.Background(), 999, input)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) && err.Error() != ErrProductNotFound.Error() {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_CreateIngredient(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	input := CreateIngredientInput{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 100.0,
		MinStock:      10.0,
	}

	ingredient, err := svc.CreateIngredient(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}
	if ingredient.Name != "Flour" {
		t.Errorf("expected name 'Flour', got '%s'", ingredient.Name)
	}
	if ingredient.Unit != "kg" {
		t.Errorf("expected unit 'kg', got '%s'", ingredient.Unit)
	}
}

func TestProductService_ListIngredients(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg"}
	mockRepo.ingredients[2] = &domain.Ingredient{ID: 2, Name: "Sugar", Unit: "kg"}

	ingredients, err := svc.ListIngredients(context.Background())
	if err != nil {
		t.Fatalf("ListIngredients failed: %v", err)
	}
	if len(ingredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(ingredients))
	}
}

func TestProductService_GetIngredient(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg"}

	ingredient, err := svc.GetIngredient(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIngredient failed: %v", err)
	}
	if ingredient.Name != "Flour" {
		t.Errorf("expected name 'Flour', got '%s'", ingredient.Name)
	}
}

func TestProductService_GetIngredient_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	_, err := svc.GetIngredient(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent ingredient")
	}
	if !errors.Is(err, ErrIngredientNotFound) && err.Error() != ErrIngredientNotFound.Error() {
		t.Errorf("expected ErrIngredientNotFound, got: %v", err)
	}
}

func TestProductService_UpdateIngredientStock(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg", StockQuantity: 50.0}

	input := UpdateStockInput{
		StockQuantity: 100.0,
	}

	ingredient, err := svc.UpdateIngredientStock(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateIngredientStock failed: %v", err)
	}
	if ingredient.StockQuantity != 100.0 {
		t.Errorf("expected stock 100.0, got %f", ingredient.StockQuantity)
	}
}

func TestProductService_UpdateIngredientStock_QuantityField(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg", StockQuantity: 50.0}

	input := UpdateStockInput{
		Quantity: 75.0,
	}

	ingredient, err := svc.UpdateIngredientStock(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateIngredientStock failed: %v", err)
	}
	if ingredient.StockQuantity != 75.0 {
		t.Errorf("expected stock 75.0, got %f", ingredient.StockQuantity)
	}
}

func TestProductService_UpdateIngredientStock_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	input := UpdateStockInput{
		StockQuantity: 100.0,
	}

	_, err := svc.UpdateIngredientStock(context.Background(), 999, input)
	if err == nil {
		t.Error("expected error for non-existent ingredient")
	}
	if !errors.Is(err, ErrIngredientNotFound) && err.Error() != ErrIngredientNotFound.Error() {
		t.Errorf("expected ErrIngredientNotFound, got: %v", err)
	}
}

func TestProductService_UpdateIngredient(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg", StockQuantity: 50.0, MinStock: 10.0}

	input := UpdateIngredientInput{
		Name:          "Sugar",
		Unit:          "g",
		StockQuantity: 200.0,
		MinStock:      20.0,
	}

	ingredient, err := svc.UpdateIngredient(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateIngredient failed: %v", err)
	}
	if ingredient.Name != "Sugar" {
		t.Errorf("expected name 'Sugar', got '%s'", ingredient.Name)
	}
	if ingredient.Unit != "g" {
		t.Errorf("expected unit 'g', got '%s'", ingredient.Unit)
	}
}

func TestProductService_DeleteIngredient(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.ingredients[1] = &domain.Ingredient{ID: 1, Name: "Flour", Unit: "kg"}

	err := svc.DeleteIngredient(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteIngredient failed: %v", err)
	}

	if _, exists := mockRepo.ingredients[1]; exists {
		t.Error("expected ingredient to be deleted")
	}
}

func TestProductService_SetProductIngredients(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Cake", Active: true}
	mockRepo.ingredients[2] = &domain.Ingredient{ID: 2, Name: "Flour", Unit: "kg"}
	mockRepo.ingredients[3] = &domain.Ingredient{ID: 3, Name: "Sugar", Unit: "kg"}

	input := SetProductIngredientsInput{
		Items: []ProductIngredientInput{
			{IngredientID: 2, Quantity: 0.5},
			{IngredientID: 3, Quantity: 0.2},
		},
	}

	err := svc.SetProductIngredients(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("SetProductIngredients failed: %v", err)
	}

	ingredients := mockRepo.productIngredients[1]
	if len(ingredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(ingredients))
	}
}

func TestProductService_SetProductIngredients_IngredientNotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Cake", Active: true}

	input := SetProductIngredientsInput{
		Items: []ProductIngredientInput{
			{IngredientID: 999, Quantity: 0.5},
		},
	}

	err := svc.SetProductIngredients(context.Background(), 1, input)
	if err == nil {
		t.Error("expected error for non-existent ingredient")
	}
}

func TestProductService_DuplicateProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	original := &domain.Product{
		ID:          1,
		Name:        "Original Product",
		Description: "Original description",
		Price:       domain.FromFloat64(10.99),
		Active:      true,
		PhotoURL:    "http://example.com/photo.jpg",
		CategoryID:  nil,
	}
	mockRepo.products[1] = original

	duplicate, err := svc.DuplicateProduct(context.Background(), 1)
	if err != nil {
		t.Fatalf("DuplicateProduct failed: %v", err)
	}
	if duplicate.Name != "Original Product (Cópia)" {
		t.Errorf("expected name 'Original Product (Cópia)', got '%s'", duplicate.Name)
	}
	if !duplicate.Active {
		t.Error("expected duplicate to be active")
	}
	if !duplicate.IsNew {
		t.Error("expected duplicate to be marked as new")
	}
	if duplicate.Featured {
		t.Error("expected duplicate not to be featured")
	}
}

func TestProductService_DuplicateProduct_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	_, err := svc.DuplicateProduct(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) && err.Error() != ErrProductNotFound.Error() {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_ArchiveProduct(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Test Product", Active: true}

	err := svc.ArchiveProduct(context.Background(), 1)
	if err != nil {
		t.Fatalf("ArchiveProduct failed: %v", err)
	}

	if mockRepo.products[1].Active {
		t.Error("expected product to be archived (inactive)")
	}
}

func TestProductService_ArchiveProduct_NotFound(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	err := svc.ArchiveProduct(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent product")
	}
	if !errors.Is(err, ErrProductNotFound) && err.Error() != ErrProductNotFound.Error() {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}

func TestProductService_UpdateProduct_WithPromotion(t *testing.T) {
	mockRepo := NewMockProductRepository()
	svc := NewProductService(mockRepo)

	now := time.Now()
	mockRepo.products[1] = &domain.Product{ID: 1, Name: "Test Product", Price: domain.FromFloat64(20.0), Active: true}

	promotionPrice := domain.FromFloat64(15.0)
	input := UpdateProductInput{
		Name:           "Test Product",
		Price:          domain.FromFloat64(20.0),
		PromotionPrice: &promotionPrice,
		PromotionStart: &now,
	}

	product, err := svc.UpdateProduct(context.Background(), 1, input)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}
	if product.PromotionPrice == nil {
		t.Error("expected promotion price to be set")
	}
	if *product.PromotionPrice != domain.FromFloat64(15.0) {
		t.Errorf("expected promotion price 15.0, got %f", product.PromotionPrice.ToFloat64())
	}
}
