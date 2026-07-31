package webhook

import (
	"context"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/consumers/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockWebhookProvider is a mock implementation of WebhookProvider for testing
type MockWebhookProvider struct {
	SentWebhooks []Webhook
	SendError    error
}

func (m *MockWebhookProvider) Send(ctx context.Context, webhook Webhook) error {
	if m.SendError != nil {
		return m.SendError
	}
	m.SentWebhooks = append(m.SentWebhooks, webhook)
	return nil
}

func (m *MockWebhookProvider) Close() error {
	return nil
}

var _ WebhookProvider = (*MockWebhookProvider)(nil)

func TestLogWebhookProvider(t *testing.T) {
	provider := NewLogWebhookProvider()
	ctx := context.Background()

	webhook := Webhook{
		URL: "https://example.com/webhook",
		Payload: map[string]interface{}{
			"event": "test",
		},
		Headers: map[string]string{
			"Template": "test",
			"Payload":  "{}",
			"EventID":  "1",
		},
	}

	err := provider.Send(ctx, webhook)
	assert.NoError(t, err)

	err = provider.Close()
	assert.NoError(t, err)
}

func TestInvitationWebhookTemplate(t *testing.T) {
	template := NewInvitationWebhookTemplate()

	url, payload, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com/webhooks/invitation", url)
	assert.Equal(t, "invitation.created", payload["event"])
}

func TestOrderCreatedWebhookTemplate(t *testing.T) {
	template := NewOrderCreatedWebhookTemplate()

	url, payload, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com/webhooks/order", url)
	assert.Equal(t, "order.created", payload["event"])
}

func TestCompanyCreatedWebhookTemplate(t *testing.T) {
	template := NewCompanyCreatedWebhookTemplate()

	url, payload, err := template.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com/webhooks/company", url)
	assert.Equal(t, "company.created", payload["event"])
}

func TestWebhookProcessor_Process(t *testing.T) {
	mockProvider := &MockWebhookProvider{}
	processor := NewWebhookProcessor(mockProvider)

	event := framework.Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{"order_id": 100},
	}

	err := processor.Process(context.Background(), event)
	assert.NoError(t, err)
	assert.Len(t, mockProvider.SentWebhooks, 1)
	assert.Equal(t, "https://api.example.com/webhooks/order", mockProvider.SentWebhooks[0].URL)
}

func TestWebhookProcessor_Process_UnknownType(t *testing.T) {
	mockProvider := &MockWebhookProvider{}
	processor := NewWebhookProcessor(mockProvider)

	event := framework.Event{
		ID:            1,
		EventType:     "unknown.event",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	err := processor.Process(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no template for event type")
}

func TestWebhookProcessor_Process_ProviderError(t *testing.T) {
	mockProvider := &MockWebhookProvider{
		SendError: assert.AnError,
	}
	processor := NewWebhookProcessor(mockProvider)

	event := framework.Event{
		ID:            1,
		EventType:     "order.created",
		AggregateType: "order",
		AggregateID:   100,
		TenantID:      1,
		Payload:       map[string]interface{}{},
	}

	err := processor.Process(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send webhook")
}

func TestWebhookProcessor_AllTemplates(t *testing.T) {
	mockProvider := &MockWebhookProvider{}
	processor := NewWebhookProcessor(mockProvider)

	events := []framework.Event{
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
		err := processor.Process(context.Background(), event)
		assert.NoError(t, err)
	}

	assert.Len(t, mockProvider.SentWebhooks, 3)
	assert.Equal(t, "https://api.example.com/webhooks/invitation", mockProvider.SentWebhooks[0].URL)
	assert.Equal(t, "https://api.example.com/webhooks/order", mockProvider.SentWebhooks[1].URL)
	assert.Equal(t, "https://api.example.com/webhooks/company", mockProvider.SentWebhooks[2].URL)
}

func TestWebhookConsumer_Framework(t *testing.T) {
	mockProvider := &MockWebhookProvider{}
	config := framework.DefaultConfig()
	config.EnableMetrics = false // Disable metrics for simpler testing

	consumer := NewWebhookConsumer(nil, "test_queue", mockProvider, config)
	require.NotNil(t, consumer)

	// Test that consumer has framework methods
	assert.NotNil(t, consumer.GetMetrics())
	assert.Equal(t, framework.StateClosed, consumer.GetCircuitBreakerState())

	// Test clear idempotency
	consumer.ClearIdempotency()

	// Test close
	err := consumer.Close()
	assert.NoError(t, err)
}
