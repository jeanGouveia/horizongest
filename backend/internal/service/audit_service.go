package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

// AuditService handles audit logging for critical operations
// FASE A.3: B17 - Centralized audit logging service
type AuditService struct {
	repo ports.AuditRepository
}

// NewAuditService creates a new audit service
func NewAuditService(repo ports.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID           uint      `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	UserID       uint      `json:"user_id"`
	CompanyID    uint      `json:"company_id"`
	Operation    string    `json:"operation"`
	Resource     string    `json:"resource"`
	ResourceID   string    `json:"resource_id,omitempty"`
	Before       string    `json:"before,omitempty"`
	After        string    `json:"after,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	CorrelationID string   `json:"correlation_id,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
}

// LogAudit logs an audit entry
func (s *AuditService) LogAudit(ctx context.Context, entry AuditEntry) error {
	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	// Convert to domain audit entry
	auditEntry := &domain.AuditLog{
		Timestamp:    entry.Timestamp,
		UserID:       entry.UserID,
		CompanyID:    entry.CompanyID,
		Operation:    entry.Operation,
		Resource:     entry.Resource,
		ResourceID:   entry.ResourceID,
		Before:       entry.Before,
		After:        entry.After,
		IPAddress:    entry.IPAddress,
		UserAgent:    entry.UserAgent,
		CorrelationID: entry.CorrelationID,
		RequestID:    entry.RequestID,
	}

	return s.repo.CreateAuditLog(ctx, auditEntry)
}

// LogUserAction logs a user action for audit
func (s *AuditService) LogUserAction(ctx context.Context, userID, companyID uint, operation, resource string, resourceID string, before, after interface{}, ip, userAgent, correlationID, requestID string) error {
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)

	entry := AuditEntry{
		UserID:       userID,
		CompanyID:    companyID,
		Operation:    operation,
		Resource:     resource,
		ResourceID:   resourceID,
		Before:       string(beforeJSON),
		After:        string(afterJSON),
		IPAddress:    ip,
		UserAgent:    userAgent,
		CorrelationID: correlationID,
		RequestID:    requestID,
	}

	return s.LogAudit(ctx, entry)
}

// GetAuditLogs retrieves audit logs for a company
func (s *AuditService) GetAuditLogs(ctx context.Context, companyID uint, limit int) ([]domain.AuditLog, error) {
	return s.repo.FindAuditLogsByCompany(ctx, companyID, limit)
}

// GetAuditLogsByUser retrieves audit logs for a specific user
func (s *AuditService) GetAuditLogsByUser(ctx context.Context, userID uint, limit int) ([]domain.AuditLog, error) {
	return s.repo.FindAuditLogsByUser(ctx, userID, limit)
}
