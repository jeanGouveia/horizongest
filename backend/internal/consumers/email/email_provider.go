package email

import "context"

// EmailProvider defines the interface for sending emails
// This allows different implementations (SMTP, SendGrid, SES, Mailgun, etc.)
type EmailProvider interface {
	// Send sends an email with the given parameters
	Send(ctx context.Context, email Email) error

	// Close closes any resources held by the provider
	Close() error
}

// Email represents an email to be sent
type Email struct {
	To      string
	Subject string
	Body    string
	HTML    bool
	Headers map[string]string
}
