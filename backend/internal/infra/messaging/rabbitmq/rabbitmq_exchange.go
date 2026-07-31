package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ExchangeManager gerencia declaração de exchanges e filas
type ExchangeManager struct {
	channel *amqp.Channel
	config  Config
}

// NewExchangeManager cria um novo gerenciador de exchanges
func NewExchangeManager(channel *amqp.Channel, config Config) *ExchangeManager {
	return &ExchangeManager{
		channel: channel,
		config:  config,
	}
}

// DeclareExchange declara a exchange principal
func (em *ExchangeManager) DeclareExchange() error {
	err := em.channel.ExchangeDeclare(
		em.config.Exchange,      // name
		em.config.ExchangeType,  // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Printf("Exchange declared: name=%s, type=%s", em.config.Exchange, em.config.ExchangeType)
	return nil
}

// DeclareQueue declara uma fila para um tenant específico
func (em *ExchangeManager) DeclareQueue(tenantID uint, eventType string) (amqp.Queue, error) {
	queueName := fmt.Sprintf("%s.%d.%s", em.config.QueuePrefix, tenantID, eventType)

	queue, err := em.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("Queue declared: name=%s, tenant_id=%d, event_type=%s", queueName, tenantID, eventType)
	return queue, nil
}

// DeclareDeadLetterQueue declara uma DLQ para eventos que falharam
func (em *ExchangeManager) DeclareDeadLetterQueue(tenantID uint) (amqp.Queue, error) {
	queueName := fmt.Sprintf("%s.%d.dlq", em.config.QueuePrefix, tenantID)

	queue, err := em.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		map[string]interface{}{
			"x-dead-letter-exchange": em.config.Exchange,
		},
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare dead letter queue: %w", err)
	}

	log.Printf("Dead letter queue declared: name=%s, tenant_id=%d", queueName, tenantID)
	return queue, nil
}

// BindQueue faz o binding de uma fila à exchange com uma routing key
func (em *ExchangeManager) BindQueue(queueName, routingKey string) error {
	err := em.channel.QueueBind(
		queueName,        // queue name
		routingKey,       // routing key
		em.config.Exchange, // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Printf("Queue bound: queue=%s, routing_key=%s, exchange=%s", queueName, routingKey, em.config.Exchange)
	return nil
}

// GetRoutingKey retorna a routing key para um evento
func (em *ExchangeManager) GetRoutingKey(aggregateType, eventType string) string {
	return fmt.Sprintf("%s.%s", aggregateType, eventType)
}
