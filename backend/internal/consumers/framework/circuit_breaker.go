package framework

import (
	"context"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu                sync.Mutex
	state             CircuitBreakerState
	failureCount      int
	threshold         int
	timeout           time.Duration
	lastFailureTime   time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     StateClosed,
		threshold: threshold,
		timeout:   timeout,
	}
}

// Execute runs the given function if the circuit is closed or half-open
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	
	// Check if we should transition from open to half-open
	if cb.state == StateOpen && time.Since(cb.lastFailureTime) > cb.timeout {
		cb.state = StateHalfOpen
		cb.failureCount = 0
	}
	
	// Fail fast if circuit is open
	if cb.state == StateOpen {
		cb.mu.Unlock()
		return &CircuitBreakerError{State: cb.state}
	}
	
	cb.mu.Unlock()
	
	// Execute the function
	err := fn()
	
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		
		// Transition to open if threshold reached
		if cb.failureCount >= cb.threshold {
			cb.state = StateOpen
		}
		
		return err
	}
	
	// Reset on success
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failureCount = 0
	}
	
	return nil
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
}

// CircuitBreakerError is returned when the circuit is open
type CircuitBreakerError struct {
	State CircuitBreakerState
}

func (e *CircuitBreakerError) Error() string {
	return "circuit breaker is open"
}
