package email

import (
	"sync"
)

// IdempotencyStore tracks processed event IDs to prevent duplicate processing
// For production, this should be replaced with a persistent store (database, Redis, etc.)
type IdempotencyStore struct {
	mu    sync.RWMutex
	ids   map[uint]bool // event_id -> true
}

// NewIdempotencyStore creates a new in-memory idempotency store
func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		ids: make(map[uint]bool),
	}
}

// IsProcessed checks if an event has already been processed
func (s *IdempotencyStore) IsProcessed(eventID uint) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ids[eventID]
}

// MarkProcessed marks an event as processed
func (s *IdempotencyStore) MarkProcessed(eventID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[eventID] = true
}

// Clear clears all processed event IDs (useful for testing)
func (s *IdempotencyStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids = make(map[uint]bool)
}
