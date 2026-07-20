package domain

import "time"

// PlatformAudit represents an audit log entry for platform-level actions
type PlatformAudit struct {
	ID            uint
	PlatformUserID *uint
	Action        string
	EntityType    string
	EntityID      *uint
	Changes       string // JSON string of changes
	IPAddress     string
	UserAgent     string
	CreatedAt     time.Time
}
