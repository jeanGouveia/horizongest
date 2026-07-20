package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// InvitationStatus represents the status of an invitation
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

// IsValid checks if the invitation status is valid
func (s InvitationStatus) IsValid() bool {
	switch s {
	case InvitationStatusPending, InvitationStatusAccepted, InvitationStatusExpired, InvitationStatusRevoked:
		return true
	default:
		return false
	}
}

// String returns the string representation of the invitation status
func (s InvitationStatus) String() string {
	return string(s)
}

// DisplayName returns the display name of the invitation status
func (s InvitationStatus) DisplayName() string {
	switch s {
	case InvitationStatusPending:
		return "Pendente"
	case InvitationStatusAccepted:
		return "Aceito"
	case InvitationStatusExpired:
		return "Expirado"
	case InvitationStatusRevoked:
		return "Revogado"
	default:
		return "Desconhecido"
	}
}

// ParseInvitationStatus parses a string to InvitationStatus
func ParseInvitationStatus(s string) (InvitationStatus, bool) {
	status := InvitationStatus(s)
	return status, status.IsValid()
}

// AllInvitationStatuses returns all valid invitation statuses
func AllInvitationStatuses() []InvitationStatus {
	return []InvitationStatus{
		InvitationStatusPending,
		InvitationStatusAccepted,
		InvitationStatusExpired,
		InvitationStatusRevoked,
	}
}

// Invitation represents an invitation to join a company
type Invitation struct {
	ID          uint
	CompanyID   uint
	Email       string
	Role        Role
	Token       string
	Status      InvitationStatus
	ExpiresAt   time.Time
	AcceptedAt  *time.Time
	CreatedBy   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsExpired checks if the invitation is expired
func (i *Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// CanBeAccepted checks if the invitation can be accepted
func (i *Invitation) CanBeAccepted() bool {
	return i.Status == InvitationStatusPending && !i.IsExpired()
}

// GenerateToken generates a secure random token for the invitation
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
