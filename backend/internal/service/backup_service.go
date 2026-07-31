package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type BackupService struct {
	dbHost       string
	dbPort       string
	dbUser       string
	dbPassword   string
	dbName       string
	backupDir    string
	platformName string // Platform name for backup filename prefix (Sprint 3.6)
}

func NewBackupService(dbHost, dbPort, dbUser, dbPassword, dbName, backupDir, platformName string) *BackupService {
	return &BackupService{
		dbHost:       dbHost,
		dbPort:       dbPort,
		dbUser:       dbUser,
		dbPassword:   dbPassword,
		dbName:       dbName,
		backupDir:    backupDir,
		platformName: platformName,
	}
}

type BackupResult struct {
	FileName  string
	Size      int64
	Path      string
	CreatedAt time.Time
}

func (s *BackupService) CreateBackup(ctx context.Context) (*BackupResult, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("BackupService.CreateBackup: criar diretório de backup: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	// Use platform name for backup filename prefix (Sprint 3.6)
	platformPrefix := s.platformName
	if platformPrefix == "" {
		platformPrefix = "platform" // Fallback if platform name is empty
	}
	fileName := fmt.Sprintf("%s_backup_%s.sql", platformPrefix, timestamp)
	filePath := fmt.Sprintf("%s/%s", s.backupDir, fileName)

	// Build mysqldump command
	args := []string{
		"-h", s.dbHost,
		"-P", s.dbPort,
		"-u", s.dbUser,
		fmt.Sprintf("-p%s", s.dbPassword),
		s.dbName,
		"--single-transaction",
		"--quick",
		"--lock-tables=false",
		"--routines",
		"--triggers",
		"--events",
	}

	// Execute mysqldump
	cmd := exec.CommandContext(ctx, "mysqldump", args...)
	cmd.Stdout = nil // We'll redirect to file

	// Create output file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("BackupService.CreateBackup: criar arquivo de backup: %w", err)
	}
	defer file.Close()

	cmd.Stdout = file

	if err := cmd.Run(); err != nil {
		// Clean up failed backup file
		os.Remove(filePath)
		return nil, fmt.Errorf("BackupService.CreateBackup: mysqldump falhou: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("BackupService.CreateBackup: obter informações do arquivo: %w", err)
	}

	return &BackupResult{
		FileName:  fileName,
		Size:      fileInfo.Size(),
		Path:      filePath,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BackupService) ListBackups(ctx context.Context) ([]BackupResult, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupResult{}, nil
		}
		return nil, fmt.Errorf("BackupService.ListBackups: ler diretório de backup: %w", err)
	}

	var backups []BackupResult
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		backups = append(backups, BackupResult{
			FileName:  entry.Name(),
			Size:      info.Size(),
			Path:      fmt.Sprintf("%s/%s", s.backupDir, entry.Name()),
			CreatedAt: info.ModTime(),
		})
	}

	return backups, nil
}

func (s *BackupService) DeleteBackup(ctx context.Context, fileName string) error {
	filePath := fmt.Sprintf("%s/%s", s.backupDir, fileName)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("BackupService.DeleteBackup: deletar backup: %w", err)
	}
	return nil
}
