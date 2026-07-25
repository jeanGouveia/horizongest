package domain

import "time"

// ImpersonationAudit records all impersonation sessions for audit purposes
type ImpersonationAudit struct {
	ID                     uint      `json:"id"`
	PlatformUserID         uint      `json:"platformUserId"`         // The platform admin who initiated impersonation
	CompanyID              uint      `json:"companyId"`              // The company being impersonated
	CompanyOwnerUserID     uint      `json:"companyOwnerUserId"`     // The owner user of the company being impersonated
	StartedAt              time.Time `json:"startedAt"`              // When impersonation started
	EndedAt                *time.Time `json:"endedAt,omitempty"`     // When impersonation ended (NULL if still active)
	IPAddress              string    `json:"ipAddress"`              // IP address of the platform admin
	UserAgent              string    `json:"userAgent"`              // User agent of the platform admin
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	DeletedAt              *time.Time `json:"deletedAt,omitempty"`
}

// IsActive returns true if the impersonation session is still active
func (ia *ImpersonationAudit) IsActive() bool {
	return ia.EndedAt == nil
}

// End marks the impersonation session as ended
func (ia *ImpersonationAudit) End() {
	now := time.Now()
	ia.EndedAt = &now
	ia.UpdatedAt = now
}
