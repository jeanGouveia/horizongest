package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

// Mock repository for testing
type mockPlatformBrandRepository struct {
	brand *domain.PlatformBrandConfig
	err   error
}

func (m *mockPlatformBrandRepository) Get(ctx context.Context) (*domain.PlatformBrandConfig, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.brand, nil
}

func (m *mockPlatformBrandRepository) Update(ctx context.Context, brand *domain.PlatformBrandConfig, updatedBy uint) error {
	if m.err != nil {
		return m.err
	}
	m.brand = brand
	return nil
}

func (m *mockPlatformBrandRepository) Initialize(ctx context.Context) error {
	return m.err
}

func TestPlatformBrandService_Get(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockPlatformBrandRepository
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockPlatformBrandRepository{
				brand: &domain.PlatformBrandConfig{
					PlatformName:      "TestPlatform",
					PlatformShortName: "Test",
					OwnerCompanyName:  "Test Company",
					Website:          "https://test.com",
					SupportEmail:      "support@test.com",
					SupportURL:        "https://support.test.com",
					Copyright:         "© 2024 Test Company",
					PrimaryColor:      "#000000",
					SecondaryColor:    "#ffffff",
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockPlatformBrandRepository{
				err: errors.New("repository error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPlatformBrandService(tt.repo)
			_, err := svc.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("PlatformBrandService.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlatformBrandService_Update(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockPlatformBrandRepository
		brand   *domain.PlatformBrandConfig
		userID  uint
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockPlatformBrandRepository{
				brand: &domain.PlatformBrandConfig{},
			},
			brand: &domain.PlatformBrandConfig{
				PlatformName:      "TestPlatform",
				PlatformShortName: "Test",
				OwnerCompanyName:  "Test Company",
				Website:          "https://test.com",
				SupportEmail:      "support@test.com",
				SupportURL:        "https://support.test.com",
				Copyright:         "© 2024 Test Company",
				PrimaryColor:      "#000000",
				SecondaryColor:    "#ffffff",
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "invalid brand - missing platform name",
			repo: &mockPlatformBrandRepository{
				brand: &domain.PlatformBrandConfig{},
			},
			brand: &domain.PlatformBrandConfig{
				PlatformShortName: "Test",
				OwnerCompanyName:  "Test Company",
				Website:          "https://test.com",
				SupportEmail:      "support@test.com",
				SupportURL:        "https://support.test.com",
				Copyright:         "© 2024 Test Company",
				PrimaryColor:      "#000000",
				SecondaryColor:    "#ffffff",
			},
			userID:  1,
			wantErr: true,
		},
		{
			name: "repository error",
			repo: &mockPlatformBrandRepository{
				err: errors.New("repository error"),
			},
			brand: &domain.PlatformBrandConfig{
				PlatformName:      "TestPlatform",
				PlatformShortName: "Test",
				OwnerCompanyName:  "Test Company",
				Website:          "https://test.com",
				SupportEmail:      "support@test.com",
				SupportURL:        "https://support.test.com",
				Copyright:         "© 2024 Test Company",
				PrimaryColor:      "#000000",
				SecondaryColor:    "#ffffff",
			},
			userID:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPlatformBrandService(tt.repo)
			err := svc.Update(context.Background(), tt.brand, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlatformBrandService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlatformBrandService_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockPlatformBrandRepository
		wantErr bool
	}{
		{
			name:    "success",
			repo:    &mockPlatformBrandRepository{},
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockPlatformBrandRepository{
				err: errors.New("repository error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPlatformBrandService(tt.repo)
			err := svc.Initialize(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("PlatformBrandService.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
