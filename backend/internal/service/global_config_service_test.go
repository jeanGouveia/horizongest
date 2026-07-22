package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// Mock repository for testing
type mockGlobalConfigRepository struct {
	config *domain.GlobalConfig
	err    error
}

func (m *mockGlobalConfigRepository) Get(ctx context.Context) (*domain.GlobalConfig, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.config, nil
}

func (m *mockGlobalConfigRepository) Update(ctx context.Context, config *domain.GlobalConfig, updatedBy uint) error {
	if m.err != nil {
		return m.err
	}
	m.config = config
	return nil
}

func (m *mockGlobalConfigRepository) Initialize(ctx context.Context) error {
	return m.err
}

func TestGlobalConfigService_Get(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockGlobalConfigRepository
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{
					DefaultTimezone: "America/Sao_Paulo",
					DefaultLocale:   "pt-BR",
					MonetaryFormat:  "BRL R$ 1.000,00",
					DateFormat:      "DD/MM/YYYY",
					TimeFormat:      "HH:mm",
					MaxUploadSizeMB: 10,
					MaxImageSizeMB:  5,
				},
			},
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockGlobalConfigRepository{
				err: errors.New("repository error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewGlobalConfigService(tt.repo)
			_, err := svc.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfigService.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalConfigService_Update(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockGlobalConfigRepository
		config  *domain.GlobalConfig
		userID  uint
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{},
			},
			config: &domain.GlobalConfig{
				DefaultTimezone:   "America/Sao_Paulo",
				DefaultLocale:     "pt-BR",
				MonetaryFormat:    "BRL R$ 1.000,00",
				DateFormat:        "DD/MM/YYYY",
				TimeFormat:        "HH:mm",
				MaxUploadSizeMB:   10,
				MaxImageSizeMB:    5,
				AllowedImageTypes: "jpg,png,webp",
				AllowedFileTypes:  "pdf,doc,xlsx",
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "invalid config - missing timezone",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{},
			},
			config: &domain.GlobalConfig{
				DefaultLocale:   "pt-BR",
				MonetaryFormat:  "BRL R$ 1.000,00",
				DateFormat:      "DD/MM/YYYY",
				TimeFormat:      "HH:mm",
				MaxUploadSizeMB: 10,
				MaxImageSizeMB:  5,
			},
			userID:  1,
			wantErr: true,
		},
		{
			name: "repository error",
			repo: &mockGlobalConfigRepository{
				err: errors.New("repository error"),
			},
			config: &domain.GlobalConfig{
				DefaultTimezone: "America/Sao_Paulo",
				DefaultLocale:   "pt-BR",
				MonetaryFormat:  "BRL R$ 1.000,00",
				DateFormat:      "DD/MM/YYYY",
				TimeFormat:      "HH:mm",
				MaxUploadSizeMB: 10,
				MaxImageSizeMB:  5,
			},
			userID:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewGlobalConfigService(tt.repo)
			err := svc.Update(context.Background(), tt.config, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfigService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGlobalConfigService_IsModuleEnabled(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockGlobalConfigRepository
		module  string
		want    bool
		wantErr bool
	}{
		{
			name: "finance enabled",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{
					EnableFinance: true,
				},
			},
			module:  "finance",
			want:    true,
			wantErr: false,
		},
		{
			name: "finance disabled",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{
					EnableFinance: false,
				},
			},
			module:  "finance",
			want:    false,
			wantErr: false,
		},
		{
			name: "unknown module",
			repo: &mockGlobalConfigRepository{
				config: &domain.GlobalConfig{},
			},
			module:  "unknown",
			want:    false,
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockGlobalConfigRepository{
				err: errors.New("repository error"),
			},
			module:  "finance",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewGlobalConfigService(tt.repo)
			got, err := svc.IsModuleEnabled(context.Background(), tt.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfigService.IsModuleEnabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GlobalConfigService.IsModuleEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalConfigService_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockGlobalConfigRepository
		wantErr bool
	}{
		{
			name:    "success",
			repo:    &mockGlobalConfigRepository{},
			wantErr: false,
		},
		{
			name: "repository error",
			repo: &mockGlobalConfigRepository{
				err: errors.New("repository error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewGlobalConfigService(tt.repo)
			err := svc.Initialize(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfigService.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
