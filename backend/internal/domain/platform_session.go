package domain

import "time"

// PlatformSession represents a platform user session
type PlatformSession struct {
	ID            uint
	PlatformUserID uint
	Token         string
	ExpiresAt     time.Time
	CreatedAt     time.Time
}

// IsExpired checks if the session is expired
func (s *PlatformSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
