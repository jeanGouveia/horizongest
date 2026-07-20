package domain

import (
	"errors"
	"time"
)

type Product struct {
	ID                     uint
	Name                   string
	Description            string
	Price                  float64
	IsComposto             bool
	Active                 bool // "Pode ser utilizado pelo negócio?"
	PhotoURL               string
	CategoryID             *uint
	CompanyID              uint // ID da empresa/tenant (obrigatório - Sprint 3)
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

	// Sprint 4 - Ficha Técnica Avançada
	Cost           float64 // custo total do produto (soma dos ingredientes)
	CMV            float64 // custo merca da venda (CMV = Cost / Price)
	Margin         float64 // margem em % (Margin = (Price - Cost) / Price)
	Profit         float64 // lucro por unidade (Profit = Price - Cost)
	SuggestedPrice float64 // preço sugerido baseado no custo e margem desejada

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

// CalculateCost calcula o custo total do produto baseado nos ingredientes
func (p *Product) CalculateCost() float64 {
	totalCost := 0.0
	for _, ingredient := range p.Ingredients {
		totalCost += ingredient.CalculateCost()
	}
	p.Cost = totalCost
	return totalCost
}

// CalculateCMV calcula o custo merca da venda
// CMV = Custo / Preço
func (p *Product) CalculateCMV() float64 {
	if p.Price == 0 {
		p.CMV = 0
		return 0
	}
	p.CMV = p.Cost / p.Price
	return p.CMV
}

// CalculateMargin calcula a margem de lucro
// Margem = (Preço - Custo) / Preço
func (p *Product) CalculateMargin() float64 {
	if p.Price == 0 {
		p.Margin = 0
		return 0
	}
	p.Margin = (p.Price - p.Cost) / p.Price
	return p.Margin
}

// CalculateProfit calcula o lucro por unidade
// Lucro = Preço - Custo
func (p *Product) CalculateProfit() float64 {
	p.Profit = p.Price - p.Cost
	return p.Profit
}

// CalculateSuggestedPrice calcula o preço sugerido baseado no custo e margem desejada
// PreçoSugerido = Custo / (1 - MargemDesejada)
func (p *Product) CalculateSuggestedPrice(desiredMargin float64) float64 {
	if desiredMargin >= 1.0 {
		desiredMargin = 0.5 // Se margem >= 100%, assume 50%
	}
	if desiredMargin < 0 {
		desiredMargin = 0.3 // Se margem negativa, assume 30%
	}
	suggestedPrice := p.Cost / (1.0 - desiredMargin)
	p.SuggestedPrice = suggestedPrice
	return suggestedPrice
}

// HasRecipe verifica se o produto tem ficha técnica
func (p *Product) HasRecipe() bool {
	return len(p.Ingredients) > 0
}

// ValidateRecipe valida a ficha técnica do produto
// Retorna erro se:
// - Produto não tem ficha técnica
// - Ingrediente não existe
// - Ingrediente está inativo
func (p *Product) ValidateRecipe() error {
	if !p.HasRecipe() {
		return errors.New("produto sem ficha técnica")
	}
	for _, ingredient := range p.Ingredients {
		if ingredient.Ingredient == nil {
			return errors.New("ingrediente não encontrado")
		}
		if !ingredient.Ingredient.Active {
			return errors.New("ingrediente inativo")
		}
	}
	return nil
}
