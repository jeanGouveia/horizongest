package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EmailConsumer consumes events from RabbitMQ and sends emails
type EmailConsumer struct {
	connection      *amqp.Connection
	queue           string
	emailProvider   EmailProvider
	idempotencyStore *IdempotencyStore
	templates       map[string]Template
}

// NewEmailConsumer creates a new EmailConsumer
func NewEmailConsumer(
	conn *amqp.Connection,
	queue string,
	emailProvider EmailProvider,
	idempotencyStore *IdempotencyStore,
) *EmailConsumer {
	templates := map[string]Template{
		"invitation.created": NewInvitationTemplate(),
		"order.created":      NewOrderCreatedTemplate(),
		"company.created":    NewCompanyCreatedTemplate(),
	}

	return &EmailConsumer{
		connection:      conn,
		queue:           queue,
		emailProvider:   emailProvider,
		idempotencyStore: idempotencyStore,
		templates:       templates,
	}
}

// Event represents a domain event from RabbitMQ
type Event struct {
	ID            uint                   `json:"event_id"`
	EventType     string                 `json:"event_type"`
	AggregateType string                 `json:"aggregate_type"`
	AggregateID   uint                   `json:"aggregate_id"`
	TenantID      uint                   `json:"tenant_id"`
	Payload       map[string]interface{} `json:"payload"`
}

// Start begins consuming events from the queue
func (c *EmailConsumer) Start(ctx context.Context) error {
	log.Printf("EmailConsumer: starting to consume from queue: %s", c.queue)

	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack (we'll ack manually)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("EmailConsumer: shutdown requested")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("EmailConsumer: channel closed")
				return nil
			}

			c.processMessage(ctx, msg)
		}
	}
}

// processMessage processes a single message from RabbitMQ
func (c *EmailConsumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	startTime := time.Now()

	// Parse event
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("EmailConsumer: failed to unmarshal event: %v", err)
		msg.Nack(false, false) // Don't requeue malformed messages
		return
	}

	log.Printf("EmailConsumer: received event id=%d, type=%s", event.ID, event.EventType)

	// Check idempotency
	if c.idempotencyStore.IsProcessed(event.ID) {
		log.Printf("EmailConsumer: event id=%d already processed, ignoring", event.ID)
		msg.Ack(false) // Ack to remove from queue
		return
	}

	// Process event
	if err := c.processEvent(ctx, event); err != nil {
		log.Printf("EmailConsumer: failed to process event id=%d: %v", event.ID, err)
		msg.Nack(false, true) // Requeue for retry
		return
	}

	// Mark as processed
	c.idempotencyStore.MarkProcessed(event.ID)

	// Ack message
	msg.Ack(false)

	duration := time.Since(startTime)
	log.Printf("EmailConsumer: event id=%d processed successfully in %v", event.ID, duration)
}

// processEvent processes a single event
func (c *EmailConsumer) processEvent(ctx context.Context, event Event) error {
	// Get template for event type
	template, ok := c.templates[event.EventType]
	if !ok {
		log.Printf("EmailConsumer: no template for event type: %s", event.EventType)
		return fmt.Errorf("no template for event type: %s", event.EventType)
	}

	// Render template
	subject, body, err := template.Render(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Determine recipient (for now, use a placeholder)
	// In production, this would come from the payload or a lookup
	to := "user@example.com"

	// Create email
	email := Email{
		To:      to,
		Subject: subject,
		Body:    body,
		HTML:    false,
		Headers: map[string]string{
			"Template": event.EventType,
			"Payload":   fmt.Sprintf("%+v", event.Payload),
			"EventID":   fmt.Sprintf("%d", event.ID),
		},
	}

	// Send email
	if err := c.emailProvider.Send(ctx, email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Close closes the consumer
func (c *EmailConsumer) Close() error {
	log.Printf("EmailConsumer: closing")
	if c.emailProvider != nil {
		return c.emailProvider.Close()
	}
	return nil
}
