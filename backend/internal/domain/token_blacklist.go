package domain

import "time"

// TokenBlacklist representa um token JWT revogado
type TokenBlacklist struct {
	Token     string    `gorm:"primaryKey;size:500"`
	RevokedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}
