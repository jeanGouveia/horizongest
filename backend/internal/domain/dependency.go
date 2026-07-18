package domain

// DependencyCheck representa o resultado de uma verificação de dependência
type DependencyCheck struct {
	CanDelete bool               `json:"canDelete"`
	Reasons   []DependencyReason `json:"reasons"`
}

type DependencyReason struct {
	Type        string `json:"type"` // "order", "product", "recipe", etc.
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProductDependencies representa dependências que impedem exclusão de produto
type ProductDependencies struct {
	Orders []struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
		Date   string `json:"date"`
	} `json:"orders"`
	Recipes []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"recipes"`
}

// IngredientDependencies representa dependências que impedem exclusão de ingrediente
type IngredientDependencies struct {
	Products []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"products"`
}

// CategoryDependencies representa dependências que impedem exclusão de categoria
type CategoryDependencies struct {
	Products []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"products"`
}
