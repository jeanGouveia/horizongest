package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeadLetterSender is the interface for sending messages to dead letter queue
type DeadLetterSender interface {
	Send(ctx context.Context, event Event, reason string, attempt int) error
}

// DeadLetterHandler handles messages that exceed retry attempts
type DeadLetterHandler struct {
	connection *amqp.Connection
	queue      string
}

var _ DeadLetterSender = (*DeadLetterHandler)(nil)

// NewDeadLetterHandler creates a new dead letter handler
func NewDeadLetterHandler(conn *amqp.Connection, queue string) *DeadLetterHandler {
	return &DeadLetterHandler{
		connection: conn,
		queue:      queue,
	}
}

// Send sends a message to the dead letter queue
func (dlh *DeadLetterHandler) Send(ctx context.Context, event Event, reason string, attempt int) error {
	channel, err := dlh.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	// Declare dead letter queue if it doesn't exist
	_, err = channel.QueueDeclare(
		dlh.queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter queue: %w", err)
	}

	// Create dead letter message
	deadLetterMsg := DeadLetterMessage{
		Event:     event,
		Reason:    reason,
		Attempt:   attempt,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(deadLetterMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal dead letter message: %w", err)
	}

	// Publish to dead letter queue
	err = channel.PublishWithContext(
		ctx,
		"",        // exchange
		dlh.queue, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish to dead letter queue: %w", err)
	}

	log.Printf("[DeadLetter] Event id=%d sent to dead letter queue: %s (reason: %s, attempt: %d)",
		event.ID, dlh.queue, reason, attempt)

	return nil
}

// DeadLetterMessage represents a message in the dead letter queue
type DeadLetterMessage struct {
	Event     Event     `json:"event"`
	Reason    string    `json:"reason"`
	Attempt   int       `json:"attempt"`
	Timestamp time.Time `json:"timestamp"`
}
