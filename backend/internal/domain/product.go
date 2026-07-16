package domain

import "time"

type Product struct {
	ID                     uint
	Name                   string
	Description            string
	Price                  float64
	IsComposto             bool
	Active                 bool // "Pode ser utilizado pelo negócio?"
	PhotoURL               string
	CategoryID             *uint
	DisplayOrder           int
	PreparationTimeMinutes int
	Featured               bool
	IsNew                  bool
	PromotionPrice         *float64
	PromotionStart         *time.Time
	PromotionEnd           *time.Time
	AvailableFrom          string
	AvailableUntil         string
	SKU                    string
	InternalNotes          string
	DeletedAt              *time.Time // "O registro foi removido logicamente"
	CreatedAt              time.Time
	UpdatedAt              time.Time

	// SEO fields para Cardápio Digital
	Slug            string
	MetaTitle       string
	MetaDescription string
	AltImage        string
	Canonical       string

	// iFood integration fields
	ExternalID    string
	MarketplaceID string
	SyncStatus    string
	LastSync      *time.Time

	// Preenchido sob demanda (não vem do banco direto)
	Ingredients []ProductIngredient
}
