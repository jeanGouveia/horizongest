package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	rateLimitKeyPrefix = "ratelimit"
)

// rateLimitKey generates a namespaced rate limit key
func rateLimitKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", Namespace, rateLimitKeyPrefix, key)
}

// RateLimiter defines the rate limiting interface
type RateLimiter interface {
	// Allow checks if a request is allowed based on the rate limit
	// Returns true if allowed, false if rate limit exceeded
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)

	// Reset resets the rate limit for a key
	Reset(ctx context.Context, key string) error

	// GetRemaining returns the remaining requests allowed in the current window
	GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error)
}

// RedisRateLimiter implements RateLimiter using Redis with sliding window algorithm
type RedisRateLimiter struct {
	client *Client
}

// NewRateLimiter creates a new Redis rate limiter
func NewRateLimiter(client *Client) RateLimiter {
	return &RedisRateLimiter{client: client}
}

// Allow checks if a request is allowed based on the rate limit
// Uses a simple token bucket approach with Redis INCR and EXPIRE
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Use INCR to atomically increment the counter
	current, err := r.client.Incr(ctx, rateLimitKey(key)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit counter for key %s: %w", key, err)
	}

	// If this is the first request (current == 1), set the expiration
	if current == 1 {
		if err := r.client.Expire(ctx, rateLimitKey(key), window).Err(); err != nil {
			return false, fmt.Errorf("failed to set expiration for rate limit key %s: %w", key, err)
		}
	}

	// Check if the limit has been exceeded
	return current <= int64(limit), nil
}

// Reset resets the rate limit for a key
func (r *RedisRateLimiter) Reset(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, rateLimitKey(key)).Err(); err != nil {
		return fmt.Errorf("failed to reset rate limit for key %s: %w", key, err)
	}
	return nil
}

// GetRemaining returns the remaining requests allowed in the current window
func (r *RedisRateLimiter) GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	current, err := r.client.Get(ctx, rateLimitKey(key)).Result()
	if err != nil {
		// Key doesn't exist, so all requests are available
		return limit, nil
	}

	var count int64
	if _, err := fmt.Sscanf(current, "%d", &count); err != nil {
		return 0, fmt.Errorf("failed to parse rate limit counter for key %s: %w", key, err)
	}

	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}
