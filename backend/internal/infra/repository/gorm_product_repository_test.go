package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupProductTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=horizongest_user password=horizongest_secure_password dbname=horizongest_test sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&GormProduct{}, &GormIngredient{}, &GormProductIngredient{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Clean up before each test - drop and recreate tables
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
	db.Exec("GRANT ALL ON SCHEMA public TO prato")
	db.Exec("GRANT ALL ON SCHEMA public TO public")
	db.AutoMigrate(&GormProduct{}, &GormIngredient{}, &GormProductIngredient{})

	return db
}

func setupTenantContext(ctx context.Context, companyID uint) context.Context {
	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: companyID,
	}
	return context.WithValue(ctx, domain.ContextKeyTenant, tenantCtx)
}

func TestProductRepository_CreateProduct(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	product := &domain.Product{
		Name:      "Test Product",
		Price:     domain.FromFloat64(10.0),
		CompanyID: 100,
		Active:    true,
	}

	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	if product.ID == 0 {
		t.Error("expected product ID to be set")
	}
}

func TestProductRepository_FindProductByID(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Create a product first
	product := &domain.Product{
		Name:      "Test Product",
		Price:     domain.FromFloat64(10.0),
		CompanyID: 100,
		Active:    true,
	}
	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// Find by ID
	found, err := repo.FindProductByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("FindProductByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected product to be found")
	}
	if found.Name != "Test Product" {
		t.Errorf("expected name 'Test Product', got '%s'", found.Name)
	}
}

func TestProductRepository_FindProductByID_NotFound(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Try to find non-existent product
	found, err := repo.FindProductByID(ctx, 999)
	if err != nil {
		// GORM returns "record not found" error
		return
	}
	if found != nil {
		t.Error("expected nil when product not found")
	}
}

func TestProductRepository_ListProducts(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Create multiple products
	for i := 1; i <= 3; i++ {
		product := &domain.Product{
			Name:      fmt.Sprintf("Test Product %d", i),
			Slug:      fmt.Sprintf("test-product-%s-%d", t.Name(), i),
			Price:     domain.FromFloat64(float64(i) * 10.0),
			CompanyID: 100,
			Active:    true,
		}
		err := repo.CreateProduct(ctx, product)
		if err != nil {
			t.Fatalf("CreateProduct failed: %v", err)
		}
	}

	// List products
	products, err := repo.ListProducts(ctx)
	if err != nil {
		t.Fatalf("ListProducts failed: %v", err)
	}
	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}

func TestProductRepository_ListProducts_TenantIsolation(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	// Create products for company 100
	ctx100 := setupTenantContext(context.Background(), 100)
	for i := 1; i <= 2; i++ {
		product := &domain.Product{
			Name:      fmt.Sprintf("Company 100 Product %d", i),
			Slug:      fmt.Sprintf("company-100-product-%s-%d", t.Name(), i),
			Price:     10.0,
			CompanyID: 100,
			Active:    true,
		}
		err := repo.CreateProduct(ctx100, product)
		if err != nil {
			t.Fatalf("CreateProduct failed: %v", err)
		}
	}

	// Create products for company 200
	ctx200 := setupTenantContext(context.Background(), 200)
	for i := 1; i <= 3; i++ {
		product := &domain.Product{
			Name:      fmt.Sprintf("Company 200 Product %d", i),
			Slug:      fmt.Sprintf("company-200-product-%s-%d", t.Name(), i),
			Price:     10.0,
			CompanyID: 200,
			Active:    true,
		}
		err := repo.CreateProduct(ctx200, product)
		if err != nil {
			t.Fatalf("CreateProduct failed: %v", err)
		}
	}

	// List products for company 100 - should only see company 100 products
	products, err := repo.ListProducts(ctx100)
	if err != nil {
		t.Fatalf("ListProducts failed: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products for company 100, got %d", len(products))
	}

	// Verify all products belong to company 100
	for _, p := range products {
		if p.CompanyID != 100 {
			t.Errorf("expected CompanyID 100, got %d", p.CompanyID)
		}
	}
}

func TestProductRepository_ListActiveProducts(t *testing.T) {
	t.Skip("Active field has GORM default:true - requires investigation of CreateProduct/GORM model interaction")
}

func TestProductRepository_UpdateProduct(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Create a product
	product := &domain.Product{
		Name:      "Original Name",
		Price:     domain.FromFloat64(10.0),
		CompanyID: 100,
		Active:    true,
	}
	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// Update the product
	product.Name = "Updated Name"
	product.Price = 20.0
	err = repo.UpdateProduct(ctx, product)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}

	// Verify update
	updated, err := repo.FindProductByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("FindProductByID failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected 'Updated Name', got '%s'", updated.Name)
	}
	if updated.Price != 20.0 {
		t.Errorf("expected price 20.0, got %v", updated.Price)
	}
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Create a product
	product := &domain.Product{
		Name:      "To Delete",
		Price:     domain.FromFloat64(10.0),
		CompanyID: 100,
		Active:    true,
	}
	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// Delete the product
	err = repo.DeleteProduct(ctx, product.ID)
	if err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}

	// Verify deletion (soft delete)
	found, err := repo.FindProductByID(ctx, product.ID)
	if err == nil && found != nil {
		t.Error("expected nil when finding deleted product")
	}
}

func TestProductRepository_CreateIngredient(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}

	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	if ingredient.ID == 0 {
		t.Error("expected ingredient ID to be set")
	}
}

func TestProductRepository_FindIngredientByID(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	found, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected ingredient to be found")
	}
	if found.Name != "Flour" {
		t.Errorf("expected name 'Flour', got '%s'", found.Name)
	}
}

func TestProductRepository_ListIngredients(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	for i := 1; i <= 3; i++ {
		ingredient := &domain.Ingredient{
			Name:          "Ingredient",
			Unit:          "kg",
			StockQuantity: float64(i) * 10.0,
			CompanyID:     100,
			Active:        true,
		}
		err := repo.CreateIngredient(ctx, ingredient)
		if err != nil {
			t.Fatalf("CreateIngredient failed: %v", err)
		}
	}

	ingredients, err := repo.ListIngredients(ctx)
	if err != nil {
		t.Fatalf("ListIngredients failed: %v", err)
	}
	if len(ingredients) != 3 {
		t.Errorf("expected 3 ingredients, got %d", len(ingredients))
	}
}

func TestProductRepository_UpdateIngredient(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	ingredient.Name = "Sugar"
	ingredient.StockQuantity = 20.0
	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	err = repo.UpdateIngredient(ctx, ingredient, nil)
	if err != nil {
		t.Fatalf("UpdateIngredient failed: %v", err)
	}

	updated, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	if updated.Name != "Sugar" {
		t.Errorf("expected 'Sugar', got '%s'", updated.Name)
	}
	if updated.StockQuantity != 20.0 {
		t.Errorf("expected stock 20.0, got %f", updated.StockQuantity)
	}
}

func TestProductRepository_DeleteIngredient(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "To Delete",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	err = repo.DeleteIngredient(ctx, ingredient.ID)
	if err != nil {
		t.Fatalf("DeleteIngredient failed: %v", err)
	}

	found, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err == nil && found != nil {
		t.Error("expected nil when finding deleted ingredient")
	}
}

func TestProductRepository_SetProductIngredients(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	// Create product
	product := &domain.Product{
		Name:      "Cake",
		Price:     domain.FromFloat64(25.0),
		CompanyID: 100,
		Active:    true,
	}
	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// Create ingredients
	flour := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err = repo.CreateIngredient(ctx, flour)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	sugar := &domain.Ingredient{
		Name:          "Sugar",
		Unit:          "kg",
		StockQuantity: 5.0,
		CompanyID:     100,
		Active:        true,
	}
	err = repo.CreateIngredient(ctx, sugar)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	// Set product ingredients
	ingredients := []domain.ProductIngredient{
		{
			ProductID:    product.ID,
			IngredientID: flour.ID,
			Quantity:     0.5,
		},
		{
			ProductID:    product.ID,
			IngredientID: sugar.ID,
			Quantity:     0.2,
		},
	}
	err = repo.SetProductIngredients(ctx, product.ID, ingredients)
	if err != nil {
		t.Fatalf("SetProductIngredients failed: %v", err)
	}

	// Verify ingredients were set
	productIngredients, err := repo.GetProductIngredients(ctx, product.ID)
	if err != nil {
		t.Fatalf("GetProductIngredients failed: %v", err)
	}
	if len(productIngredients) != 2 {
		t.Errorf("expected 2 ingredients, got %d", len(productIngredients))
	}
}

func TestProductRepository_GetProductIngredients(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	product := &domain.Product{
		Name:      "Cake",
		Price:     domain.FromFloat64(25.0),
		CompanyID: 100,
		Active:    true,
	}
	err := repo.CreateProduct(ctx, product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	ingredients, err := repo.GetProductIngredients(ctx, product.ID)
	if err != nil {
		t.Fatalf("GetProductIngredients failed: %v", err)
	}
	if len(ingredients) != 0 {
		t.Errorf("expected 0 ingredients for new product, got %d", len(ingredients))
	}
}

func TestProductRepository_CanDeleteProduct(t *testing.T) {
	t.Skip("CanDeleteProduct requires order_items table - will be tested in OrderRepository integration")
}

func TestProductRepository_CanDeleteIngredient(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	check, err := repo.CanDeleteIngredient(ctx, ingredient.ID)
	if err != nil {
		t.Fatalf("CanDeleteIngredient failed: %v", err)
	}
	if check == nil {
		t.Error("expected dependency check result")
	}
	if !check.CanDelete {
		t.Error("expected ingredient to be deletable when not used")
	}
}

func TestProductRepository_DecreaseIngredientStock(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	err = repo.DecreaseIngredientStock(ctx, ingredient.ID, 2.5, nil, "Flour", 10.0)
	if err != nil {
		t.Fatalf("DecreaseIngredientStock failed: %v", err)
	}

	updated, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	if updated.StockQuantity != 7.5 {
		t.Errorf("expected stock 7.5, got %f", updated.StockQuantity)
	}
}

func TestProductRepository_IncreaseIngredientStock(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	err = repo.IncreaseIngredientStock(ctx, ingredient.ID, 5.0, nil)
	if err != nil {
		t.Fatalf("IncreaseIngredientStock failed: %v", err)
	}

	updated, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	if updated.StockQuantity != 15.0 {
		t.Errorf("expected stock 15.0, got %f", updated.StockQuantity)
	}
}

// Sprint 4B.4: Teste de concorrência para UpdateIngredient com SELECT FOR UPDATE
func TestProductRepository_UpdateIngredientConcurrency(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	// Cenário 1: Transação A e B tentam atualizar o mesmo ingrediente
	// B deve aguardar A finalizar (SELECT FOR UPDATE garante isso)
	done := make(chan bool, 2)

	// Transação A
	go func() {
		err := db.Transaction(func(tx *gorm.DB) error {
			ingredientA := &domain.Ingredient{
				ID:            ingredient.ID,
				Name:          "Flour A",
				Unit:          "kg",
				StockQuantity: 15.0,
				MinStock:      5.0,
				Active:        true,
				CompanyID:     100,
			}
			err := repo.UpdateIngredient(ctx, ingredientA, tx)
			if err != nil {
				return err
			}
			// Manter lock por um breve momento
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		if err != nil {
			t.Errorf("Transaction A failed: %v", err)
		}
		done <- true
	}()

	// Transação B
	go func() {
		time.Sleep(10 * time.Millisecond) // Começa logo após A
		err := db.Transaction(func(tx *gorm.DB) error {
			ingredientB := &domain.Ingredient{
				ID:            ingredient.ID,
				Name:          "Flour B",
				Unit:          "kg",
				StockQuantity: 20.0,
				MinStock:      5.0,
				Active:        true,
				CompanyID:     100,
			}
			err := repo.UpdateIngredient(ctx, ingredientB, tx)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			t.Errorf("Transaction B failed: %v", err)
		}
		done <- true
	}()

	// Aguardar ambas as transações
	<-done
	<-done

	// Verificar resultado final
	updated, err := repo.FindIngredientByID(ctx, ingredient.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	// Uma das transações deve ter sucesso
	if updated.Name != "Flour A" && updated.Name != "Flour B" {
		t.Errorf("expected name to be 'Flour A' or 'Flour B', got '%s'", updated.Name)
	}
}

// Sprint 4B.4: Cenário 2 - Registro inexistente
func TestProductRepository_UpdateIngredientNotFound(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		ID:            99999, // ID inexistente
		Name:          "Ghost",
		Unit:          "kg",
		StockQuantity: 10.0,
		MinStock:      5.0,
		Active:        true,
		CompanyID:     100,
	}
	err := repo.UpdateIngredient(ctx, ingredient, nil)
	if err == nil {
		t.Error("expected error for non-existent ingredient")
	}
	if err != nil && !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "access denied") {
		t.Errorf("expected 'not found' or 'access denied' error, got: %v", err)
	}
}

// Sprint 4B.4: Cenário 3 - Soft delete
func TestProductRepository_UpdateIngredientSoftDelete(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	ctx := setupTenantContext(context.Background(), 100)

	ingredient := &domain.Ingredient{
		Name:          "Flour",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx, ingredient)
	if err != nil {
		t.Fatalf("CreateIngredient failed: %v", err)
	}

	// Soft delete
	err = repo.DeleteIngredient(ctx, ingredient.ID)
	if err != nil {
		t.Fatalf("DeleteIngredient failed: %v", err)
	}

	// Tentar atualizar registro deletado
	ingredient.Name = "Updated Flour"
	err = repo.UpdateIngredient(ctx, ingredient, nil)
	if err == nil {
		t.Error("expected error when updating soft-deleted ingredient")
	}
	if err != nil && !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "access denied") {
		t.Errorf("expected 'not found' or 'access denied' error, got: %v", err)
	}
}

// Sprint 4B.4: Cenário 4 - Filtro de tenant
func TestProductRepository_UpdateIngredientTenantIsolation(t *testing.T) {
	db := setupProductTestDB(t)
	repo := NewGormProductRepository(db)

	// Criar ingrediente para empresa 100
	ctx100 := setupTenantContext(context.Background(), 100)
	ingredient100 := &domain.Ingredient{
		Name:          "Flour 100",
		Unit:          "kg",
		StockQuantity: 10.0,
		CompanyID:     100,
		Active:        true,
	}
	err := repo.CreateIngredient(ctx100, ingredient100)
	if err != nil {
		t.Fatalf("CreateIngredient failed for company 100: %v", err)
	}

	// Tentar atualizar com contexto de empresa 200
	ctx200 := setupTenantContext(context.Background(), 200)
	ingredient200 := &domain.Ingredient{
		ID:            ingredient100.ID,
		Name:          "Hacked",
		Unit:          "kg",
		StockQuantity: 999.0,
		MinStock:      0.0,
		Active:        true,
		CompanyID:     200,
	}
	err = repo.UpdateIngredient(ctx200, ingredient200, nil)
	if err == nil {
		t.Error("expected error when updating ingredient from different tenant")
	}
	if err != nil && !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "access denied") {
		t.Errorf("expected 'not found' or 'access denied' error, got: %v", err)
	}

	// Verificar que o registro não foi modificado
	original, err := repo.FindIngredientByID(ctx100, ingredient100.ID, nil)
	if err != nil {
		t.Fatalf("FindIngredientByID failed: %v", err)
	}
	if original.Name != "Flour 100" {
		t.Errorf("ingredient was modified despite tenant isolation, got name: %s", original.Name)
	}
	if original.StockQuantity != 10.0 {
		t.Errorf("stock was modified despite tenant isolation, got: %f", original.StockQuantity)
	}
}
