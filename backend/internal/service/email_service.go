package service

import (
	"fmt"
	"log"
)

// EmailService handles sending emails
type EmailService struct {
	enabled      bool
	from         string
	platformName string // Platform name for email templates (Sprint 3.6)
}

// NewEmailService creates a new email service
func NewEmailService(enabled bool, from, platformName string) *EmailService {
	return &EmailService{
		enabled:      enabled,
		from:         from,
		platformName: platformName,
	}
}

// SendWelcomeEmail sends a welcome email to a new company owner with their temporary password
func (s *EmailService) SendWelcomeEmail(to, name, companyName, tempPassword string) error {
	if !s.enabled {
		log.Printf("[EMAIL] Email sending disabled. Would send welcome email to %s for company %s", to, companyName)
		return nil
	}

	// Use dynamic platform name from brand config (Sprint 3.6)
	platformName := s.platformName
	if platformName == "" {
		platformName = "HorizonGest" // Fallback
	}

	subject := fmt.Sprintf("Bem-vindo ao %s - Sua conta foi criada", platformName)
	_ = fmt.Sprintf(`Olá %s,

Sua empresa "%s" foi criada com sucesso no %s.

Abaixo estão suas credenciais de acesso temporárias:
- E-mail: %s
- Senha temporária: %s

Por favor, faça login e altere sua senha imediatamente por motivos de segurança.

Atenciosamente,
Equipe %s`, name, companyName, platformName, to, tempPassword, platformName)

	// Email sending is not implemented in this version (Sprint 3.4)
	// SMTP implementation is deferred to future sprint
	// Email content is logged for debugging purposes only
	log.Printf("[EMAIL] To: %s, Subject: %s", to, subject)

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, name, resetLink string) error {
	if !s.enabled {
		log.Printf("[EMAIL] Email sending disabled. Would send password reset email to %s", to)
		return nil
	}

	// Use dynamic platform name from brand config (Sprint 3.6)
	platformName := s.platformName
	if platformName == "" {
		platformName = "HorizonGest" // Fallback
	}

	subject := fmt.Sprintf("Redefinição de Senha - %s", platformName)
	body := fmt.Sprintf(`Olá %s,

Recebemos uma solicitação para redefinir sua senha no %s.

Clique no link abaixo para redefinir sua senha:
%s

Se você não solicitou esta redefinição, ignore este email.

Atenciosamente,
Equipe %s`, name, platformName, resetLink, platformName)

	log.Printf("[EMAIL] To: %s, Subject: %s, Body: %s", to, subject, body)

	return nil
}
