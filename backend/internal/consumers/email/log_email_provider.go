package email

import (
	"context"
	"log"
	"strings"
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
	// FASE A: Mask email in logs
	emailMask := maskEmail(email.To)
	log.Printf("[EMAIL] To: %s", emailMask)
	log.Printf("[EMAIL] Subject: %s", email.Subject)
	log.Printf("[EMAIL] Template: %s", email.Headers["Template"])
	log.Printf("[EMAIL] Payload: %s", email.Headers["Payload"])
	log.Printf("[EMAIL] EventID: %s", email.Headers["EventID"])
	log.Printf("[EMAIL] Body Length: %d bytes", len(email.Body))
	return nil
}

// maskEmail masks an email address for logging (FASE A)
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	username := parts[0]
	domain := parts[1]

	// Mask username: show first 2 chars, rest as *
	if len(username) > 2 {
		username = username[:2] + strings.Repeat("*", len(username)-2)
	} else {
		username = strings.Repeat("*", len(username))
	}

	// Mask domain: show first char, rest as *
	domainParts := strings.Split(domain, ".")
	if len(domainParts) >= 2 {
		domainName := domainParts[0]
		tld := strings.Join(domainParts[1:], ".")
		if len(domainName) > 1 {
			domainName = domainName[:1] + strings.Repeat("*", len(domainName)-1)
		}
		return username + "@" + domainName + "." + tld
	}

	return username + "@" + strings.Repeat("*", len(domain))
}

// Close is a no-op for LogEmailProvider
func (p *LogEmailProvider) Close() error {
	log.Printf("[EMAIL] LogEmailProvider closed")
	return nil
}

var _ EmailProvider = (*LogEmailProvider)(nil)
