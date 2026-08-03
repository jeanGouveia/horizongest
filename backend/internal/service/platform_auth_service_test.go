package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// MockPlatformUserRepository for testing
type MockPlatformUserRepository struct {
	users       map[uint]*domain.PlatformUser
	findByEmail *domain.PlatformUser
	findByID    *domain.PlatformUser
	createError error
	updateError error
}

func NewMockPlatformUserRepository() *MockPlatformUserRepository {
	return &MockPlatformUserRepository{
		users: make(map[uint]*domain.PlatformUser),
	}
}

func (m *MockPlatformUserRepository) Create(ctx context.Context, user *domain.PlatformUser) error {
	if m.createError != nil {
		return m.createError
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *MockPlatformUserRepository) FindByEmail(ctx context.Context, email string) (*domain.PlatformUser, error) {
	if m.findByEmail != nil {
		return m.findByEmail, nil
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockPlatformUserRepository) FindByID(ctx context.Context, id uint) (*domain.PlatformUser, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return m.users[id], nil
}

func (m *MockPlatformUserRepository) Update(ctx context.Context, user *domain.PlatformUser) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.users[user.ID] = user
	return nil
}

// MockPlatformSessionRepository for testing
type MockPlatformSessionRepository struct {
	sessions    map[string]*domain.PlatformSession
	createError error
	findError   error
	deleteError error
}

func NewMockPlatformSessionRepository() *MockPlatformSessionRepository {
	return &MockPlatformSessionRepository{
		sessions: make(map[string]*domain.PlatformSession),
	}
}

func (m *MockPlatformSessionRepository) Create(ctx context.Context, session *domain.PlatformSession) error {
	if m.createError != nil {
		return m.createError
	}
	m.sessions[session.Token] = session
	return nil
}

func (m *MockPlatformSessionRepository) FindByToken(ctx context.Context, token string) (*domain.PlatformSession, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	return m.sessions[token], nil
}

func (m *MockPlatformSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	delete(m.sessions, token)
	return nil
}

func (m *MockPlatformSessionRepository) DeleteByPlatformUserID(ctx context.Context, platformUserID uint) error {
	for token, session := range m.sessions {
		if session.PlatformUserID == platformUserID {
			delete(m.sessions, token)
		}
	}
	return nil
}

func TestPlatformAuthService_Login_Success(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	user := &domain.PlatformUser{
		ID:           1,
		Name:         "Admin User",
		Email:        "admin@platform.com",
		PasswordHash: "$2a$10$test", // Mock hash
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}
	userRepo.users[1] = user

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	// Mock password verification by setting a known hash
	user.PasswordHash = hashPassword("password123")

	result, err := svc.Login(context.Background(), PlatformLoginInput{
		Email:    "admin@platform.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if result == nil {
		t.Fatal("Login returned nil result")
	}
	if result.Token == "" {
		t.Error("Login returned empty token")
	}
	if result.User == nil {
		t.Error("Login returned nil user")
	}
	if result.User.Email != "admin@platform.com" {
		t.Errorf("expected email admin@platform.com, got %s", result.User.Email)
	}
}

func TestPlatformAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	user := &domain.PlatformUser{
		ID:           1,
		Name:         "Admin User",
		Email:        "admin@platform.com",
		PasswordHash: hashPassword("password123"),
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}
	userRepo.users[1] = user

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, err := svc.Login(context.Background(), PlatformLoginInput{
		Email:    "admin@platform.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
	if !errors.Is(err, ErrPlatformInvalidCredentials) && err.Error() != ErrPlatformInvalidCredentials.Error() {
		t.Errorf("expected ErrPlatformInvalidCredentials, got: %v", err)
	}
}

func TestPlatformAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, err := svc.Login(context.Background(), PlatformLoginInput{
		Email:    "nonexistent@platform.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
	if !errors.Is(err, ErrPlatformInvalidCredentials) && err.Error() != ErrPlatformInvalidCredentials.Error() {
		t.Errorf("expected ErrPlatformInvalidCredentials, got: %v", err)
	}
}

func TestPlatformAuthService_Login_InactiveUser(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	user := &domain.PlatformUser{
		ID:           1,
		Name:         "Admin User",
		Email:        "admin@platform.com",
		PasswordHash: hashPassword("password123"),
		Role:         domain.PlatformRoleAdmin,
		Active:       false,
	}
	userRepo.users[1] = user

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, err := svc.Login(context.Background(), PlatformLoginInput{
		Email:    "admin@platform.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("expected error for inactive user, got nil")
	}
	if !errors.Is(err, ErrPlatformUserInactive) && err.Error() != ErrPlatformUserInactive.Error() {
		t.Errorf("expected ErrPlatformUserInactive, got: %v", err)
	}
}

func TestPlatformAuthService_Logout(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	token := "test-token"
	session := &domain.PlatformSession{
		PlatformUserID: 1,
		Token:          token,
		ExpiresAt:      time.Now().Add(time.Hour),
	}
	sessionRepo.sessions[token] = session

	err := svc.Logout(context.Background(), token)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Verify session was deleted
	if _, exists := sessionRepo.sessions[token]; exists {
		t.Error("session should be deleted after logout")
	}
}

func TestPlatformAuthService_ValidateToken(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	// Generate a valid token
	user := &domain.PlatformUser{
		ID:   1,
		Role: domain.PlatformRoleAdmin,
	}
	token, sessionID, err := svc.generateToken(user.ID, user.Role)
	if err != nil {
		t.Fatalf("generateToken failed: %v", err)
	}

	// FASE A: Create a session in the mock repository (B4 - Platform Session Validation)
	session := &domain.PlatformSession{
		ID:             uint(sessionID),
		PlatformUserID: user.ID,
		ExpiresAt:      time.Now().Add(time.Hour),
		Token:          token,
	}
	if err := sessionRepo.Create(context.Background(), session); err != nil {
		t.Fatalf("sessionRepo.Create failed: %v", err)
	}

	userID, role, err := svc.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if userID != user.ID {
		t.Errorf("expected userID %d, got %d", user.ID, userID)
	}
	if role != user.Role {
		t.Errorf("expected role %v, got %v", user.Role, role)
	}
}

func TestPlatformAuthService_ValidateToken_Invalid(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, _, err := svc.ValidateToken(context.Background(), "invalid-token")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestPlatformAuthService_CreatePlatformUser(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	user, err := svc.CreatePlatformUser(context.Background(), "New User", "new@platform.com", "password123", domain.PlatformRoleSupport)
	if err != nil {
		t.Fatalf("CreatePlatformUser failed: %v", err)
	}
	if user == nil {
		t.Fatal("CreatePlatformUser returned nil user")
	}
	if user.Email != "new@platform.com" {
		t.Errorf("expected email new@platform.com, got %s", user.Email)
	}
	if user.Role != domain.PlatformRoleSupport {
		t.Errorf("expected role Support, got %v", user.Role)
	}
}

func TestPlatformAuthService_CreatePlatformUser_EmailExists(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	existingUser := &domain.PlatformUser{
		ID:    1,
		Email: "existing@platform.com",
	}
	userRepo.users[1] = existingUser
	userRepo.findByEmail = existingUser

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, err := svc.CreatePlatformUser(context.Background(), "New User", "existing@platform.com", "password123", domain.PlatformRoleSupport)
	if err == nil {
		t.Error("expected error for existing email, got nil")
	}
	if !errors.Is(err, ErrPlatformEmailAlreadyExists) && err.Error() != ErrPlatformEmailAlreadyExists.Error() {
		t.Errorf("expected ErrPlatformEmailAlreadyExists, got: %v", err)
	}
}

func TestPlatformAuthService_GetPlatformUser(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	user := &domain.PlatformUser{
		ID:    1,
		Name:  "Admin User",
		Email: "admin@platform.com",
		Role:  domain.PlatformRoleAdmin,
	}
	userRepo.users[1] = user
	userRepo.findByID = user

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	output, err := svc.GetPlatformUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPlatformUser failed: %v", err)
	}
	if output == nil {
		t.Fatal("GetPlatformUser returned nil output")
	}
	if output.Email != "admin@platform.com" {
		t.Errorf("expected email admin@platform.com, got %s", output.Email)
	}
}

func TestPlatformAuthService_GetPlatformUser_NotFound(t *testing.T) {
	userRepo := NewMockPlatformUserRepository()
	sessionRepo := NewMockPlatformSessionRepository()

	svc := NewPlatformAuthService(userRepo, sessionRepo, "test-secret", time.Hour, bcrypt.DefaultCost)

	_, err := svc.GetPlatformUser(context.Background(), 999)
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
	if !errors.Is(err, ErrPlatformUserNotFound) && err.Error() != ErrPlatformUserNotFound.Error() {
		t.Errorf("expected ErrPlatformUserNotFound, got: %v", err)
	}
}
