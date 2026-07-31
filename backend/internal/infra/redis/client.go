package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	rediscmd "github.com/redis/go-redis/v9"
)

const (
	// Namespace is the prefix for all Redis keys to avoid collisions
	Namespace = "hg"
)

// Config holds Redis configuration
type Config struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
}

// HealthStatus represents the health status of Redis
type HealthStatus struct {
	Healthy    bool
	Latency    time.Duration
	Connected  bool
	DB         int
	ClientName string
}

// Client wraps go-redis client with additional functionality
type Client struct {
	*rediscmd.Client
	config Config
	mu     sync.RWMutex
}

// NewClient creates a new Redis client with startup validation
func NewClient(cfg Config) (*Client, error) {
	rdb := rediscmd.NewClient(&rediscmd.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	})

	client := &Client{
		Client: rdb,
		config: cfg,
	}

	// Startup validation
	if err := client.StartupValidation(context.Background()); err != nil {
		return nil, fmt.Errorf("Redis startup validation failed: %w", err)
	}

	return client, nil
}

// StartupValidation performs comprehensive startup validation
func (c *Client) StartupValidation(ctx context.Context) error {
	// Test basic connectivity
	if err := c.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test write operation
	testKey := "healthcheck:startup"
	if err := c.Set(ctx, testKey, "ok", 10*time.Second).Err(); err != nil {
		return fmt.Errorf("write test failed: %w", err)
	}

	// Test read operation
	val, err := c.Get(ctx, testKey).Result()
	if err != nil {
		return fmt.Errorf("read test failed: %w", err)
	}
	if val != "ok" {
		return fmt.Errorf("read test returned unexpected value: %s", val)
	}

	// Cleanup test key
	if err := c.Del(ctx, testKey).Err(); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	return nil
}

// Close closes the Redis connection gracefully
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Client == nil {
		return nil
	}

	return c.Client.Close()
}

// HealthCheck checks if Redis is healthy and returns detailed status
func (c *Client) HealthCheck(ctx context.Context) (*HealthStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Client == nil {
		return &HealthStatus{Healthy: false, Connected: false}, nil
	}

	start := time.Now()
	err := c.Ping(ctx).Err()
	latency := time.Since(start)

	if err != nil {
		return &HealthStatus{
			Healthy:   false,
			Latency:   latency,
			Connected: false,
		}, nil
	}

	return &HealthStatus{
		Healthy:    true,
		Latency:    latency,
		Connected:  true,
		DB:         c.config.DB,
		ClientName: "horizongest",
	}, nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
