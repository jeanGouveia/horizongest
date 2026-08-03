package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

// JWTKey represents a JWT signing key with metadata
type JWTKey struct {
	ID        string    // Key identifier (kid)
	Secret    string    // The actual secret
	CreatedAt time.Time // When the key was created
	ExpiresAt time.Time // When the key expires (for grace period)
	Active    bool      // Whether this is the active key
}

// JWTKeyStore manages JWT keys for rotation
type JWTKeyStore struct {
	mu            sync.RWMutex
	activeKey     *JWTKey
	previousKeys  []*JWTKey
	gracePeriod   time.Duration // How long previous keys remain valid
	keyExpiration time.Duration // How long until a key should be rotated
}

// NewJWTKeyStore creates a new JWT key store
func NewJWTKeyStore(initialSecret string, gracePeriod, keyExpiration time.Duration) (*JWTKeyStore, error) {
	store := &JWTKeyStore{
		gracePeriod:   gracePeriod,
		keyExpiration: keyExpiration,
	}

	// Generate key ID for initial secret
	keyID, err := generateKeyID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key ID: %w", err)
	}

	now := time.Now()
	store.activeKey = &JWTKey{
		ID:        keyID,
		Secret:    initialSecret,
		CreatedAt: now,
		ExpiresAt: now.Add(keyExpiration),
		Active:    true,
	}

	return store, nil
}

// GetActiveKey returns the current active signing key
func (ks *JWTKeyStore) GetActiveKey() *JWTKey {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	return ks.activeKey
}

// GetKeyByID returns a key by its ID (kid)
func (ks *JWTKeyStore) GetKeyByID(kid string) (*JWTKey, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	// Check active key
	if ks.activeKey.ID == kid {
		return ks.activeKey, true
	}

	// Check previous keys
	for _, key := range ks.previousKeys {
		if key.ID == kid {
			// Check if key is still within grace period
			if time.Now().Before(key.ExpiresAt) {
				return key, true
			}
			return nil, false
		}
	}

	return nil, false
}

// RotateKey creates a new active key and moves the current one to previous keys
func (ks *JWTKeyStore) RotateKey() error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	// Generate new key
	newSecret, err := generateSecureSecret(32)
	if err != nil {
		return fmt.Errorf("failed to generate new secret: %w", err)
	}

	newKeyID, err := generateKeyID()
	if err != nil {
		return fmt.Errorf("failed to generate key ID: %w", err)
	}

	now := time.Now()

	// Move current active key to previous keys
	if ks.activeKey != nil {
		ks.activeKey.Active = false
		ks.activeKey.ExpiresAt = now.Add(ks.gracePeriod)
		ks.previousKeys = append([]*JWTKey{ks.activeKey}, ks.previousKeys...)
	}

	// Set new active key
	ks.activeKey = &JWTKey{
		ID:        newKeyID,
		Secret:    newSecret,
		CreatedAt: now,
		ExpiresAt: now.Add(ks.keyExpiration),
		Active:    true,
	}

	// Clean up expired previous keys
	ks.cleanupExpiredKeys()

	return nil
}

// cleanupExpiredKeys removes keys that are past their grace period
func (ks *JWTKeyStore) cleanupExpiredKeys() {
	now := time.Now()
	validKeys := make([]*JWTKey, 0, len(ks.previousKeys))

	for _, key := range ks.previousKeys {
		if now.Before(key.ExpiresAt) {
			validKeys = append(validKeys, key)
		}
	}

	ks.previousKeys = validKeys
}

// RevokeKey revokes a specific key by ID
func (ks *JWTKeyStore) RevokeKey(kid string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	// Cannot revoke active key
	if ks.activeKey.ID == kid {
		return fmt.Errorf("cannot revoke active key")
	}

	// Remove from previous keys
	for i, key := range ks.previousKeys {
		if key.ID == kid {
			ks.previousKeys = append(ks.previousKeys[:i], ks.previousKeys[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("key not found")
}

// GetAllKeys returns all keys (for debugging/monitoring)
func (ks *JWTKeyStore) GetAllKeys() []*JWTKey {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	keys := make([]*JWTKey, 0, len(ks.previousKeys)+1)
	keys = append(keys, ks.activeKey)
	keys = append(keys, ks.previousKeys...)

	return keys
}

// generateSecureSecret generates a cryptographically secure random secret
func generateSecureSecret(bytes int) (string, error) {
	randomBytes := make([]byte, bytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// generateKeyID generates a unique key identifier
func generateKeyID() (string, error) {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return fmt.Sprintf("key-%x", randomBytes), nil
}
