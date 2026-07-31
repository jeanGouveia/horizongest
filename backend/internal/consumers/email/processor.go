package email

import (
	"context"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/consumers/framework"
)

// EmailProcessor implements the framework.Processor interface for email processing
// This is where the email-specific business logic lives
type EmailProcessor struct {
	emailProvider EmailProvider
	templates     map[string]Template
}

// NewEmailProcessor creates a new email processor
func NewEmailProcessor(emailProvider EmailProvider) *EmailProcessor {
	templates := map[string]Template{
		"invitation.created": NewInvitationTemplate(),
		"order.created":      NewOrderCreatedTemplate(),
		"company.created":    NewCompanyCreatedTemplate(),
	}

	return &EmailProcessor{
		emailProvider: emailProvider,
		templates:     templates,
	}
}

// Process processes an event by transforming it into an email and sending it
// This is the only business logic in the email consumer
func (p *EmailProcessor) Process(ctx context.Context, event framework.Event) error {
	// Get template for event type
	template, ok := p.templates[event.EventType]
	if !ok {
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
			"Payload":  fmt.Sprintf("%+v", event.Payload),
			"EventID":  fmt.Sprintf("%d", event.ID),
		},
	}

	// Send email
	if err := p.emailProvider.Send(ctx, email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Close closes the email processor
func (p *EmailProcessor) Close() error {
	if p.emailProvider != nil {
		return p.emailProvider.Close()
	}
	return nil
}
