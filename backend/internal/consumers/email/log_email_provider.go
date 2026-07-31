package email

import (
	"context"
	"log"
)

// LogEmailProvider is a mock implementation that logs emails instead of sending them
// This is used for development and testing
type LogEmailProvider struct{}

// NewLogEmailProvider creates a new LogEmailProvider
func NewLogEmailProvider() *LogEmailProvider {
	return &LogEmailProvider{}
}

// Send logs the email details instead of actually sending it
func (p *LogEmailProvider) Send(ctx context.Context, email Email) error {
	log.Printf("[EMAIL] To: %s", email.To)
	log.Printf("[EMAIL] Subject: %s", email.Subject)
	log.Printf("[EMAIL] Template: %s", email.Headers["Template"])
	log.Printf("[EMAIL] Payload: %s", email.Headers["Payload"])
	log.Printf("[EMAIL] EventID: %s", email.Headers["EventID"])
	log.Printf("[EMAIL] Body Length: %d bytes", len(email.Body))
	return nil
}

// Close is a no-op for LogEmailProvider
func (p *LogEmailProvider) Close() error {
	log.Printf("[EMAIL] LogEmailProvider closed")
	return nil
}

var _ EmailProvider = (*LogEmailProvider)(nil)
