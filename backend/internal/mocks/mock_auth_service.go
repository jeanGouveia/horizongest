package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

// MockAuthService is a mock implementation of service.AuthServiceInterface
type MockAuthService struct {
	ValidateTokenResult *service.JWTClaims
	ValidateTokenError   error
	LogoutError         error
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) Login(ctx context.Context, input service.LoginInput) (*service.LoginResult, error) {
	return &service.LoginResult{}, nil
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	return m.LogoutError
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*service.JWTClaims, error) {
	return m.ValidateTokenResult, m.ValidateTokenError
}

func (m *MockAuthService) UpdateProfile(ctx context.Context, userID uint, input service.UpdateProfileInput) (*domain.User, error) {
	return &domain.User{}, nil
}

func (m *MockAuthService) ChangePassword(ctx context.Context, userID uint, input service.ChangePasswordInput) error {
	return nil
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, input service.RequestPasswordResetInput) error {
	return nil
}

func (m *MockAuthService) ResetPassword(ctx context.Context, input service.ResetPasswordInput) error {
	return nil
}
