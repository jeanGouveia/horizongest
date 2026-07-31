package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/pg"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

// ─── GORM Model ────────────────────────────────────────────────────────────

type GormOutboxEvent struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"`
	AggregateType string     `gorm:"not null;size:100;index:idx_outbox_aggregate,priority:1"`
	AggregateID   uint       `gorm:"not null;index:idx_outbox_aggregate,priority:2"`
	EventType     string     `gorm:"not null;size:100"`
	EventVersion  string     `gorm:"not null;size:20;default:'1.0'"`
	Payload       string     `gorm:"not null;type:text"`
	TenantID      uint       `gorm:"not null;index:idx_outbox_tenant_status,priority:1"`
	Status        string     `gorm:"not null;size:20;default:'pending';index:idx_outbox_tenant_status,priority:2"`
	Priority      int        `gorm:"not null;default:5;index:idx_outbox_priority,priority:1"`
	Attempts      int        `gorm:"not null;default:0"`
	AvailableAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_outbox_available_at;index:idx_outbox_priority,priority:2"`
	ProcessedAt   *time.Time `gorm:"index:idx_outbox_processed_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	LastError     *string    `gorm:"type:text"`
}

func (GormOutboxEvent) TableName() string { return "outbox_events" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.OutboxRepository = (*GormOutboxRepository)(nil)

type GormOutboxRepository struct {
	db *gorm.DB
}

func NewGormOutboxRepository(db *gorm.DB) *GormOutboxRepository {
	return &GormOutboxRepository{db: db}
}

// getDB returns the transaction if provided, otherwise returns the default DB
func (r *GormOutboxRepository) getDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

// applyTenantFilter applies tenant filtering for outbox events (uses tenant_id instead of company_id)
func (r *GormOutboxRepository) applyTenantFilter(ctx context.Context, db *gorm.DB) *gorm.DB {
	tenantCtx, ok := domain.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query as-is
		return db
	}
	return db.Where("tenant_id = ?", tenantCtx.GetCompanyID())
}

// applyTenantFilterWithID applies tenant filtering for outbox events by ID (uses tenant_id instead of company_id)
func (r *GormOutboxRepository) applyTenantFilterWithID(ctx context.Context, db *gorm.DB, id uint) *gorm.DB {
	tenantCtx, ok := domain.GetTenantContextFromContext(ctx)
	if !ok {
		// No tenant context, return query with only ID
		return db.Where("id = ?", id)
	}
	return db.Where("id = ? AND tenant_id = ?", id, tenantCtx.GetCompanyID())
}

// Create stores a new outbox event
// If tx is provided, the operation is executed within that transaction
func (r *GormOutboxRepository) Create(ctx context.Context, event *domain.OutboxEvent, tx *gorm.DB) error {
	// Auto-fill TenantID from tenant context if not set
	if event.TenantID == 0 {
		companyID, err := GetCompanyIDFromContext(ctx)
		if err != nil {
			return fmt.Errorf("OutboxRepository.Create: %w", err)
		}
		event.TenantID = companyID
	}

	// Validate event before creation
	if !event.IsValid() {
		return fmt.Errorf("OutboxRepository.Create: invalid event data")
	}

	gEvent := outboxDomainToGorm(event)
	db := r.getDB(ctx, tx)

	// Handle unique constraint violation for idempotency
	err := db.Create(&gEvent).Error
	if err != nil {
		if pg.IsUniqueViolation(err) {
			return fmt.Errorf("OutboxRepository.Create: event already exists for this aggregate (aggregate_type=%s, aggregate_id=%d, event_type=%s)",
				event.AggregateType, event.AggregateID, event.EventType)
		}
		return fmt.Errorf("OutboxRepository.Create: %w", err)
	}

	event.ID = gEvent.ID
	event.CreatedAt = gEvent.CreatedAt
	return nil
}

// FindByID retrieves an outbox event by its ID
func (r *GormOutboxRepository) FindByID(ctx context.Context, id uint) (*domain.OutboxEvent, error) {
	var gEvent GormOutboxEvent
	query := r.applyTenantFilterWithID(ctx, r.db, id)
	err := query.First(&gEvent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("OutboxRepository.FindByID: %w", err)
	}

	return outboxGormToDomain(&gEvent), nil
}

// FindPendingEvents retrieves events that are pending processing
// Ordered by priority (critical first) and available_at (oldest first)
func (r *GormOutboxRepository) FindPendingEvents(ctx context.Context, tenantID uint, limit int) ([]*domain.OutboxEvent, error) {
	var gEvents []GormOutboxEvent

	// Use explicit tenant_id filter for this method since tenantID is passed as parameter
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ? AND available_at <= ?", tenantID, string(domain.OutboxStatusPending), time.Now()).
		Order("priority ASC, available_at ASC").
		Limit(limit)

	err := query.Find(&gEvents).Error
	if err != nil {
		return nil, fmt.Errorf("OutboxRepository.FindPendingEvents: %w", err)
	}

	events := make([]*domain.OutboxEvent, len(gEvents))
	for i, g := range gEvents {
		events[i] = outboxGormToDomain(&g)
	}

	return events, nil
}

// UpdateStatus changes the status of an outbox event
func (r *GormOutboxRepository) UpdateStatus(ctx context.Context, id uint, status domain.OutboxStatus) error {
	query := r.applyTenantFilterWithID(ctx, r.db, id)
	err := query.Model(&GormOutboxEvent{}).
		Update("status", string(status)).Error

	if err != nil {
		return fmt.Errorf("OutboxRepository.UpdateStatus: %w", err)
	}

	return nil
}

// UpdateStatusWithOptimisticLock changes the status of an outbox event with optimistic locking
// Returns true if the update was successful (event was in the expected status)
// Returns false if the event was already in a different status (concurrent modification)
func (r *GormOutboxRepository) UpdateStatusWithOptimisticLock(ctx context.Context, id uint, expectedStatus, newStatus domain.OutboxStatus) (bool, error) {
	query := r.applyTenantFilterWithID(ctx, r.db, id)
	result := query.Model(&GormOutboxEvent{}).
		Where("status = ?", string(expectedStatus)).
		Update("status", string(newStatus))

	if result.Error != nil {
		return false, fmt.Errorf("OutboxRepository.UpdateStatusWithOptimisticLock: %w", result.Error)
	}

	// If no rows were affected, the event was not in the expected status
	return result.RowsAffected > 0, nil
}

// IncrementAttempts increments the attempt counter and sets the last error and available_at for retry
func (r *GormOutboxRepository) IncrementAttempts(ctx context.Context, id uint, errorMsg string, availableAt time.Time) error {
	query := r.applyTenantFilterWithID(ctx, r.db, id)
	err := query.Model(&GormOutboxEvent{}).
		Updates(map[string]interface{}{
			"attempts":     gorm.Expr("attempts + 1"),
			"last_error":   errorMsg,
			"status":       string(domain.OutboxStatusPending),
			"available_at": availableAt,
		}).Error

	if err != nil {
		return fmt.Errorf("OutboxRepository.IncrementAttempts: %w", err)
	}

	return nil
}

// MarkAsCompleted marks an event as successfully processed
func (r *GormOutboxRepository) MarkAsCompleted(ctx context.Context, id uint) error {
	now := time.Now()
	query := r.applyTenantFilterWithID(ctx, r.db, id)
	err := query.Model(&GormOutboxEvent{}).
		Updates(map[string]interface{}{
			"status":       string(domain.OutboxStatusCompleted),
			"processed_at": &now,
		}).Error

	if err != nil {
		return fmt.Errorf("OutboxRepository.MarkAsCompleted: %w", err)
	}

	return nil
}

// FindByAggregate retrieves all events for a specific aggregate
func (r *GormOutboxRepository) FindByAggregate(ctx context.Context, aggregateType string, aggregateID uint) ([]*domain.OutboxEvent, error) {
	var gEvents []GormOutboxEvent

	query := r.applyTenantFilter(ctx, r.db)
	err := query.WithContext(ctx).
		Where("aggregate_type = ? AND aggregate_id = ?", aggregateType, aggregateID).
		Order("created_at DESC").
		Find(&gEvents).Error

	if err != nil {
		return nil, fmt.Errorf("OutboxRepository.FindByAggregate: %w", err)
	}

	events := make([]*domain.OutboxEvent, len(gEvents))
	for i, g := range gEvents {
		events[i] = outboxGormToDomain(&g)
	}

	return events, nil
}

// DeleteOldCompletedEvents removes completed events older than a specified duration
func (r *GormOutboxRepository) DeleteOldCompletedEvents(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)

	query := r.applyTenantFilter(ctx, r.db)
	result := query.WithContext(ctx).
		Where("status = ? AND processed_at < ?", string(domain.OutboxStatusCompleted), cutoff).
		Delete(&GormOutboxEvent{})

	if result.Error != nil {
		return 0, fmt.Errorf("OutboxRepository.DeleteOldCompletedEvents: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ─── Mappers ──────────────────────────────────────────────────────────────────

func outboxDomainToGorm(d *domain.OutboxEvent) GormOutboxEvent {
	return GormOutboxEvent{
		ID:            d.ID,
		AggregateType: d.AggregateType,
		AggregateID:   d.AggregateID,
		EventType:     d.EventType,
		EventVersion:  d.EventVersion,
		Payload:       d.Payload,
		TenantID:      d.TenantID,
		Status:        string(d.Status),
		Priority:      d.Priority,
		Attempts:      d.Attempts,
		AvailableAt:   d.AvailableAt,
		ProcessedAt:   d.ProcessedAt,
		CreatedAt:     d.CreatedAt,
		LastError:     d.LastError,
	}
}

func outboxGormToDomain(g *GormOutboxEvent) *domain.OutboxEvent {
	return &domain.OutboxEvent{
		ID:            g.ID,
		AggregateType: g.AggregateType,
		AggregateID:   g.AggregateID,
		EventType:     g.EventType,
		EventVersion:  g.EventVersion,
		Payload:       g.Payload,
		TenantID:      g.TenantID,
		Status:        domain.OutboxStatus(g.Status),
		Priority:      g.Priority,
		Attempts:      g.Attempts,
		AvailableAt:   g.AvailableAt,
		ProcessedAt:   g.ProcessedAt,
		CreatedAt:     g.CreatedAt,
		LastError:     g.LastError,
	}
}
