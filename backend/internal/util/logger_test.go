package util

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLogger_LogEntryStructure(t *testing.T) {
	logger := NewLogger("test-service", "test")
	
	fields := map[string]interface{}{
		"request_id":     "req-123",
		"correlation_id": "corr-456",
		"company_id":     uint(789),
		"user_id":        uint(101),
		"operation":      "test_operation",
		"duration_ms":    int64(100),
		"error":          "test error",
	}
	
	logger.Info("test message", fields)
	
	// The logger outputs to stdout, so we can't easily capture it here
	// But we can test the LogEntry structure
	entry := LogEntry{
		Timestamp:    "2024-01-01T00:00:00Z",
		Level:        LevelInfo,
		Message:      "test message",
		RequestID:    "req-123",
		CorrelationID: "corr-456",
		CompanyID:    789,
		UserID:       101,
		Service:      "test-service",
		Operation:    "test_operation",
		DurationMs:   100,
		Error:        "test error",
		Environment:  "test",
	}
	
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal log entry: %v", err)
	}
	
	var unmarshaled LogEntry
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}
	
	if unmarshaled.RequestID != "req-123" {
		t.Errorf("expected request_id req-123, got %s", unmarshaled.RequestID)
	}
	if unmarshaled.CorrelationID != "corr-456" {
		t.Errorf("expected correlation_id corr-456, got %s", unmarshaled.CorrelationID)
	}
	if unmarshaled.CompanyID != 789 {
		t.Errorf("expected company_id 789, got %d", unmarshaled.CompanyID)
	}
	if unmarshaled.UserID != 101 {
		t.Errorf("expected user_id 101, got %d", unmarshaled.UserID)
	}
	if unmarshaled.Operation != "test_operation" {
		t.Errorf("expected operation test_operation, got %s", unmarshaled.Operation)
	}
	if unmarshaled.DurationMs != 100 {
		t.Errorf("expected duration_ms 100, got %d", unmarshaled.DurationMs)
	}
	if unmarshaled.Error != "test error" {
		t.Errorf("expected error test error, got %s", unmarshaled.Error)
	}
}

func TestLogger_LogLevels(t *testing.T) {
	logger := NewLogger("test-service", "test")
	
	logger.Debug("debug message", nil)
	logger.Info("info message", nil)
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)
	
	// Fatal would exit, so we can't test it here
}

func TestLogger_GlobalLogger(t *testing.T) {
	InitLogger("global-service", "global")
	
	logger := GetLogger()
	if logger == nil {
		t.Error("expected global logger to be initialized")
	}
	
	if logger.serviceName != "global-service" {
		t.Errorf("expected service name global-service, got %s", logger.serviceName)
	}
}

func TestLogger_ConvenienceFunctions(t *testing.T) {
	InitLogger("conv-service", "test")
	
	LogDebug("debug", nil)
	LogInfo("info", nil)
	LogWarn("warn", nil)
	LogError("error", nil)
}

func TestLogEntry_RequiredFields(t *testing.T) {
	entry := LogEntry{
		Timestamp:   "2024-01-01T00:00:00Z",
		Level:       LevelInfo,
		Message:     "test",
		Service:     "test-service",
		Environment: "test",
	}
	
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal log entry: %v", err)
	}
	
	jsonStr := string(jsonBytes)
	
	// Check required fields are present
	requiredFields := []string{"timestamp", "level", "message", "service", "environment"}
	for _, field := range requiredFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("required field %s not found in JSON output", field)
		}
	}
}
