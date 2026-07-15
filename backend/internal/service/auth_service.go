package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var (
	ErrEmailAlreadyExists = errors.New("e-mail já cadastrado")
	ErrInvalidCredentials = errors.New("e-mail ou senha inválidos")
)

// JWTClaims é exportado para que o middleware possa usar o tipo.
type JWTClaims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo   ports.UserRepository
	secret     []byte
	expiry     time.Duration
	bcryptCost int
}

func NewAuthService(userRepo ports.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-troque-em-producao"
	}
	return &AuthService{
		userRepo:   userRepo,
		secret:     []byte(secret),
		expiry:     24 * time.Hour,
		bcryptCost: bcrypt.DefaultCost,
	}
}

// --- Register ---

type RegisterInput struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	existing, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("Register: hash: %w", err)
	}

	user := &domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
	}
	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}
	return user, nil
}

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
		return nil, fmt.Errorf("Login: %w", err)
	}
	// Mesmo erro para e-mail inexistente e senha errada (evita user enumeration)
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
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
		return nil, fmt.Errorf("UpdateProfile: %w", err)
	}
	if existing != nil && existing.ID != userID {
		return nil, ErrEmailAlreadyExists
	}

	// Buscar usuário atual
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UpdateProfile: %w", err)
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Atualizar campos
	user.Name = input.Name
	user.Email = input.Email

	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("UpdateProfile: %w", err)
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
		return fmt.Errorf("ChangePassword: %w", err)
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
		return fmt.Errorf("ChangePassword: hash: %w", err)
	}

	// Atualizar senha
	user.PasswordHash = string(hash)
	if err = s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("ChangePassword: %w", err)
	}
	return nil
}

// --- ValidateToken ---
// Retorna *JWTClaims (exportado) para o middleware extrair UserID, Email e Name.

func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo inesperado: %v", t.Header["alg"])
			}
			return s.secret, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("ValidateToken: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("ValidateToken: claims inválidos")
	}
	return claims, nil
}

// --- helper privado ---

func (s *AuthService) generateJWT(user *domain.User) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "pratoOnline",
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("generateJWT: %w", err)
	}
	return signed, nil
}
