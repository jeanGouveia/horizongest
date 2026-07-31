# Sprint 5C.2 - Event Publisher Infrastructure - Desenho Arquitetural

**Data:** 30/07/2026
**Objetivo:** Projetar a arquitetura completa do Event Publisher + RabbitMQ

---

## 1. Arquitetura em Camadas

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           PRESENTATION LAYER                             │
│                        (cmd/server/main.go)                              │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                     Dependency Injection                          │   │
│  │  - Repositories                                                    │   │
│  │  - Services                                                        │   │
│  │  - Publishers                                                      │   │
│  │  - Dispatcher                                                       │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                          APPLICATION LAYER                              │
│                        (internal/service/)                              │
│                                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐│
│  │ OrderService │  │ProductService│  │StockMovement │  │ EventDispatcher││
│  │              │  │              │  │   Service     │  │               ││
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘│
│         │                 │                 │                 │         │
│         │                 │                 │                 │         │
│         └─────────────────┴─────────────────┴─────────────────┘         │
│                           ↓                                               │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                    OutboxRepository (interface)                   │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                   EventPublisher (interface)                      │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                            DOMAIN LAYER                                  │
│                         (internal/domain/)                               │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                       OutboxEvent                                  │   │
│  │  - ID, AggregateType, AggregateID                                 │   │
│  │  - EventType, EventVersion, Payload                                │   │
│  │  - TenantID, Status, Priority                                      │   │
│  │  - Attempts, AvailableAt, ProcessedAt                             │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                       OutboxStatus                                 │   │
│  │  - Pending, Processing, Completed, Failed                         │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                         INFRASTRUCTURE LAYER                              │
│                        (internal/infra/)                                 │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                    GormOutboxRepository                           │   │
│  │                    (infra/repository/)                             │   │
│  │  - Create, FindByID, FindPendingEvents                            │   │
│  │  - UpdateStatus, IncrementAttempts, MarkAsCompleted               │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                    RabbitMQPublisher                               │   │
│  │                 (infra/messaging/rabbitmq/)                         │   │
│  │  - Connect, Reconnect                                              │   │
│  │  - Declare Exchanges, Queues, Bindings                             │   │
│  │  - Publish, Publisher Confirm                                     │   │
│  │  - Retry Logic                                                     │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                         EXTERNAL SYSTEMS                                 │
│                                                                           │
│  ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐  │
│  │   PostgreSQL     │    │     RabbitMQ     │    │   (Future)       │  │
│  │                  │    │                  │    │   - Email         │  │
│  │   outbox_events  │    │   Exchanges      │    │   - Webhook       │  │
│  │                  │    │   Queues         │    │   - iFood         │  │
│  └──────────────────┘    └──────────────────┘    └──────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Fluxo de Publicação de Eventos

### Fluxo Normal (Happy Path)

```
┌─────────────┐
│   Service   │
│ (e.g., Order)│
└──────┬──────┘
       │ 1. Inicia transação
       ↓
┌─────────────────┐
│   Transaction   │
│   (gorm.DB)     │
└──────┬──────────┘
       │ 2. Cria entidade (order)
       ↓
┌─────────────────┐
│ OrderRepository │
│   .Create()     │
└──────┬──────────┘
       │ 3. Cria evento na mesma transação
       ↓
┌─────────────────┐
│ OutboxRepository│
│   .Create()     │
└──────┬──────────┘
       │ 4. Commit transação
       ↓
┌─────────────────┐
│   PostgreSQL    │
│   outbox_events │
│ (status=pending)│
└─────────────────┘
       │ 5. Dispatcher busca eventos pendentes
       ↓
┌─────────────────┐
│   Dispatcher    │
│  (background)   │
└──────┬──────────┘
       │ 6. Busca batch de eventos
       ↓
┌─────────────────┐
│ OutboxRepository│
│.FindPending()   │
└──────┬──────────┘
       │ 7. Publica eventos
       ↓
┌─────────────────┐
│ EventPublisher  │
│   .Publish()    │
└──────┬──────────┘
       │ 8. Envia para RabbitMQ
       ↓
┌─────────────────┐
│   RabbitMQ      │
│   Exchange      │
└──────┬──────────┘
       │ 9. Publisher Confirm
       ↓
┌─────────────────┐
│   Dispatcher    │
└──────┬──────────┘
       │ 10. Marca como completed
       ↓
┌─────────────────┐
│ OutboxRepository│
│.MarkAsCompleted│
└─────────────────┘
```

### Fluxo de Erro com Retry

```
┌─────────────────┐
│   Dispatcher    │
└──────┬──────────┘
       │ 1. Publica evento
       ↓
┌─────────────────┐
│ EventPublisher  │
│   .Publish()    │
└──────┬──────────┘
       │ 2. Erro na publicação
       ↓
┌─────────────────┐
│   RabbitMQ      │
│  (connection    │
│   error)        │
└─────────────────┘
       │ 3. Incrementa tentativas
       ↓
┌─────────────────┐
│ OutboxRepository│
│.IncrementAttempts│
└──────┬──────────┘
       │ 4. Agenda retry (available_at)
       ↓
┌─────────────────┐
│   PostgreSQL    │
│ (status=failed, │
│  attempts++)    │
└─────────────────┘
       │ 5. Próximo ciclo do dispatcher
       ↓
┌─────────────────┐
│   Dispatcher    │
│ (retry após     │
│  backoff)       │
└─────────────────┘
```

### Fluxo de Dead Letter

```
┌─────────────────┐
│   Dispatcher    │
└──────┬──────────┘
       │ 1. Verifica attempts >= max
       ↓
┌─────────────────┐
│   PostgreSQL    │
│ (attempts >=    │
│  retry_count)   │
└──────┬──────────┘
       │ 2. Move para Dead Letter
       ↓
┌─────────────────┐
│ OutboxRepository│
│ (move para DLQ  │
│  ou marca como  │
│  dead_letter)   │
└─────────────────┘
       │ 3. Alerta/monitoramento
       ↓
┌─────────────────┐
│   Observability│
│ (log de DLQ)    │
└─────────────────┘
```

---

## 3. Estrutura de Pacotes

```
internal/
├── ports/
│   ├── outbox_repository.go          # Já existe
│   └── event_publisher.go            # NOVO: Interface EventPublisher
│
├── domain/
│   └── outbox_event.go               # Já existe
│
├── service/
│   ├── order_service.go              # Já existe
│   ├── stock_movement_service.go     # Já existe
│   └── event_dispatcher.go           # NOVO: Dispatcher
│
└── infra/
    ├── repository/
    │   └── gorm_outbox_repository.go # Já existe
    │
    └── messaging/
        └── rabbitmq/
            ├── rabbitmq_publisher.go       # NOVO: Implementação RabbitMQ
            ├── rabbitmq_config.go          # NOVO: Configuração
            ├── rabbitmq_connection.go      # NOVO: Gerenciamento de conexão
            ├── rabbitmq_exchange.go        # NOVO: Declaração de exchanges/filas
            └── rabbitmq_publisher_test.go  # NOVO: Testes
```

---

## 4. Interfaces e Implementações

### 4.1 EventPublisher (Interface)

**Local:** `internal/ports/event_publisher.go`

```go
package ports

import (
    "context"
    "github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// EventPublisher define a interface para publicação de eventos
// Implementa o padrão Port Adapter para desacoplar domínio de infraestrutura
type EventPublisher interface {
    // Publish publica um único evento
    Publish(ctx context.Context, event domain.OutboxEvent) error

    // PublishBatch publica múltiplos eventos em batch
    PublishBatch(ctx context.Context, events []domain.OutboxEvent) error

    // Close fecha a conexão com o message broker
    Close() error
}
```

### 4.2 RabbitMQPublisher (Implementação)

**Local:** `internal/infra/messaging/rabbitmq/rabbitmq_publisher.go`

```go
package rabbitmq

import (
    "context"
    "amqp"
    "github.com/jeanGouveia/horizongest/backend/internal/domain"
    "github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type RabbitMQPublisher struct {
    config     Config
    connection *amqp.Connection
    channel    *amqp.Channel
}

func NewRabbitMQPublisher(config Config) *RabbitMQPublisher {
    // inicialização
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
    // publicação com publisher confirm
}

func (p *RabbitMQPublisher) PublishBatch(ctx context.Context, events []domain.OutboxEvent) error {
    // publicação em batch
}

func (p *RabbitMQPublisher) Close() error {
    // graceful shutdown
}

var _ ports.EventPublisher = (*RabbitMQPublisher)(nil)
```

### 4.3 EventDispatcher (Service)

**Local:** `internal/service/event_dispatcher.go`

```go
package service

import (
    "context"
    "time"
    "github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type EventDispatcher struct {
    outboxRepo  ports.OutboxRepository
    publisher   ports.EventPublisher
    config      DispatcherConfig
}

type DispatcherConfig struct {
    Interval         time.Duration
    BatchSize        int
    RetryCount       int
    RetryBackoff     time.Duration
    PublisherTimeout time.Duration
}

func NewEventDispatcher(
    outboxRepo ports.OutboxRepository,
    publisher ports.EventPublisher,
    config DispatcherConfig,
) *EventDispatcher {
    // inicialização
}

func (d *EventDispatcher) Start(ctx context.Context) {
    // loop de processamento
}

func (d *EventDispatcher) Shutdown() {
    // graceful shutdown
}
```

---

## 5. Quem Conhece Quem

### Regras de Dependência

```
✅ PERMITIDO:
- Service → Domain
- Service → Ports (interfaces)
- Infra → Ports (interfaces)
- Infra → Domain
- Presentation → Service
- Presentation → Infra

❌ PROIBIDO:
- Domain → Infra
- Domain → Service
- Ports → Infra
- Ports → Domain (apenas tipos)
- Service → Infra (apenas via interfaces)
```

### Matriz de Dependências

| Camada | Domain | Ports | Service | Infra | External |
|-------|--------|-------|---------|-------|----------|
| Domain | - | ✅ | ❌ | ❌ | ❌ |
| Ports | ✅ | - | ❌ | ❌ | ❌ |
| Service | ✅ | ✅ | - | ❌ | ❌ |
| Infra | ✅ | ✅ | ❌ | - | ✅ |
| Presentation | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 6. Inversão de Dependência

### Diagrama de Inversão

```
┌─────────────────────────────────────────────────────────────┐
│                     HIGH LEVEL                               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ OrderService │    │ Dispatcher   │    │   Service    │  │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘  │
│         │                   │                   │          │
│         └───────────────────┴───────────────────┘          │
│                           ↓                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              EventPublisher (interface)              │  │
│  │              (internal/ports/)                         │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           ↑
                           │
┌─────────────────────────────────────────────────────────────┐
│                     LOW LEVEL                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           RabbitMQPublisher (implementation)          │  │
│  │           (internal/infra/messaging/rabbitmq/)       │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Princípio
- **High Level** (Service, Dispatcher) depende de abstrações (interfaces)
- **Low Level** (RabbitMQPublisher) implementa abstrações
- Dependência aponta para interfaces, não implementações

---

## 7. Substituibilidade

### Cenário 1: Trocar RabbitMQ por Redis Streams

```go
// Criar nova implementação
type RedisStreamsPublisher struct {
    client *redis.Client
}

func (p *RedisStreamsPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
    // publicação no Redis Streams
}

// Implementa a mesma interface
var _ ports.EventPublisher = (*RedisStreamsPublisher)(nil)

// Injeção de dependência (main.go)
publisher := infra.NewRedisStreamsPublisher(redisConfig)
dispatcher := service.NewEventDispatcher(outboxRepo, publisher, config)
```

### Cenário 2: Trocar RabbitMQ por Kafka

```go
// Criar nova implementação
type KafkaPublisher struct {
    producer *kafka.Producer
}

func (p *KafkaPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
    // publicação no Kafka
}

// Implementa a mesma interface
var _ ports.EventPublisher = (*KafkaPublisher)(nil)

// Injeção de dependência (main.go)
publisher := infra.NewKafkaPublisher(kafkaConfig)
dispatcher := service.NewEventDispatcher(outboxRepo, publisher, config)
```

### Cenário 3: Mock para Testes

```go
// Mock para testes unitários
type MockEventPublisher struct {
    PublishedEvents []domain.OutboxEvent
    PublishError    error
}

func (m *MockEventPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
    if m.PublishError != nil {
        return m.PublishError
    }
    m.PublishedEvents = append(m.PublishedEvents, event)
    return nil
}

// Implementa a mesma interface
var _ ports.EventPublisher = (*MockEventPublisher)(nil)

// Uso em testes
mockPublisher := &MockEventPublisher{}
dispatcher := service.NewEventDispatcher(outboxRepo, mockPublisher, config)
```

---

## 8. Configuração

### Estrutura de Configuração

```go
// internal/infra/messaging/rabbitmq/rabbitmq_config.go
package rabbitmq

type Config struct {
    URL             string
    Exchange        string
    ExchangeType    string
    QueuePrefix     string
    RetryCount      int
    PublisherTimeout time.Duration
    ReconnectDelay  time.Duration
}

// internal/service/event_dispatcher.go
type DispatcherConfig struct {
    Interval         time.Duration
    BatchSize        int
    RetryCount       int
    RetryBackoff     time.Duration
    PublisherTimeout time.Duration
}
```

### Environment Variables

```bash
# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=horizongest.events
RABBITMQ_EXCHANGE_TYPE=topic
RABBITMQ_QUEUE_PREFIX=horizongest
RABBITMQ_RETRY_COUNT=3
RABBITMQ_PUBLISHER_TIMEOUT=10s
RABBITMQ_RECONNECT_DELAY=5s

# Dispatcher
DISPATCHER_INTERVAL=5s
DISPATCHER_BATCH_SIZE=50
DISPATCHER_RETRY_COUNT=5
DISPATCHER_RETRY_BACKOFF=30s
```

---

## 9. Observabilidade

### Eventos de Log

```go
// Dispatcher iniciado
log.Info("event_dispatcher_started", 
    "interval", config.Interval,
    "batch_size", config.BatchSize,
)

// Evento encontrado
log.Debug("event_found",
    "event_id", event.ID,
    "event_type", event.EventType,
    "tenant_id", event.TenantID,
)

// Evento publicado
log.Info("event_published",
    "event_id", event.ID,
    "event_type", event.EventType,
    "tenant_id", event.TenantID,
    "duration_ms", duration,
)

// Retry
log.Warn("event_retry",
    "event_id", event.ID,
    "attempt", event.Attempts,
    "error", err,
    "next_retry_at", event.AvailableAt,
)

// Dead Letter
log.Error("event_dead_letter",
    "event_id", event.ID,
    "attempts", event.Attempts,
    "last_error", event.LastError,
)

// Reconnect
log.Info("rabbitmq_reconnected",
    "url", config.URL,
    "downtime_ms", downtime,
)

// Publisher Confirm
log.Debug("publisher_confirm",
    "event_id", event.ID,
    "ack", ack,
)
```

---

## 10. Conclusão

### Arquitetura Proposta

- ✅ **Clean Architecture** mantida
- ✅ **DDD** respeitado
- ✅ **Ports & Adapters** implementado
- ✅ **Inversão de Dependência** garantida
- ✅ **Substituibilidade** total (RabbitMQ ↔ Redis ↔ Kafka)
- ✅ **Testabilidade** alta (mocks fáceis)
- ✅ **Observabilidade** planejada
- ✅ **Graceful Shutdown** incluído

### Próximos Passos

1. ✅ FASE 1: Auditoria - CONCLUÍDA
2. ✅ FASE 2: Desenho Arquitetural - CONCLUÍDA
3. ⏭️ FASE 3: Event Publisher Interface
4. ⏭️ FASE 4: RabbitMQ Adapter
5. ⏭️ FASE 5: Dispatcher
6. ⏭️ FASE 6: Configuração
7. ⏭️ FASE 7: Observabilidade
8. ⏭️ FASE 8: Testes
9. ⏭️ FASE 9: Validação Arquitetural
10. ⏭️ FASE 10: Relatório Final

---

**Status:** ✅ **DESENHO ARQUITETURAL CONCLUÍDO - PRONTO PARA FASE 3**
