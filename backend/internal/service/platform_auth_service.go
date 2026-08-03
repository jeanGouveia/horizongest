package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/security"
)

var (
	ErrPlatformEmailAlreadyExists = errors.New("platform email already exists")
	ErrPlatformInvalidCredentials = errors.New("platform invalid credentials")
	ErrPlatformUserNotFound       = errors.New("platform user not found")
	ErrPlatformUserInactive       = errors.New("platform user is inactive")
)

type PlatformAuthService struct {
	platformUserRepo    PlatformUserRepository
	platformSessionRepo PlatformSessionRepository
	keyStore            *security.JWTKeyStore // FASE A: JWT key store for rotation
	sessionDuration     time.Duration
	bcryptCost          int
}

type PlatformUserRepository interface {
	Create(ctx context.Context, user *domain.PlatformUser) error
	FindByEmail(ctx context.Context, email string) (*domain.PlatformUser, error)
	FindByID(ctx context.Context, id uint) (*domain.PlatformUser, error)
	Update(ctx context.Context, user *domain.PlatformUser) error
}

type PlatformSessionRepository interface {
	Create(ctx context.Context, session *domain.PlatformSession) error
	FindByToken(ctx context.Context, token string) (*domain.PlatformSession, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByPlatformUserID(ctx context.Context, platformUserID uint) error
}

func NewPlatformAuthService(
	platformUserRepo PlatformUserRepository,
	platformSessionRepo PlatformSessionRepository,
	jwtSecret string,
	sessionDuration time.Duration,
	bcryptCost int,
) *PlatformAuthService {
	if jwtSecret == "" {
		panic("JWT_PLATFORM_SECRET environment variable is required but not set")
	}
	// FASE A: Initialize JWT key store for rotation
	keyStore, err := security.NewJWTKeyStore(jwtSecret, 24*time.Hour, 30*24*time.Hour)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize JWT key store: %v", err))
	}
	return &PlatformAuthService{
		platformUserRepo:    platformUserRepo,
		platformSessionRepo: platformSessionRepo,
		keyStore:            keyStore,
		sessionDuration:     sessionDuration,
		bcryptCost:          bcryptCost,
	}
}

type PlatformLoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type PlatformLoginOutput struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      *PlatformUserOutput
}

type PlatformUserOutput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Login authenticates a platform user and returns a JWT token
func (s *PlatformAuthService) Login(ctx context.Context, input PlatformLoginInput) (*PlatformLoginOutput, error) {
	// Find user by email
	user, err := s.platformUserRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("PlatformAuthService.Login: buscar usuário: %w", err)
	}
	if user == nil {
		return nil, ErrPlatformInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return nil, ErrPlatformUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrPlatformInvalidCredentials
	}

	// Generate JWT token
	token, expiresAt, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("PlatformAuthService.Login: gerar token: %w", err)
	}

	// Create session
	session := &domain.PlatformSession{
		PlatformUserID: user.ID,
		Token:          token,
		ExpiresAt:      time.Unix(expiresAt, 0),
		CreatedAt:      time.Now(),
	}
	if err := s.platformSessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("PlatformAuthService.Login: criar sessão: %w", err)
	}

	return &PlatformLoginOutput{
		Token:     token,
		ExpiresAt: expiresAt,
		User: &PlatformUserOutput{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.String(),
		},
	}, nil
}

// Logout invalidates a platform user's session
func (s *PlatformAuthService) Logout(ctx context.Context, token string) error {
	if err := s.platformSessionRepo.DeleteByToken(ctx, token); err != nil {
		return fmt.Errorf("PlatformAuthService.Logout: deletar sessão: %w", err)
	}
	return nil
}

// ValidateToken validates a JWT token and returns the platform user ID and role
func (s *PlatformAuthService) ValidateToken(ctx context.Context, token string) (uint, domain.PlatformRole, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("PlatformAuthService.ValidateToken: algoritmo inesperado: %v", t.Header["alg"])
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
			return nil, fmt.Errorf("PlatformAuthService.ValidateToken: chave não encontrada ou expirada: %s", kid)
		}
		return []byte(key.Secret), nil
	})

	if err != nil {
		return 0, domain.PlatformRoleSupport, fmt.Errorf("PlatformAuthService.ValidateToken: validar token: %w", err)
	}

	// Extract user ID and role
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, domain.PlatformRoleSupport, errors.New("invalid token: missing subject")
	}

	userID := uint(userIDFloat)

	roleStr, ok := claims["role"].(string)
	if !ok {
		return 0, domain.PlatformRoleSupport, errors.New("invalid token: missing role")
	}

	role, ok := domain.ParsePlatformRole(roleStr)
	if !ok {
		return 0, domain.PlatformRoleSupport, errors.New("invalid token: invalid role")
	}

	// FASE A: Check if session exists and is valid
	session, err := s.platformSessionRepo.FindByToken(ctx, token)
	if err != nil {
		return 0, domain.PlatformRoleSupport, fmt.Errorf("PlatformAuthService.ValidateToken: verificar sessão: %w", err)
	}
	if session == nil {
		return 0, domain.PlatformRoleSupport, errors.New("session not found or revoked")
	}
	// Verify session belongs to the same user
	if session.PlatformUserID != userID {
		return 0, domain.PlatformRoleSupport, errors.New("session user mismatch")
	}
	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return 0, domain.PlatformRoleSupport, errors.New("session expired")
	}

	return userID, role, nil
}

// generateToken generates a JWT token for a platform user
func (s *PlatformAuthService) generateToken(userID uint, role domain.PlatformRole) (string, int64, error) {
	expiresAt := time.Now().Add(s.sessionDuration)

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role.String(),
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
		"type": "platform",
	}

	// FASE A: Get active key and add kid header
	activeKey := s.keyStore.GetActiveKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = activeKey.ID
	tokenString, err := token.SignedString([]byte(activeKey.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// CreatePlatformUser creates a new platform user (only for PlatformAdmin)
func (s *PlatformAuthService) CreatePlatformUser(ctx context.Context, name, email, password string, role domain.PlatformRole) (*domain.PlatformUser, error) {
	// Check if email already exists
	existing, err := s.platformUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("PlatformAuthService.CreatePlatformUser: verificar usuário existente: %w", err)
	}
	if existing != nil {
		return nil, ErrPlatformEmailAlreadyExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("PlatformAuthService.CreatePlatformUser: gerar hash: %w", err)
	}

	// Create user
	user := &domain.PlatformUser{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.platformUserRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("PlatformAuthService.CreatePlatformUser: criar usuário: %w", err)
	}

	return user, nil
}

// GetPlatformUser retrieves a platform user by ID
func (s *PlatformAuthService) GetPlatformUser(ctx context.Context, userID uint) (*PlatformUserOutput, error) {
	user, err := s.platformUserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("PlatformAuthService.GetPlatformUser: buscar usuário: %w", err)
	}
	if user == nil {
		return nil, ErrPlatformUserNotFound
	}

	return &PlatformUserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.String(),
	}, nil
}
