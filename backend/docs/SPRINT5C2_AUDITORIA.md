# Sprint 5C.2 - Event Publisher Infrastructure - Auditoria

**Data:** 30/07/2026
**Objetivo:** Mapear a infraestrutura existente antes de implementar Event Publisher + RabbitMQ

---

## 1. OutboxRepository - Implementação Atual

### Localização
- **Interface:** `internal/ports/outbox_repository.go`
- **Implementação:** `internal/infra/repository/gorm_outbox_repository.go`
- **Domain Model:** `internal/domain/outbox_event.go`
- **Migration:** `migrations/00035_create_outbox_events.sql`

### Interface
```go
type OutboxRepository interface {
    Create(ctx context.Context, event *domain.OutboxEvent, tx *gorm.DB) error
    FindByID(ctx context.Context, id uint) (*domain.OutboxEvent, error)
    FindPendingEvents(ctx context.Context, tenantID uint, limit int) ([]*domain.OutboxEvent, error)
    UpdateStatus(ctx context.Context, id uint, status domain.OutboxStatus) error
    IncrementAttempts(ctx context.Context, id uint, error string) error
    MarkAsCompleted(ctx context.Context, id uint) error
    FindByAggregate(ctx context.Context, aggregateType string, aggregateID uint) ([]*domain.OutboxEvent, error)
    DeleteOldCompletedEvents(ctx context.Context, olderThan time.Duration) (int64, error)
}
```

### Características
- ✅ Suporta transações via parâmetro `tx *gorm.DB`
- ✅ Tenant filtering via `tenant_id`
- ✅ Auto-fill de `tenant_id` do contexto
- ✅ Isolamento multi-tenant em todos os métodos
- ✅ Índices otimizados para workload do Outbox Pattern
- ✅ Unique constraint para idempotência

---

## 2. Persistência de Eventos

### Tabela: `outbox_events`

| Coluna | Tipo | Descrição |
|--------|------|-----------|
| `id` | INTEGER PRIMARY KEY AUTOINCREMENT | ID do evento |
| `aggregate_type` | VARCHAR(100) NOT NULL | Tipo do agregado |
| `aggregate_id` | INTEGER NOT NULL | ID do agregado |
| `event_type` | VARCHAR(100) NOT NULL | Tipo do evento |
| `event_version` | VARCHAR(20) NOT NULL DEFAULT '1.0' | Versão do schema |
| `payload` | TEXT NOT NULL | JSON do evento |
| `tenant_id` | INTEGER NOT NULL | ID do tenant |
| `status` | VARCHAR(20) NOT NULL DEFAULT 'pending' | Status |
| `priority` | INTEGER NOT NULL DEFAULT 5 | Prioridade (1-10) |
| `attempts` | INTEGER NOT NULL DEFAULT 0 | Tentativas |
| `available_at` | DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP | Disponibilidade |
| `processed_at` | DATETIME | Processamento |
| `created_at` | DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP | Criação |
| `last_error`️ TEXT | Último erro |

### Índices
- `idx_outbox_tenant_status` - (tenant_id, status)
- `idx_outbox_available_at` - (available_at)
- `idx_outbox_aggregate` - (aggregate_type, aggregate_id)
- `idx_outbox_processed_at` - (processed_at)
- `idx_outbox_priority` - (priority, available_at)

### Unique Constraint
- `UNIQUE (aggregate_type, aggregate_id, event_type)` - Idempotência

---

## 3. Transaction Propagation

### Implementação Atual

**Pattern usado:** Transação explícita via `tx *gorm.DB`

**Exemplo em `stock_movement_service.go`:**
```go
err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 1. Buscar ingrediente com SELECT FOR UPDATE
    ingredient, err := s.productRepo.FindIngredientByIDForUpdate(ctx, input.IngredientID, tx)
    
    // 2. Criar movimentação DENTRO da transação
    if err := s.stockMovementRepo.Create(ctx, movement, tx); err != nil {
        return err
    }
    
    // 3. Atualizar estoque DENTRO da transação
    if err := s.productRepo.UpdateIngredient(ctx, ingredient, tx); err != nil {
        return err
    }
    
    return nil
})
```

**Mecanismo no OutboxRepository:**
```go
func (r *GormOutboxRepository) getDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
    if tx != nil {
        return tx.WithContext(ctx)
    }
    return r.db.WithContext(ctx)
}
```

### Conclusão
- ✅ Transaction propagation está implementada corretamente
- ✅ Padrão é consistente: passar `tx *gorm.DB` para métodos que precisam participar da transação
- ✅ Services recebem `db *gorm.DB` no construtor para iniciar transações
- ✅ Atomicidade garantida quando tx é fornecido

---

## 4. Local do Dispatcher

### Recomendação

**Local:** `cmd/server/main.go`

**Ponto de inicialização:** Após inicialização de todos os serviços e repositories

**Como goroutine separada:**
```go
// Após inicialização de serviços
outboxRepo := repository.NewGormOutboxRepository(db)
eventPublisher := infra.NewRabbitMQPublisher(rabbitMQConfig)
dispatcher := service.NewEventDispatcher(outboxRepo, eventPublisher, dispatcherConfig)

// Iniciar dispatcher em background
ctx, cancel := context.WithCancel(context.Background())
go dispatcher.Start(ctx)

// Graceful shutdown
defer cancel()
dispatcher.Shutdown()
```

**Alternativa:** Package dedicado `internal/infra/dispatcher/`

---

## 5. Local da Infraestrutura RabbitMQ

### Estrutura Recomendada

```
internal/
├── ports/
│   └── event_publisher.go          # Interface EventPublisher
├── infra/
│   └── messaging/
│       └── rabbitmq/
│           ├── rabbitmq_publisher.go    # Implementação RabbitMQ
│           ├── rabbitmq_config.go       # Configuração
│           ├── rabbitmq_connection.go   # Gerenciamento de conexão
│           └── rabbitmq_exchange.go     # Declaração de exchanges/filas
└── service/
    └── event_dispatcher.go         # Dispatcher (ou internal/infra/dispatcher/)
```

### Interface EventPublisher (ports)
```go
// internal/ports/event_publisher.go
package ports

import (
    "context"
    "github.com/jeanGouveia/horizongest/backend/internal/domain"
)

type EventPublisher interface {
    Publish(ctx context.Context, event domain.OutboxEvent) error
    PublishBatch(ctx context.Context, events []domain.OutboxEvent) error
    Close() error
}
```

### Implementação RabbitMQ (infra)
```go
// internal/infra/messaging/rabbitmq/rabbitmq_publisher.go
package rabbitmq

import (
    "github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type RabbitMQPublisher struct {
    // conexão, canal, configurações
}

func NewRabbitMQPublisher(config Config) *RabbitMQPublisher {
    // inicialização
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
    // publicação no RabbitMQ
}

var _ ports.EventPublisher = (*RabbitMQPublisher)(nil)
```

---

## 6. Respostas às Perguntas da Auditoria

### Qual package?
- **Interface:** `internal/ports`
- **Implementação RabbitMQ:** `internal/infra/messaging/rabbitmq`
- **Dispatcher:** `internal/service` ou `internal/infra/dispatcher`

### Qual interface?
- **EventPublisher** em `internal/ports/event_publisher.go`
- Métodos: `Publish()`, `PublishBatch()`, `Close()`

### Qual adapter?
- **RabbitMQPublisher** em `internal/infra/messaging/rabbitmq/rabbitmq_publisher.go`
- Implementa a interface `EventPublisher`
- Encapsula toda a lógica RabbitMQ (conexão, reconexão, exchanges, filas, bindings)

### Como manter a inversão de dependência?
```
Service (domain/application)
    ↓ depende de
EventPublisher (interface em ports)
    ↓ implementado por
RabbitMQPublisher (adapter em infra)
```

- Domínio/application conhece apenas a interface `EventPublisher`
- Implementação RabbitMQ está em `infra` (camada externa)
- Injeção de dependência via construtor
- Fácil substituir por outro broker (Redis Streams, Kafka)

---

## 7. Diagrama Arquitetural Proposto

```
┌─────────────────────────────────────────────────────────────┐
│                      Application Layer                       │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │ OrderService │    │ ProductService│   │  Dispatcher  │ │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘ │
│         │                   │                   │          │
│         ↓                   ↓                   ↓          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              OutboxRepository (ports)                │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                       Domain Layer                           │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              OutboxEvent (domain)                    │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  ┌──────────────────┐    ┌──────────────────────────────┐  │
│  │ GormOutboxRepo   │    │   EventPublisher (ports)     │  │
│  │   (repository)   │    └──────────────┬───────────────┘  │
│  └──────────────────┘                   │                  │
│                                         ↓                   │
│                          ┌──────────────────────────────┐  │
│                          │   RabbitMQPublisher          │  │
│                          │   (messaging/rabbitmq)       │  │
│                          └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ↓
                    ┌─────────────────┐
                    │   RabbitMQ      │
                    │   (External)    │
                    └─────────────────┘
```

### Fluxo de Dependências
1. **Service** → depende de → **OutboxRepository** (interface em ports)
2. **Service** → depende de → **EventPublisher** (interface em ports)
3. **Dispatcher** → depende de → **OutboxRepository** (interface em ports)
4. **Dispatcher** → depende de → **EventPublisher** (interface em ports)
5. **GormOutboxRepository** → implementa → **OutboxRepository** (infra)
6. **RabbitMQPublisher** → implementa → **EventPublisher** (infra)

### Regras de Dependência
- ✅ Nenhuma dependência de infra para domain
- ✅ Nenhuma dependência de infra para application
- ✅ Domain/application depende apenas de interfaces (ports)
- ✅ Implementações concretas ficam em infra
- ✅ Inversão de dependência mantida

---

## 8. Conclusão da Auditoria

### Status Atual
- ✅ OutboxRepository implementado corretamente
- ✅ Transaction propagation funcionando
- ✅ Tenant isolation implementado
- ✅ Índices otimizados
- ✅ Estrutura preparada para Dispatcher

### Próximos Passos
1. Criar interface `EventPublisher` em `internal/ports`
2. Implementar `RabbitMQPublisher` em `internal/infra/messaging/rabbitmq`
3. Implementar `Dispatcher` em `internal/service` ou `internal/infra/dispatcher`
4. Adicionar configurações RabbitMQ
5. Inicializar Dispatcher em `cmd/server/main.go`
6. Adicionar observabilidade (logs)
7. Criar testes unitários e integração

### Risco Técnico
**Classificação:** ✅ **BAIXO**

Infraestrutura está sólida. Não há impedimentos para implementação.

---

**Status:** ✅ **AUDITORIA CONCLUÍDA - PRONTO PARA FASE 2**
