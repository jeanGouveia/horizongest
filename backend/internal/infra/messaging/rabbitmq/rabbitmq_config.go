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
}

// DefaultConfig retorna configurações padrão para desenvolvimento
func DefaultConfig() Config {
	return Config{
		URL:                    "amqp://guest:guest@localhost:5672/",
		Exchange:               "horizongest.events",
		ExchangeType:           "topic",
		QueuePrefix:            "horizongest",
		RetryCount:             3,
		PublisherTimeout:       10 * time.Second,
		ReconnectDelay:         5 * time.Second,
		EnablePublisherConfirm: true,
	}
}
