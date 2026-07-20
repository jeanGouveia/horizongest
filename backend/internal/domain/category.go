package domain

import "time"

type Category struct {
	ID           uint
	Name         string
	Description  string
	DisplayOrder int
	Active       bool
	CompanyID    uint // ID da empresa/tenant (obrigatório - Sprint 3)
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
