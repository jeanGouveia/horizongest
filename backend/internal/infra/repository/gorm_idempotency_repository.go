package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/pg"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

// ─── GORM Model ────────────────────────────────────────────────────────────

type GormIdempotencyKey struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Key       string `gorm:"not null;index:idx_idempotency_lookup,priority:2"`
	CompanyID uint   `gorm:"not null;index:idx_idempotency_lookup,priority:1"`

	RequestHash   string `gorm:"not null;size:64"`
	RequestParams string `gorm:"not null;type:text"`

	Status string `gorm:"not null;default:'processing';index:idx_idempotency_status"`

	ResponseStatusCode int
	ResponseHeaders    string `gorm:"type:text"`
	ResponseBody       string `gorm:"type:text"`

	ErrorMessage string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_idempotency_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (GormIdempotencyKey) TableName() string { return "idempotency_keys" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.IdempotencyRepository = (*GormIdempotencyRepository)(nil)

type GormIdempotencyRepository struct {
	db *gorm.DB
}

func NewGormIdempotencyRepository(db *gorm.DB) *GormIdempotencyRepository {
	return &GormIdempotencyRepository{db: db}
}

// Check verifica se uma chave de idempotência existe e retorna o resultado
func (r *GormIdempotencyRepository) Check(ctx context.Context, companyID uint, key string, requestHash string) (*domain.IdempotencyResult, error) {
	var gKey GormIdempotencyKey
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND key = ?", companyID, key).
		First(&gKey).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Chave não existe - pode prosseguir com execução
		return &domain.IdempotencyResult{
			Found:           false,
			Replayable:      false,
			PayloadMismatch: false,
			Processing:      false,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("IdempotencyRepository.Check: %w", err)
	}

	// Chave existe - analisar estado
	result := &domain.IdempotencyResult{
		Found:           true,
		ExistingRecord:  gormToDomain(&gKey),
		PayloadMismatch: gKey.RequestHash != requestHash,
		Processing:      gKey.Status == string(domain.IdempotencyStatusProcessing),
	}

	// Só é replayable se status = succeeded E payload match
	if gKey.Status == string(domain.IdempotencyStatusSucceeded) && !result.PayloadMismatch {
		result.Replayable = true
		result.Response = &domain.IdempotencyResponse{
			StatusCode: gKey.ResponseStatusCode,
			Headers:    parseHeadersJSON(gKey.ResponseHeaders),
			Body:       gKey.ResponseBody,
		}
	}

	return result, nil
}

// CreateOrGet cria um novo registro de idempotência com status "processing"
// ou retorna o existente se já foi criado por outra request (Stripe approach)
// Usa INSERT ... ON CONFLICT DO NOTHING para eliminar race condition
func (r *GormIdempotencyRepository) CreateOrGet(ctx context.Context, record *domain.IdempotencyKey) (*domain.IdempotencyKey, error) {
	gKey := domainToGorm(record)

	// INSERT ... ON CONFLICT DO NOTHING ... RETURNING
	// Se o registro já existe, retorna NULL (DO NOTHING não retorna nada)
	err := r.db.WithContext(ctx).
		Exec(`
			INSERT INTO idempotency_keys 
				(key, company_id, request_hash, request_params, status)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT (company_id, key) 
			DO NOTHING
		`, gKey.Key, gKey.CompanyID, gKey.RequestHash, gKey.RequestParams, gKey.Status).Error

	if err != nil && !pg.IsUniqueViolation(err) {
		return nil, fmt.Errorf("IdempotencyRepository.CreateOrGet: %w", err)
	}

	// Após INSERT, buscar o registro (pode ter sido criado por nós ou por outro request)
	var existingKey GormIdempotencyKey
	err = r.db.WithContext(ctx).
		Where("company_id = ? AND key = ?", gKey.CompanyID, gKey.Key).
		First(&existingKey).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Isso não deveria acontecer se INSERT sucedeu
		return nil, fmt.Errorf("IdempotencyRepository.CreateOrGet: registro não encontrado após INSERT")
	}

	if err != nil {
		return nil, fmt.Errorf("IdempotencyRepository.CreateOrGet: buscar após INSERT: %w", err)
	}

	return gormToDomain(&existingKey), nil
}

// UpdateSuccess atualiza o registro para status "succeeded" com a resposta
func (r *GormIdempotencyRepository) UpdateSuccess(ctx context.Context, id uint, response *domain.IdempotencyResponse) error {
	headersJSON, err := json.Marshal(response.Headers)
	if err != nil {
		return fmt.Errorf("IdempotencyRepository.UpdateSuccess: marshal headers: %w", err)
	}

	updates := map[string]interface{}{
		"status":               string(domain.IdempotencyStatusSucceeded),
		"response_status_code": response.StatusCode,
		"response_headers":     string(headersJSON),
		"response_body":        response.Body,
	}

	err = r.db.WithContext(ctx).
		Model(&GormIdempotencyKey{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("IdempotencyRepository.UpdateSuccess: %w", err)
	}

	return nil
}

// UpdateFailure atualiza o registro para status "failed" com o erro
func (r *GormIdempotencyRepository) UpdateFailure(ctx context.Context, id uint, errorMessage string) error {
	updates := map[string]interface{}{
		"status":        string(domain.IdempotencyStatusFailed),
		"error_message": errorMessage,
	}

	err := r.db.WithContext(ctx).
		Model(&GormIdempotencyKey{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("IdempotencyRepository.UpdateFailure: %w", err)
	}

	return nil
}

// DeleteExpired remove registros expirados (para job de limpeza)
func (r *GormIdempotencyRepository) DeleteExpired(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Delete(&GormIdempotencyKey{})

	if result.Error != nil {
		return 0, fmt.Errorf("IdempotencyRepository.DeleteExpired: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

// ComputeRequestHash calcula SHA-256 do payload normalizado
func ComputeRequestHash(body []byte) string {
	// Normalizar JSON: remover whitespace e ordenar chaves
	var normalized map[string]interface{}
	if err := json.Unmarshal(body, &normalized); err == nil {
		// Se for JSON válido, normalizar
		normalizedBytes, _ := json.Marshal(normalized)
		hash := sha256.Sum256(normalizedBytes)
		return hex.EncodeToString(hash[:])
	}
	// Se não for JSON válido, hash direto
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
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

func domainToGorm(d *domain.IdempotencyKey) GormIdempotencyKey {
	return GormIdempotencyKey{
		Key:                d.Key,
		CompanyID:          d.CompanyID,
		RequestHash:        d.RequestHash,
		RequestParams:      d.RequestParams,
		Status:             string(d.Status),
		ResponseStatusCode: d.ResponseStatusCode,
		ResponseHeaders:    d.ResponseHeaders,
		ResponseBody:       d.ResponseBody,
		ErrorMessage:       d.ErrorMessage,
	}
}

func gormToDomain(g *GormIdempotencyKey) *domain.IdempotencyKey {
	return &domain.IdempotencyKey{
		ID:                 g.ID,
		Key:                g.Key,
		CompanyID:          g.CompanyID,
		RequestHash:        g.RequestHash,
		RequestParams:      g.RequestParams,
		Status:             domain.IdempotencyStatus(g.Status),
		ResponseStatusCode: g.ResponseStatusCode,
		ResponseHeaders:    g.ResponseHeaders,
		ResponseBody:       g.ResponseBody,
		ErrorMessage:       g.ErrorMessage,
		CreatedAt:          g.CreatedAt,
		UpdatedAt:          g.UpdatedAt,
	}
}
