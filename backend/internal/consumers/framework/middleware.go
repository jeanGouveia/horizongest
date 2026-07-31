package framework

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Middleware is a function that processes a message before/after the handler
type Middleware func(next Handler) Handler

// Handler is the function that processes a message
type Handler func(ctx context.Context, event Event) error

// Chain chains multiple middlewares together
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// LoggingMiddleware logs message processing
func LoggingMiddleware(consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			startTime := time.Now()
			log.Printf("[%s] Processing event id=%d, type=%s", consumerName, event.ID, event.EventType)

			err := next(ctx, event)

			duration := time.Since(startTime)
			if err != nil {
				log.Printf("[%s] Event id=%d failed in %v: %v", consumerName, event.ID, duration, err)
			} else {
				log.Printf("[%s] Event id=%d processed successfully in %v", consumerName, event.ID, duration)
			}

			return err
		}
	}
}

// IdempotencyMiddleware checks if an event has already been processed
func IdempotencyMiddleware(store *IdempotencyStore, consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			if store.IsProcessed(event.ID) {
				log.Printf("[%s] Event id=%d already processed, ignoring", consumerName, event.ID)
				return &IdempotencyError{EventID: event.ID}
			}

			err := next(ctx, event)

			if err == nil {
				store.MarkProcessed(event.ID)
			}

			return err
		}
	}
}

// RetryMiddleware adds retry logic with exponential backoff
func RetryMiddleware(config RetryConfig, consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			return Retry(ctx, config, func() error {
				return next(ctx, event)
			})
		}
	}
}

// CircuitBreakerMiddleware adds circuit breaker protection
func CircuitBreakerMiddleware(cb *CircuitBreaker, consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			err := cb.Execute(ctx, func() error {
				return next(ctx, event)
			})

			// Check if error is circuit breaker error
			var circuitBreakerErr *CircuitBreakerError
			if err != nil && errors.As(err, &circuitBreakerErr) {
				log.Printf("[%s] Circuit breaker is open, skipping event id=%d", consumerName, event.ID)
			}

			return err
		}
	}
}

// TimeoutMiddleware adds timeout to operations
func TimeoutMiddleware(timeout time.Duration, consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return next(ctx, event)
		}
	}
}

// DeadLetterMiddleware sends failed messages to dead letter queue
func DeadLetterMiddleware(dlh DeadLetterSender, maxAttempts int, consumerName string) Middleware {
	attemptMap := make(map[uint]int)

	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			attemptMap[event.ID]++
			attempt := attemptMap[event.ID]

			err := next(ctx, event)

			if err != nil && attempt >= maxAttempts {
				log.Printf("[%s] Event id=%d exceeded max attempts (%d), sending to dead letter",
					consumerName, event.ID, maxAttempts)

				dlhErr := dlh.Send(ctx, event, err.Error(), attempt)
				if dlhErr != nil {
					log.Printf("[%s] Failed to send event id=%d to dead letter: %v", consumerName, event.ID, dlhErr)
				}

				delete(attemptMap, event.ID)

				// Return DeadLetterError to signal that message was sent to DLQ
				// This tells the consumer to Ack (not Nack) to avoid RabbitMQ re-delivery
				return &DeadLetterError{
					EventID: event.ID,
					Reason:  err.Error(),
				}
			}

			if err == nil {
				delete(attemptMap, event.ID)
			}

			return err
		}
	}
}

// MetricsMiddleware tracks metrics
func MetricsMiddleware(metrics MetricsCollector, consumerName string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, event Event) error {
			metrics.IncrementReceived()

			startTime := time.Now()
			err := next(ctx, event)
			duration := time.Since(startTime)

			if err != nil {
				if _, ok := err.(*IdempotencyError); ok {
					metrics.IncrementIgnored()
				} else {
					metrics.IncrementFailed()
				}
			} else {
				metrics.IncrementProcessed(duration)
			}

			return err
		}
	}
}

// IdempotencyError is returned when an event is a duplicate
type IdempotencyError struct {
	EventID uint
}

func (e *IdempotencyError) Error() string {
	return fmt.Sprintf("event id=%d already processed", e.EventID)
}

// DeadLetterError is returned when an event has been sent to dead letter
type DeadLetterError struct {
	EventID uint
	Reason  string
}

func (e *DeadLetterError) Error() string {
	return fmt.Sprintf("event id=%d sent to dead letter: %s", e.EventID, e.Reason)
}

// ParseMessage parses a RabbitMQ message into an Event
func ParseMessage(msg amqp.Delivery) (Event, error) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return Event{}, fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return event, nil
}
