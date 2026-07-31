package framework

import "time"

// Config holds configuration for a consumer
type Config struct {
	// Queue name to consume from
	Queue string

	// Consumer name for logging
	ConsumerName string

	// Retry configuration
	MaxRetries      int
	InitialRetryDelay time.Duration
	MaxRetryDelay    time.Duration
	RetryMultiplier  float64

	// Timeout configuration
	OperationTimeout time.Duration

	// Circuit breaker configuration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout  time.Duration

	// Dead letter configuration
	DeadLetterQueue string
	MaxRetryAttempts int

	// Metrics configuration
	EnableMetrics bool
	MetricsPrefix string
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		MaxRetries:           3,
		InitialRetryDelay:    1 * time.Second,
		MaxRetryDelay:        30 * time.Second,
		RetryMultiplier:      2.0,
		OperationTimeout:     30 * time.Second,
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,
		MaxRetryAttempts:     3,
		EnableMetrics:       true,
		MetricsPrefix:       "consumer",
	}
}
