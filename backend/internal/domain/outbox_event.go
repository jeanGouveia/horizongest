package domain

import "time"

// OutboxStatus represents the processing status of an outbox event
type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusCompleted OutboxStatus = "completed"
	OutboxStatusFailed    OutboxStatus = "failed"
)

// OutboxEvent represents a domain event stored in the outbox table
// This implements the Outbox Pattern for reliable event publishing
type OutboxEvent struct {
	ID            uint
	AggregateType string        // 'order', 'product', 'ingredient', etc.
	AggregateID   uint          // ID of the aggregate (ex: order_id)
	EventType     string        // 'order.created', 'product.updated', etc.
	EventVersion  string        // Version of the event schema (default: '1.0')
	Payload       string        // JSON payload of the event
	TenantID      uint          // company_id for multi-tenant isolation
	Status        OutboxStatus  // 'pending', 'processing', 'completed', 'failed'
	Priority      int           // 1=critical, 5=normal, 10=low priority
	Attempts      int           // Number of processing attempts
	AvailableAt   time.Time     // When the event becomes available for processing
	ProcessedAt   *time.Time    // When the event was successfully processed
	CreatedAt     time.Time
	LastError     *string       // Last error message (if any)
}

// IsValid checks if the outbox event has valid required fields
func (e *OutboxEvent) IsValid() bool {
	if e.AggregateType == "" {
		return false
	}
	if e.AggregateID == 0 {
		return false
	}
	if e.EventType == "" {
		return false
	}
	if e.EventVersion == "" {
		return false
	}
	if e.Payload == "" {
		return false
	}
	if e.TenantID == 0 {
		return false
	}
	if e.Status == "" {
		return false
	}
	if e.Priority < 1 || e.Priority > 10 {
		return false
	}
	return true
}
