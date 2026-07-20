package domain

import "time"

// User é a entidade central. Sem tags GORM aqui — domínio puro.
// Sprint 3: CompanyID e Role são obrigatórios (NOT NULL) - não existe mais usuário sem empresa.
type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool       // "Pode ser utilizado pelo negócio?"
	CompanyID    uint       // ID da empresa/tenant (obrigatório - Sprint 3)
	Role         Role       // Papel do usuário na empresa (obrigatório - Sprint 3)
	DeletedAt    *time.Time // "O registro foi removido logicamente"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
