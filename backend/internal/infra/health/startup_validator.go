package health

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jeanGouveia/horizongest/backend/internal/config"
	"github.com/jeanGouveia/horizongest/backend/internal/util"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// StartupValidator validates all critical dependencies on startup
// FASE A.4: Startup Validation - Fail-fast on critical dependency failure
type StartupValidator struct {
	cfg      *config.Config
	db       *pgx.Conn
	redis    *redis.Client
	rabbitMQ *amqp091.Connection
}

// NewStartupValidator creates a new startup validator
func NewStartupValidator(cfg *config.Config) *StartupValidator {
	return &StartupValidator{
		cfg: cfg,
	}
}

// Validate validates all critical dependencies
// FASE A.4: Abort startup if any critical dependency is unavailable
func (v *StartupValidator) Validate(ctx context.Context) error {
	logger := util.GetLogger()

	logger.Info("Starting dependency validation", nil)

	// Validate Database
	if err := v.validateDatabase(ctx); err != nil {
		logger.Fatal("Database validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("database validation failed: %w", err)
	}

	// Validate Redis
	if err := v.validateRedis(ctx); err != nil {
		logger.Fatal("Redis validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("redis validation failed: %w", err)
	}

	// Validate RabbitMQ
	if err := v.validateRabbitMQ(ctx); err != nil {
		logger.Fatal("RabbitMQ validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("rabbitmq validation failed: %w", err)
	}

	// Validate JWT Keys
	if err := v.validateJWTKeys(); err != nil {
		logger.Fatal("JWT keys validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("jwt keys validation failed: %w", err)
	}

	// Validate Storage
	if err := v.validateStorage(); err != nil {
		logger.Fatal("Storage validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("storage validation failed: %w", err)
	}

	logger.Info("All dependencies validated successfully", nil)
	return nil
}

// validateDatabase validates database connectivity
// FASE A.4: Ensure database is accessible and responsive
func (v *StartupValidator) validateDatabase(ctx context.Context) error {
	logger := util.GetLogger()

	logger.Info("Validating database connection", nil)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, v.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// Test query
	var result int
	err = conn.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("database query returned unexpected result")
	}

	v.db = conn
	logger.Info("Database validation passed", nil)
	return nil
}

// validateRedis validates Redis connectivity
// FASE A.4: Ensure Redis is accessible and responsive
func (v *StartupValidator) validateRedis(ctx context.Context) error {
	logger := util.GetLogger()

	logger.Info("Validating Redis connection", nil)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts, err := redis.ParseURL(v.cfg.RedisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	defer client.Close()

	// Test connection
	result, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	if result != "PONG" {
		return fmt.Errorf("redis ping returned unexpected result: %s", result)
	}

	v.redis = client
	logger.Info("Redis validation passed", nil)
	return nil
}

// validateRabbitMQ validates RabbitMQ connectivity
// FASE A.4: Ensure RabbitMQ is accessible and responsive
func (v *StartupValidator) validateRabbitMQ(ctx context.Context) error {
	logger := util.GetLogger()

	logger.Info("Validating RabbitMQ connection", nil)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := amqp091.Dial(v.cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Test connection
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer channel.Close()

	v.rabbitMQ = conn
	logger.Info("RabbitMQ validation passed", nil)
	return nil
}

// validateJWTKeys validates JWT configuration
// FASE A.4: Ensure JWT keys are properly configured
func (v *StartupValidator) validateJWTKeys() error {
	logger := util.GetLogger()

	logger.Info("Validating JWT configuration", nil)

	if v.cfg.JWTSecret == "" {
		return fmt.Errorf("JWT secret is not configured")
	}

	if len(v.cfg.JWTSecret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters")
	}

	if v.cfg.JWTAccessTokenExpiry < time.Minute {
		return fmt.Errorf("JWT access token expiry must be at least 1 minute")
	}

	if v.cfg.JWTRefreshTokenExpiry < time.Hour {
		return fmt.Errorf("JWT refresh token expiry must be at least 1 hour")
	}

	logger.Info("JWT validation passed", nil)
	return nil
}

// validateStorage validates storage configuration
// FASE A.4: Ensure storage is accessible
func (v *StartupValidator) validateStorage() error {
	logger := util.GetLogger()

	logger.Info("Validating storage configuration", nil)

	if v.cfg.StoragePath == "" {
		return fmt.Errorf("storage path is not configured")
	}

	// Note: Actual filesystem validation is done in the storage service
	// This just validates configuration

	logger.Info("Storage validation passed", nil)
	return nil
}
