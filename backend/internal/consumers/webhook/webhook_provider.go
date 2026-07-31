package webhook

import "context"

// WebhookProvider defines the interface for sending webhooks
// This allows different implementations (HTTP, async, retry, etc.)
type WebhookProvider interface {
	// Send sends a webhook with the given parameters
	Send(ctx context.Context, webhook Webhook) error

	// Close closes any resources held by the provider
	Close() error
}

// Webhook represents a webhook to be sent
type Webhook struct {
	URL     string
	Payload map[string]interface{}
	Headers map[string]string
}
