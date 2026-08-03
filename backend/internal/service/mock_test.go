package service

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	users             map[uint]*domain.User
	findByEmailResult *domain.User
	findByEmailError  error
	updateError       error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) CreateBootstrapOwner(ctx context.Context, user *domain.User, companyID uint) error {
	user.CompanyID = companyID
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	return m.users[id], nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.findByEmailResult != nil {
		return m.findByEmailResult, m.findByEmailError
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.users[user.ID] = user
	return nil
}

type MockCompanyRepository struct {
	companies map[uint]*domain.Company
}

func NewMockCompanyRepository() *MockCompanyRepository {
	return &MockCompanyRepository{
		companies: make(map[uint]*domain.Company),
	}
}

func (m *MockCompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	return m.companies[id], nil
}

func (m *MockCompanyRepository) FindBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	for _, company := range m.companies {
		if company.Slug == slug {
			return company, nil
		}
	}
	return nil, nil
}

func (m *MockCompanyRepository) List(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	for _, company := range m.companies {
		companies = append(companies, *company)
	}
	return companies, nil
}

func (m *MockCompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id uint) error {
	delete(m.companies, id)
	return nil
}

type MockTokenBlacklistRepository struct {
	blacklisted map[string]bool
	addError    error
	checkError  error
}

func NewMockTokenBlacklistRepository() *MockTokenBlacklistRepository {
	return &MockTokenBlacklistRepository{
		blacklisted: make(map[string]bool),
	}
}

func (m *MockTokenBlacklistRepository) Add(ctx context.Context, entry *domain.TokenBlacklist) error {
	if m.addError != nil {
		return m.addError
	}
	m.blacklisted[entry.Token] = true
	return nil
}

func (m *MockTokenBlacklistRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	if m.checkError != nil {
		return false, m.checkError
	}
	return m.blacklisted[token], nil
}

func (m *MockTokenBlacklistRepository) CleanExpired(ctx context.Context) error {
	return nil
}

type MockPasswordResetRepository struct {
	tokens map[string]*domain.PasswordResetToken
}

func NewMockPasswordResetRepository() *MockPasswordResetRepository {
	return &MockPasswordResetRepository{
		tokens: make(map[string]*domain.PasswordResetToken),
	}
}

func (m *MockPasswordResetRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	m.tokens[token.Token] = token
	return nil
}

func (m *MockPasswordResetRepository) FindByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	return m.tokens[token], nil
}

func (m *MockPasswordResetRepository) FindByTokenForUpdate(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	return m.tokens[token], nil
}

func (m *MockPasswordResetRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.PasswordResetToken, error) {
	var tokens []*domain.PasswordResetToken
	for _, token := range m.tokens {
		if token.UserID == userID {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}

func (m *MockPasswordResetRepository) MarkAsUsed(ctx context.Context, token *domain.PasswordResetToken) error {
	token.Used = true
	m.tokens[token.Token] = token
	return nil
}

func (m *MockPasswordResetRepository) Delete(ctx context.Context, token *domain.PasswordResetToken) error {
	delete(m.tokens, token.Token)
	return nil
}

func (m *MockPasswordResetRepository) DeleteExpired(ctx context.Context) error {
	return nil
}
