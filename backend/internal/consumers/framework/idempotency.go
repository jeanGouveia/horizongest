package framework

import (
	"context"
	"sync"
)

// IdempotencyChecker defines the interface for idempotency operations
type IdempotencyChecker interface {
	// IsProcessed checks if an event has already been processed
	IsProcessed(ctx context.Context, eventID uint) (bool, error)

	// MarkProcessed marks an event as processed
	MarkProcessed(ctx context.Context, eventID uint) error
}

// IdempotencyStore tracks processed event IDs to prevent duplicate processing
// This is a shared implementation used by all consumers (in-memory)
type IdempotencyStore struct {
	mu  sync.RWMutex
	ids map[uint]bool // event_id -> true
}

// NewIdempotencyStore creates a new in-memory idempotency store
func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		ids: make(map[uint]bool),
	}
}

// IsProcessed checks if an event has already been processed
func (s *IdempotencyStore) IsProcessed(ctx context.Context, eventID uint) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ids[eventID], nil
}

// MarkProcessed marks an event as processed
func (s *IdempotencyStore) MarkProcessed(ctx context.Context, eventID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[eventID] = true
	return nil
}

// Clear clears all processed event IDs (useful for testing)
func (s *IdempotencyStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids = make(map[uint]bool)
}
