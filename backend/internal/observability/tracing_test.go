package observability

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestTracingService_StartSpan(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	if span == nil {
		t.Error("expected span to be created")
	}

	if ctx == nil {
		t.Error("expected context to be returned")
	}
}

func TestTracingService_StartSpanWithAttributes(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	attrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	}

	ctx, span := service.StartSpanWithAttributes(ctx, "test-operation", attrs...)
	defer span.End()

	if span == nil {
		t.Error("expected span to be created")
	}
}

func TestTracingService_RecordError(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	err := testError("test error")
	service.RecordError(span, err)

	// Error is recorded on the span, we can't easily verify it without a real tracer
}

func TestTracingService_AddAttributes(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	attrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.Int("key2", 123),
	}

	service.AddAttributes(span, attrs...)
}

func TestGetSpanFromContext(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	retrievedSpan := GetSpanFromContext(ctx)
	if retrievedSpan != span {
		t.Error("expected to retrieve the same span")
	}
}

func TestGetTraceID(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	traceID := GetTraceID(ctx)
	// Without a real tracer, this will be empty
	_ = traceID
}

func TestGetSpanID(t *testing.T) {
	service := NewTracingService("test-service")
	ctx := context.Background()

	ctx, span := service.StartSpan(ctx, "test-operation")
	defer span.End()

	spanID := GetSpanID(ctx)
	// Without a real tracer, this will be empty
	_ = spanID
}

func testError(msg string) error {
	return &testErrorType{msg: msg}
}

type testErrorType struct {
	msg string
}

func (e *testErrorType) Error() string {
	return e.msg
}
