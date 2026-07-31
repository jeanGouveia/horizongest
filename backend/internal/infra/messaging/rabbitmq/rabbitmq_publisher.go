package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implementa a interface EventPublisher usando RabbitMQ
type RabbitMQPublisher struct {
	connection      *Connection
	exchangeManager *ExchangeManager
	config          Config
}

// NewRabbitMQPublisher cria um novo publisher RabbitMQ
func NewRabbitMQPublisher(config Config) (*RabbitMQPublisher, error) {
	// Criar conexão
	conn, err := NewConnection(config)
	if err != nil {
		return nil, err
	}

	// Obter canal
	channel, err := conn.GetChannel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Criar gerenciador de exchanges
	exchangeManager := NewExchangeManager(channel, config)

	// Declarar exchange
	if err := exchangeManager.DeclareExchange(); err != nil {
		conn.Close()
		return nil, err
	}

	publisher := &RabbitMQPublisher{
		connection:      conn,
		exchangeManager: exchangeManager,
		config:          config,
	}

	log.Printf("RabbitMQPublisher initialized successfully")
	return publisher, nil
}

// Publish publica um único evento no RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
	return p.PublishBatch(ctx, []domain.OutboxEvent{event})
}

// PublishBatch publica múltiplos eventos em batch
func (p *RabbitMQPublisher) PublishBatch(ctx context.Context, events []domain.OutboxEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Verificar conexão
	if p.connection.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed")
	}

	// Obter canal
	channel, err := p.connection.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Publicar cada evento
	for _, event := range events {
		if err := p.publishSingle(ctx, channel, event); err != nil {
			return fmt.Errorf("failed to publish event id=%d: %w", event.ID, err)
		}
	}

	log.Printf("Published %d events successfully", len(events))
	return nil
}

// publishSingle publica um único evento sem retry de negócio
// Retry de negócio é responsabilidade do Dispatcher
// Este método apenas tenta publicar uma vez e falha rápido se houver erro
func (p *RabbitMQPublisher) publishSingle(ctx context.Context, channel *amqp.Channel, event domain.OutboxEvent) error {
	routingKey := p.exchangeManager.GetRoutingKey(event.AggregateType, event.EventType)

	// Preparar mensagem
	message := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Headers: amqp.Table{
			"event_id":       event.ID,
			"event_type":     event.EventType,
			"event_version":  event.EventVersion,
			"aggregate_type": event.AggregateType,
			"aggregate_id":   event.AggregateID,
			"tenant_id":      event.TenantID,
		},
		Body: []byte(event.Payload),
	}

	// Contexto com timeout
	publishCtx, cancel := context.WithTimeout(ctx, p.config.PublisherTimeout)
	defer cancel()

	// Publicar (sem retry - retry é responsabilidade do Dispatcher)
	err := channel.PublishWithContext(
		publishCtx,
		p.config.Exchange, // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		message,           // message
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// Publisher confirm se habilitado
	if p.config.EnablePublisherConfirm {
		ack := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		select {
		case confirm := <-ack:
			if !confirm.Ack {
				return fmt.Errorf("publisher confirm: message not acknowledged")
			}
		case <-publishCtx.Done():
			return fmt.Errorf("publisher confirm timeout")
		}
	}

	log.Printf("Event published: id=%d, type=%s, routing_key=%s", event.ID, event.EventType, routingKey)
	return nil
}

// Close fecha a conexão com o RabbitMQ
func (p *RabbitMQPublisher) Close() error {
	if p.connection != nil {
		return p.connection.Close()
	}
	return nil
}

// Ensure RabbitMQPublisher implements EventPublisher
var _ ports.EventPublisher = (*RabbitMQPublisher)(nil)

// MockEventPublisher é um mock para testes
type MockEventPublisher struct {
	PublishedEvents []domain.OutboxEvent
	PublishError    error
}

func (m *MockEventPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, event)
	return nil
}

func (m *MockEventPublisher) PublishBatch(ctx context.Context, events []domain.OutboxEvent) error {
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
}

func (m *MockEventPublisher) Close() error {
	return nil
}

// Ensure MockEventPublisher implements EventPublisher
var _ ports.EventPublisher = (*MockEventPublisher)(nil)
