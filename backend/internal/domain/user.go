package domain

import "time"

// User é a entidade central. Sem tags GORM aqui — domínio puro.
type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool       // "Pode ser utilizado pelo negócio?"
	CompanyID    *uint      // ID da empresa/tenant (null para compatibilidade com Core V1)
	Role         *Role      // Papel do usuário na empresa (null para compatibilidade com Core V1)
	DeletedAt    *time.Time // "O registro foi removido logicamente"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
