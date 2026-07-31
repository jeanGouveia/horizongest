package redis

import (
	"context"
	"fmt"
	"log"
	"time"
)

// FallbackCache wraps a Cache with graceful degradation
// If Redis fails, it logs the error and returns cache miss (allowing fallback to PostgreSQL)
type FallbackCache struct {
	cache Cache
}

// NewFallbackCache creates a new cache with graceful degradation
func NewFallbackCache(cache Cache) Cache {
	return &FallbackCache{cache: cache}
}

// Get retrieves a value from cache, returning error as cache miss on Redis failure
func (f *FallbackCache) Get(ctx context.Context, key string, dest interface{}) error {
	err := f.cache.Get(ctx, key, dest)
	if err != nil {
		// Log Redis error but treat as cache miss
		log.Printf("[RedisFallback] Cache get failed for key %s: %v (treating as cache miss)", key, err)
		return fmt.Errorf("key not found: %s", key)
	}
	return nil
}

// Set stores a value in cache, ignoring errors on Redis failure
func (f *FallbackCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := f.cache.Set(ctx, key, value, ttl)
	if err != nil {
		// Log Redis error but don't fail the operation
		log.Printf("[RedisFallback] Cache set failed for key %s: %v (cache disabled)", key, err)
		return nil
	}
	return nil
}

// Delete removes a value from cache, ignoring errors on Redis failure
func (f *FallbackCache) Delete(ctx context.Context, key string) error {
	err := f.cache.Delete(ctx, key)
	if err != nil {
		log.Printf("[RedisFallback] Cache delete failed for key %s: %v (ignoring)", key, err)
		return nil
	}
	return nil
}

// Exists checks if a key exists in cache, returning false on Redis failure
func (f *FallbackCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := f.cache.Exists(ctx, key)
	if err != nil {
		log.Printf("[RedisFallback] Cache exists check failed for key %s: %v (treating as not exists)", key, err)
		return false, nil
	}
	return exists, nil
}

// TTL returns the remaining time to live of a key, returning 0 on Redis failure
func (f *FallbackCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := f.cache.TTL(ctx, key)
	if err != nil {
		log.Printf("[RedisFallback] Cache TTL check failed for key %s: %v (returning 0)", key, err)
		return 0, nil
	}
	return ttl, nil
}

// Invalidate removes multiple keys from cache, ignoring errors on Redis failure
func (f *FallbackCache) Invalidate(ctx context.Context, keys ...string) error {
	err := f.cache.Invalidate(ctx, keys...)
	if err != nil {
		log.Printf("[RedisFallback] Cache invalidate failed for %d keys: %v (ignoring)", len(keys), err)
		return nil
	}
	return nil
}

// SetNX sets a value only if the key does not exist, returning false on Redis failure
func (f *FallbackCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	result, err := f.cache.SetNX(ctx, key, value, ttl)
	if err != nil {
		log.Printf("[RedisFallback] Cache SetNX failed for key %s: %v (returning false)", key, err)
		return false, nil
	}
	return result, nil
}

// FallbackLockManager wraps a LockManager with graceful degradation
// If Redis fails, it logs the error and returns success (allowing operation without lock)
type FallbackLockManager struct {
	lock LockManager
}

// NewFallbackLockManager creates a new lock manager with graceful degradation
func NewFallbackLockManager(lock LockManager) LockManager {
	return &FallbackLockManager{lock: lock}
}

// Acquire acquires a lock, returning success on Redis failure (degraded mode)
func (f *FallbackLockManager) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	acquired, err := f.lock.Acquire(ctx, key, ttl)
	if err != nil {
		log.Printf("[RedisFallback] Lock acquire failed for key %s: %v (proceeding without lock - degraded mode)", key, err)
		return true, nil // Allow operation to proceed without lock
	}
	return acquired, nil
}

// Release releases a lock, ignoring errors on Redis failure
func (f *FallbackLockManager) Release(ctx context.Context, key string) error {
	err := f.lock.Release(ctx, key)
	if err != nil {
		log.Printf("[RedisFallback] Lock release failed for key %s: %v (ignoring)", key, err)
		return nil
	}
	return nil
}

// TryAcquireWithRetry tries to acquire a lock with retry logic, returning success on Redis failure
func (f *FallbackLockManager) TryAcquireWithRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (bool, error) {
	acquired, err := f.lock.TryAcquireWithRetry(ctx, key, ttl, maxRetries, retryDelay)
	if err != nil {
		log.Printf("[RedisFallback] Lock acquire with retry failed for key %s: %v (proceeding without lock - degraded mode)", key, err)
		return true, nil // Allow operation to proceed without lock
	}
	return acquired, nil
}

// FallbackRateLimiter wraps a RateLimiter with graceful degradation
// If Redis fails, it logs the error and allows all requests (degraded mode)
type FallbackRateLimiter struct {
	limiter RateLimiter
}

// NewFallbackRateLimiter creates a new rate limiter with graceful degradation
func NewFallbackRateLimiter(limiter RateLimiter) RateLimiter {
	return &FallbackRateLimiter{limiter: limiter}
}

// Allow checks if a request is allowed, allowing all on Redis failure (degraded mode)
func (f *FallbackRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	allowed, err := f.limiter.Allow(ctx, key, limit, window)
	if err != nil {
		log.Printf("[RedisFallback] Rate limit check failed for key %s: %v (allowing request - degraded mode)", key, err)
		return true, nil // Allow request in degraded mode
	}
	return allowed, nil
}

// Reset resets the rate limit, ignoring errors on Redis failure
func (f *FallbackRateLimiter) Reset(ctx context.Context, key string) error {
	err := f.limiter.Reset(ctx, key)
	if err != nil {
		log.Printf("[RedisFallback] Rate limit reset failed for key %s: %v (ignoring)", key, err)
		return nil
	}
	return nil
}

// GetRemaining returns the remaining requests, returning limit on Redis failure
func (f *FallbackRateLimiter) GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	remaining, err := f.limiter.GetRemaining(ctx, key, limit, window)
	if err != nil {
		log.Printf("[RedisFallback] Rate limit get remaining failed for key %s: %v (returning limit - degraded mode)", key, err)
		return limit, nil
	}
	return remaining, nil
}

// FallbackIdempotencyChecker wraps an IdempotencyChecker with graceful degradation
// If Redis fails, it logs the error and allows processing (may have duplicates in degraded mode)
type FallbackIdempotencyChecker struct {
	checker IdempotencyChecker
}

// NewFallbackIdempotencyChecker creates a new idempotency checker with graceful degradation
func NewFallbackIdempotencyChecker(checker IdempotencyChecker) IdempotencyChecker {
	return &FallbackIdempotencyChecker{checker: checker}
}

// IsProcessed checks if an event has been processed, returning false on Redis failure
func (f *FallbackIdempotencyChecker) IsProcessed(ctx context.Context, eventID uint) (bool, error) {
	processed, err := f.checker.IsProcessed(ctx, eventID)
	if err != nil {
		log.Printf("[RedisFallback] Idempotency check failed for event %d: %v (allowing processing - degraded mode)", eventID, err)
		return false, nil // Allow processing (may have duplicates)
	}
	return processed, nil
}

// MarkProcessed marks an event as processed, ignoring errors on Redis failure
func (f *FallbackIdempotencyChecker) MarkProcessed(ctx context.Context, eventID uint) error {
	err := f.checker.MarkProcessed(ctx, eventID)
	if err != nil {
		log.Printf("[RedisFallback] Idempotency mark failed for event %d: %v (ignoring - degraded mode)", eventID, err)
		return nil
	}
	return nil
}

// IdempotencyChecker interface for fallback wrapper
type IdempotencyChecker interface {
	IsProcessed(ctx context.Context, eventID uint) (bool, error)
	MarkProcessed(ctx context.Context, eventID uint) error
}
