package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          uint
	OrderNumber int // Número comercial do pedido (sequencial por empresa)
	Status      OrderStatus
	TotalPrice  float64
	Notes       string
	CompanyID   uint       // ID da empresa/tenant (obrigatório - Sprint 3)
	DeletedAt   *time.Time // "O registro foi removido logicamente"
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Preenchido sob demanda
	Items []OrderItem
}
