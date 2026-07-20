package domain

import "time"

// Company representa uma empresa/tenant no sistema multi-tenant.
// Cada tenant opera de forma isolada, com seus próprios dados e configurações.
type Company struct {
	ID          uint
	Name        string
	Slug        string // Identificador único para URL (ex: "empresa-xyz")
	Description string
	Active      bool // "A empresa está ativa no sistema?"

	// Configurações do tenant
	LogoURL        string
	PrimaryColor   string // Cor primária para theming
	SecondaryColor string // Cor secundária para theming

	// Business Engine - Identidade funcional (Sprint 3)
	BusinessType BusinessType // Tipo de negócio (restaurant, bakery, etc.)
	Locale       string       // Localização (ex: "pt-BR", "en-US")
	Currency     string       // Moeda (ex: "BRL", "USD")
	Timezone     string       // Fuso horário (ex: "America/Sao_Paulo")

	// Metadados
	DeletedAt *time.Time // "O registro foi removido logicamente"
	CreatedAt time.Time
	UpdatedAt time.Time

	// Plan and Status (Sprint 3.2)
	PlanID      *uint      `json:"plan_id,omitempty" gorm:"index"`
	Status      string     `json:"status" gorm:"size:20;default:'active'"` // active, trial, suspended, cancelled
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
}
