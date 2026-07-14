package domain

import "time"

// OrderItem representa um produto dentro de um pedido.
// Todos os campos de snapshot são congelados no momento da venda e nunca alterados.
// Princípio #4: Histórico é imutável.
type OrderItem struct {
	ID                 uint
	OrderID            uint
	ProductID          uint
	Quantity           float64
	UnitPrice          float64    // snapshot do preço no momento do pedido
	ProductName        string     // snapshot do nome do produto no momento do pedido
	ProductDescription string     // snapshot da descrição do produto no momento do pedido
	ProductIsComposto  bool       // snapshot da flag is_composto no momento do pedido
	DeletedAt          *time.Time // "O registro foi removido logicamente"

	// Preenchido em joins (leitura) - apenas para navegação, não para dados históricos
	Product *Product
}
