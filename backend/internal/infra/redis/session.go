package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rediscmd "github.com/redis/go-redis/v9"
)

const (
	sessionKeyPrefix = "session"
)

// sessionKey generates a namespaced session key
func sessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s:%s", Namespace, sessionKeyPrefix, sessionID)
}

// SessionStore defines the session storage interface
type SessionStore interface {
	// Get retrieves a session by ID
	Get(ctx context.Context, sessionID string, dest interface{}) error

	// Set stores a session with TTL
	Set(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error

	// Delete removes a session
	Delete(ctx context.Context, sessionID string) error

	// Exists checks if a session exists
	Exists(ctx context.Context, sessionID string) (bool, error)

	// Refresh extends the TTL of a session
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error

	// Clear removes all sessions for a user
	Clear(ctx context.Context, userID string) error
}

// RedisSessionStore implements SessionStore using Redis
type RedisSessionStore struct {
	client *Client
}

// NewSessionStore creates a new Redis session store
func NewSessionStore(client *Client) SessionStore {
	return &RedisSessionStore{client: client}
}

// Get retrieves a session by ID
func (s *RedisSessionStore) Get(ctx context.Context, sessionID string, dest interface{}) error {
	val, err := s.client.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		if err == rediscmd.Nil {
			return fmt.Errorf("session not found: %s", sessionID)
		}
		return fmt.Errorf("failed to get session %s: %w", sessionID, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal session %s: %w", sessionID, err)
	}

	return nil
}

// Set stores a session with TTL
func (s *RedisSessionStore) Set(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error {
	sessionData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session %s: %w", sessionID, err)
	}

	if err := s.client.Set(ctx, sessionKey(sessionID), sessionData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set session %s: %w", sessionID, err)
	}

	return nil
}

// Delete removes a session
func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	if err := s.client.Del(ctx, sessionKey(sessionID)).Err(); err != nil {
		return fmt.Errorf("failed to delete session %s: %w", sessionID, err)
	}
	return nil
}

// Exists checks if a session exists
func (s *RedisSessionStore) Exists(ctx context.Context, sessionID string) (bool, error) {
	result, err := s.client.Exists(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check session existence %s: %w", sessionID, err)
	}
	return result > 0, nil
}

// Refresh extends the TTL of a session
func (s *RedisSessionStore) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	if err := s.client.Expire(ctx, sessionKey(sessionID), ttl).Err(); err != nil {
		return fmt.Errorf("failed to refresh session %s: %w", sessionID, err)
	}
	return nil
}

// Clear removes all sessions for a user
// This assumes sessions are stored with a pattern like "session:{userID}:{sessionID}"
func (s *RedisSessionStore) Clear(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s:%s:%s:*", Namespace, sessionKeyPrefix, userID)

	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	keys := make([]string, 0)

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan sessions for user %s: %w", userID, err)
	}

	if len(keys) > 0 {
		if err := s.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to clear sessions for user %s: %w", userID, err)
		}
	}

	return nil
}
