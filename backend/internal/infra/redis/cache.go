package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rediscmd "github.com/redis/go-redis/v9"
)

const (
	cacheKeyPrefix = "cache"
)

// cacheKey generates a namespaced cache key
func cacheKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", Namespace, cacheKeyPrefix, key)
}

// Cache defines the cache interface
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error

	// Set stores a value in cache with optional TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// TTL returns the remaining time to live of a key
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Invalidate removes multiple keys from cache
	Invalidate(ctx context.Context, keys ...string) error

	// SetNX sets a value only if the key does not exist (atomic)
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *Client
}

// NewCache creates a new Redis cache
func NewCache(client *Client) Cache {
	return &RedisCache{client: client}
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, cacheKey(key)).Result()
	if err != nil {
		if err == rediscmd.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Set stores a value in cache with optional TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := c.client.Set(ctx, cacheKey(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, cacheKey(key)).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, cacheKey(key)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return result > 0, nil
}

// TTL returns the remaining time to live of a key
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, cacheKey(key)).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}
	return ttl, nil
}

// Invalidate removes multiple keys from cache
func (c *RedisCache) Invalidate(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	namespacedKeys := make([]string, len(keys))
	for i, key := range keys {
		namespacedKeys[i] = cacheKey(key)
	}

	if err := c.client.Del(ctx, namespacedKeys...).Err(); err != nil {
		return fmt.Errorf("failed to invalidate keys: %w", err)
	}
	return nil
}

// SetNX sets a value only if the key does not exist (atomic)
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	result, err := c.client.SetNX(ctx, cacheKey(key), data, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key %s with NX: %w", key, err)
	}

	return result, nil
}
