package ports

import (
	"context"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type IdempotencyRepository interface {
	// Check verifica se uma chave de idempotência existe e retorna o resultado
	Check(ctx context.Context, companyID uint, key string, requestHash string) (*domain.IdempotencyResult, error)

	// CreateOrGet cria um novo registro de idempotência com status "processing"
	// ou retorna o existente se já foi criado por outra request (elimina race condition)
	CreateOrGet(ctx context.Context, record *domain.IdempotencyKey) (*domain.IdempotencyKey, error)

	// UpdateSuccess atualiza o registro para status "succeeded" com a resposta
	UpdateSuccess(ctx context.Context, id uint, response *domain.IdempotencyResponse) error

	// UpdateFailure atualiza o registro para status "failed" com o erro
	UpdateFailure(ctx context.Context, id uint, errorMessage string) error

	// DeleteExpired remove registros expirados (para job de limpeza)
	DeleteExpired(ctx context.Context, cutoff time.Time) (int64, error)
}
