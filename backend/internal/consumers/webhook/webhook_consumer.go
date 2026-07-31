package webhook

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/consumers/framework"
	amqp "github.com/rabbitmq/amqp091-go"
)

// WebhookConsumer is a thin wrapper around the framework's BaseConsumer
// It only provides webhook-specific configuration and delegates to the framework
type WebhookConsumer struct {
	baseConsumer *framework.BaseConsumer
}

// NewWebhookConsumer creates a new webhook consumer using the framework
func NewWebhookConsumer(
	conn *amqp.Connection,
	queue string,
	webhookProvider WebhookProvider,
	config framework.Config,
) *WebhookConsumer {
	processor := NewWebhookProcessor(webhookProvider)
	config.ConsumerName = "WebhookConsumer"
	config.Queue = queue

	return &WebhookConsumer{
		baseConsumer: framework.NewBaseConsumer(conn, config, processor),
	}
}

// Start begins consuming events from the queue
func (c *WebhookConsumer) Start(ctx context.Context) error {
	return c.baseConsumer.Start(ctx)
}

// Close closes the consumer
func (c *WebhookConsumer) Close() error {
	return c.baseConsumer.Close()
}

// GetMetrics returns the current metrics snapshot
func (c *WebhookConsumer) GetMetrics() framework.MetricsSnapshot {
	return c.baseConsumer.GetMetrics()
}

// ResetMetrics resets all metrics
func (c *WebhookConsumer) ResetMetrics() {
	c.baseConsumer.ResetMetrics()
}

// GetCircuitBreakerState returns the current circuit breaker state
func (c *WebhookConsumer) GetCircuitBreakerState() framework.CircuitBreakerState {
	return c.baseConsumer.GetCircuitBreakerState()
}

// ResetCircuitBreaker resets the circuit breaker to closed state
func (c *WebhookConsumer) ResetCircuitBreaker() {
	c.baseConsumer.ResetCircuitBreaker()
}

// ClearIdempotency clears all processed event IDs (useful for testing)
func (c *WebhookConsumer) ClearIdempotency() {
	c.baseConsumer.ClearIdempotency()
}
