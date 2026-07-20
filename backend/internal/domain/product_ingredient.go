package domain

import "time"

// ProductIngredient define a ficha técnica:
// qual ingrediente e em qual quantidade compõe um produto.
// Exemplo: Pizza Margherita usa 0.150 kg de Queijo Mozzarella.
type ProductIngredient struct {
	ID           uint
	ProductID    uint
	IngredientID uint
	Quantity     float64 // quantidade consumida por unidade do produto
	// Sprint 4 - Ficha Técnica Avançada
	Loss      float64    // perda em % (ex: 10% = 0.10)
	Yield     float64    // rendimento em % (ex: 90% = 0.90)
	UnitCost  float64    // custo unitário do ingrediente (calculado)
	TotalCost float64    // custo total = quantity * unitCost / yield
	DeletedAt *time.Time // "O registro foi removido logicamente"

	// Preenchido em joins (leitura)
	Ingredient *Ingredient
}

// CalculateCost calcula o custo total do ingrediente no produto
// Custo = (Quantidade * CustoUnitário) / Rendimento
func (pi *ProductIngredient) CalculateCost() float64 {
	if pi.Yield == 0 {
		pi.Yield = 1.0 // Se não definido, assume 100%
	}
	cost := (pi.Quantity * pi.UnitCost) / pi.Yield
	pi.TotalCost = cost
	return cost
}

// GetEffectiveQuantity retorna a quantidade efetiva considerando perdas
// QuantidadeEfetiva = Quantidade / (1 - Perda)
func (pi *ProductIngredient) GetEffectiveQuantity() float64 {
	if pi.Loss >= 1.0 {
		pi.Loss = 0.0 // Se perda >= 100%, assume 0%
	}
	effectiveQuantity := pi.Quantity / (1.0 - pi.Loss)
	return effectiveQuantity
}
