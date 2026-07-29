package domain

import "time"

// IdempotencyStatus representa o status de uma operação idempotente
type IdempotencyStatus string

const (
	IdempotencyStatusProcessing IdempotencyStatus = "processing"
	IdempotencyStatusSucceeded  IdempotencyStatus = "succeeded"
	IdempotencyStatusFailed     IdempotencyStatus = "failed"
)

// IdempotencyKey representa um registro de idempotência
type IdempotencyKey struct {
	ID        uint
	Key       string // Chave de idempotência (UUID v4)
	CompanyID uint

	// Hash do payload request para detectar reutilização com payload diferente
	RequestHash string

	// Parâmetros da requisição (method, path, query params) em JSON
	RequestParams string

	// Status da operação
	Status IdempotencyStatus

	// Resposta HTTP armazenada para replay
	ResponseStatusCode int
	ResponseHeaders    string // JSON headers
	ResponseBody       string // JSON body

	// Erro se a operação falhou
	ErrorMessage string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IdempotencyResult representa o resultado de uma verificação de idempotência
type IdempotencyResult struct {
	// Found indica se a chave existe
	Found bool

	// Replayable indica se a resposta pode ser reutilizada
	// (só é true se status = succeeded)
	Replayable bool

	// PayloadMismatch indica se a chave foi reutilizada com payload diferente
	PayloadMismatch bool

	// Processing indica se a operação ainda está em execução
	Processing bool

	// Response contém a resposta armazenada (se replayable)
	Response *IdempotencyResponse

	// ExistingRecord contém o registro existente (para debugging)
	ExistingRecord *IdempotencyKey
}

// IdempotencyResponse representa uma resposta HTTP armazenada
type IdempotencyResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}
