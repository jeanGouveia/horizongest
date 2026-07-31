package framework

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestIdempotencyStore tests the idempotency store
func TestIdempotencyStore(t *testing.T) {
	ctx := context.Background()
	store := NewIdempotencyStore()

	// Initially not processed
	processed, err := store.IsProcessed(ctx, 1)
	assert.NoError(t, err)
	assert.False(t, processed)

	// Mark as processed
	err = store.MarkProcessed(ctx, 1)
	assert.NoError(t, err)

	processed, err = store.IsProcessed(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, processed)

	// Different ID not processed
	processed, err = store.IsProcessed(ctx, 2)
	assert.NoError(t, err)
	assert.False(t, processed)

	// Clear
	store.Clear()
	processed, err = store.IsProcessed(ctx, 1)
	assert.NoError(t, err)
	assert.False(t, processed)
}

// TestRetry tests the retry logic with exponential backoff
func TestRetry(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	// Test success on first attempt
	attempts := 0
	err := Retry(ctx, config, func() error {
		attempts++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)

	// Test success after retries
	attempts = 0
	err = Retry(ctx, config, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, attempts)

	// Test failure after max retries
	attempts = 0
	err = Retry(ctx, config, func() error {
		attempts++
		return errors.New("persistent error")
	})
	assert.Error(t, err)
	assert.Equal(t, 4, attempts) // 1 initial + 3 retries
	assert.Contains(t, err.Error(), "failed after 4 attempts")
}

// TestExponentialBackoff tests the exponential backoff calculation
func TestExponentialBackoff(t *testing.T) {
	config := RetryConfig{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	// Test backoff for various attempts
	tests := []struct {
		attempt     int
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{0, 0, 0},
		{1, 10 * time.Millisecond, 10 * time.Millisecond},
		{2, 20 * time.Millisecond, 20 * time.Millisecond},
		{3, 40 * time.Millisecond, 40 * time.Millisecond},
		{4, 80 * time.Millisecond, 80 * time.Millisecond},
		{5, 100 * time.Millisecond, 100 * time.Millisecond}, // Capped at MaxDelay
	}

	for _, tt := range tests {
		delay := ExponentialBackoff(tt.attempt, config.InitialDelay, config.Multiplier, config.MaxDelay)
		assert.True(t, delay >= tt.expectedMin && delay <= tt.expectedMax,
			"attempt %d: expected delay between %v and %v, got %v", tt.attempt, tt.expectedMin, tt.expectedMax, delay)
	}
}

// TestCircuitBreaker tests the circuit breaker
func TestCircuitBreaker(t *testing.T) {
	ctx := context.Background()
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	// Initially closed
	assert.Equal(t, StateClosed, cb.State())

	// Successful execution
	err := cb.Execute(ctx, func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.State())

	// Failures trigger circuit opening (need threshold failures to open)
	// With threshold=3, we need 3 failures to open the circuit
	for i := 0; i < 3; i++ {
		err := cb.Execute(ctx, func() error {
			return errors.New("error")
		})
		assert.Error(t, err)
	}

	// Circuit should be open after threshold failures
	assert.Equal(t, StateOpen, cb.State())

	// Circuit breaker error when open
	err = cb.Execute(ctx, func() error {
		return nil
	})
	assert.Error(t, err)
	assert.IsType(t, &CircuitBreakerError{}, err)

	// Wait for timeout to transition to half-open
	time.Sleep(150 * time.Millisecond)

	// Check state - should transition to half-open on next Execute call
	// The transition happens when Execute is called, not automatically
	err = cb.Execute(ctx, func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.State()) // Success in half-open closes circuit

	// Reset
	cb.Reset()
	assert.Equal(t, StateClosed, cb.State())
}

// TestCircuitBreakerMiddleware tests the circuit breaker middleware
func TestCircuitBreakerMiddleware(t *testing.T) {
	ctx := context.Background()
	cb := NewCircuitBreaker(2, 100*time.Millisecond)
	middleware := CircuitBreakerMiddleware(cb, "TestConsumer")

	handler := middleware(func(ctx context.Context, event Event) error {
		return errors.New("error")
	})

	event := Event{ID: 1, EventType: "test"}

	// First failure
	err := handler(ctx, event)
	assert.Error(t, err)

	// Second failure - circuit should open
	err = handler(ctx, event)
	assert.Error(t, err)

	// Third call - circuit breaker error
	err = handler(ctx, event)
	assert.Error(t, err)
	assert.IsType(t, &CircuitBreakerError{}, err)
}

// TestTimeoutMiddleware tests the timeout middleware
func TestTimeoutMiddleware(t *testing.T) {
	middleware := TimeoutMiddleware(10*time.Millisecond, "TestConsumer")

	// Test timeout - handler should respect the context cancellation
	handler := middleware(func(ctx context.Context, event Event) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			return nil
		}
	})

	event := Event{ID: 1, EventType: "test"}
	err := handler(context.Background(), event)
	assert.Error(t, err)

	// Test no timeout
	handler = middleware(func(ctx context.Context, event Event) error {
		return nil
	})

	event = Event{ID: 2, EventType: "test"}
	err = handler(context.Background(), event)
	assert.NoError(t, err)
}

// TestIdempotencyMiddleware tests the idempotency middleware
func TestIdempotencyMiddleware(t *testing.T) {
	ctx := context.Background()
	store := NewIdempotencyStore()
	middleware := IdempotencyMiddleware(store, "TestConsumer")

	handler := middleware(func(ctx context.Context, event Event) error {
		return nil
	})

	event := Event{ID: 1, EventType: "test"}

	// First call - success
	err := handler(ctx, event)
	assert.NoError(t, err)

	processed, err := store.IsProcessed(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, processed)

	// Second call - idempotency error
	err = handler(ctx, event)
	assert.Error(t, err)
	assert.IsType(t, &IdempotencyError{}, err)
}

// TestRetryMiddleware tests the retry middleware
func TestRetryMiddleware(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxRetries:   2,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     50 * time.Millisecond,
		Multiplier:   2.0,
	}
	middleware := RetryMiddleware(config, "TestConsumer")

	attempts := 0
	handler := middleware(func(ctx context.Context, event Event) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	event := Event{ID: 1, EventType: "test"}
	err := handler(ctx, event)
	assert.NoError(t, err)
	assert.Equal(t, 2, attempts)
}

// TestDeadLetterMiddleware tests the dead letter middleware
func TestDeadLetterMiddleware(t *testing.T) {
	ctx := context.Background()

	// Mock dead letter sender
	dlh := &mockDeadLetterSender{
		sentEvents: make([]DeadLetterMessage, 0),
	}

	middleware := DeadLetterMiddleware(dlh, 2, "TestConsumer")

	attempts := 0
	handler := middleware(func(ctx context.Context, event Event) error {
		attempts++
		return errors.New("persistent error")
	})

	// Use same event ID to test attempt counting
	event1 := Event{ID: 1, EventType: "test"}

	// First attempt - should not send to DLQ (attempt=1, maxAttempts=2)
	err := handler(ctx, event1)
	assert.Error(t, err)
	assert.Equal(t, 1, attempts)
	assert.Len(t, dlh.sentEvents, 0)

	// Second attempt - should send to DLQ (attempt=2 >= maxAttempts=2)
	err = handler(ctx, event1)
	assert.Error(t, err)
	assert.Equal(t, 2, attempts)
	assert.Len(t, dlh.sentEvents, 1)
	// After DLQ, it should return DeadLetterError
	var deadLetterErr *DeadLetterError
	assert.ErrorAs(t, err, &deadLetterErr)
}

// TestMetricsMiddleware tests the metrics middleware
func TestMetricsMiddleware(t *testing.T) {
	ctx := context.Background()
	metrics := NewInMemoryMetrics()
	middleware := MetricsMiddleware(metrics, "TestConsumer")

	event := Event{ID: 1, EventType: "test"}

	// Test successful processing
	handler := middleware(func(ctx context.Context, event Event) error {
		return nil
	})
	err := handler(ctx, event)
	assert.NoError(t, err)

	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(1), snapshot.EventsReceived)
	assert.Equal(t, uint64(1), snapshot.EventsProcessed)
	assert.Equal(t, uint64(0), snapshot.EventsFailed)

	// Test failed processing
	metrics.Reset()
	handler = middleware(func(ctx context.Context, event Event) error {
		return errors.New("error")
	})
	err = handler(ctx, event)
	assert.Error(t, err)

	snapshot = metrics.GetSnapshot()
	assert.Equal(t, uint64(1), snapshot.EventsReceived)
	assert.Equal(t, uint64(0), snapshot.EventsProcessed)
	assert.Equal(t, uint64(1), snapshot.EventsFailed)
}

// TestLoggingMiddleware tests the logging middleware
func TestLoggingMiddleware(t *testing.T) {
	ctx := context.Background()
	middleware := LoggingMiddleware("TestConsumer")

	handler := middleware(func(ctx context.Context, event Event) error {
		return nil
	})

	event := Event{ID: 1, EventType: "test"}
	err := handler(ctx, event)
	assert.NoError(t, err)
}

// TestMiddlewareChain tests the middleware chain
func TestMiddlewareChain(t *testing.T) {
	ctx := context.Background()

	callOrder := []string{}

	m1 := func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			callOrder = append(callOrder, "m1-before")
			err := next(ctx, event)
			callOrder = append(callOrder, "m1-after")
			return err
		}
	}

	m2 := func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			callOrder = append(callOrder, "m2-before")
			err := next(ctx, event)
			callOrder = append(callOrder, "m2-after")
			return err
		}
	}

	handler := Chain(m1, m2)(func(ctx context.Context, event Event) error {
		callOrder = append(callOrder, "handler")
		return nil
	})

	event := Event{ID: 1, EventType: "test"}
	err := handler(ctx, event)
	assert.NoError(t, err)

	// Verify order: m1-before, m2-before, handler, m2-after, m1-after
	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	assert.Equal(t, expected, callOrder)
}

// TestInMemoryMetrics tests the in-memory metrics implementation
func TestInMemoryMetrics(t *testing.T) {
	metrics := NewInMemoryMetrics()

	// Test increments
	metrics.IncrementReceived()
	metrics.IncrementProcessed(100 * time.Millisecond)
	metrics.IncrementIgnored()
	metrics.IncrementFailed()
	metrics.IncrementDeadLetter()
	metrics.IncrementCircuitBreakerTrip()

	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(1), snapshot.EventsReceived)
	assert.Equal(t, uint64(1), snapshot.EventsProcessed)
	assert.Equal(t, uint64(1), snapshot.EventsIgnored)
	assert.Equal(t, uint64(1), snapshot.EventsFailed)
	assert.Equal(t, uint64(1), snapshot.DeadLetterSent)
	assert.Equal(t, uint64(1), snapshot.CircuitBreakerTrips)
	assert.Equal(t, 100*time.Millisecond, snapshot.AvgProcessingTime)

	// Test reset
	metrics.Reset()
	snapshot = metrics.GetSnapshot()
	assert.Equal(t, uint64(0), snapshot.EventsReceived)
	assert.Equal(t, uint64(0), snapshot.EventsProcessed)
}

// TestMetricsCollectorInterface tests that InMemoryMetrics implements MetricsCollector
func TestMetricsCollectorInterface(t *testing.T) {
	var metrics MetricsCollector = NewInMemoryMetrics()
	assert.NotNil(t, metrics)

	metrics.IncrementReceived()
	metrics.IncrementProcessed(50 * time.Millisecond)

	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(1), snapshot.EventsReceived)
	assert.Equal(t, uint64(1), snapshot.EventsProcessed)

	metrics.Reset()
	snapshot = metrics.GetSnapshot()
	assert.Equal(t, uint64(0), snapshot.EventsReceived)
}

// Mock implementations for testing

type mockDeadLetterSender struct {
	sentEvents []DeadLetterMessage
}

func (m *mockDeadLetterSender) Send(ctx context.Context, event Event, reason string, attempt int) error {
	m.sentEvents = append(m.sentEvents, DeadLetterMessage{
		Event:     event,
		Reason:    reason,
		Attempt:   attempt,
		Timestamp: time.Now(),
	})
	return nil
}
