package domain

import "time"

// ErrorResponse representa a resposta de erro padronizada
type ErrorResponse struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"requestId"`
}

// NewErrorResponse cria uma nova resposta de erro padronizada
func NewErrorResponse(code, message string, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	}
}

// WithDetails adiciona detalhes ao erro
func (e *ErrorResponse) WithDetails(details map[string]interface{}) *ErrorResponse {
	e.Details = details
	return e
}
