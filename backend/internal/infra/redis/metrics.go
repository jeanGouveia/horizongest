package redis

import (
	"context"
	"time"
)

// RedisMetrics defines the metrics interface for Redis operations
type RedisMetrics interface {
	// Cache operations
	IncrementCacheHit()
	IncrementCacheMiss()
	RecordCacheOperation(operation string, duration time.Duration)
	
	// Lock operations
	IncrementLockAcquired()
	IncrementLockReleased()
	IncrementLockFailed()
	RecordLockOperation(operation string, duration time.Duration)
	
	// Rate limit operations
	IncrementRateLimitAllowed()
	IncrementRateLimitDenied()
	RecordRateLimitCheck(duration time.Duration)
	
	// Idempotency operations
	IncrementIdempotencyHit()
	IncrementIdempotencyMiss()
	RecordIdempotencyCheck(duration time.Duration)
	
	// Health check operations
	RecordHealthCheck(duration time.Duration, healthy bool)
}

// NoOpRedisMetrics is a no-op implementation for when metrics are not needed
type NoOpRedisMetrics struct{}

func (m *NoOpRedisMetrics) IncrementCacheHit() {}
func (m *NoOpRedisMetrics) IncrementCacheMiss() {}
func (m *NoOpRedisMetrics) RecordCacheOperation(operation string, duration time.Duration) {}
func (m *NoOpRedisMetrics) IncrementLockAcquired() {}
func (m *NoOpRedisMetrics) IncrementLockReleased() {}
func (m *NoOpRedisMetrics) IncrementLockFailed() {}
func (m *NoOpRedisMetrics) RecordLockOperation(operation string, duration time.Duration) {}
func (m *NoOpRedisMetrics) IncrementRateLimitAllowed() {}
func (m *NoOpRedisMetrics) IncrementRateLimitDenied() {}
func (m *NoOpRedisMetrics) RecordRateLimitCheck(duration time.Duration) {}
func (m *NoOpRedisMetrics) IncrementIdempotencyHit() {}
func (m *NoOpRedisMetrics) IncrementIdempotencyMiss() {}
func (m *NoOpRedisMetrics) RecordIdempotencyCheck(duration time.Duration) {}
func (m *NoOpRedisMetrics) RecordHealthCheck(duration time.Duration, healthy bool) {}

// RedisCacheWithMetrics wraps a Cache with metrics
type RedisCacheWithMetrics struct {
	cache   Cache
	metrics RedisMetrics
}

// NewRedisCacheWithMetrics creates a new Redis cache with metrics
func NewRedisCacheWithMetrics(cache Cache, metrics RedisMetrics) Cache {
	if metrics == nil {
		return cache
	}
	return &RedisCacheWithMetrics{
		cache:   cache,
		metrics: metrics,
	}
}

func (c *RedisCacheWithMetrics) Get(ctx context.Context, key string, dest interface{}) error {
	start := time.Now()
	err := c.cache.Get(ctx, key, dest)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("get", duration)
	if err != nil {
		c.metrics.IncrementCacheMiss()
	} else {
		c.metrics.IncrementCacheHit()
	}
	
	return err
}

func (c *RedisCacheWithMetrics) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	err := c.cache.Set(ctx, key, value, ttl)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("set", duration)
	return err
}

func (c *RedisCacheWithMetrics) Delete(ctx context.Context, key string) error {
	start := time.Now()
	err := c.cache.Delete(ctx, key)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("delete", duration)
	return err
}

func (c *RedisCacheWithMetrics) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	exists, err := c.cache.Exists(ctx, key)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("exists", duration)
	if err == nil && exists {
		c.metrics.IncrementCacheHit()
	} else {
		c.metrics.IncrementCacheMiss()
	}
	
	return exists, err
}

func (c *RedisCacheWithMetrics) TTL(ctx context.Context, key string) (time.Duration, error) {
	start := time.Now()
	ttl, err := c.cache.TTL(ctx, key)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("ttl", duration)
	return ttl, err
}

func (c *RedisCacheWithMetrics) Invalidate(ctx context.Context, keys ...string) error {
	start := time.Now()
	err := c.cache.Invalidate(ctx, keys...)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("invalidate", duration)
	return err
}

func (c *RedisCacheWithMetrics) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	start := time.Now()
	result, err := c.cache.SetNX(ctx, key, value, ttl)
	duration := time.Since(start)
	
	c.metrics.RecordCacheOperation("setnx", duration)
	return result, err
}

// RedisLockManagerWithMetrics wraps a LockManager with metrics
type RedisLockManagerWithMetrics struct {
	lock    LockManager
	metrics RedisMetrics
}

// NewRedisLockManagerWithMetrics creates a new Redis lock manager with metrics
func NewRedisLockManagerWithMetrics(lock LockManager, metrics RedisMetrics) LockManager {
	if metrics == nil {
		return lock
	}
	return &RedisLockManagerWithMetrics{
		lock:    lock,
		metrics: metrics,
	}
}

func (l *RedisLockManagerWithMetrics) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	start := time.Now()
	acquired, err := l.lock.Acquire(ctx, key, ttl)
	duration := time.Since(start)
	
	l.metrics.RecordLockOperation("acquire", duration)
	if err != nil {
		l.metrics.IncrementLockFailed()
	} else if acquired {
		l.metrics.IncrementLockAcquired()
	} else {
		l.metrics.IncrementLockFailed()
	}
	
	return acquired, err
}

func (l *RedisLockManagerWithMetrics) Release(ctx context.Context, key string) error {
	start := time.Now()
	err := l.lock.Release(ctx, key)
	duration := time.Since(start)
	
	l.metrics.RecordLockOperation("release", duration)
	if err == nil {
		l.metrics.IncrementLockReleased()
	} else {
		l.metrics.IncrementLockFailed()
	}
	
	return err
}

func (l *RedisLockManagerWithMetrics) TryAcquireWithRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (bool, error) {
	start := time.Now()
	acquired, err := l.lock.TryAcquireWithRetry(ctx, key, ttl, maxRetries, retryDelay)
	duration := time.Since(start)
	
	l.metrics.RecordLockOperation("try_acquire_with_retry", duration)
	if err != nil {
		l.metrics.IncrementLockFailed()
	} else if acquired {
		l.metrics.IncrementLockAcquired()
	} else {
		l.metrics.IncrementLockFailed()
	}
	
	return acquired, err
}

// RedisRateLimiterWithMetrics wraps a RateLimiter with metrics
type RedisRateLimiterWithMetrics struct {
	limiter RateLimiter
	metrics RedisMetrics
}

// NewRedisRateLimiterWithMetrics creates a new Redis rate limiter with metrics
func NewRedisRateLimiterWithMetrics(limiter RateLimiter, metrics RedisMetrics) RateLimiter {
	if metrics == nil {
		return limiter
	}
	return &RedisRateLimiterWithMetrics{
		limiter: limiter,
		metrics: metrics,
	}
}

func (r *RedisRateLimiterWithMetrics) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	start := time.Now()
	allowed, err := r.limiter.Allow(ctx, key, limit, window)
	duration := time.Since(start)
	
	r.metrics.RecordRateLimitCheck(duration)
	if err != nil {
		return false, err
	}
	
	if allowed {
		r.metrics.IncrementRateLimitAllowed()
	} else {
		r.metrics.IncrementRateLimitDenied()
	}
	
	return allowed, nil
}

func (r *RedisRateLimiterWithMetrics) Reset(ctx context.Context, key string) error {
	return r.limiter.Reset(ctx, key)
}

func (r *RedisRateLimiterWithMetrics) GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	return r.limiter.GetRemaining(ctx, key, limit, window)
}
