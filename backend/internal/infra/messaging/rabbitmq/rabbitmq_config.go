package rabbitmq

import "time"

// Config define as configurações do RabbitMQ Publisher
type Config struct {
	// URL é a URL de conexão do RabbitMQ
	// Exemplo: amqp://guest:guest@localhost:5672/
	URL string

	// Exchange é o nome da exchange onde os eventos serão publicados
	Exchange string

	// ExchangeType é o tipo da exchange (topic, direct, fanout)
	ExchangeType string

	// QueuePrefix é o prefixo para nomes de filas
	QueuePrefix string

	// RetryCount é o número máximo de tentativas de publicação
	RetryCount int

	// PublisherTimeout é o timeout para publicação de cada evento
	PublisherTimeout time.Duration

	// ReconnectDelay é o delay entre tentativas de reconexão
	ReconnectDelay time.Duration

	// EnablePublisherConfirm habilita publisher confirms
	EnablePublisherConfirm bool

	// Sprint 5D.3 - Performance Hardening: DLQ and Prefetch Configuration

	// DLQEnabled habilita Dead Letter Queue
	DLQEnabled bool

	// DLQName é o nome da Dead Letter Queue
	DLQName string

	// DLQTTL é o TTL da DLQ em milissegundos
	DLQTTL int64

	// DLQMaxRetries é o número máximo de retries antes de enviar para DLQ
	DLQMaxRetries int

	// PrefetchCount é o número de mensagens que um consumer pode receber antes de ack
	PrefetchCount int
}

// DefaultConfig retorna configurações padrão para desenvolvimento
func DefaultConfig() Config {
	return Config{
		URL:                    "amqp://guest:guest@localhost:5672/",
		Exchange:               "horizongest.events",
		ExchangeType:           "topic",
		QueuePrefix:            "horizongest",
		RetryCount:             3,
		PublisherTimeout:       30 * time.Second, // Sprint 5D.3: Increased from 10s to 30s
		ReconnectDelay:         5 * time.Second,
		EnablePublisherConfirm: true,
		// Sprint 5D.3 - Performance Hardening: DLQ and Prefetch
		DLQEnabled:    true,
		DLQName:       "horizongest.dlq",
		DLQTTL:        86400000, // 24 hours in milliseconds
		DLQMaxRetries: 3,
		PrefetchCount: 10, // Optimal for most workloads
	}
}
