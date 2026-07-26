package mocks

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockUserRepository is a mock implementation of ports.UserRepository
type MockUserRepository struct {
	Users             map[uint]*domain.User
	CreateError       error
	FindError         error
	UpdateError       error
	DeleteError       error
	FindByEmailResult *domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if user.ID == 0 {
		user.ID = uint(len(m.Users) + 1)
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.Users[id], nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	if m.FindByEmailResult != nil {
		return m.FindByEmailResult, nil
	}
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range m.Users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Users, id)
	return nil
}
