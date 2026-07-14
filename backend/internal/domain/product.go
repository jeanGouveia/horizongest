package domain

import "time"

type Product struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	IsComposto  bool
	Active      bool       // "Pode ser utilizado pelo negócio?"
	DeletedAt   *time.Time // "O registro foi removido logicamente"
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Preenchido sob demanda (não vem do banco direto)
	Ingredients []ProductIngredient
}
