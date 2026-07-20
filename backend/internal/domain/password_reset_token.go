package domain

import "time"

// PasswordResetToken represents a temporary token for password recovery
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (PasswordResetToken) TableName() string { return "password_reset_tokens" }
