package service

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadDispatcherConfigFromEnv carrega configurações do Dispatcher de variáveis de ambiente
func LoadDispatcherConfigFromEnv() DispatcherConfig {
	_ = godotenv.Load()

	return DispatcherConfig{
		Interval:         getEnvDuration("DISPATCHER_INTERVAL", 5*time.Second),
		BatchSize:        getEnvInt("DISPATCHER_BATCH_SIZE", 50),
		RetryCount:       getEnvInt("DISPATCHER_RETRY_COUNT", 5),
		RetryBackoff:     getEnvDuration("DISPATCHER_RETRY_BACKOFF", 30*time.Second),
		PublisherTimeout: getEnvDuration("RABBITMQ_PUBLISHER_TIMEOUT", 10*time.Second),
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
