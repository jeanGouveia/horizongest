package config

import (
	"os"
	"testing"
)

func TestEnvValidator_Validate(t *testing.T) {
	// Save original env values
	origEnv := make(map[string]string)
	envVars := []string{"DATABASE_URL", "REDIS_URL", "RABBITMQ_URL", "JWT_SECRET"}
	for _, v := range envVars {
		origEnv[v] = os.Getenv(v)
	}
	
	// Restore after test
	defer func() {
		for k, v := range origEnv {
			os.Setenv(k, v)
		}
	}()
	
	t.Run("Required variables present", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		os.Setenv("REDIS_URL", "redis://localhost")
		os.Setenv("RABBITMQ_URL", "amqp://localhost")
		os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long")
		
		validator := NewEnvValidator()
		validator.Require("DATABASE_URL")
		validator.Require("REDIS_URL")
		validator.Require("RABBITMQ_URL")
		validator.Require("JWT_SECRET")
		
		err := validator.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
	
	t.Run("Required variable missing", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		os.Setenv("REDIS_URL", "redis://localhost")
		os.Setenv("RABBITMQ_URL", "amqp://localhost")
		os.Setenv("JWT_SECRET", "test-secret-key-32-chars-long")
		
		validator := NewEnvValidator()
		validator.Require("DATABASE_URL")
		validator.Require("REDIS_URL")
		validator.Require("RABBITMQ_URL")
		validator.Require("JWT_SECRET")
		
		err := validator.Validate()
		if err == nil {
			t.Error("expected error for missing required variable")
		}
	})
	
	t.Run("Insecure secret value", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		os.Setenv("REDIS_URL", "redis://localhost")
		os.Setenv("RABBITMQ_URL", "amqp://localhost")
		os.Setenv("JWT_SECRET", "changeme")
		
		validator := NewEnvValidator()
		validator.Require("DATABASE_URL")
		validator.Require("REDIS_URL")
		validator.Require("RABBITMQ_URL")
		validator.Require("JWT_SECRET")
		
		err := validator.Validate()
		if err == nil {
			t.Error("expected error for insecure secret value")
		}
	})
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")
	
	result := GetEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("expected test_value, got %s", result)
	}
	
	result = GetEnv("NON_EXISTENT", "default")
	if result != "default" {
		t.Errorf("expected default, got %s", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	
	result := GetEnvInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
	
	result = GetEnvInt("NON_EXISTENT", 10)
	if result != 10 {
		t.Errorf("expected 10, got %d", result)
	}
}

func TestGetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL_TRUE", "true")
	os.Setenv("TEST_BOOL_FALSE", "false")
	os.Setenv("TEST_BOOL_1", "1")
	defer func() {
		os.Unsetenv("TEST_BOOL_TRUE")
		os.Unsetenv("TEST_BOOL_FALSE")
		os.Unsetenv("TEST_BOOL_1")
	}()
	
	if !GetEnvBool("TEST_BOOL_TRUE", false) {
		t.Error("expected true")
	}
	
	if GetEnvBool("TEST_BOOL_FALSE", true) {
		t.Error("expected false")
	}
	
	if !GetEnvBool("TEST_BOOL_1", false) {
		t.Error("expected true")
	}
	
	if !GetEnvBool("NON_EXISTENT", true) {
		t.Error("expected default true")
	}
}

func TestValidateProductionEnv(t *testing.T) {
	// Save original env
	origEnv := os.Getenv("ENVIRONMENT")
	defer os.Setenv("ENVIRONMENT", origEnv)
	
	t.Run("Production with all secrets", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		os.Setenv("REDIS_URL", "redis://localhost")
		os.Setenv("RABBITMQ_URL", "amqp://localhost")
		os.Setenv("JWT_SECRET", "secure-secret-key-32-chars-long")
		
		err := ValidateProductionEnv()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		
		// Cleanup
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("JWT_SECRET")
	})
	
	t.Run("Production missing secret", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Unsetenv("DATABASE_URL")
		
		err := ValidateProductionEnv()
		if err == nil {
			t.Error("expected error for missing secret")
		}
	})
	
	t.Run("Production insecure secret", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("DATABASE_URL", "postgres://localhost/test")
		os.Setenv("REDIS_URL", "redis://localhost")
		os.Setenv("RABBITMQ_URL", "amqp://localhost")
		os.Setenv("JWT_SECRET", "changeme")
		
		err := ValidateProductionEnv()
		if err == nil {
			t.Error("expected error for insecure secret")
		}
		
		// Cleanup
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("JWT_SECRET")
	})
	
	t.Run("Development environment", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "development")
		
		err := ValidateProductionEnv()
		if err != nil {
			t.Errorf("expected no error in development, got %v", err)
		}
	})
}
