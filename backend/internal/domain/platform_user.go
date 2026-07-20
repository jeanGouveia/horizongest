package domain

import "time"

// PlatformUser represents a user at the platform level (PlatformAdmin, PlatformSupport)
// Platform users manage companies and have no CompanyID
type PlatformUser struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool
	Role         PlatformRole
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
