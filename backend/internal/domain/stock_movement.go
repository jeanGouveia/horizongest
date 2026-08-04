package domain

import "time"

// StockMovementType representa o tipo de movimentação de estoque
type StockMovementType string

const (
	StockMovementEntry     StockMovementType = "entry"     // Entrada (compras, ajuste positivo)
	StockMovementExit      StockMovementType = "exit"      // Saída (produção, ajuste negativo)
	StockMovementAdjust    StockMovementType = "adjust"    // Ajuste manual
	StockMovementInventory StockMovementType = "inventory" // Inventário
)

// StockMovement representa uma movimentação de estoque
type StockMovement struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	CompanyID     uint              `json:"companyId" gorm:"not null;index:idx_stock_movements_company"`
	IngredientID  uint              `json:"ingredientId" gorm:"not null;index:idx_stock_movements_ingredient"`
	Type          StockMovementType `json:"type" gorm:"not null;size:20"`
	Quantity      float64           `json:"quantity" gorm:"not null"` // Positivo para entrada, negativo para saída
	PreviousStock float64           `json:"previousStock" gorm:"not null"`
	NewStock      float64           `json:"newStock" gorm:"not null"`
	Reason        string            `json:"reason" gorm:"size:500"`
	ReferenceType string            `json:"referenceType" gorm:"size:50"` // "purchase", "order", "adjustment", "inventory", "ifood_order"
	ReferenceID   uint              `json:"referenceId" gorm:"index"`
	PerformedBy   *uint             `json:"performedBy" gorm:"index"` // User ID (nullable for system operations like iFood webhooks)
	PerformedAt   time.Time         `json:"performedAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt     time.Time         `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time        `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Ingredient *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
	Performer  *User       `json:"performer,omitempty" gorm:"foreignKey:PerformedBy"`
}

// TableName especifica o nome da tabela
func (StockMovement) TableName() string {
	return "stock_movements"
}

// StockInventory representa um inventário de estoque
type StockInventory struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	CompanyID     uint       `json:"companyId" gorm:"not null;index:idx_stock_inventories_company"`
	InventoryDate time.Time  `json:"inventoryDate" gorm:"not null;index"`
	Status        string     `json:"status" gorm:"not null;size:20;default:'draft'"` // draft, completed, cancelled
	Notes         string     `json:"notes" gorm:"size:1000"`
	PerformedBy   uint       `json:"performedBy" gorm:"not null"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Items []StockInventoryItem `json:"items,omitempty" gorm:"foreignKey:InventoryID"`
}

// TableName especifica o nome da tabela
func (StockInventory) TableName() string {
	return "stock_inventories"
}

// StockInventoryItem representa um item de inventário
type StockInventoryItem struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	InventoryID   uint       `json:"inventoryId" gorm:"not null;index:idx_stock_inventory_items_inventory"`
	IngredientID  uint       `json:"ingredientId" gorm:"not null;index:idx_stock_inventory_items_ingredient"`
	ExpectedStock float64    `json:"expectedStock" gorm:"not null"` // Estoque esperado (sistema)
	ActualStock   float64    `json:"actualStock" gorm:"not null"`   // Estoque real (contado)
	Difference    float64    `json:"difference" gorm:"not null"`    // Diferença (actual - expected)
	Reason        string     `json:"reason" gorm:"size:500"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Ingredient *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
}

// TableName especifica o nome da tabela
func (StockInventoryItem) TableName() string {
	return "stock_inventory_items"
}
