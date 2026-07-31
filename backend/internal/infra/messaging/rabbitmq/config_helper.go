package rabbitmq

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfigFromEnv carrega configurações do RabbitMQ de variáveis de ambiente
func LoadConfigFromEnv() Config {
	// Carregar .env se existir
	_ = godotenv.Load()

	return Config{
		URL:                    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		Exchange:               getEnv("RABBITMQ_EXCHANGE", "horizongest.events"),
		ExchangeType:           getEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		QueuePrefix:            getEnv("RABBITMQ_QUEUE_PREFIX", "horizongest"),
		RetryCount:             getEnvInt("RABBITMQ_RETRY_COUNT", 3),
		PublisherTimeout:       getEnvDuration("RABBITMQ_PUBLISHER_TIMEOUT", 10*time.Second),
		ReconnectDelay:         getEnvDuration("RABBITMQ_RECONNECT_DELAY", 5*time.Second),
		EnablePublisherConfirm: true,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
