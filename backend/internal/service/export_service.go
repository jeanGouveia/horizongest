package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type ExportService struct {
	companyRepo ports.CompanyRepository
	userRepo    ports.UserRepository
	exportDir   string
}

func NewExportService(companyRepo ports.CompanyRepository, userRepo ports.UserRepository, exportDir string) *ExportService {
	return &ExportService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
		exportDir:   exportDir,
	}
}

type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
)

type ExportResult struct {
	FileName  string
	Format    ExportFormat
	Size      int64
	Path      string
	CreatedAt time.Time
}

func (s *ExportService) ExportCompanies(ctx context.Context, format ExportFormat) (*ExportResult, error) {
	// Get all companies
	companies, err := s.companyRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("ExportCompanies: failed to list companies: %w", err)
	}

	// Ensure export directory exists
	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return nil, fmt.Errorf("ExportCompanies: failed to create export directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var fileName, filePath string

	switch format {
	case ExportFormatCSV:
		fileName = fmt.Sprintf("companies_export_%s.csv", timestamp)
		filePath = fmt.Sprintf("%s/%s", s.exportDir, fileName)
		if err := s.exportCompaniesToCSV(companies, filePath); err != nil {
			return nil, err
		}
	case ExportFormatJSON:
		fileName = fmt.Sprintf("companies_export_%s.json", timestamp)
		filePath = fmt.Sprintf("%s/%s", s.exportDir, fileName)
		if err := s.exportCompaniesToJSON(companies, filePath); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("ExportCompanies: failed to get file info: %w", err)
	}

	return &ExportResult{
		FileName:  fileName,
		Format:    format,
		Size:      fileInfo.Size(),
		Path:      filePath,
		CreatedAt: time.Now(),
	}, nil
}

func (s *ExportService) exportCompaniesToCSV(companies []domain.Company, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Slug", "Description", "BusinessType", "Locale", "Currency", "Timezone", "Active", "Status", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, company := range companies {
		row := []string{
			fmt.Sprintf("%d", company.ID),
			company.Name,
			company.Slug,
			company.Description,
			string(company.BusinessType),
			company.Locale,
			company.Currency,
			company.Timezone,
			fmt.Sprintf("%t", company.Active),
			"", // Status field - will be populated from domain
			company.CreatedAt.Format(time.RFC3339),
			company.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func (s *ExportService) exportCompaniesToJSON(companies []domain.Company, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(companies); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (s *ExportService) ExportUsers(ctx context.Context, format ExportFormat) (*ExportResult, error) {
	// Get all users
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("ExportUsers: failed to list users: %w", err)
	}

	// Ensure export directory exists
	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return nil, fmt.Errorf("ExportUsers: failed to create export directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var fileName, filePath string

	switch format {
	case ExportFormatCSV:
		fileName = fmt.Sprintf("users_export_%s.csv", timestamp)
		filePath = fmt.Sprintf("%s/%s", s.exportDir, fileName)
		if err := s.exportUsersToCSV(users, filePath); err != nil {
			return nil, err
		}
	case ExportFormatJSON:
		fileName = fmt.Sprintf("users_export_%s.json", timestamp)
		filePath = fmt.Sprintf("%s/%s", s.exportDir, fileName)
		if err := s.exportUsersToJSON(users, filePath); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("ExportUsers: failed to get file info: %w", err)
	}

	return &ExportResult{
		FileName:  fileName,
		Format:    format,
		Size:      fileInfo.Size(),
		Path:      filePath,
		CreatedAt: time.Now(),
	}, nil
}

func (s *ExportService) exportUsersToCSV(users []*domain.User, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Email", "CompanyID", "Role", "Active", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, user := range users {
		row := []string{
			fmt.Sprintf("%d", user.ID),
			user.Name,
			user.Email,
			fmt.Sprintf("%d", user.CompanyID),
			string(user.Role),
			fmt.Sprintf("%t", user.Active),
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func (s *ExportService) exportUsersToJSON(users []*domain.User, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(users); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
