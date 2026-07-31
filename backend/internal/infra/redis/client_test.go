package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConfig tests the Config struct
func TestConfig(t *testing.T) {
	cfg := Config{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 6379, cfg.Port)
	assert.Equal(t, 0, cfg.DB)
	assert.Equal(t, 10, cfg.PoolSize)
}

// TestHealthStatus tests the HealthStatus struct
func TestHealthStatus(t *testing.T) {
	status := HealthStatus{
		Healthy:    true,
		Latency:    10 * time.Millisecond,
		Connected:  true,
		DB:         0,
		ClientName: "horizongest",
	}

	assert.True(t, status.Healthy)
	assert.True(t, status.Connected)
	assert.Equal(t, 0, status.DB)
	assert.Equal(t, "horizongest", status.ClientName)
}

// TestNewClient_InvalidConfig tests NewClient with invalid configuration
func TestNewClient_InvalidConfig(t *testing.T) {
	cfg := Config{
		Host:         "invalid-host",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolTimeout:  1 * time.Second,
	}

	_, err := NewClient(cfg)
	assert.Error(t, err)
}

// TestGetConfig tests GetConfig method
func TestGetConfig(t *testing.T) {
	cfg := Config{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	// Create a mock client without actual connection
	client := &Client{
		Client: nil,
		config: cfg,
	}

	retrievedCfg := client.GetConfig()
	assert.Equal(t, cfg.Host, retrievedCfg.Host)
	assert.Equal(t, cfg.Port, retrievedCfg.Port)
	assert.Equal(t, cfg.DB, retrievedCfg.DB)
}

// TestClose tests Close method
func TestClose(t *testing.T) {
	client := &Client{
		Client: nil,
		config: Config{},
	}

	err := client.Close()
	assert.NoError(t, err)
}

// TestHealthCheck_NilClient tests HealthCheck with nil client
func TestHealthCheck_NilClient(t *testing.T) {
	client := &Client{
		Client: nil,
		config: Config{},
	}

	status, err := client.HealthCheck(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.Healthy)
	assert.False(t, status.Connected)
}
