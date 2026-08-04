package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

// ShutdownManager handles graceful shutdown of application components
// FASE A.4: Graceful Shutdown - Ensure clean shutdown without data loss
type ShutdownManager struct {
	shutdownTimeout time.Duration
	shutdownFuncs   []ShutdownFunc
	wg              sync.WaitGroup
	mu              sync.Mutex
}

// ShutdownFunc represents a function to be called during shutdown
type ShutdownFunc func(ctx context.Context) error

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		shutdownTimeout: timeout,
		shutdownFuncs:   make([]ShutdownFunc, 0),
	}
}

// Register registers a shutdown function
// FASE A.4: Register components for graceful shutdown
func (m *ShutdownManager) Register(name string, fn ShutdownFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.shutdownFuncs = append(m.shutdownFuncs, func(ctx context.Context) error {
		logger := util.GetLogger()
		logger.Info(fmt.Sprintf("Shutting down %s", name), nil)
		start := time.Now()
		err := fn(ctx)
		duration := time.Since(start)
		
		if err != nil {
			logger.Error(fmt.Sprintf("Error shutting down %s", name), map[string]interface{}{
				"error": err.Error(),
				"duration_ms": duration.Milliseconds(),
			})
		} else {
			logger.Info(fmt.Sprintf("Successfully shut down %s", name), map[string]interface{}{
				"duration_ms": duration.Milliseconds(),
			})
		}
		
		return err
	})
}

// WaitForShutdown waits for shutdown signal and executes shutdown functions
// FASE A.4: Graceful shutdown with signal handling
func (m *ShutdownManager) WaitForShutdown() error {
	logger := util.GetLogger()
	
	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Wait for signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", map[string]interface{}{
		"signal": sig.String(),
	})
	
	// Create shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()
	
	// Execute shutdown functions in reverse order (LIFO)
	m.mu.Lock()
	funcs := make([]ShutdownFunc, len(m.shutdownFuncs))
	copy(funcs, m.shutdownFuncs)
	m.mu.Unlock()
	
	// Execute in reverse order
	for i := len(funcs) - 1; i >= 0; i-- {
		m.wg.Add(1)
		go func(fn ShutdownFunc) {
			defer m.wg.Done()
			fn(ctx)
		}(funcs[i])
	}
	
	// Wait for all shutdown functions to complete
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		logger.Info("All components shut down successfully", nil)
		return nil
	case <-ctx.Done():
		logger.Error("Shutdown timeout exceeded", nil)
		return fmt.Errorf("shutdown timeout exceeded after %v", m.shutdownTimeout)
	}
}

// Shutdown manually triggers shutdown
// FASE A.4: Manual shutdown trigger
func (m *ShutdownManager) Shutdown() error {
	logger := util.GetLogger()
	
	logger.Info("Manual shutdown triggered", nil)
	
	ctx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()
	
	m.mu.Lock()
	funcs := make([]ShutdownFunc, len(m.shutdownFuncs))
	copy(funcs, m.shutdownFuncs)
	m.mu.Unlock()
	
	// Execute in reverse order
	for i := len(funcs) - 1; i >= 0; i-- {
		if err := funcs[i](ctx); err != nil {
			logger.Error("Error during manual shutdown", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
	
	return nil
}
