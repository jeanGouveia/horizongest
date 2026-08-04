package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

// EnvValidator validates environment variables on startup
// FASE A.4: Environment Validation - Fail-fast on missing required variables
type EnvValidator struct {
	requiredVars map[string]bool
	optionalVars map[string]bool
}

// NewEnvValidator creates a new environment validator
func NewEnvValidator() *EnvValidator {
	return &EnvValidator{
		requiredVars: make(map[string]bool),
		optionalVars: make(map[string]bool),
	}
}

// Require marks a variable as required
func (v *EnvValidator) Require(name string) *EnvValidator {
	v.requiredVars[name] = true
	return v
}

// Optional marks a variable as optional
func (v *EnvValidator) Optional(name string) *EnvValidator {
	v.optionalVars[name] = true
	return v
}

// Validate validates all required variables
// FASE A.4: Fail-fast on missing required environment variables
func (v *EnvValidator) Validate() error {
	logger := util.GetLogger()
	
	missing := []string{}
	
	// Check required variables
	for name := range v.requiredVars {
		value := os.Getenv(name)
		if value == "" {
			missing = append(missing, name)
		} else {
			// Validate no secret is empty or default
			if strings.Contains(strings.ToLower(name), "secret") || 
			   strings.Contains(strings.ToLower(name), "key") ||
			   strings.Contains(strings.ToLower(name), "password") {
				if value == "changeme" || value == "secret" || value == "" {
					return fmt.Errorf("environment variable %s has insecure default value", name)
				}
			}
		}
	}
	
	if len(missing) > 0 {
		logger.Fatal("Missing required environment variables", map[string]interface{}{
			"missing": missing,
		})
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	
	logger.Info("Environment validation passed", nil)
	return nil
}

// GetEnv gets an environment variable with fallback
func GetEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvRequired gets a required environment variable
// FASE A.4: Fail-fast if required variable is missing
func GetEnvRequired(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", name)
	}
	return value, nil
}

// GetEnvInt gets an environment variable as integer
func GetEnvInt(name string, defaultValue int) int {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetEnvBool gets an environment variable as boolean
func GetEnvBool(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}

// ValidateProductionEnv validates production-specific requirements
// FASE A.4: Ensure production environment is properly configured
func ValidateProductionEnv() error {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		requiredSecrets := []string{
			"JWT_SECRET",
			"DATABASE_URL",
			"REDIS_URL",
			"RABBITMQ_URL",
		}
		
		for _, secret := range requiredSecrets {
			value := os.Getenv(secret)
			if value == "" {
				return fmt.Errorf("production environment requires %s to be set", secret)
			}
			
			// Check for insecure defaults
			if strings.Contains(strings.ToLower(secret), "secret") {
				if value == "changeme" || value == "secret" || value == "default" {
					return fmt.Errorf("production environment variable %s has insecure value", secret)
				}
			}
		}
	}
	return nil
}
