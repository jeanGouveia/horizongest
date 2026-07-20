package service

import (
	"fmt"
	"log"
)

// EmailService handles sending emails
type EmailService struct {
	enabled bool
	from    string
}

// NewEmailService creates a new email service
func NewEmailService(enabled bool, from string) *EmailService {
	return &EmailService{
		enabled: enabled,
		from:    from,
	}
}

// SendWelcomeEmail sends a welcome email to a new company owner with their temporary password
func (s *EmailService) SendWelcomeEmail(to, name, companyName, tempPassword string) error {
	if !s.enabled {
		log.Printf("[EMAIL] Email sending disabled. Would send welcome email to %s for company %s", to, companyName)
		return nil
	}

	subject := "Bem-vindo ao PratoOnline - Sua conta foi criada"
	_ = fmt.Sprintf(`Olá %s,

Sua empresa "%s" foi criada com sucesso no PratoOnline.

Abaixo estão suas credenciais de acesso temporárias:
- E-mail: %s
- Senha temporária: %s

Por favor, faça login e altere sua senha imediatamente por motivos de segurança.

Atenciosamente,
Equipe PratoOnline`, name, companyName, to, tempPassword)

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

	subject := "Redefinição de Senha - PratoOnline"
	body := fmt.Sprintf(`Olá %s,

Recebemos uma solicitação para redefinir sua senha no PratoOnline.

Clique no link abaixo para redefinir sua senha:
%s

Se você não solicitou esta redefinição, ignore este email.

Atenciosamente,
Equipe PratoOnline`, name, resetLink)

	log.Printf("[EMAIL] To: %s, Subject: %s, Body: %s", to, subject, body)

	return nil
}
