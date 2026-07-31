package ports

import (
	"context"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"gorm.io/gorm"
)

// OutboxRepository defines the interface for outbox event persistence
// This repository is used to store domain events reliably within the same transaction
// as the main business operation, implementing the Outbox Pattern
type OutboxRepository interface {
	// Create stores a new outbox event
	// If tx is provided, the operation is executed within that transaction
	// This ensures atomicity between the business operation and event creation
	Create(ctx context.Context, event *domain.OutboxEvent, tx *gorm.DB) error

	// FindByID retrieves an outbox event by its ID
	FindByID(ctx context.Context, id uint) (*domain.OutboxEvent, error)

	// FindPendingEvents retrieves events that are pending processing
	// Ordered by priority (critical first) and available_at (oldest first)
	// Limit specifies the maximum number of events to retrieve
	FindPendingEvents(ctx context.Context, tenantID uint, limit int) ([]*domain.OutboxEvent, error)

	// UpdateStatus changes the status of an outbox event
	// Used to transition events between pending -> processing -> completed/failed
	UpdateStatus(ctx context.Context, id uint, status domain.OutboxStatus) error

	// UpdateStatusWithOptimisticLock changes the status of an outbox event with optimistic locking
	// Returns true if the update was successful (event was in the expected status)
	// Returns false if the event was already in a different status (concurrent modification)
	// Used to prevent race conditions when multiple dispatchers process events concurrently
	UpdateStatusWithOptimisticLock(ctx context.Context, id uint, expectedStatus, newStatus domain.OutboxStatus) (bool, error)

	// IncrementAttempts increments the attempt counter and sets the last error and available_at for retry
	// Used when an event processing fails and needs to be retried
	IncrementAttempts(ctx context.Context, id uint, error string, availableAt time.Time) error

	// MarkAsCompleted marks an event as successfully processed
	// Sets status to completed and records the processed_at timestamp
	MarkAsCompleted(ctx context.Context, id uint) error

	// FindByAggregate retrieves all events for a specific aggregate
	// Useful for auditing and debugging
	FindByAggregate(ctx context.Context, aggregateType string, aggregateID uint) ([]*domain.OutboxEvent, error)

	// DeleteOldCompletedEvents removes completed events older than a specified duration
	// Used for cleanup/maintenance to prevent table bloat
	DeleteOldCompletedEvents(ctx context.Context, olderThan time.Duration) (int64, error)
}
