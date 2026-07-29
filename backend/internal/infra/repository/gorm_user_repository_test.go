package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
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
	err = db.AutoMigrate(&GormUserModel{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Clean up before each test - drop and recreate schema
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
	db.Exec("GRANT ALL ON SCHEMA public TO prato")
	db.Exec("GRANT ALL ON SCHEMA public TO public")
	db.AutoMigrate(&GormUserModel{})

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user := &domain.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleOwner,
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user ID to be set")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user := &domain.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleOwner,
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected user to be found")
	}
	if found.Name != "Test User" {
		t.Errorf("expected name 'Test User', got '%s'", found.Name)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	found, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if found != nil {
		t.Error("expected nil when user not found")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user := &domain.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleOwner,
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected user to be found")
	}
	if found.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", found.Email)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	found, err := repo.FindByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected nil when user not found")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user := &domain.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleOwner,
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	user.Name = "Updated Name"
	user.Active = false
	err = repo.Update(context.Background(), user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.FindByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
	if updated.Active {
		t.Error("expected user to be inactive")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	for i := 1; i <= 3; i++ {
		user := &domain.User{
			Name:         "Test User",
			Email:        fmt.Sprintf("test%d@example.com", i),
			PasswordHash: "hashedpassword",
			Active:       true,
			CompanyID:    100,
			Role:         domain.RoleOwner,
		}
		err := repo.Create(context.Background(), user)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	users, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user1 := &domain.User{
		Name:         "Test User 1",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleOwner,
	}
	err := repo.Create(context.Background(), user1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	user2 := &domain.User{
		Name:         "Test User 2",
		Email:        "test@example.com", // Same email
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleAdmin,
	}
	err = repo.Create(context.Background(), user2)
	if err == nil {
		t.Error("expected error when creating user with duplicate email")
	}
}

func TestUserRepository_RoleParsing(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewGormUserRepository(db)

	user := &domain.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Active:       true,
		CompanyID:    100,
		Role:         domain.RoleManager,
	}
	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.FindByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Role != domain.RoleManager {
		t.Errorf("expected role Manager, got %v", found.Role)
	}
}
