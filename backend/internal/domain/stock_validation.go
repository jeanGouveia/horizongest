package domain

// StockValidationRequest representa o payload para validação de estoque
type StockValidationRequest struct {
	Items []OrderItem `json:"items"`
}

// StockValidationResponse representa o resultado da validação
type StockValidationResponse struct {
	Valid             bool                     `json:"valid"`
	InsufficientStock []InsufficientIngredient `json:"insufficientStock,omitempty"`
}

// InsufficientIngredient representa um ingrediente com estoque insuficiente
type InsufficientIngredient struct {
	IngredientID   uint    `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Required       float64 `json:"required"`
	Available      float64 `json:"available"`
	Shortage       float64 `json:"shortage"`
	Unit           string  `json:"unit"`
}
