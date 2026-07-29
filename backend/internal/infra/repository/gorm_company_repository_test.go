package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCompanyTestDB(t *testing.T) *gorm.DB {
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
	err = db.AutoMigrate(&GormCompanyModel{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Clean up before each test - drop and recreate schema
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
	db.Exec("GRANT ALL ON SCHEMA public TO prato")
	db.Exec("GRANT ALL ON SCHEMA public TO public")
	db.AutoMigrate(&GormCompanyModel{})

	return db
}

func TestCompanyRepository_Create(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-create",
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}

	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if company.ID == 0 {
		t.Error("expected company ID to be set")
	}
}

func TestCompanyRepository_FindByID(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected company to be found")
	}
	if found.Name != "Test Company" {
		t.Errorf("expected name 'Test Company', got '%s'", found.Name)
	}
}

func TestCompanyRepository_FindByID_NotFound(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	found, err := repo.FindByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected nil when company not found")
	}
}

func TestCompanyRepository_FindBySlug(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindBySlug(context.Background(), company.Slug)
	if err != nil {
		t.Fatalf("FindBySlug failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected company to be found")
	}
	if found.Name != "Test Company" {
		t.Errorf("expected name 'Test Company', got '%s'", found.Name)
	}
}

func TestCompanyRepository_FindBySlug_NotFound(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	found, err := repo.FindBySlug(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("FindBySlug failed: %v", err)
	}
	if found != nil {
		t.Error("expected nil when company not found")
	}
}

func TestCompanyRepository_List(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	for i := 1; i <= 3; i++ {
		company := &domain.Company{
			Name:         fmt.Sprintf("Test Company %d", i),
			Slug:         fmt.Sprintf("test-company-%d", i),
			Description:  "A test company",
			Active:       true,
			BusinessType: domain.BusinessTypeRestaurant,
			Locale:       "pt-BR",
			Currency:     "BRL",
			Timezone:     "America/Sao_Paulo",
		}
		err := repo.Create(context.Background(), company)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	companies, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(companies) != 3 {
		t.Errorf("expected 3 companies, got %d", len(companies))
	}
}

func TestCompanyRepository_Update(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	company.Name = "Updated Company"
	company.Description = "Updated description"
	company.Active = false
	err = repo.Update(context.Background(), company)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.FindByID(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Name != "Updated Company" {
		t.Errorf("expected name 'Updated Company', got '%s'", updated.Name)
	}
	if updated.Active {
		t.Error("expected company to be inactive")
	}
}

func TestCompanyRepository_Delete(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify soft delete
	found, err := repo.FindByID(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected nil after soft delete")
	}
}

func TestCompanyRepository_Create_DuplicateSlug(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company1 := &domain.Company{
		Name:         "Test Company 1",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	company2 := &domain.Company{
		Name:         "Test Company 2",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())), // Same slug
		Description:  "Another test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err = repo.Create(context.Background(), company2)
	if err == nil {
		t.Error("expected error when creating company with duplicate slug")
	}
}

func TestCompanyRepository_BusinessTypeParsing(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeBakery,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.BusinessType != domain.BusinessTypeBakery {
		t.Errorf("expected BusinessType Bakery, got %v", found.BusinessType)
	}
}

func TestCompanyRepository_Timestamps(t *testing.T) {
	db := setupCompanyTestDB(t)
	repo := NewGormCompanyRepository(db)

	company := &domain.Company{
		Name:         "Test Company",
		Slug:         "test-company-" + t.Name() + "-" + strconv.Itoa(int(time.Now().Unix())),
		Description:  "A test company",
		Active:       true,
		BusinessType: domain.BusinessTypeRestaurant,
		Locale:       "pt-BR",
		Currency:     "BRL",
		Timezone:     "America/Sao_Paulo",
	}
	err := repo.Create(context.Background(), company)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if company.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if company.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	// Wait a moment to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	company.Name = "Updated Company"
	err = repo.Update(context.Background(), company)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.FindByID(context.Background(), company.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.UpdatedAt.Equal(company.CreatedAt) {
		t.Error("expected UpdatedAt to be different from CreatedAt after update")
	}
}
