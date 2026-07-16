package domain

import "time"

type Category struct {
	ID           uint
	Name         string
	Description  string
	DisplayOrder int
	Active       bool
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
