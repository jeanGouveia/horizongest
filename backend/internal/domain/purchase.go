package domain

import "time"

// Supplier representa um fornecedor
type Supplier struct {
	ID        uint       `json:"id"`
	CompanyID uint       `json:"companyId" gorm:"not null;index"`
	Name      string     `json:"name" gorm:"not null"`
	CNPJ      string     `json:"cnpj"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Address   string     `json:"address"`
	City      string     `json:"city"`
	State     string     `json:"state"`
	ZipCode   string     `json:"zipCode"`
	Notes     string     `json:"notes"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

// TableName especifica o nome da tabela
func (Supplier) TableName() string {
	return "suppliers"
}

// PurchaseOrderStatus representa o status de um pedido de compra
type PurchaseOrderStatus string

const (
	PurchaseOrderDraft     PurchaseOrderStatus = "draft"
	PurchaseOrderSent      PurchaseOrderStatus = "sent"
	PurchaseOrderConfirmed PurchaseOrderStatus = "confirmed"
	PurchaseOrderReceived  PurchaseOrderStatus = "received"
	PurchaseOrderCancelled PurchaseOrderStatus = "cancelled"
)

// PurchaseOrder representa um pedido de compra
type PurchaseOrder struct {
	ID           uint                `json:"id"`
	CompanyID    uint                `json:"companyId" gorm:"not null;index"`
	SupplierID   uint                `json:"supplierId" gorm:"not null;index"`
	OrderNumber  string              `json:"orderNumber" gorm:"uniqueIndex"`
	Status       PurchaseOrderStatus `json:"status" gorm:"not null;default:'draft'"`
	OrderDate    time.Time           `json:"orderDate" gorm:"not null"`
	ExpectedDate *time.Time          `json:"expectedDate"`
	ReceivedDate *time.Time          `json:"receivedDate"`
	Subtotal     Money               `json:"subtotal" gorm:"default:0"`
	Tax          Money               `json:"tax" gorm:"default:0"`
	Discount     Money               `json:"discount" gorm:"default:0"`
	Total        Money               `json:"total" gorm:"default:0"`
	Notes        string              `json:"notes"`
	CreatedBy    uint                `json:"createdBy" gorm:"not null"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	DeletedAt    *time.Time          `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Supplier *Supplier           `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Items    []PurchaseOrderItem `json:"items,omitempty" gorm:"foreignKey:PurchaseOrderID"`
}

// TableName especifica o nome da tabela
func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem representa um item de pedido de compra
type PurchaseOrderItem struct {
	ID              uint       `json:"id"`
	PurchaseOrderID uint       `json:"purchaseOrderId" gorm:"not null;index"`
	IngredientID    uint       `json:"ingredientId" gorm:"not null;index"`
	Quantity        float64    `json:"quantity" gorm:"not null"`
	Unit            string     `json:"unit" gorm:"not null"`
	UnitPrice       Money      `json:"unitPrice" gorm:"not null"`
	Subtotal        Money      `json:"subtotal" gorm:"default:0"`
	ReceivedQty     float64    `json:"receivedQty" gorm:"default:0"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Ingredient *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
}

// TableName especifica o nome da tabela
func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// PurchaseReceiving representa um recebimento de compra
type PurchaseReceiving struct {
	ID              uint       `json:"id"`
	PurchaseOrderID uint       `json:"purchaseOrderId" gorm:"not null;index"`
	ReceivedDate    time.Time  `json:"receivedDate" gorm:"not null"`
	ReceivedBy      uint       `json:"receivedBy" gorm:"not null"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Items []PurchaseReceivingItem `json:"items,omitempty" gorm:"foreignKey:PurchaseReceivingID"`
}

// TableName especifica o nome da tabela
func (PurchaseReceiving) TableName() string {
	return "purchase_receivings"
}

// PurchaseReceivingItem representa um item de recebimento
type PurchaseReceivingItem struct {
	ID                  uint       `json:"id"`
	PurchaseReceivingID uint       `json:"purchaseReceivingId" gorm:"not null;index"`
	PurchaseOrderItemID uint       `json:"purchaseOrderItemId" gorm:"not null;index"`
	IngredientID        uint       `json:"ingredientId" gorm:"not null;index"`
	Quantity            float64    `json:"quantity" gorm:"not null"`
	Unit                string     `json:"unit" gorm:"not null"`
	UnitPrice           Money      `json:"unitPrice" gorm:"not null"`
	Subtotal            Money      `json:"subtotal" gorm:"default:0"`
	Notes               string     `json:"notes"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Ingredient *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
}

// TableName especifica o nome da tabela
func (PurchaseReceivingItem) TableName() string {
	return "purchase_receiving_items"
}
