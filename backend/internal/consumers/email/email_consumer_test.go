package email

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEmailProvider is a mock implementation of EmailProvider for testing
type MockEmailProvider struct {
	SentEmails []Email
	SendError  error
}

func (m *MockEmailProvider) Send(ctx context.Context, email Email) error {
	if m.SendError != nil {
		return m.SendError
	}
	m.SentEmails = append(m.SentEmails, email)
	return nil
}

func (m *MockEmailProvider) Close() error {
	return nil
}

var _ EmailProvider = (*MockEmailProvider)(nil)

func TestIdempotencyStore(t *testing.T) {
	store := NewIdempotencyStore()

	// Initially not processed
	assert.False(t, store.IsProcessed(1))

	// Mark as processed
	store.MarkProcessed(1)
	assert.True(t, store.IsProcessed(1))

	// Different ID not processed
	assert.False(t, store.IsProcessed(2))

	// Clear
	store.Clear()
	assert.False(t, store.IsProcessed(1))
}

func TestLogEmailProvider(t *testing.T) {
	provider := NewLogEmailProvider()
	ctx := context.Background()

	email := Email{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
		Headers: map[string]string{
			"Template": "test",
			"Payload":  "{}",
			"EventID":  "1",
		},
	}

	err := provider.Send(ctx, email)
	assert.NoError(t, err)

	err = provider.Close()
	assert.NoError(t, err)
}

func TestInvitationTemplate(t *testing.T) {
	template := NewInvitationTemplate()

	subject, body, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "You're Invited!", subject)
	assert.Equal(t, "Welcome to HorizonGest. You have been invited to join our platform.", body)
}

func TestOrderCreatedTemplate(t *testing.T) {
	template := NewOrderCreatedTemplate()

	subject, body, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Order Confirmation", subject)
	assert.Equal(t, "Thank you for your order. Your order has been received and is being processed.", body)
}

func TestCompanyCreatedTemplate(t *testing.T) {
	template := NewCompanyCreatedTemplate()

	subject, body, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Welcome to HorizonGest", subject)
	assert.Equal(t, "Your company has been successfully created. Welcome aboard!", body)
}

func TestEmailConsumer_ProcessEvent(t *testing.T) {
	mockProvider := &MockEmailProvider{}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	event := Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{"order_id": 100},
	}

	err := consumer.processEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Len(t, mockProvider.SentEmails, 1)
	assert.Equal(t, "Order Confirmation", mockProvider.SentEmails[0].Subject)
}

func TestEmailConsumer_ProcessEvent_UnknownType(t *testing.T) {
	mockProvider := &MockEmailProvider{}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	event := Event{
		ID:            1,
		EventType:     "unknown.event",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	err := consumer.processEvent(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no template for event type")
}

func TestEmailConsumer_ProcessEvent_ProviderError(t *testing.T) {
	mockProvider := &MockEmailProvider{
		SendError: assert.AnError,
	}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	event := Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	err := consumer.processEvent(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestEmailConsumer_Idempotency(t *testing.T) {
	mockProvider := &MockEmailProvider{}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	event := Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	// First processing
	err := consumer.processEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Len(t, mockProvider.SentEmails, 1)

	// Mark as processed
	idempotencyStore.MarkProcessed(event.ID)

	// Second processing should be ignored (handled in processMessage, but we test the store)
	assert.True(t, idempotencyStore.IsProcessed(event.ID))
}

func TestEmailConsumer_AllTemplates(t *testing.T) {
	mockProvider := &MockEmailProvider{}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	events := []Event{
		{
			ID:            1,
			EventType:     "invitation.created",
			AggregateType: "invitation",
			AggregateID:   1,
			TenantID:      1,
			Payload:       map[string]interface{}{},
		},
		{
			ID:            2,
			EventType:     "order.created",
			AggregateType: "order",
			AggregateID:   100,
			TenantID:      1,
			Payload:       map[string]interface{}{},
		},
		{
			ID:            3,
			EventType:     "company.created",
			AggregateType: "company",
			AggregateID:   1,
			TenantID:      1,
			Payload:       map[string]interface{}{},
		},
	}

	for _, event := range events {
		err := consumer.processEvent(context.Background(), event)
		assert.NoError(t, err)
	}

	assert.Len(t, mockProvider.SentEmails, 3)
	assert.Equal(t, "You're Invited!", mockProvider.SentEmails[0].Subject)
	assert.Equal(t, "Order Confirmation", mockProvider.SentEmails[1].Subject)
	assert.Equal(t, "Welcome to HorizonGest", mockProvider.SentEmails[2].Subject)
}

func TestEmailConsumer_ProcessMessage(t *testing.T) {
	mockProvider := &MockEmailProvider{}
	idempotencyStore := NewIdempotencyStore()
	consumer := NewEmailConsumer(nil, "test_queue", mockProvider, idempotencyStore)

	event := Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	body, err := json.Marshal(event)
	require.NoError(t, err)

	msg := amqp.Delivery{
		Body: body,
	}

	// This would normally be called in the consume loop
	// For testing, we simulate the processing
	consumer.processMessage(context.Background(), msg)

	// Give it a moment to process
	time.Sleep(10 * time.Millisecond)

	// Verify event was processed
	assert.True(t, idempotencyStore.IsProcessed(event.ID))
}
