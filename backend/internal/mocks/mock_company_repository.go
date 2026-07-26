package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockCompanyRepository is a mock implementation of ports.CompanyRepository
type MockCompanyRepository struct {
	Companies   map[uint]*domain.Company
	CreateError error
	FindError   error
	UpdateError error
	DeleteError error
}

func NewMockCompanyRepository() *MockCompanyRepository {
	return &MockCompanyRepository{
		Companies: make(map[uint]*domain.Company),
	}
}

func (m *MockCompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if company.ID == 0 {
		company.ID = uint(len(m.Companies) + 1)
	}
	m.Companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Companies[id], nil
}

func (m *MockCompanyRepository) FindBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	for _, company := range m.Companies {
		if company.Slug == slug {
			return company, nil
		}
	}
	return nil, nil
}

func (m *MockCompanyRepository) List(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	for _, company := range m.Companies {
		companies = append(companies, *company)
	}
	return companies, nil
}

func (m *MockCompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Companies, id)
	return nil
}
