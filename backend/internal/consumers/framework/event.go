package framework

// Event represents a domain event from RabbitMQ
// This is a shared struct used by all consumers
type Event struct {
	ID            uint                   `json:"event_id"`
	EventType     string                 `json:"event_type"`
	AggregateType string                 `json:"aggregate_type"`
	AggregateID   uint                   `json:"aggregate_id"`
	TenantID      uint                   `json:"tenant_id"`
	Payload       map[string]interface{} `json:"payload"`
}
