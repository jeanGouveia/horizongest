package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	lockKeyPrefix = "lock"
)

// lockKey generates a namespaced lock key
func lockKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", Namespace, lockKeyPrefix, key)
}

// LockManager defines the distributed lock interface
type LockManager interface {
	// Acquire acquires a lock with the given key and TTL
	// Returns true if lock was acquired, false otherwise
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Release releases a lock with the given key
	Release(ctx context.Context, key string) error

	// TryAcquireWithRetry tries to acquire a lock with retry logic
	TryAcquireWithRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (bool, error)
}

// RedisLockManager implements LockManager using Redis
type RedisLockManager struct {
	client *Client
}

// NewLockManager creates a new Redis lock manager
func NewLockManager(client *Client) LockManager {
	return &RedisLockManager{client: client}
}

// Acquire acquires a lock with the given key and TTL
func (l *RedisLockManager) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use SET NX EX for atomic lock acquisition
	result, err := l.client.SetNX(ctx, lockKey(key), "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock for key %s: %w", key, err)
	}
	return result, nil
}

// Release releases a lock with the given key
func (l *RedisLockManager) Release(ctx context.Context, key string) error {
	if err := l.client.Del(ctx, lockKey(key)).Err(); err != nil {
		return fmt.Errorf("failed to release lock for key %s: %w", key, err)
	}
	return nil
}

// TryAcquireWithRetry tries to acquire a lock with retry logic
func (l *RedisLockManager) TryAcquireWithRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		acquired, err := l.Acquire(ctx, key, ttl)
		if err != nil {
			return false, err
		}
		if acquired {
			return true, nil
		}

		// Wait before retrying
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryDelay):
			continue
		}
	}

	return false, nil
}
