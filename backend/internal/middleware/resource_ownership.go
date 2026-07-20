package middleware

import (
	"errors"
)

// ValidateCompanyOwnership is a helper function to validate company ownership
// This can be used directly in handlers for custom validation logic
// Sprint 3.4 - Security Hardening
func ValidateCompanyOwnership(userCompanyID, resourceCompanyID uint) error {
	if userCompanyID != resourceCompanyID {
		return errors.New("access denied - resource belongs to another company")
	}
	return nil
}
