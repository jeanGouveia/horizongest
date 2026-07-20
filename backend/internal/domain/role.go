package domain

// Role represents the user's role within a company.
// Nullable for Core V1 compatibility (users without CompanyID have Role == null).
type Role string

const (
	RoleOwner    Role = "owner"    // Full access to everything, can manage other users
	RoleAdmin    Role = "admin"    // Full access except cannot alter Owner role
	RoleManager  Role = "manager"  // Can manage orders, products, and view reports
	RoleEmployee Role = "employee" // Limited operational access
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleManager, RoleEmployee:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// DisplayName returns the human-readable display name
func (r Role) DisplayName() string {
	switch r {
	case RoleOwner:
		return "Proprietário"
	case RoleAdmin:
		return "Administrador"
	case RoleManager:
		return "Gerente"
	case RoleEmployee:
		return "Funcionário"
	default:
		return "Desconhecido"
	}
}

// AllRoles returns all available roles
func AllRoles() []Role {
	return []Role{
		RoleOwner,
		RoleAdmin,
		RoleManager,
		RoleEmployee,
	}
}

// ParseRole parses a string into a Role
func ParseRole(s string) (Role, bool) {
	r := Role(s)
	if r.IsValid() {
		return r, true
	}
	return "", false
}
