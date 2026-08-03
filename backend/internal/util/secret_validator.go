package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	ErrSecretTooShort      = errors.New("secret must be at least 32 bytes (256 bits)")
	ErrSecretLowEntropy    = errors.New("secret has insufficient entropy (use random base64)")
	ErrSecretInsecureValue = errors.New("secret contains insecure placeholder value")
)

// insecurePlaceholders are values that should never be used in production
var insecurePlaceholders = []string{
	"troque-este-valor",
	"change-in-production",
	"your-",
	"secret-key",
	"placeholder",
	"example",
	"test",
}

// ValidateJWTSecret validates a JWT secret for minimum security requirements
func ValidateJWTSecret(secret string, isProduction bool) error {
	if secret == "" {
		return errors.New("secret cannot be empty")
	}

	// Check minimum length (32 bytes = 256 bits)
	if len(secret) < 32 {
		return ErrSecretTooShort
	}

	// In production, check for insecure placeholders
	if isProduction {
		lowerSecret := strings.ToLower(secret)
		for _, placeholder := range insecurePlaceholders {
			if strings.Contains(lowerSecret, placeholder) {
				return ErrSecretInsecureValue
			}
		}
	}

	// Estimate entropy
	entropy := estimateEntropy(secret)
	minEntropy := 160 // 160 bits of entropy (NIST recommendation for 2024)
	if entropy < minEntropy {
		return fmt.Errorf("%w (estimated: %d bits, minimum: %d bits)", ErrSecretLowEntropy, entropy, minEntropy)
	}

	return nil
}

// estimateEntropy estimates the entropy of a string in bits
func estimateEntropy(s string) int {
	if len(s) == 0 {
		return 0
	}

	// Count character set size
	charSet := make(map[rune]bool)
	for _, c := range s {
		charSet[c] = true
	}

	charSetSize := len(charSet)

	// Calculate entropy based on character set size
	// entropy = length * log2(charSet_size)
	entropy := float64(len(s)) * log2(float64(charSetSize))

	return int(entropy)
}

// log2 calculates base-2 logarithm
func log2(x float64) float64 {
	const ln2 = 0.6931471805599453
	return 1.4426950408889634 * logNatural(x) // log2(x) = ln(x) / ln(2)
}

func logNatural(x float64) float64 {
	// Simple approximation for natural logarithm
	// For production, use math.Log from standard library
	if x <= 0 {
		return 0
	}
	if x == 1 {
		return 0
	}

	// Taylor series approximation around 1
	result := 0.0
	term := (x - 1) / (x + 1)
	termSquared := term * term
	for i := 0; i < 20; i++ {
		power := 2*float64(i) + 1
		result += term / power
		term *= termSquared
	}
	return 2 * result
}

// GenerateSecureSecret generates a cryptographically secure random secret
func GenerateSecureSecret(bytes int) (string, error) {
	if bytes < 32 {
		bytes = 32 // Minimum 32 bytes
	}

	randomBytes := make([]byte, bytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// GenerateSecureSecretHex generates a cryptographically secure random hex secret
func GenerateSecureSecretHex(bytes int) (string, error) {
	if bytes < 32 {
		bytes = 32 // Minimum 32 bytes
	}

	randomBytes := make([]byte, bytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return fmt.Sprintf("%x", randomBytes), nil
}

// GenerateSecureSecretAlphaNum generates a cryptographically secure alphanumeric secret
func GenerateSecureSecretAlphaNum(length int) (string, error) {
	if length < 32 {
		length = 32 // Minimum 32 characters
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
