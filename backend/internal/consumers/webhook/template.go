package webhook

// WebhookTemplate defines the interface for webhook templates
type WebhookTemplate interface {
	// Render renders the template with the given data
	// Returns URL and payload
	Render(data interface{}) (url string, payload map[string]interface{}, err error)
}

// InvitationWebhookTemplate handles invitation.created events
type InvitationWebhookTemplate struct{}

// NewInvitationWebhookTemplate creates a new InvitationWebhookTemplate
func NewInvitationWebhookTemplate() *InvitationWebhookTemplate {
	return &InvitationWebhookTemplate{}
}

// Render renders the invitation webhook
func (t *InvitationWebhookTemplate) Render(data interface{}) (string, map[string]interface{}, error) {
	url := "https://api.example.com/webhooks/invitation"
	payload := map[string]interface{}{
		"event": "invitation.created",
		"data":  data,
	}
	return url, payload, nil
}

// OrderCreatedWebhookTemplate handles order.created events
type OrderCreatedWebhookTemplate struct{}

// NewOrderCreatedWebhookTemplate creates a new OrderCreatedWebhookTemplate
func NewOrderCreatedWebhookTemplate() *OrderCreatedWebhookTemplate {
	return &OrderCreatedWebhookTemplate{}
}

// Render renders the order confirmation webhook
func (t *OrderCreatedWebhookTemplate) Render(data interface{}) (string, map[string]interface{}, error) {
	url := "https://api.example.com/webhooks/order"
	payload := map[string]interface{}{
		"event": "order.created",
		"data":  data,
	}
	return url, payload, nil
}

// CompanyCreatedWebhookTemplate handles company.created events
type CompanyCreatedWebhookTemplate struct{}

// NewCompanyCreatedWebhookTemplate creates a new CompanyCreatedWebhookTemplate
func NewCompanyCreatedWebhookTemplate() *CompanyCreatedWebhookTemplate {
	return &CompanyCreatedWebhookTemplate{}
}

// Render renders the company welcome webhook
func (t *CompanyCreatedWebhookTemplate) Render(data interface{}) (string, map[string]interface{}, error) {
	url := "https://api.example.com/webhooks/company"
	payload := map[string]interface{}{
		"event": "company.created",
		"data":  data,
	}
	return url, payload, nil
}
