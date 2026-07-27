package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// Helper function to create a bcrypt hash
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// TestAuthService_Login_Success tests successful login flow
func TestAuthService_Login_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	result, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password",
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
	if result.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", result.User.Email)
	}

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(result.Token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, claims.Email)
	}
	if claims.CompanyID != user.CompanyID {
		t.Errorf("expected CompanyID %d, got %d", user.CompanyID, claims.CompanyID)
	}
	if claims.Issuer != "TestPlatform" {
		t.Errorf("expected Issuer TestPlatform, got %s", claims.Issuer)
	}
}

// TestAuthService_Login_InactiveUser tests that inactive users cannot login
func TestAuthService_Login_InactiveUser(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "inactive@example.com",
		Name:         "Inactive User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       false,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "inactive@example.com",
		Password: "password",
	})

	if err == nil {
		t.Error("expected error for inactive user, got nil")
	}
	if !errors.Is(err, errors.New("usuário desativado. Entre em contato com o administrador.")) && err.Error() != "usuário desativado. Entre em contato com o administrador." {
		t.Errorf("expected inactive user error, got: %v", err)
	}
}

// TestAuthService_Logout tests logout functionality
func TestAuthService_Logout(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	// First login to get a token
	result, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Logout
	err = svc.Logout(context.Background(), result.Token)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Verify token is blacklisted
	blacklisted, err := mockBlacklistRepo.IsBlacklisted(context.Background(), result.Token)
	if err != nil {
		t.Fatalf("IsBlacklisted failed: %v", err)
	}
	if !blacklisted {
		t.Error("token should be blacklisted after logout")
	}

	// Verify blacklisted token cannot be validated
	_, err = svc.ValidateToken(context.Background(), result.Token)
	if err == nil {
		t.Error("expected error when validating blacklisted token")
	}
}

// TestAuthService_ValidateToken tests token validation
func TestAuthService_ValidateToken(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	// Generate token
	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	// Verify token is valid
	claims, err := svc.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims == nil {
		t.Fatal("ValidateToken returned nil claims")
	}
	if claims.UserID != user.ID {
		t.Errorf("expected UserID %d, got %d", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, claims.Email)
	}
}

// TestAuthService_ValidateToken_InvalidSecret tests that tokens with wrong secret are rejected
func TestAuthService_ValidateToken_InvalidSecret(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	// Create a service with different secret
	wrongSecretSvc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "wrong-secret", "TestPlatform")

	_, err = wrongSecretSvc.ValidateToken(context.Background(), token)
	if err == nil {
		t.Fatal("expected error when validating token with wrong secret")
	}
}

// TestAuthService_Login_WrongPassword tests login with wrong password
func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) && err.Error() != ErrInvalidCredentials.Error() {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

// TestAuthService_Login_UserNotFound tests login with non-existent user
func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password",
	})

	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) && err.Error() != ErrInvalidCredentials.Error() {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

// TestAuthService_ValidateToken_Blacklisted tests that blacklisted tokens are rejected
func TestAuthService_ValidateToken_Blacklisted(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	token, err := svc.generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT failed: %v", err)
	}

	// Blacklist the token
	mockBlacklistRepo.blacklisted[token] = true

	_, err = svc.ValidateToken(context.Background(), token)
	if err == nil {
		t.Error("expected error when validating blacklisted token")
	}
}

// TestAuthService_generateJWT_DynamicIssuer tests JWT generation with dynamic issuer
func TestAuthService_generateJWT_DynamicIssuer(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	tests := []struct {
		name    string
		issuer  string
		user    *domain.User
		wantErr bool
	}{
		{
			name:   "success with custom issuer",
			issuer: "TestPlatform",
			user: &domain.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name:   "success with empty issuer (fallback)",
			issuer: "",
			user: &domain.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", tt.issuer)

			token, err := svc.generateJWT(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.generateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Error("AuthService.generateJWT() returned empty token")
			}

			// Verify issuer in token
			if !tt.wantErr && token != "" {
				parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				if err != nil {
					t.Fatalf("failed to parse token: %v", err)
				}
				claims, ok := parsedToken.Claims.(*JWTClaims)
				if !ok {
					t.Fatal("failed to extract claims")
				}
				expectedIssuer := tt.issuer
				if expectedIssuer == "" {
					expectedIssuer = "platform"
				}
				if claims.Issuer != expectedIssuer {
					t.Errorf("expected issuer %s, got %s", expectedIssuer, claims.Issuer)
				}
			}
		})
	}
}

// TestAuthService_UpdateProfile tests profile update
func TestAuthService_UpdateProfile(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	updated, err := svc.UpdateProfile(context.Background(), 1, UpdateProfileInput{
		Name:  "Updated Name",
		Email: "updated@example.com",
	})

	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
	if updated.Email != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got '%s'", updated.Email)
	}
}

// TestAuthService_UpdateProfile_EmailAlreadyExists tests profile update with existing email
func TestAuthService_UpdateProfile_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user1 := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	user2 := &domain.User{
		ID:           2,
		Email:        "other@example.com",
		Name:         "Other User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user1
	mockUserRepo.users[2] = user2

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	_, err := svc.UpdateProfile(context.Background(), 1, UpdateProfileInput{
		Name:  "Updated Name",
		Email: "other@example.com",
	})

	if err == nil {
		t.Error("expected error for existing email, got nil")
	}
	if !errors.Is(err, ErrEmailAlreadyExists) && err.Error() != ErrEmailAlreadyExists.Error() {
		t.Errorf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

// TestAuthService_ChangePassword tests password change
func TestAuthService_ChangePassword(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ChangePassword(context.Background(), 1, ChangePasswordInput{
		CurrentPassword: "password",
		NewPassword:     "newpassword",
	})

	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// Verify password was changed
	updatedUser := mockUserRepo.users[1]
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("newpassword"))
	if err != nil {
		t.Errorf("password was not updated correctly: %v", err)
	}
}

// TestAuthService_ChangePassword_WrongCurrentPassword tests password change with wrong current password
func TestAuthService_ChangePassword_WrongCurrentPassword(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ChangePassword(context.Background(), 1, ChangePasswordInput{
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword",
	})

	if err == nil {
		t.Error("expected error for wrong current password, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) && err.Error() != ErrInvalidCredentials.Error() {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

// TestAuthService_RequestPasswordReset tests password reset request
func TestAuthService_RequestPasswordReset(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.RequestPasswordReset(context.Background(), RequestPasswordResetInput{
		Email: "test@example.com",
	})

	if err != nil {
		t.Fatalf("RequestPasswordReset failed: %v", err)
	}

	// Verify a token was created
	tokens, err := mockPasswordResetRepo.FindByUserID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByUserID failed: %v", err)
	}
	if len(tokens) == 0 {
		t.Error("expected password reset token to be created")
	}
}

// TestAuthService_RequestPasswordReset_UserNotFound tests password reset for non-existent user (should not reveal)
func TestAuthService_RequestPasswordReset_UserNotFound(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.RequestPasswordReset(context.Background(), RequestPasswordResetInput{
		Email: "nonexistent@example.com",
	})

	if err != nil {
		t.Fatalf("RequestPasswordReset should not error for non-existent user: %v", err)
	}
}

// TestAuthService_ResetPassword tests password reset with valid token
func TestAuthService_ResetPassword(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	resetToken := &domain.PasswordResetToken{
		UserID:    1,
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	mockPasswordResetRepo.tokens["valid-token"] = resetToken

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ResetPassword(context.Background(), ResetPasswordInput{
		Token:       "valid-token",
		NewPassword: "newpassword",
	})

	if err != nil {
		t.Fatalf("ResetPassword failed: %v", err)
	}

	// Verify password was changed
	updatedUser := mockUserRepo.users[1]
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("newpassword"))
	if err != nil {
		t.Errorf("password was not reset correctly: %v", err)
	}

	// Verify token was marked as used
	if !resetToken.Used {
		t.Error("token should be marked as used")
	}
}

// TestAuthService_ResetPassword_InvalidToken tests password reset with invalid token
func TestAuthService_ResetPassword_InvalidToken(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ResetPassword(context.Background(), ResetPasswordInput{
		Token:       "invalid-token",
		NewPassword: "newpassword",
	})

	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
	if !errors.Is(err, ErrInvalidResetToken) && err.Error() != ErrInvalidResetToken.Error() {
		t.Errorf("expected ErrInvalidResetToken, got: %v", err)
	}
}

// TestAuthService_ResetPassword_ExpiredToken tests password reset with expired token
func TestAuthService_ResetPassword_ExpiredToken(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	resetToken := &domain.PasswordResetToken{
		UserID:    1,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Used:      false,
	}
	mockPasswordResetRepo.tokens["expired-token"] = resetToken

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ResetPassword(context.Background(), ResetPasswordInput{
		Token:       "expired-token",
		NewPassword: "newpassword",
	})

	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
	if !errors.Is(err, ErrInvalidResetToken) && err.Error() != ErrInvalidResetToken.Error() {
		t.Errorf("expected ErrInvalidResetToken, got: %v", err)
	}
}

// TestAuthService_ResetPassword_AlreadyUsedToken tests password reset with already used token
func TestAuthService_ResetPassword_AlreadyUsedToken(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}
	mockUserRepo.users[1] = user

	resetToken := &domain.PasswordResetToken{
		UserID:    1,
		Token:     "used-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      true,
	}
	mockPasswordResetRepo.tokens["used-token"] = resetToken

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	err := svc.ResetPassword(context.Background(), ResetPasswordInput{
		Token:       "used-token",
		NewPassword: "newpassword",
	})

	if err == nil {
		t.Error("expected error for already used token, got nil")
	}
	if !errors.Is(err, ErrResetTokenAlreadyUsed) && err.Error() != ErrResetTokenAlreadyUsed.Error() {
		t.Errorf("expected ErrResetTokenAlreadyUsed, got: %v", err)
	}
}

// TestAuthService_generateJWTWithImpersonation tests JWT generation with impersonation claims
func TestAuthService_generateJWTWithImpersonation(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockCompanyRepo := NewMockCompanyRepository()
	mockBlacklistRepo := NewMockTokenBlacklistRepository()
	mockPasswordResetRepo := NewMockPasswordResetRepository()

	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashPassword("password"),
		CompanyID:    1,
		Role:         domain.RoleOwner,
		Active:       true,
	}

	svc := NewAuthService(mockUserRepo, mockCompanyRepo, mockBlacklistRepo, mockPasswordResetRepo, "test-secret", "TestPlatform")

	token, err := svc.generateJWTWithImpersonation(user, true, 999)
	if err != nil {
		t.Fatalf("generateJWTWithImpersonation failed: %v", err)
	}

	// Verify impersonation claims
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if !claims.IsImpersonating {
		t.Error("expected IsImpersonating to be true")
	}
	if claims.OriginalPlatformUserID != 999 {
		t.Errorf("expected OriginalPlatformUserID 999, got %d", claims.OriginalPlatformUserID)
	}
}
