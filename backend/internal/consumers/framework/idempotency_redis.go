package framework

import (
	"context"
	"fmt"
	"time"

	rediscache "github.com/jeanGouveia/horizongest/backend/internal/infra/redis"
)

const (
	idempotencyKeyPrefix = "idempotency"
)

// idempotencyKey generates a namespaced idempotency key
func idempotencyKey(eventID uint) string {
	return fmt.Sprintf("%s:%s:event:%d", rediscache.Namespace, idempotencyKeyPrefix, eventID)
}

// RedisIdempotencyStore tracks processed event IDs using Redis to prevent duplicate processing
// This is a shared implementation used by all consumers with persistent storage
type RedisIdempotencyStore struct {
	cache rediscache.Cache
	ttl   time.Duration
}

// NewRedisIdempotencyStore creates a new Redis-based idempotency store
func NewRedisIdempotencyStore(cache rediscache.Cache, ttl time.Duration) *RedisIdempotencyStore {
	return &RedisIdempotencyStore{
		cache: cache,
		ttl:   ttl,
	}
}

// IsProcessed checks if an event has already been processed
func (s *RedisIdempotencyStore) IsProcessed(ctx context.Context, eventID uint) (bool, error) {
	key := idempotencyKey(eventID)
	exists, err := s.cache.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check if event %d is processed: %w", eventID, err)
	}
	return exists, nil
}

// MarkProcessed marks an event as processed (atomic operation using SetNX)
func (s *RedisIdempotencyStore) MarkProcessed(ctx context.Context, eventID uint) error {
	key := idempotencyKey(eventID)

	// Use SetNX for atomic operation - only sets if key doesn't exist
	alreadyProcessed, err := s.cache.SetNX(ctx, key, true, s.ttl)
	if err != nil {
		return fmt.Errorf("failed to mark event %d as processed: %w", eventID, err)
	}

	// If alreadyProcessed is false, the event was already marked
	if !alreadyProcessed {
		return nil // Not an error, just already processed
	}

	return nil
}

// Clear clears all processed event IDs for a specific prefix
// This is useful for testing or manual cleanup
func (s *RedisIdempotencyStore) Clear(ctx context.Context) error {
	// Note: This would require scanning keys, which is expensive
	// For production, consider using a different approach or manual key management
	return fmt.Errorf("clear operation not implemented for Redis idempotency store")
}
