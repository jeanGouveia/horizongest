package framework

import (
	"sync"
	"time"
)

// MetricsCollector is the interface for collecting consumer metrics
// This allows different implementations (in-memory, Prometheus, OpenTelemetry, Datadog, etc.)
type MetricsCollector interface {
	// IncrementReceived increments the events received counter
	IncrementReceived()

	// IncrementProcessed increments the events processed counter
	IncrementProcessed(duration time.Duration)

	// IncrementIgnored increments the events ignored counter
	IncrementIgnored()

	// IncrementFailed increments the events failed counter
	IncrementFailed()

	// IncrementDeadLetter increments the dead letter sent counter
	IncrementDeadLetter()

	// IncrementCircuitBreakerTrip increments the circuit breaker trip counter
	IncrementCircuitBreakerTrip()

	// GetSnapshot returns a snapshot of current metrics
	GetSnapshot() MetricsSnapshot

	// Reset resets all metrics
	Reset()
}

// MetricsSnapshot represents a snapshot of metrics at a point in time
type MetricsSnapshot struct {
	EventsReceived      uint64
	EventsProcessed     uint64
	EventsIgnored       uint64
	EventsFailed        uint64
	DeadLetterSent      uint64
	AvgProcessingTime   time.Duration
	CircuitBreakerTrips uint64
}

// InMemoryMetrics is an in-memory implementation of MetricsCollector
// This is the default implementation used by the framework
type InMemoryMetrics struct {
	mu sync.RWMutex

	// Counters
	EventsReceived  uint64
	EventsProcessed uint64
	EventsIgnored   uint64
	EventsFailed    uint64
	DeadLetterSent  uint64

	// Timing
	TotalProcessingTime time.Duration

	// Circuit breaker
	CircuitBreakerTrips uint64
}

// NewInMemoryMetrics creates a new in-memory metrics instance
func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{}
}

// IncrementReceived increments the events received counter
func (m *InMemoryMetrics) IncrementReceived() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventsReceived++
}

// IncrementProcessed increments the events processed counter
func (m *InMemoryMetrics) IncrementProcessed(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventsProcessed++
	m.TotalProcessingTime += duration
}

// IncrementIgnored increments the events ignored counter
func (m *InMemoryMetrics) IncrementIgnored() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventsIgnored++
}

// IncrementFailed increments the events failed counter
func (m *InMemoryMetrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventsFailed++
}

// IncrementDeadLetter increments the dead letter sent counter
func (m *InMemoryMetrics) IncrementDeadLetter() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DeadLetterSent++
}

// IncrementCircuitBreakerTrip increments the circuit breaker trip counter
func (m *InMemoryMetrics) IncrementCircuitBreakerTrip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CircuitBreakerTrips++
}

// GetSnapshot returns a snapshot of current metrics
func (m *InMemoryMetrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgProcessingTime time.Duration
	if m.EventsProcessed > 0 {
		avgProcessingTime = m.TotalProcessingTime / time.Duration(m.EventsProcessed)
	}

	return MetricsSnapshot{
		EventsReceived:      m.EventsReceived,
		EventsProcessed:     m.EventsProcessed,
		EventsIgnored:       m.EventsIgnored,
		EventsFailed:        m.EventsFailed,
		DeadLetterSent:      m.DeadLetterSent,
		AvgProcessingTime:   avgProcessingTime,
		CircuitBreakerTrips: m.CircuitBreakerTrips,
	}
}

// Reset resets all metrics
func (m *InMemoryMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventsReceived = 0
	m.EventsProcessed = 0
	m.EventsIgnored = 0
	m.EventsFailed = 0
	m.DeadLetterSent = 0
	m.TotalProcessingTime = 0
	m.CircuitBreakerTrips = 0
}

// NewMetrics creates a new metrics instance (alias for NewInMemoryMetrics for backward compatibility)
// Deprecated: Use NewInMemoryMetrics instead
func NewMetrics() *InMemoryMetrics {
	return NewInMemoryMetrics()
}
