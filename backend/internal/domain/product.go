package domain

import (
	"errors"
	"math"
	"time"
)

type Product struct {
	ID                     uint
	Name                   string
	Description            string
	Price                  Money
	IsComposto             bool
	Active                 bool // "Pode ser utilizado pelo negócio?"
	PhotoURL               string
	CategoryID             *uint
	CompanyID              uint // ID da empresa/tenant (obrigatório - Sprint 3)
	DisplayOrder           int
	PreparationTimeMinutes int
	Featured               bool
	IsNew                  bool
	PromotionPrice         *Money
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
	Cost           Money   // custo total do produto (soma dos ingredientes)
	CMV            Money   // custo merca da venda (CMV = Cost / Price)
	Margin         float64 // margem em % (Margin = (Price - Cost) / Price)
	Profit         Money   // lucro por unidade (Profit = Price - Cost)
	SuggestedPrice Money   // preço sugerido baseado no custo e margem desejada

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
func (p *Product) CalculateCost() Money {
	totalCost := Money(0)
	for _, ingredient := range p.Ingredients {
		totalCost = totalCost.Add(ingredient.CalculateCost())
	}
	p.Cost = totalCost
	return totalCost
}

// CalculateCMV calcula o custo merca da venda
// CMV = Custo / Preço (retorna como float64 percentual)
func (p *Product) CalculateCMV() float64 {
	if p.Price.IsZero() {
		p.CMV = Money(0)
		return 0
	}
	cmv := float64(p.Cost) / float64(p.Price)
	p.CMV = Money(math.Round(cmv * 100)) // Armazena como centavos de percentual
	return cmv
}

// CalculateMargin calcula a margem de lucro
// Margem = (Preço - Custo) / Preço (retorna como float64 percentual)
func (p *Product) CalculateMargin() float64 {
	if p.Price.IsZero() {
		p.Margin = 0
		return 0
	}
	profit := float64(p.Price.Sub(p.Cost)) / 100.0
	price := p.Price.ToFloat64()
	p.Margin = profit / price
	return p.Margin
}

// CalculateProfit calcula o lucro por unidade
// Lucro = Preço - Custo
func (p *Product) CalculateProfit() Money {
	p.Profit = p.Price.Sub(p.Cost)
	return p.Profit
}

// CalculateSuggestedPrice calcula o preço sugerido baseado no custo e margem desejada
// PreçoSugerido = Custo / (1 - MargemDesejada)
func (p *Product) CalculateSuggestedPrice(desiredMargin float64) Money {
	if desiredMargin >= 1.0 {
		desiredMargin = 0.5 // Se margem >= 100%, assume 50%
	}
	if desiredMargin < 0 {
		desiredMargin = 0.3 // Se margem negativa, assume 30%
	}
	costFloat := p.Cost.ToFloat64()
	suggestedPrice := costFloat / (1.0 - desiredMargin)
	p.SuggestedPrice = FromFloat64(suggestedPrice)
	return p.SuggestedPrice
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
