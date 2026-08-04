package config

import (
	"fmt"
	"time"
)

// Config holds application configuration
// FASE A.4: Centralized configuration with validation
type Config struct {
	// Server
	ServerPort         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration

	// Database
	DatabaseURL             string
	DatabaseMaxOpenConns    int
	DatabaseMaxIdleConns    int
	DatabaseConnMaxLifetime time.Duration
	DatabaseConnMaxIdleTime time.Duration

	// Redis
	RedisURL          string
	RedisPoolSize     int
	RedisMinIdleConns int
	RedisDialTimeout  time.Duration
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
	RedisPoolTimeout  time.Duration

	// RabbitMQ
	RabbitMQURL               string
	RabbitMQExchange          string
	RabbitMQQueue             string
	RabbitMQRoutingKey        string
	RabbitMQPrefetchCount     int
	RabbitMQReconnectInterval time.Duration

	// JWT
	JWTSecret             string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration

	// Storage
	StoragePath         string
	StorageMaxFileSize  int64
	StorageAllowedTypes []string

	// Environment
	Environment string
	LogLevel    string

	// Observability
	EnableTracing   bool
	EnableMetrics   bool
	EnableAudit     bool
	TracingEndpoint string
	MetricsEndpoint string

	// Shutdown
	ShutdownTimeout time.Duration
}

// LoadConfig loads configuration from environment variables
// FASE A.4: Load and validate configuration with defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		// Server
		ServerPort:         GetEnv("SERVER_PORT", "8080"),
		ServerReadTimeout:  time.Duration(GetEnvInt("SERVER_READ_TIMEOUT", 30)) * time.Second,
		ServerWriteTimeout: time.Duration(GetEnvInt("SERVER_WRITE_TIMEOUT", 30)) * time.Second,
		ServerIdleTimeout:  time.Duration(GetEnvInt("SERVER_IDLE_TIMEOUT", 120)) * time.Second,

		// Database
		DatabaseURL:             GetEnv("DATABASE_URL", "postgres://localhost:5432/horizongest"),
		DatabaseMaxOpenConns:    GetEnvInt("DATABASE_MAX_OPEN_CONNS", 25),
		DatabaseMaxIdleConns:    GetEnvInt("DATABASE_MAX_IDLE_CONNS", 5),
		DatabaseConnMaxLifetime: time.Duration(GetEnvInt("DATABASE_CONN_MAX_LIFETIME", 3600)) * time.Second,
		DatabaseConnMaxIdleTime: time.Duration(GetEnvInt("DATABASE_CONN_MAX_IDLE_TIME", 300)) * time.Second,

		// Redis
		RedisURL:          GetEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPoolSize:     GetEnvInt("REDIS_POOL_SIZE", 10),
		RedisMinIdleConns: GetEnvInt("REDIS_MIN_IDLE_CONNS", 2),
		RedisDialTimeout:  time.Duration(GetEnvInt("REDIS_DIAL_TIMEOUT", 5)) * time.Second,
		RedisReadTimeout:  time.Duration(GetEnvInt("REDIS_READ_TIMEOUT", 3)) * time.Second,
		RedisWriteTimeout: time.Duration(GetEnvInt("REDIS_WRITE_TIMEOUT", 3)) * time.Second,
		RedisPoolTimeout:  time.Duration(GetEnvInt("REDIS_POOL_TIMEOUT", 4)) * time.Second,

		// RabbitMQ
		RabbitMQURL:               GetEnv("RABBITMQ_URL", "amqp://localhost:5672"),
		RabbitMQExchange:          GetEnv("RABBITMQ_EXCHANGE", "horizongest"),
		RabbitMQQueue:             GetEnv("RABBITMQ_QUEUE", "events"),
		RabbitMQRoutingKey:        GetEnv("RABBITMQ_ROUTING_KEY", "events"),
		RabbitMQPrefetchCount:     GetEnvInt("RABBITMQ_PREFETCH_COUNT", 10),
		RabbitMQReconnectInterval: time.Duration(GetEnvInt("RABBITMQ_RECONNECT_INTERVAL", 5)) * time.Second,

		// JWT
		JWTSecret:             GetEnv("JWT_SECRET", "changeme"),
		JWTAccessTokenExpiry:  time.Duration(GetEnvInt("JWT_ACCESS_TOKEN_EXPIRY", 900)) * time.Second,     // 15 minutes
		JWTRefreshTokenExpiry: time.Duration(GetEnvInt("JWT_REFRESH_TOKEN_EXPIRY", 604800)) * time.Second, // 7 days

		// Storage
		StoragePath:         GetEnv("STORAGE_PATH", "./uploads"),
		StorageMaxFileSize:  int64(GetEnvInt("STORAGE_MAX_FILE_SIZE", 10485760)), // 10MB
		StorageAllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},

		// Environment
		Environment: GetEnv("ENVIRONMENT", "development"),
		LogLevel:    GetEnv("LOG_LEVEL", "info"),

		// Observability
		EnableTracing:   GetEnvBool("ENABLE_TRACING", true),
		EnableMetrics:   GetEnvBool("ENABLE_METRICS", true),
		EnableAudit:     GetEnvBool("ENABLE_AUDIT", true),
		TracingEndpoint: GetEnv("TRACING_ENDPOINT", "http://localhost:4318"),
		MetricsEndpoint: GetEnv("METRICS_ENDPOINT", ":9090"),

		// Shutdown
		ShutdownTimeout: time.Duration(GetEnvInt("SHUTDOWN_TIMEOUT", 30)) * time.Second,
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
// FASE A.4: Validate configuration values
func (c *Config) Validate() error {
	// Validate production environment
	if err := ValidateProductionEnv(); err != nil {
		return err
	}

	// Validate required fields
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Validate timeouts
	if c.ServerReadTimeout < time.Second {
		return fmt.Errorf("server read timeout must be at least 1 second")
	}
	if c.ServerWriteTimeout < time.Second {
		return fmt.Errorf("server write timeout must be at least 1 second")
	}

	// Validate JWT
	if c.JWTAccessTokenExpiry < time.Minute {
		return fmt.Errorf("JWT access token expiry must be at least 1 minute")
	}
	if c.JWTRefreshTokenExpiry < time.Hour {
		return fmt.Errorf("JWT refresh token expiry must be at least 1 hour")
	}

	// Validate storage
	if c.StorageMaxFileSize < 1024 {
		return fmt.Errorf("storage max file size must be at least 1KB")
	}

	return nil
}
