package webhook

import (
	"context"
	"log"
)

// LogWebhookProvider is a mock implementation that logs webhooks instead of sending them
// This is used for development and testing
type LogWebhookProvider struct{}

// NewLogWebhookProvider creates a new LogWebhookProvider
func NewLogWebhookProvider() *LogWebhookProvider {
	return &LogWebhookProvider{}
}

// Send logs the webhook details instead of actually sending it
func (p *LogWebhookProvider) Send(ctx context.Context, webhook Webhook) error {
	log.Printf("[WEBHOOK] URL: %s", webhook.URL)
	log.Printf("[WEBHOOK] Template: %s", webhook.Headers["Template"])
	log.Printf("[WEBHOOK] Payload: %s", webhook.Headers["Payload"])
	log.Printf("[WEBHOOK] EventID: %s", webhook.Headers["EventID"])
	log.Printf("[WEBHOOK] Payload Size: %d bytes", len(webhook.Payload))
	return nil
}

// Close is a no-op for LogWebhookProvider
func (p *LogWebhookProvider) Close() error {
	log.Printf("[WEBHOOK] LogWebhookProvider closed")
	return nil
}

var _ WebhookProvider = (*LogWebhookProvider)(nil)
