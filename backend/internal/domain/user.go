package domain

import "time"

// User é a entidade central. Sem tags GORM aqui — domínio puro.
type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Active       bool       // "Pode ser utilizado pelo negócio?"
	DeletedAt    *time.Time // "O registro foi removido logicamente"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
