package webhook

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/consumers/framework"
)

// WebhookProcessor implements the framework.Processor interface for webhook processing
// This is where the webhook-specific business logic lives
type WebhookProcessor struct {
	webhookProvider WebhookProvider
	templates       map[string]WebhookTemplate
}

// NewWebhookProcessor creates a new webhook processor
func NewWebhookProcessor(webhookProvider WebhookProvider) *WebhookProcessor {
	templates := map[string]WebhookTemplate{
		"invitation.created": NewInvitationWebhookTemplate(),
		"order.created":      NewOrderCreatedWebhookTemplate(),
		"company.created":    NewCompanyCreatedWebhookTemplate(),
	}

	return &WebhookProcessor{
		webhookProvider: webhookProvider,
		templates:       templates,
	}
}

// Process processes an event by transforming it into a webhook and sending it
// This is the only business logic in the webhook consumer
func (p *WebhookProcessor) Process(ctx context.Context, event framework.Event) error {
	// Get template for event type
	template, ok := p.templates[event.EventType]
	if !ok {
		return fmt.Errorf("no template for event type: %s", event.EventType)
	}

	// Render template
	url, payload, err := template.Render(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Create webhook
	webhook := Webhook{
		URL:     url,
		Payload: payload,
		Headers: map[string]string{
			"Template": event.EventType,
			"Payload":   fmt.Sprintf("%+v", event.Payload),
			"EventID":   fmt.Sprintf("%d", event.ID),
		},
	}

	// Send webhook
	if err := p.webhookProvider.Send(ctx, webhook); err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}

	return nil
}

// Close closes the webhook processor
func (p *WebhookProcessor) Close() error {
	if p.webhookProvider != nil {
		return p.webhookProvider.Close()
	}
	return nil
}
