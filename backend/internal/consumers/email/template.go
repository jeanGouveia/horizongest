package email

// Template defines the interface for email templates
type Template interface {
	// Render renders the template with the given data
	// Returns subject and body
	Render(data interface{}) (subject string, body string, err error)
}

// InvitationTemplate handles invitation.created events
type InvitationTemplate struct{}

// NewInvitationTemplate creates a new InvitationTemplate
func NewInvitationTemplate() *InvitationTemplate {
	return &InvitationTemplate{}
}

// Render renders the invitation email
func (t *InvitationTemplate) Render(data interface{}) (string, string, error) {
	subject := "You're Invited!"
	body := "Welcome to HorizonGest. You have been invited to join our platform."
	return subject, body, nil
}

// OrderCreatedTemplate handles order.created events
type OrderCreatedTemplate struct{}

// NewOrderCreatedTemplate creates a new OrderCreatedTemplate
func NewOrderCreatedTemplate() *OrderCreatedTemplate {
	return &OrderCreatedTemplate{}
}

// Render renders the order confirmation email
func (t *OrderCreatedTemplate) Render(data interface{}) (string, string, error) {
	subject := "Order Confirmation"
	body := "Thank you for your order. Your order has been received and is being processed."
	return subject, body, nil
}

// CompanyCreatedTemplate handles company.created events
type CompanyCreatedTemplate struct{}

// NewCompanyCreatedTemplate creates a new CompanyCreatedTemplate
func NewCompanyCreatedTemplate() *CompanyCreatedTemplate {
	return &CompanyCreatedTemplate{}
}

// Render renders the company welcome email
func (t *CompanyCreatedTemplate) Render(data interface{}) (string, string, error) {
	subject := "Welcome to HorizonGest"
	body := "Your company has been successfully created. Welcome aboard!"
	return subject, body, nil
}
