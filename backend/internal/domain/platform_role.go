package domain

// PlatformRole represents the role of a platform user
type PlatformRole string

const (
	PlatformRoleAdmin   PlatformRole = "PlatformAdmin"   // Full access to platform
	PlatformRoleSupport PlatformRole = "PlatformSupport" // Support access, limited permissions
)

// IsValid checks if the role is valid
func (r PlatformRole) IsValid() bool {
	switch r {
	case PlatformRoleAdmin, PlatformRoleSupport:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r PlatformRole) String() string {
	return string(r)
}

// DisplayName returns the human-readable display name
func (r PlatformRole) DisplayName() string {
	switch r {
	case PlatformRoleAdmin:
		return "Administrador da Plataforma"
	case PlatformRoleSupport:
		return "Suporte da Plataforma"
	default:
		return "Desconhecido"
	}
}

// AllPlatformRoles returns all available platform roles
func AllPlatformRoles() []PlatformRole {
	return []PlatformRole{
		PlatformRoleAdmin,
		PlatformRoleSupport,
	}
}

// ParsePlatformRole parses a string into a PlatformRole
func ParsePlatformRole(s string) (PlatformRole, bool) {
	r := PlatformRole(s)
	if r.IsValid() {
		return r, true
	}
	return "", false
}
