package ports

import (
	"context"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// EventPublisher define a interface para publicação de eventos
// Implementa o padrão Port Adapter para desacoplar domínio de infraestrutura
// O domínio conhece apenas esta interface, não RabbitMQ ou outros brokers
type EventPublisher interface {
	// Publish publica um único evento no message broker
	// Retorna erro se a publicação falhar
	Publish(ctx context.Context, event domain.OutboxEvent) error

	// PublishBatch publica múltiplos eventos em batch
	// Mais eficiente que publicações individuais
	PublishBatch(ctx context.Context, events []domain.OutboxEvent) error

	// Close fecha a conexão com o message broker
	// Deve ser chamado durante graceful shutdown
	Close() error
}
