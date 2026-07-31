package framework

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Processor is the interface that specific consumers must implement
// This is where the business logic for each consumer lives
type Processor interface {
	// Process processes an event and returns an error if processing fails
	// The processor should only:
	// 1. Transform the event payload into the appropriate message format
	// 2. Call the provider to send the message
	Process(ctx context.Context, event Event) error

	// Close closes any resources held by the processor
	Close() error
}

// BaseConsumer is the framework's base consumer that handles all common logic
// Specific consumers (email, webhook, ifood, etc.) should embed this and provide a Processor
type BaseConsumer struct {
	connection *amqp.Connection
	config     Config
	processor  Processor

	// Framework components
	idempotencyStore  *IdempotencyStore
	circuitBreaker    *CircuitBreaker
	deadLetterHandler *DeadLetterHandler
	metrics           MetricsCollector

	// Middleware chain
	handler Handler
}

// NewBaseConsumer creates a new base consumer with the framework components
func NewBaseConsumer(
	conn *amqp.Connection,
	config Config,
	processor Processor,
) *BaseConsumer {
	// Initialize framework components
	idempotencyStore := NewIdempotencyStore()
	circuitBreaker := NewCircuitBreaker(config.CircuitBreakerThreshold, config.CircuitBreakerTimeout)

	var deadLetterHandler *DeadLetterHandler
	if config.DeadLetterQueue != "" {
		deadLetterHandler = NewDeadLetterHandler(conn, config.DeadLetterQueue)
	}

	metrics := NewInMemoryMetrics()

	// Build middleware chain
	handler := func(ctx context.Context, event Event) error {
		return processor.Process(ctx, event)
	}

	// Apply middlewares in order:
	// 1. Logging - logs all processing
	// 2. Idempotency - removes duplicates immediately
	// 3. CircuitBreaker - fails fast if circuit is open
	// 4. Retry - retries with exponential backoff
	// 5. Timeout - enforces timeout within retry
	// 6. Metrics - tracks metrics
	middlewares := []Middleware{
		LoggingMiddleware(config.ConsumerName),
		IdempotencyMiddleware(idempotencyStore, config.ConsumerName),
		CircuitBreakerMiddleware(circuitBreaker, config.ConsumerName),
		RetryMiddleware(RetryConfig{
			MaxRetries:   config.MaxRetries,
			InitialDelay: config.InitialRetryDelay,
			MaxDelay:     config.MaxRetryDelay,
			Multiplier:   config.RetryMultiplier,
		}, config.ConsumerName),
		TimeoutMiddleware(config.OperationTimeout, config.ConsumerName),
	}

	if config.EnableMetrics {
		middlewares = append(middlewares, MetricsMiddleware(metrics, config.ConsumerName))
	}

	if deadLetterHandler != nil {
		middlewares = append(middlewares, DeadLetterMiddleware(deadLetterHandler, config.MaxRetryAttempts, config.ConsumerName))
	}

	// Chain all middlewares
	handler = Chain(middlewares...)(handler)

	return &BaseConsumer{
		connection:        conn,
		config:            config,
		processor:         processor,
		idempotencyStore:  idempotencyStore,
		circuitBreaker:    circuitBreaker,
		deadLetterHandler: deadLetterHandler,
		metrics:           metrics,
		handler:           handler,
	}
}

// Start begins consuming events from the queue
func (c *BaseConsumer) Start(ctx context.Context) error {
	log.Printf("[%s] Starting to consume from queue: %s", c.config.ConsumerName, c.config.Queue)

	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		c.config.Queue, // queue
		"",             // consumer
		false,          // auto-ack (we'll ack manually)
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Shutdown requested", c.config.ConsumerName)
			return nil
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("[%s] Channel closed", c.config.ConsumerName)
				return nil
			}

			c.processMessage(ctx, msg)
		}
	}
}

// processMessage processes a single message from RabbitMQ
func (c *BaseConsumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	// Parse event
	event, err := ParseMessage(msg)
	if err != nil {
		log.Printf("[%s] Failed to parse message: %v", c.config.ConsumerName, err)
		msg.Nack(false, false) // Don't requeue malformed messages
		return
	}

	// Process through middleware chain
	// This includes internal retry, circuit breaker, dead letter, etc.
	err = c.handler(ctx, event)

	// Handle ack/nack based on error type
	if err != nil {
		if _, isIdempotencyErr := err.(*IdempotencyError); isIdempotencyErr {
			// Duplicate event - ack to remove from queue
			msg.Ack(false)
		} else if _, isCircuitBreakerErr := err.(*CircuitBreakerError); isCircuitBreakerErr {
			// Circuit breaker open - nack without requeue (don't retry via RabbitMQ)
			msg.Nack(false, false)
		} else if _, isDeadLetterErr := err.(*DeadLetterError); isDeadLetterErr {
			// Sent to dead letter - ack to remove from queue (don't retry via RabbitMQ)
			msg.Ack(false)
		} else {
			// Other error after internal retry exhausted - ack to remove from queue
			// Internal retry already attempted, don't retry via RabbitMQ to avoid duplicate retry
			msg.Ack(false)
		}
		return
	}

	// Success - ack message
	msg.Ack(false)
}

// Close closes the consumer and its resources
func (c *BaseConsumer) Close() error {
	log.Printf("[%s] Closing", c.config.ConsumerName)

	if c.processor != nil {
		if err := c.processor.Close(); err != nil {
			return err
		}
	}

	return nil
}

// GetMetrics returns the current metrics snapshot
func (c *BaseConsumer) GetMetrics() MetricsSnapshot {
	return c.metrics.GetSnapshot()
}

// ResetMetrics resets all metrics
func (c *BaseConsumer) ResetMetrics() {
	c.metrics.Reset()
}

// GetCircuitBreakerState returns the current circuit breaker state
func (c *BaseConsumer) GetCircuitBreakerState() CircuitBreakerState {
	return c.circuitBreaker.State()
}

// ResetCircuitBreaker resets the circuit breaker to closed state
func (c *BaseConsumer) ResetCircuitBreaker() {
	c.circuitBreaker.Reset()
}

// ClearIdempotency clears all processed event IDs (useful for testing)
func (c *BaseConsumer) ClearIdempotency() {
	c.idempotencyStore.Clear()
}
