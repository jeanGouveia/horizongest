package shutdown

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestShutdownManager_Register(t *testing.T) {
	manager := NewShutdownManager(30 * time.Second)
	
	called := false
	manager.Register("test", func(ctx context.Context) error {
		called = true
		return nil
	})
	
	if !called {
		t.Error("shutdown function was not called")
	}
}

func TestShutdownManager_Shutdown(t *testing.T) {
	manager := NewShutdownManager(30 * time.Second)
	
	called := false
	manager.Register("test", func(ctx context.Context) error {
		called = true
		return nil
	})
	
	err := manager.Shutdown()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if !called {
		t.Error("shutdown function was not called")
	}
}

func TestShutdownManager_ShutdownWithError(t *testing.T) {
	manager := NewShutdownManager(30 * time.Second)
	
	manager.Register("test", func(ctx context.Context) error {
		return errors.New("test error")
	})
	
	err := manager.Shutdown()
	if err == nil {
		t.Error("expected error")
	}
}

func TestShutdownManager_ShutdownTimeout(t *testing.T) {
	manager := NewShutdownManager(100 * time.Millisecond)
	
	manager.Register("slow", func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	
	err := manager.Shutdown()
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestShutdownManager_ShutdownOrder(t *testing.T) {
	manager := NewShutdownManager(30 * time.Second)
	
	order := []string{}
	
	manager.Register("first", func(ctx context.Context) error {
		order = append(order, "first")
		return nil
	})
	
	manager.Register("second", func(ctx context.Context) error {
		order = append(order, "second")
		return nil
	})
	
	manager.Register("third", func(ctx context.Context) error {
		order = append(order, "third")
		return nil
	})
	
	manager.Shutdown()
	
	// Should be in reverse order (LIFO)
	expected := []string{"third", "second", "first"}
	if len(order) != len(expected) {
		t.Errorf("expected %d calls, got %d", len(expected), len(order))
	}
	
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected %s at position %d, got %s", v, i, order[i])
		}
	}
}
