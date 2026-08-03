package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

// EventDispatcher processa eventos da Outbox e publica no message broker
type EventDispatcher struct {
	outboxRepo  ports.OutboxRepository
	publisher   ports.EventPublisher
	config      DispatcherConfig
	wg          sync.WaitGroup
	shutdownCtx context.Context
	shutdown    context.CancelFunc
	running     bool
	mu          sync.RWMutex
}

// DispatcherConfig configura o comportamento do Dispatcher
type DispatcherConfig struct {
	Interval         time.Duration // Intervalo entre ciclos de processamento
	BatchSize        int           // Número máximo de eventos por batch
	RetryCount       int           // Número máximo de tentativas por evento
	RetryBackoff     time.Duration // Backoff entre tentativas
	PublisherTimeout time.Duration // Timeout para publicação de cada evento
}

// DefaultDispatcherConfig retorna configurações padrão
func DefaultDispatcherConfig() DispatcherConfig {
	return DispatcherConfig{
		Interval:         5 * time.Second,
		BatchSize:        50,
		RetryCount:       5,
		RetryBackoff:     30 * time.Second,
		PublisherTimeout: 10 * time.Second,
	}
}

// NewEventDispatcher cria um novo EventDispatcher
func NewEventDispatcher(
	outboxRepo ports.OutboxRepository,
	publisher ports.EventPublisher,
	config DispatcherConfig,
) *EventDispatcher {
	return &EventDispatcher{
		outboxRepo: outboxRepo,
		publisher:  publisher,
		config:     config,
	}
}

// Start inicia o dispatcher em background
func (d *EventDispatcher) Start(ctx context.Context) {
	d.mu.Lock()
	if d.running {
		d.mu.Unlock()
		log.Printf("EventDispatcher already running")
		return
	}
	d.running = true
	d.mu.Unlock()

	d.shutdownCtx, d.shutdown = context.WithCancel(ctx)

	d.wg.Add(1)
	go d.run()

	log.Printf("EventDispatcher started: interval=%v, batch_size=%d, retry_count=%d",
		d.config.Interval, d.config.BatchSize, d.config.RetryCount)
}

// run é o loop principal do dispatcher
func (d *EventDispatcher) run() {
	defer d.wg.Done()

	ticker := time.NewTicker(d.config.Interval)
	defer ticker.Stop()

	log.Printf("EventDispatcher running")

	for {
		select {
		case <-d.shutdownCtx.Done():
			log.Printf("EventDispatcher shutdown requested")
			return
		case <-ticker.C:
			d.processBatch()
		}
	}
}

// processBatch processa um batch de eventos pendentes
func (d *EventDispatcher) processBatch() {
	ctx, cancel := context.WithTimeout(d.shutdownCtx, 30*time.Second)
	defer cancel()

	// Buscar tenant ID do contexto (para multi-tenancy)
	// Por enquanto, processamos todos os tenants em loop
	// TODO: Implementar processamento por tenant específico

	// Buscar eventos pendentes
	events, err := d.outboxRepo.FindPendingEvents(ctx, 0, d.config.BatchSize)
	if err != nil {
		log.Printf("EventDispatcher: failed to find pending events: %v", err)
		return
	}

	if len(events) == 0 {
		log.Printf("EventDispatcher: no pending events found")
		return
	}

	log.Printf("EventDispatcher: processing %d pending events", len(events))

	// Processar cada evento
	for _, event := range events {
		if err := d.processEvent(ctx, event); err != nil {
			log.Printf("EventDispatcher: failed to process event id=%d: %v", event.ID, err)
		}
	}
}

// processEvent processa um único evento
func (d *EventDispatcher) processEvent(ctx context.Context, event *domain.OutboxEvent) error {
	// Verificar se atingiu limite de tentativas
	if event.Attempts >= d.config.RetryCount {
		log.Printf("EventDispatcher: event id=%d reached max attempts (%d), marking as dead letter",
			event.ID, event.Attempts)
		return d.handleDeadLetter(ctx, event)
	}

	// Atualizar status para processing COM LOCK OTIMISTA
	// Isso previne que dois dispatchers processem o mesmo evento simultaneamente
	locked, err := d.outboxRepo.UpdateStatusWithOptimisticLock(
		ctx,
		event.ID,
		domain.OutboxStatusPending,
		domain.OutboxStatusProcessing,
	)
	if err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}
	if !locked {
		// Outro dispatcher já pegou este evento
		log.Printf("EventDispatcher: event id=%d already being processed by another dispatcher", event.ID)
		return nil
	}

	// Publicar evento
	publishCtx, cancel := context.WithTimeout(ctx, d.config.PublisherTimeout)
	defer cancel()

	if err := d.publisher.Publish(publishCtx, *event); err != nil {
		log.Printf("EventDispatcher: failed to publish event id=%d: %v", event.ID, err)
		return d.handlePublishError(ctx, event, err)
	}

	// Sucesso: marcar como completed
	if err := d.outboxRepo.MarkAsCompleted(ctx, event.ID); err != nil {
		log.Printf("EventDispatcher: failed to mark event id=%d as completed: %v", event.ID, err)
		return fmt.Errorf("failed to mark as completed: %w", err)
	}

	log.Printf("EventDispatcher: event id=%d published and marked as completed", event.ID)
	return nil
}

// handlePublishError trata erros de publicação
func (d *EventDispatcher) handlePublishError(ctx context.Context, event *domain.OutboxEvent, publishErr error) error {
	// Incrementar tentativas
	errorMsg := publishErr.Error()

	// Calcular próximo retry com backoff exponencial
	backoff := time.Duration(event.Attempts+1) * d.config.RetryBackoff
	nextRetry := time.Now().Add(backoff)

	if err := d.outboxRepo.IncrementAttempts(ctx, event.ID, errorMsg, nextRetry); err != nil {
		return fmt.Errorf("failed to increment attempts: %w", err)
	}

	log.Printf("EventDispatcher: event id=%d will retry at %v (attempt %d/%d)",
		event.ID, nextRetry, event.Attempts+1, d.config.RetryCount)

	return publishErr
}

// handleDeadLetter trata eventos que atingiram o limite de tentativas
func (d *EventDispatcher) handleDeadLetter(ctx context.Context, event *domain.OutboxEvent) error {
	// TODO: Mover para tabela de dead letter
	// Por enquanto, apenas logar erro
	lastError := ""
	if event.LastError != nil {
		lastError = *event.LastError
	}
	log.Printf("EventDispatcher: DEAD LETTER - event id=%d, type=%s, tenant_id=%d, attempts=%d, last_error=%s",
		event.ID, event.EventType, event.TenantID, event.Attempts, lastError)

	// Marcar como failed permanentemente
	if err := d.outboxRepo.UpdateStatus(ctx, event.ID, domain.OutboxStatusFailed); err != nil {
		return fmt.Errorf("failed to mark as failed: %w", err)
	}

	return nil
}

// Shutdown para o dispatcher gracefulmente
func (d *EventDispatcher) Shutdown() {
	d.mu.Lock()
	if !d.running {
		d.mu.Unlock()
		return
	}
	d.running = false
	d.mu.Unlock()

	log.Printf("EventDispatcher shutting down...")

	// Cancelar contexto
	if d.shutdown != nil {
		d.shutdown()
	}

	// Aguardar conclusão do loop
	d.wg.Wait()

	// Fechar publisher
	if d.publisher != nil {
		if err := d.publisher.Close(); err != nil {
			log.Printf("EventDispatcher: error closing publisher: %v", err)
		}
	}

	log.Printf("EventDispatcher shutdown complete")
}

// IsRunning retorna true se o dispatcher estiver rodando
func (d *EventDispatcher) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}
