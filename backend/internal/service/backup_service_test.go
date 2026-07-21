package service

import (
	"testing"
)

// Test for BackupService focusing on dynamic platform name in filename
// This is critical for branding support (Sprint 3.7)

func TestBackupService_CreateBackup_DynamicPlatformName(t *testing.T) {
	tests := []struct {
		name         string
		platformName string
	}{
		{
			name:         "success with custom platform name",
			platformName: "TestPlatform",
		},
		{
			name:         "success with empty platform name (fallback)",
			platformName: "",
		},
		{
			name:         "success with HorizonGest",
			platformName: "HorizonGest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewBackupService(
				"localhost",
				"3306",
				"root",
				"",
				"testdb",
				"./test_backups",
				tt.platformName,
			)

			// Note: We're not actually running the backup since it requires mysqldump
			// This test verifies the service can be created with dynamic platform name
			if svc == nil {
				t.Error("NewBackupService() returned nil")
			}
		})
	}
}

func TestBackupService_platformName_Fallback(t *testing.T) {
	tests := []struct {
		name         string
		platformName string
		expected     string
	}{
		{
			name:         "custom platform name",
			platformName: "TestPlatform",
			expected:     "TestPlatform",
		},
		{
			name:         "empty platform name (no fallback in constructor)",
			platformName: "",
			expected:     "", // Constructor doesn't apply fallback, CreateBackup does
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewBackupService(
				"localhost",
				"3306",
				"root",
				"",
				"testdb",
				"./test_backups",
				tt.platformName,
			)

			if svc.platformName != tt.expected {
				t.Errorf("BackupService.platformName = %v, want %v", svc.platformName, tt.expected)
			}
		})
	}
}
