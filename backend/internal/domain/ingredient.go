package domain

import "time"

type Ingredient struct {
	ID            uint
	Name          string
	Unit          string     // "kg", "L", "un", "g", "ml"
	StockQuantity float64    // quantidade atual em estoque
	MinStock      float64    // alerta de estoque mínimo (opcional)
	Active        bool       // "Pode ser utilizado pelo negócio?"
	CompanyID     *uint      // ID da empresa/tenant (null para compatibilidade com Core V1)
	DeletedAt     *time.Time // "O registro foi removido logicamente"
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
