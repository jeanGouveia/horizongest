package security

import (
	"testing"
	"time"
)

func TestJWTKeyStore_NewJWTKeyStore(t *testing.T) {
	store, err := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)
	if err != nil {
		t.Fatalf("NewJWTKeyStore failed: %v", err)
	}

	activeKey := store.GetActiveKey()
	if activeKey == nil {
		t.Fatal("expected active key, got nil")
	}
	if activeKey.Secret != "test-secret" {
		t.Errorf("expected secret test-secret, got %s", activeKey.Secret)
	}
	if !activeKey.Active {
		t.Error("expected active key to be marked as active")
	}
	if activeKey.ID == "" {
		t.Error("expected key ID to be non-empty")
	}
}

func TestJWTKeyStore_GetKeyByID_Active(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	activeKey := store.GetActiveKey()
	key, found := store.GetKeyByID(activeKey.ID)
	if !found {
		t.Error("expected to find active key by ID")
	}
	if key != activeKey {
		t.Error("expected to get the same active key")
	}
}

func TestJWTKeyStore_GetKeyByID_NotFound(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	_, found := store.GetKeyByID("non-existent")
	if found {
		t.Error("expected not to find non-existent key")
	}
}

func TestJWTKeyStore_RotateKey(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	oldActiveKey := store.GetActiveKey()
	oldKid := oldActiveKey.ID

	err := store.RotateKey()
	if err != nil {
		t.Fatalf("RotateKey failed: %v", err)
	}

	// Verify new active key
	newActiveKey := store.GetActiveKey()
	if newActiveKey.ID == oldKid {
		t.Error("expected new active key to have different ID")
	}
	if newActiveKey.Secret == oldActiveKey.Secret {
		t.Error("expected new active key to have different secret")
	}

	// Verify old key is now in previous keys
	oldKey, found := store.GetKeyByID(oldKid)
	if !found {
		t.Error("expected to find old key in previous keys")
	}
	if oldKey.Active {
		t.Error("expected old key to be marked as inactive")
	}
}

func TestJWTKeyStore_RevokeKey(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	// Rotate to create a previous key
	store.RotateKey()
	oldActiveKey := store.GetActiveKey()
	store.RotateKey()

	// Revoke the previous key
	err := store.RevokeKey(oldActiveKey.ID)
	if err != nil {
		t.Fatalf("RevokeKey failed: %v", err)
	}

	// Verify key was revoked
	_, found := store.GetKeyByID(oldActiveKey.ID)
	if found {
		t.Error("expected revoked key to not be found")
	}
}

func TestJWTKeyStore_RevokeKey_Active(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	activeKey := store.GetActiveKey()
	err := store.RevokeKey(activeKey.ID)
	if err == nil {
		t.Error("expected error when trying to revoke active key")
	}
}

func TestJWTKeyStore_RevokeKey_NotFound(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	err := store.RevokeKey("non-existent")
	if err == nil {
		t.Error("expected error when revoking non-existent key")
	}
}

func TestJWTKeyStore_GetAllKeys(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 30*24*time.Hour, 90*24*time.Hour)

	keys := store.GetAllKeys()
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}

	// Rotate to add previous keys
	store.RotateKey()
	keys = store.GetAllKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys after rotation, got %d", len(keys))
	}
}

func TestJWTKeyStore_PreviousKeyExpiration(t *testing.T) {
	store, _ := NewJWTKeyStore("test-secret", 1*time.Hour, 90*24*time.Hour)

	// Rotate to create a previous key
	store.RotateKey()
	oldActiveKey := store.GetActiveKey()
	store.RotateKey()

	// Manually expire the previous key
	store.mu.Lock()
	for _, key := range store.previousKeys {
		if key.ID == oldActiveKey.ID {
			key.ExpiresAt = time.Now().Add(-1 * time.Hour)
			break
		}
	}
	store.mu.Unlock()

	// Try to get the expired key
	_, found := store.GetKeyByID(oldActiveKey.ID)
	if found {
		t.Error("expected expired previous key to not be found")
	}
}
