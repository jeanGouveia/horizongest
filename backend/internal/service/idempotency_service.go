package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrIdempotencyKeyRequired = errors.New("idempotency_key é obrigatório")
	ErrIdempotencyKeyReuse    = errors.New("idempotency_key já foi utilizada com payload diferente")
	ErrIdempotencyProcessing  = errors.New("operação com esta chave ainda está em processamento")
)

// IdempotencyService gerencia a lógica de idempotência
type IdempotencyService struct {
	idempotencyRepo ports.IdempotencyRepository
}

func NewIdempotencyService(idempotencyRepo ports.IdempotencyRepository) *IdempotencyService {
	return &IdempotencyService{
		idempotencyRepo: idempotencyRepo,
	}
}

// RequestParams captura parâmetros da requisição para contexto
type RequestParams struct {
	Method  string
	Path    string
	Headers map[string]string
	Query   map[string]string
}

// CheckAndCreate verifica idempotência e cria registro se necessário
// Usa INSERT ... ON CONFLICT DO NOTHING para eliminar race condition (Stripe approach)
// Retorna:
// - result: resultado da verificação (pode conter response para replay)
// - recordID: ID do registro criado (para atualização posterior)
// - err: erro se não for possível prosseguir
func (s *IdempotencyService) CheckAndCreate(
	ctx context.Context,
	companyID uint,
	idempotencyKey string,
	body []byte,
	params RequestParams,
) (*domain.IdempotencyResult, uint, error) {
	if idempotencyKey == "" {
		return nil, 0, ErrIdempotencyKeyRequired
	}

	// Calcular hash do payload
	requestHash := repository.ComputeRequestHash(body)

	// Serializar parâmetros da requisição
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, 0, fmt.Errorf("IdempotencyService.CheckOrInit: serializar parâmetros: %w", err)
	}

	// Criar registro ou obter existente (Stripe approach: INSERT ON CONFLICT)
	record := &domain.IdempotencyKey{
		Key:           idempotencyKey,
		CompanyID:     companyID,
		RequestHash:   requestHash,
		RequestParams: string(paramsJSON),
		Status:        domain.IdempotencyStatusProcessing,
	}

	existingRecord, err := s.idempotencyRepo.CreateOrGet(ctx, record)
	if err != nil {
		return nil, 0, fmt.Errorf("IdempotencyService.CheckOrInit: criar/obter registro: %w", err)
	}

	// Verificar payload mismatch
	if existingRecord.RequestHash != requestHash {
		log.Printf("[Idempotency] Payload mismatch: key=%s, company_id=%d", idempotencyKey, companyID)
		return nil, 0, ErrIdempotencyKeyReuse
	}

	// Se já está succeeded, replay
	if existingRecord.Status == domain.IdempotencyStatusSucceeded {
		log.Printf("[Idempotency] Replay: key=%s, company_id=%d", idempotencyKey, companyID)
		return &domain.IdempotencyResult{
			Found:      true,
			Replayable: true,
			Response: &domain.IdempotencyResponse{
				StatusCode: existingRecord.ResponseStatusCode,
				Headers:    parseHeadersJSON(existingRecord.ResponseHeaders),
				Body:       existingRecord.ResponseBody,
			},
		}, 0, nil
	}

	// Se está processing, outro request está executando
	if existingRecord.Status == domain.IdempotencyStatusProcessing {
		log.Printf("[Idempotency] Still processing: key=%s, company_id=%d", idempotencyKey, companyID)
		return nil, 0, ErrIdempotencyProcessing
	}

	// Se está failed, pode tentar novamente (reusar o registro)
	if existingRecord.Status == domain.IdempotencyStatusFailed {
		log.Printf("[Idempotency] Retrying failed: key=%s, company_id=%d", idempotencyKey, companyID)
		// Reusar o registro existente, não criar novo
		return &domain.IdempotencyResult{
			Found:      false,
			Replayable: false,
			Processing: false,
		}, existingRecord.ID, nil
	}

	// Nós criamos o registro (status=processing)
	log.Printf("[Idempotency] Created: key=%s, company_id=%d, id=%d", idempotencyKey, companyID, existingRecord.ID)

	return &domain.IdempotencyResult{
		Found:      false,
		Replayable: false,
		Processing: false,
	}, existingRecord.ID, nil
}

// RecordSuccess registra sucesso da operação
func (s *IdempotencyService) RecordSuccess(
	ctx context.Context,
	recordID uint,
	statusCode int,
	headers http.Header,
	body []byte,
) error {
	// Converter headers para map[string]string
	headersMap := make(map[string]string)
	for k, v := range headers {
		if len(v) > 0 {
			headersMap[k] = v[0]
		}
	}

	response := &domain.IdempotencyResponse{
		StatusCode: statusCode,
		Headers:    headersMap,
		Body:       string(body),
	}

	if err := s.idempotencyRepo.UpdateSuccess(ctx, recordID, response); err != nil {
		return fmt.Errorf("IdempotencyService.RecordSuccess: registrar sucesso: %w", err)
	}

	log.Printf("[Idempotency] Success recorded: record_id=%d, status=%d", recordID, statusCode)
	return nil
}

// RecordFailure registra falha da operação
func (s *IdempotencyService) RecordFailure(ctx context.Context, recordID uint, errorMessage string) error {
	if err := s.idempotencyRepo.UpdateFailure(ctx, recordID, errorMessage); err != nil {
		return fmt.Errorf("IdempotencyService.RecordFailure: registrar falha: %w", err)
	}

	log.Printf("[Idempotency] Failure recorded: record_id=%d, error=%s", recordID, errorMessage)
	return nil
}

// ReadBody lê e retorna o body da requisição (permite múltiplas leituras)
func ReadBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Restaurar body para que possa ser lido novamente
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// ExtractRequestParams extrai parâmetros da requisição
func ExtractRequestParams(r *http.Request) RequestParams {
	params := RequestParams{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}

	// Headers relevantes
	for k, v := range r.Header {
		if len(v) > 0 {
			params.Headers[k] = v[0]
		}
	}

	// Query params
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			params.Query[k] = v[0]
		}
	}

	return params
}

func parseHeadersJSON(headersJSON string) map[string]string {
	if headersJSON == "" {
		return nil
	}
	var headers map[string]string
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return nil
	}
	return headers
}
