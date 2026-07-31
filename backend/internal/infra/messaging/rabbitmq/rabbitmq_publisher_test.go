package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

func TestMockEventPublisher_Publish(t *testing.T) {
	mock := &MockEventPublisher{}

	event := domain.OutboxEvent{
		ID:            1,
		AggregateType: "order",
		AggregateID:   100,
		EventType:     "order.created",
		EventVersion:  "1.0",
		Payload:       `{"order_id":100}`,
		TenantID:      1,
		Status:        domain.OutboxStatusPending,
	}

	err := mock.Publish(context.Background(), event)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mock.PublishedEvents) != 1 {
		t.Fatalf("Expected 1 published event, got %d", len(mock.PublishedEvents))
	}

	if mock.PublishedEvents[0].ID != event.ID {
		t.Fatalf("Expected event ID %d, got %d", event.ID, mock.PublishedEvents[0].ID)
	}
}

func TestMockEventPublisher_PublishBatch(t *testing.T) {
	mock := &MockEventPublisher{}

	events := []domain.OutboxEvent{
		{
			ID:            1,
			AggregateType: "order",
			AggregateID:   100,
			EventType:     "order.created",
			EventVersion:  "1.0",
			Payload:       `{"order_id":100}`,
			TenantID:      1,
			Status:        domain.OutboxStatusPending,
		},
		{
			ID:            2,
			AggregateType: "product",
			AggregateID:   200,
			EventType:     "product.created",
			EventVersion:  "1.0",
			Payload:       `{"product_id":200}`,
			TenantID:      1,
			Status:        domain.OutboxStatusPending,
		},
	}

	err := mock.PublishBatch(context.Background(), events)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mock.PublishedEvents) != 2 {
		t.Fatalf("Expected 2 published events, got %d", len(mock.PublishedEvents))
	}
}

func TestMockEventPublisher_PublishError(t *testing.T) {
	mock := &MockEventPublisher{
		PublishError: &testError{message: "publish failed"},
	}

	event := domain.OutboxEvent{
		ID:            1,
		AggregateType: "order",
		AggregateID:   100,
		EventType:     "order.created",
		EventVersion:  "1.0",
		Payload:       `{"order_id":100}`,
		TenantID:      1,
		Status:        domain.OutboxStatusPending,
	}

	err := mock.Publish(context.Background(), event)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if len(mock.PublishedEvents) != 0 {
		t.Fatalf("Expected 0 published events on error, got %d", len(mock.PublishedEvents))
	}
}

func TestMockEventPublisher_Close(t *testing.T) {
	mock := &MockEventPublisher{}

	err := mock.Close()
	if err != nil {
		t.Fatalf("Expected no error on close, got %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.URL != "amqp://guest:guest@localhost:5672/" {
		t.Fatalf("Expected default URL, got %s", config.URL)
	}

	if config.Exchange != "horizongest.events" {
		t.Fatalf("Expected default exchange, got %s", config.Exchange)
	}

	if config.ExchangeType != "topic" {
		t.Fatalf("Expected default exchange type, got %s", config.ExchangeType)
	}

	if config.RetryCount != 3 {
		t.Fatalf("Expected default retry count 3, got %d", config.RetryCount)
	}

	if config.PublisherTimeout != 10*time.Second {
		t.Fatalf("Expected default publisher timeout 10s, got %v", config.PublisherTimeout)
	}

	if config.ReconnectDelay != 5*time.Second {
		t.Fatalf("Expected default reconnect delay 5s, got %v", config.ReconnectDelay)
	}

	if !config.EnablePublisherConfirm {
		t.Fatal("Expected publisher confirm enabled by default")
	}
}

// testError é um erro simples para testes
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
