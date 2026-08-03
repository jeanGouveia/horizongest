package service

import (
	"context"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// MockAuditRepository is a mock implementation of AuditRepository for testing
type MockAuditRepository struct {
	auditLogs []domain.AuditLog
}

func (m *MockAuditRepository) CreateAuditLog(ctx context.Context, audit *domain.AuditLog) error {
	m.auditLogs = append(m.auditLogs, *audit)
	return nil
}

func (m *MockAuditRepository) FindAuditLogsByCompany(ctx context.Context, companyID uint, limit int) ([]domain.AuditLog, error) {
	var result []domain.AuditLog
	for _, log := range m.auditLogs {
		if log.CompanyID == companyID {
			result = append(result, log)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockAuditRepository) FindAuditLogsByUser(ctx context.Context, userID uint, limit int) ([]domain.AuditLog, error) {
	var result []domain.AuditLog
	for _, log := range m.auditLogs {
		if log.UserID == userID {
			result = append(result, log)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func TestAuditService_LogAudit(t *testing.T) {
	mockRepo := &MockAuditRepository{}
	service := NewAuditService(mockRepo)

	entry := AuditEntry{
		UserID:        1,
		CompanyID:     123,
		Operation:     "create",
		Resource:      "product",
		ResourceID:    "456",
		IPAddress:     "192.168.1.1",
		UserAgent:     "test-agent",
		CorrelationID: "corr-123",
		RequestID:     "req-456",
	}

	err := service.LogAudit(context.Background(), entry)
	if err != nil {
		t.Fatalf("failed to log audit: %v", err)
	}

	if len(mockRepo.auditLogs) != 1 {
		t.Errorf("expected 1 audit log, got %d", len(mockRepo.auditLogs))
	}

	log := mockRepo.auditLogs[0]
	if log.UserID != 1 {
		t.Errorf("expected user_id 1, got %d", log.UserID)
	}
	if log.CompanyID != 123 {
		t.Errorf("expected company_id 123, got %d", log.CompanyID)
	}
	if log.Operation != "create" {
		t.Errorf("expected operation create, got %s", log.Operation)
	}
	if log.Resource != "product" {
		t.Errorf("expected resource product, got %s", log.Resource)
	}
}

func TestAuditService_LogUserAction(t *testing.T) {
	mockRepo := &MockAuditRepository{}
	service := NewAuditService(mockRepo)

	before := map[string]interface{}{"name": "old"}
	after := map[string]interface{}{"name": "new"}

	err := service.LogUserAction(
		context.Background(),
		1, 123, "update", "product", "456",
		before, after,
		"192.168.1.1", "test-agent", "corr-123", "req-456",
	)

	if err != nil {
		t.Fatalf("failed to log user action: %v", err)
	}

	if len(mockRepo.auditLogs) != 1 {
		t.Errorf("expected 1 audit log, got %d", len(mockRepo.auditLogs))
	}

	log := mockRepo.auditLogs[0]
	if log.Operation != "update" {
		t.Errorf("expected operation update, got %s", log.Operation)
	}
	if log.Before == "" {
		t.Error("expected before to be set")
	}
	if log.After == "" {
		t.Error("expected after to be set")
	}
}

func TestAuditService_GetAuditLogsByCompany(t *testing.T) {
	mockRepo := &MockAuditRepository{}
	service := NewAuditService(mockRepo)

	// Add some audit logs
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    1,
		CompanyID: 123,
		Operation: "create",
		Resource:  "product",
		Timestamp: time.Now(),
	})
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    2,
		CompanyID: 123,
		Operation: "update",
		Resource:  "product",
		Timestamp: time.Now(),
	})
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    3,
		CompanyID: 456,
		Operation: "delete",
		Resource:  "product",
		Timestamp: time.Now(),
	})

	logs, err := service.GetAuditLogs(context.Background(), 123, 10)
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("expected 2 audit logs for company 123, got %d", len(logs))
	}
}

func TestAuditService_GetAuditLogsByUser(t *testing.T) {
	mockRepo := &MockAuditRepository{}
	service := NewAuditService(mockRepo)

	// Add some audit logs
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    1,
		CompanyID: 123,
		Operation: "create",
		Resource:  "product",
		Timestamp: time.Now(),
	})
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    1,
		CompanyID: 123,
		Operation: "update",
		Resource:  "product",
		Timestamp: time.Now(),
	})
	mockRepo.auditLogs = append(mockRepo.auditLogs, domain.AuditLog{
		UserID:    2,
		CompanyID: 123,
		Operation: "delete",
		Resource:  "product",
		Timestamp: time.Now(),
	})

	logs, err := service.GetAuditLogsByUser(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("expected 2 audit logs for user 1, got %d", len(logs))
	}
}
