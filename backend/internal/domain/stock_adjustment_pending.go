package domain

import "time"

// StockAdjustmentStatus representa o status de um ajuste de estoque pendente
type StockAdjustmentStatus string

const (
	StockAdjustmentStatusPending  StockAdjustmentStatus = "pending"
	StockAdjustmentStatusApproved StockAdjustmentStatus = "approved"
	StockAdjustmentStatusRejected StockAdjustmentStatus = "rejected"
)

// StockAdjustmentPending identifica ajustes de estoque que precisam de aprovação manual
// ou análise antes de serem aplicados. Usado para estornos de estoque por cancelamento
// de pedidos, permitindo auditoria e validação humana.
// Todos os campos de snapshot (IngredientName, IngredientUnit) são congelados no momento
// do cancelamento e nunca alterados (princípio #4: Histórico é imutável).
type StockAdjustmentPending struct {
	ID              uint                  `json:"id"`
	OrderID         uint                  `json:"order_id"`
	IngredientID    uint                  `json:"ingredient_id"`
	Quantity        float64               `json:"quantity"`     // Quantidade que poderia ser devolvida ao estoque
	OrderStatus     string                `json:"order_status"` // Status do pedido no momento do cancelamento (para contexto)
	Status          StockAdjustmentStatus `json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	ProcessedAt     *time.Time            `json:"processed_at,omitempty"`     // Quando foi aprovado/rejeitado (null se pending)
	ProcessedBy     *uint                 `json:"processed_by,omitempty"`     // ID do usuário que aprovou/rejeitou (null se pending)
	ProcessingNotes string                `json:"processing_notes,omitempty"` // Observações do operador ao aprovar/rejeitar
	IngredientName  string                `json:"ingredient_name"`            // snapshot do nome do ingrediente no momento do cancelamento
	IngredientUnit  string                `json:"ingredient_unit"`            // snapshot da unidade do ingrediente no momento do cancelamento
	DeletedAt       *time.Time            `json:"deleted_at,omitempty"`       // "O registro foi removido logicamente"
}
