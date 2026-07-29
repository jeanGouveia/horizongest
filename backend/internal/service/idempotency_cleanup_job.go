package service

import (
	"context"
	"log"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

// IdempotencyCleanupJob limpa chaves de idempotência expiradas
// TTL: 48 horas (configurável via environment)
// Frequência: 1 hora
type IdempotencyCleanupJob struct {
	db              *gorm.DB
	idempotencyRepo ports.IdempotencyRepository
	ttl             time.Duration
	interval        time.Duration
	stopChan        chan struct{}
}

func NewIdempotencyCleanupJob(db *gorm.DB, idempotencyRepo ports.IdempotencyRepository, ttlHours int) *IdempotencyCleanupJob {
	ttl := time.Duration(ttlHours) * time.Hour
	if ttl == 0 {
		ttl = 48 * time.Hour // default
	}

	return &IdempotencyCleanupJob{
		db:              db,
		idempotencyRepo: idempotencyRepo,
		ttl:             ttl,
		interval:        1 * time.Hour,
		stopChan:        make(chan struct{}),
	}
}

// Start inicia o job de limpeza em background
func (j *IdempotencyCleanupJob) Start() {
	log.Printf("[IdempotencyCleanup] Iniciando job com TTL=%v, intervalo=%v", j.ttl, j.interval)

	ticker := time.NewTicker(j.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				j.cleanup(context.Background())
			case <-j.stopChan:
				ticker.Stop()
				log.Printf("[IdempotencyCleanup] Job parado")
				return
			}
		}
	}()
}

// Stop para o job de limpeza
func (j *IdempotencyCleanupJob) Stop() {
	close(j.stopChan)
}

// cleanup executa a limpeza de chaves expiradas
func (j *IdempotencyCleanupJob) cleanup(ctx context.Context) {
	cutoff := time.Now().Add(-j.ttl)

	deleted, err := j.idempotencyRepo.DeleteExpired(ctx, cutoff)
	if err != nil {
		log.Printf("[IdempotencyCleanup] Erro ao limpar chaves expiradas: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("[IdempotencyCleanup] %d chaves expiradas removidas (cutoff: %s)", deleted, cutoff.Format(time.RFC3339))
	}
}

// RunOnce executa a limpeza uma única vez (útil para testes ou execução manual)
func (j *IdempotencyCleanupJob) RunOnce(ctx context.Context) error {
	cutoff := time.Now().Add(-j.ttl)
	deleted, err := j.idempotencyRepo.DeleteExpired(ctx, cutoff)
	if err != nil {
		return err
	}
	log.Printf("[IdempotencyCleanup] RunOnce: %d chaves expiradas removidas (cutoff: %s)", deleted, cutoff.Format(time.RFC3339))
	return nil
}
