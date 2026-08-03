package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"github.com/jeanGouveia/horizongest/backend/internal/security"
	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

var (
	ErrEmailAlreadyExists    = errors.New("e-mail já cadastrado")
	ErrInvalidCredentials    = errors.New("e-mail ou senha inválidos")
	ErrInvalidResetToken     = errors.New("token de recuperação inválido ou expirado")
	ErrResetTokenAlreadyUsed = errors.New("token de recuperação já foi utilizado")
)

// JWTClaims é exportado para que o middleware possa usar o tipo.
type JWTClaims struct {
	UserID    uint   `json:"uid"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CompanyID uint   `json:"cid"`
	// Impersonation claims
	IsImpersonating        bool `json:"imp"`
	OriginalPlatformUserID uint `json:"opuid,omitempty"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo          ports.UserRepository
	companyRepo       ports.CompanyRepository
	tokenBlacklist    ports.TokenBlacklistRepository
	passwordResetRepo ports.PasswordResetRepository
	keyStore          *security.JWTKeyStore // FASE A: JWT key store for rotation
	expiry            time.Duration
	bcryptCost        int
	issuer            string // JWT issuer (platform name) - Sprint 3.6
}

func NewAuthService(userRepo ports.UserRepository, companyRepo ports.CompanyRepository, tokenBlacklist ports.TokenBlacklistRepository, passwordResetRepo ports.PasswordResetRepository, jwtSecret string, issuer string) *AuthService {
	if jwtSecret == "" {
		panic("JWT_TENANT_SECRET environment variable is required but not set")
	}
	// FASE A: Initialize JWT key store for rotation
	keyStore, err := security.NewJWTKeyStore(jwtSecret, 24*time.Hour, 30*24*time.Hour)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize JWT key store: %v", err))
	}
	return &AuthService{
		userRepo:          userRepo,
		companyRepo:       companyRepo,
		tokenBlacklist:    tokenBlacklist,
		passwordResetRepo: passwordResetRepo,
		keyStore:          keyStore,
		expiry:            24 * time.Hour,
		bcryptCost:        bcrypt.DefaultCost,
		issuer:            issuer, // JWT issuer from platform brand (Sprint 3.6)
	}
}

// --- Register (REMOVED - Sprint 3) ---
// Public registration has been removed. Companies are now created by platform administrators only.
// This method is kept for reference but should not be used.

// --- Login ---

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResult struct {
	Token string
	User  *domain.User
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login: buscar usuário: %w", err)
	}
	// Mesmo erro para e-mail inexistente e senha errada (evita user enumeration)
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("usuário desativado. Entre em contato com o administrador.")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login: gerar token: %w", err)
	}
	return &LoginResult{Token: token, User: user}, nil
}

// --- UpdateProfile ---

type UpdateProfileInput struct {
	Name  string `json:"name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, input UpdateProfileInput) (*domain.User, error) {
	// Verificar se o e-mail já está em uso por outro usuário
	existing, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("AuthService.UpdateProfile: verificar e-mail existente: %w", err)
	}
	if existing != nil && existing.ID != userID {
		return nil, ErrEmailAlreadyExists
	}

	// Buscar usuário atual
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("AuthService.UpdateProfile: buscar usuário: %w", err)
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Atualizar campos
	user.Name = input.Name
	user.Email = input.Email
	// CompanyID não pode ser alterado pelo próprio usuário - apenas via convite ou endpoints administrativos

	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("AuthService.UpdateProfile: atualizar usuário: %w", err)
	}
	return user, nil
}

// --- ChangePassword ---

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=6"`
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uint, input ChangePasswordInput) error {
	// Buscar usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("AuthService.ChangePassword: buscar usuário: %w", err)
	}
	if user == nil {
		return errors.New("usuário não encontrado")
	}

	// Verificar senha atual
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash da nova senha
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), s.bcryptCost)
	if err != nil {
		return fmt.Errorf("AuthService.ChangePassword: gerar hash: %w", err)
	}

	// Atualizar senha
	user.PasswordHash = string(hash)
	if err = s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("AuthService.ChangePassword: atualizar usuário: %w", err)
	}
	return nil
}

// --- ValidateToken ---
// Retorna *JWTClaims (exportado) para o middleware extrair UserID, Email e Name.

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (*JWTClaims, error) {
	// Sprint 4A: Remover log de JWT bruto por segurança
	// log.Printf("[FORENSIC] ValidateToken - JWT bruto recebido: %s", tokenStr)

	// Check if token is blacklisted no banco
	blacklisted, err := s.tokenBlacklist.IsBlacklisted(ctx, tokenStr)
	if err != nil {
		return nil, fmt.Errorf("AuthService.ValidateToken: verificar blacklist: %w", err)
	}
	if blacklisted {
		return nil, errors.New("token was revoked")
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("AuthService.ValidateToken: algoritmo inesperado: %v", t.Header["alg"])
			}
			// FASE A: Extract kid from header and resolve key
			kid, ok := t.Header["kid"].(string)
			if !ok {
				// Fallback to active key for backward compatibility
				activeKey := s.keyStore.GetActiveKey()
				return []byte(activeKey.Secret), nil
			}
			// Resolve key by kid
			key, found := s.keyStore.GetKeyByID(kid)
			if !found {
				return nil, fmt.Errorf("AuthService.ValidateToken: chave não encontrada ou expirada: %s", kid)
			}
			return []byte(key.Secret), nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("AuthService.ValidateToken: validar token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("AuthService.ValidateToken: claims inválidos")
	}

	// Sprint 4A: Remover log de claims sensíveis por segurança
	// log.Printf("[FORENSIC] ValidateToken - Claims validados - UserID: %d, CompanyID: %d, Email: %s, Name: %s, Issuer: %s, Subject: %s, IsImpersonating: %v",
	//	claims.UserID, claims.CompanyID, claims.Email, claims.Name, claims.Issuer, claims.Subject, claims.IsImpersonating)

	return claims, nil
}

// --- Logout ---

func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
	// Extrair expiration do token para saber quando expira
	claims, err := s.parseTokenClaims(tokenStr)
	if err != nil {
		return fmt.Errorf("AuthService.Logout: extrair claims: %w", err)
	}

	// Persistir no banco
	entry := &domain.TokenBlacklist{
		Token:     tokenStr,
		RevokedAt: time.Now(),
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if err := s.tokenBlacklist.Add(ctx, entry); err != nil {
		return fmt.Errorf("AuthService.Logout: adicionar à blacklist: %w", err)
	}

	return nil
}

// parseTokenClaims extrai os claims de um token sem validar a blacklist
func (s *AuthService) parseTokenClaims(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("AuthService.parseTokenClaims: algoritmo inesperado: %v", t.Header["alg"])
			}
			// FASE A: Extract kid from header and resolve key
			kid, ok := t.Header["kid"].(string)
			if !ok {
				// Fallback to active key for backward compatibility
				activeKey := s.keyStore.GetActiveKey()
				return []byte(activeKey.Secret), nil
			}
			// Resolve key by kid
			key, found := s.keyStore.GetKeyByID(kid)
			if !found {
				return nil, fmt.Errorf("AuthService.parseTokenClaims: chave não encontrada ou expirada: %s", kid)
			}
			return []byte(key.Secret), nil
		},
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, fmt.Errorf("AuthService.parseTokenClaims: analisar token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("AuthService.parseTokenClaims: claims inválidos")
	}
	return claims, nil
}

// --- helper privado ---

func (s *AuthService) generateJWT(user *domain.User) (string, error) {
	return s.generateJWTWithImpersonation(user, false, 0)
}

// generateJWTWithImpersonation generates a JWT token with optional impersonation claims
func (s *AuthService) generateJWTWithImpersonation(user *domain.User, isImpersonating bool, originalPlatformUserID uint) (string, error) {
	now := time.Now()
	// Use dynamic issuer from platform brand (Sprint 3.6)
	issuer := s.issuer
	if issuer == "" {
		issuer = "platform" // Fallback if issuer is empty
	}
	claims := JWTClaims{
		UserID:                 user.ID,
		Email:                  user.Email,
		Name:                   user.Name,
		CompanyID:              user.CompanyID,
		IsImpersonating:        isImpersonating,
		OriginalPlatformUserID: originalPlatformUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// FASE A: Mask email in logs
	emailMask := util.MaskEmail(claims.Email)
	log.Printf("[DEBUG] generateJWTWithImpersonation - Claims: UserID=%d, CompanyID=%d, Email=%s, Name=%s, IsImpersonating=%v, OriginalPlatformUserID=%d",
		claims.UserID, claims.CompanyID, emailMask, claims.Name, claims.IsImpersonating, claims.OriginalPlatformUserID)

	// FASE A: Get active key and add kid header
	activeKey := s.keyStore.GetActiveKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = activeKey.ID
	signed, err := token.SignedString([]byte(activeKey.Secret))
	if err != nil {
		return "", fmt.Errorf("AuthService.generateJWTWithImpersonation: assinar token: %w", err)
	}
	return signed, nil
}

// --- Password Recovery ---

type RequestPasswordResetInput struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, input RequestPasswordResetInput) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("AuthService.RequestPasswordReset: buscar usuário: %w", err)
	}
	if user == nil {
		// Don't reveal if email exists or not (security)
		return nil
	}

	// Generate secure token using time and random string
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("AuthService.RequestPasswordReset: gerar token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token with 1 hour expiration
	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return fmt.Errorf("AuthService.RequestPasswordReset: criar token: %w", err)
	}

	// Email is sent via EmailService
	// Sprint 3.4 - Security Hardening: Remove sensitive token from logs
	log.Printf("[AUTH] Password reset requested for user ID: %d", user.ID)

	return nil
}

type ResetPasswordInput struct {
	Token       string `json:"token"       validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (s *AuthService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	// FASE A: Use FindByTokenForUpdate with SELECT FOR UPDATE to prevent race condition
	resetToken, err := s.passwordResetRepo.FindByTokenForUpdate(ctx, input.Token)
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword: buscar token: %w", err)
	}
	if resetToken == nil {
		return ErrInvalidResetToken
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return ErrInvalidResetToken
	}

	// Check if token was already used (double-check after lock)
	if resetToken.Used {
		return ErrResetTokenAlreadyUsed
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword: buscar usuário: %w", err)
	}
	if user == nil {
		return ErrInvalidResetToken
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), s.bcryptCost)
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword: gerar hash: %w", err)
	}

	// Update user password
	user.PasswordHash = string(hash)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("AuthService.ResetPassword: atualizar usuário: %w", err)
	}

	// Mark token as used
	resetToken.Used = true
	if err := s.passwordResetRepo.MarkAsUsed(ctx, resetToken); err != nil {
		return fmt.Errorf("AuthService.ResetPassword: marcar token como usado: %w", err)
	}

	return nil
}
