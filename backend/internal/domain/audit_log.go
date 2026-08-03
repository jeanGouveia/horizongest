package domain

import "time"

// AuditLog represents an audit log entry for tracking critical operations
// FASE A.3: B17 - Audit log domain model
type AuditLog struct {
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
