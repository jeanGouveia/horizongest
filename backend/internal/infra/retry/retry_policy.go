package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

// RetryPolicy defines retry behavior
// FASE A.4: Retry Policies - Consistent retry logic for external services
type RetryPolicy struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxJitter       time.Duration
	Retryable      func(error) bool
}

// RetryableFunc represents a function that can be retried
type RetryableFunc func(ctx context.Context) error

// NewRetryPolicy creates a new retry policy with sensible defaults
func NewRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		MaxJitter:       100 * time.Millisecond,
		Retryable:      DefaultRetryable,
	}
}

// DefaultRetryable determines if an error is retryable
// FASE A.4: Default retryable logic
func DefaultRetryable(err error) bool {
	if err == nil {
		return false
	}
	
	// Add specific error types that are retryable
	// For now, retry all errors except context cancellation
	return err != context.Canceled && err != context.DeadlineExceeded
}

// Execute executes a function with retry policy
// FASE A.4: Execute function with retry logic
func (p *RetryPolicy) Execute(ctx context.Context, name string, fn RetryableFunc) error {
	logger := util.GetLogger()
	
	var lastErr error
	interval := p.InitialInterval
	
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if attempt > 0 {
			logger.Info(fmt.Sprintf("Retrying %s (attempt %d/%d)", name, attempt, p.MaxAttempts), map[string]interface{}{
				"attempt": attempt,
				"max_attempts": p.MaxAttempts,
				"error": lastErr.Error(),
			})
			
			// Add jitter to prevent thundering herd
			jitter := time.Duration(float64(p.MaxJitter) * (2.0*randFloat64() - 1.0))
			if jitter < 0 {
				jitter = -jitter
			}
			
			select {
			case <-time.After(interval + jitter):
			case <-ctx.Done():
				return ctx.Err()
			}
			
			// Exponential backoff
			interval = time.Duration(float64(interval) * p.Multiplier)
			if interval > p.MaxInterval {
				interval = p.MaxInterval
			}
		}
		
		err := fn(ctx)
		if err == nil {
			if attempt > 0 {
				logger.Info(fmt.Sprintf("%s succeeded after %d attempts", name, attempt+1), map[string]interface{}{
					"attempts": attempt + 1,
				})
			}
			return nil
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !p.Retryable(err) {
			logger.Error(fmt.Sprintf("%s failed with non-retryable error", name), map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}
	}
	
	logger.Error(fmt.Sprintf("%s failed after %d attempts", name, p.MaxAttempts), map[string]interface{}{
		"error": lastErr.Error(),
		"attempts": p.MaxAttempts,
	})
	
	return fmt.Errorf("%s failed after %d attempts: %w", name, p.MaxAttempts, lastErr)
}

// ExecuteWithBackoff executes a function with custom backoff
// FASE A.4: Execute with custom backoff strategy
func (p *RetryPolicy) ExecuteWithBackoff(ctx context.Context, name string, fn RetryableFunc, backoff func(attempt int) time.Duration) error {
	logger := util.GetLogger()
	
	var lastErr error
	
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if attempt > 0 {
			backoffDuration := backoff(attempt)
			logger.Info(fmt.Sprintf("Retrying %s (attempt %d/%d)", name, attempt, p.MaxAttempts), map[string]interface{}{
				"attempt": attempt,
				"max_attempts": p.MaxAttempts,
				"backoff": backoffDuration.String(),
				"error": lastErr.Error(),
			})
			
			select {
			case <-time.After(backoffDuration):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		
		err := fn(ctx)
		if err == nil {
			if attempt > 0 {
				logger.Info(fmt.Sprintf("%s succeeded after %d attempts", name, attempt+1), map[string]interface{}{
					"attempts": attempt + 1,
				})
			}
			return nil
		}
		
		lastErr = err
		
		if !p.Retryable(err) {
			return err
		}
	}
	
	return fmt.Errorf("%s failed after %d attempts: %w", name, p.MaxAttempts, lastErr)
}

// randFloat64 returns a random float64 between 0 and 1
func randFloat64() float64 {
	return float64(time.Now().UnixNano()) / float64(1<<63-1)
}

// Predefined retry policies for different scenarios
var (
	// RedisRetryPolicy - Fast retries for Redis
	RedisRetryPolicy = &RetryPolicy{
		MaxAttempts:     5,
		InitialInterval: 50 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      1.5,
		MaxJitter:       50 * time.Millisecond,
		Retryable:      DefaultRetryable,
	}
	
	// RabbitMQRetryPolicy - Moderate retries for RabbitMQ
	RabbitMQRetryPolicy = &RetryPolicy{
		MaxAttempts:     3,
		InitialInterval: 200 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		MaxJitter:       100 * time.Millisecond,
		Retryable:      DefaultRetryable,
	}
	
	// StorageRetryPolicy - Conservative retries for storage
	StorageRetryPolicy = &RetryPolicy{
		MaxAttempts:     2,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		MaxJitter:       50 * time.Millisecond,
		Retryable:      DefaultRetryable,
	}
	
	// DatabaseRetryPolicy - No retries for database (use connection pool)
	DatabaseRetryPolicy = &RetryPolicy{
		MaxAttempts:     1,
		InitialInterval: 0,
		MaxInterval:     0,
		Multiplier:      0,
		MaxJitter:       0,
		Retryable:      func(error) bool { return false },
	}
)
