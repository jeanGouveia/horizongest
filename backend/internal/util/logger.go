package util

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

// LogEntry represents a structured log entry
// FASE A.3: B16 - Enhanced with all required fields
type LogEntry struct {
	Timestamp     string                 `json:"timestamp"`
	Level         LogLevel               `json:"level"`
	Message       string                 `json:"message"`
	RequestID     string                 `json:"request_id,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	CompanyID     uint                   `json:"company_id,omitempty"`
	UserID        uint                   `json:"user_id,omitempty"`
	Service       string                 `json:"service"`
	Operation     string                 `json:"operation,omitempty"`
	DurationMs    int64                  `json:"duration_ms,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Environment   string                 `json:"environment"`
	Fields        map[string]interface{} `json:"fields,omitempty"`
}

// Logger is a structured logger that outputs JSON format
type Logger struct {
	serviceName string
	env         string
}

// NewLogger creates a new structured logger
func NewLogger(serviceName, env string) *Logger {
	return &Logger{
		serviceName: serviceName,
		env:         env,
	}
}

// log writes a structured log entry
// FASE A.3: B16 - Enhanced to include all required fields
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       level,
		Message:     message,
		Service:     l.serviceName,
		Environment: l.env,
		Fields:      fields,
	}

	// Extract common fields from the fields map if present
	if fields != nil {
		if requestID, ok := fields["request_id"]; ok {
			entry.RequestID = requestID.(string)
			delete(fields, "request_id")
		}
		if correlationID, ok := fields["correlation_id"]; ok {
			entry.CorrelationID = correlationID.(string)
			delete(fields, "correlation_id")
		}
		if companyID, ok := fields["company_id"]; ok {
			entry.CompanyID = companyID.(uint)
			delete(fields, "company_id")
		}
		if userID, ok := fields["user_id"]; ok {
			entry.UserID = userID.(uint)
			delete(fields, "user_id")
		}
		if operation, ok := fields["operation"]; ok {
			entry.Operation = operation.(string)
			delete(fields, "operation")
		}
		if durationMs, ok := fields["duration_ms"]; ok {
			entry.DurationMs = durationMs.(int64)
			delete(fields, "duration_ms")
		}
		if err, ok := fields["error"]; ok {
			entry.Error = err.(string)
			delete(fields, "error")
		}
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal log entry: %v", err)
		return
	}

	// Output to stdout
	log.Println(string(jsonBytes))

	// Exit on fatal
	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log(LevelDebug, message, fields)
}

// Info logs an info message
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(LevelInfo, message, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(LevelWarn, message, fields)
}

// Error logs an error message
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(LevelError, message, fields)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	l.log(LevelFatal, message, fields)
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// For simplicity, return the same logger
	// In a more sophisticated implementation, this would return a child logger
	return l
}

// Global logger instance
var globalLogger *Logger

// InitLogger initializes the global logger
func InitLogger(serviceName, env string) {
	globalLogger = NewLogger(serviceName, env)
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewLogger("horizongest", "development")
	}
	return globalLogger
}

// Convenience functions using the global logger
func LogDebug(message string, fields map[string]interface{}) {
	GetLogger().Debug(message, fields)
}

func LogInfo(message string, fields map[string]interface{}) {
	GetLogger().Info(message, fields)
}

func LogWarn(message string, fields map[string]interface{}) {
	GetLogger().Warn(message, fields)
}

func LogError(message string, fields map[string]interface{}) {
	GetLogger().Error(message, fields)
}

func LogFatal(message string, fields map[string]interface{}) {
	GetLogger().Fatal(message, fields)
}
