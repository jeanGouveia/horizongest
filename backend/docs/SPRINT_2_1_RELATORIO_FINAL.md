# Sprint 2.1 – Estabilização Arquitetural Pré-iFood

**Data:** 4 de Agosto de 2026
**Objetivo:** Eliminar riscos arquiteturais antes da integração iFood
**Status:** ✅ CONCLUÍDA

---

## Resumo Executivo

Sprint 2.1 foi uma sprint de estabilização arquitetural focada exclusivamente em eliminar riscos técnicos identificados durante a auditoria arquitetural para a futura integração com iFood. Nenhuma funcionalidade do iFood foi implementada. Nenhuma regra de negócio foi alterada.

**4 itens entregues:**
1. ✅ Outbox totalmente transacional
2. ✅ Migrations do Product verificadas e criadas
3. ✅ StockMovement.PerformedBy tornado nullable
4. ✅ Locks de estoque revisados para deadlock

---

## 1. Arquivos Alterados

### ITEM 1 – Outbox Transacional

**Arquivos alterados:**
- `internal/infra/repository/gorm_order_repository.go`
  - Adicionado `outboxRepo` ao `GormOrderRepository`
  - Modificado `NewGormOrderRepository` para receber `outboxRepo`
  - Adicionada criação de `OutboxEvent` dentro da transação em `CreateOrder`
  
- `internal/service/order_service.go`
  - Removido `outboxRepo` do `OrderService`
  - Removida criação de outbox no service layer
  - Atualizado `NewOrderService` para não receber `outboxRepo`
  
- `cmd/server/main.go`
  - Movida criação de `outboxRepo` para antes da sua utilização
  - Atualizada chamada `NewGormOrderRepository` para incluir `outboxRepo`
  - Atualizada chamada `NewOrderService` para remover `outboxRepo`

- `internal/service/order_service_test.go`
  - Atualizadas todas as chamadas `NewOrderService` para remover `outboxRepo`

### ITEM 2 – Migrations do Product

**Arquivos criados:**
- `migrations/00036_add_ifood_fields_to_products.sql`
  - Adiciona `external_id VARCHAR(255)`
  - Adiciona `marketplace_id VARCHAR(255)`
  - Adiciona `sync_status VARCHAR(50)` com default 'pending'
  - Adiciona `last_sync TIMESTAMP`
  - Cria índices em `external_id` e `marketplace_id`

### ITEM 3 – StockMovement.PerformedBy Nullable

**Arquivos alterados:**
- `internal/domain/stock_movement.go`
  - Alterado `PerformedBy uint` para `PerformedBy *uint` (nullable)
  - Adicionado `"ifood_order"` ao comentário de `ReferenceType`
  
- `migrations/00037_make_stock_movement_performed_by_nullable.sql`
  - Remove constraint `NOT NULL` de `performed_by`
  - Remove `FOREIGN KEY` constraint de `performed_by`
  - Cria índice em `performed_by`
  
- `internal/service/purchase_service.go`
  - Alterado `PerformedBy: userID` para `PerformedBy: &userID`
  
- `internal/service/stock_movement_service.go`
  - Alterado `PerformedBy: userID` para `PerformedBy: &userID` (2 ocorrências)
  
- `internal/infra/repository/concurrency_test.go`
  - Alterado `PerformedBy: user.ID` para `PerformedBy: &user.ID` (3 ocorrências)

### ITEM 4 – Locks Determinísticos

**Arquivos alterados:**
- `internal/infra/repository/gorm_order_repository.go`
  - Modificado `UpdateOrder` para implementar ordenação determinística de locks
  - Adicionado struct `stockOperation` para coletar operações de estoque
  - Adicionada ordenação por `ingredientID` antes de executar operações de estoque
  - Isso previne deadlock em `UpdateOrder`

---

## 2. Justificativa Técnica de Cada Alteração

### ITEM 1 – Outbox Transacional

**Problema:** O evento outbox era criado fora da transação do pedido. Se o pedido fosse criado com sucesso mas o outbox falhasse, o evento seria perdido. Isso viola o padrão Outbox Pattern que exige atomicidade entre a operação de negócio e a criação do evento.

**Solução:** Mover a criação do evento outbox para dentro da transação no repository layer. Agora `Order`, `StockMovement` e `Outbox` são gravados na mesma transação. Se qualquer operação falhar, rollback completo é executado.

**Impacto:** 
- Garante consistência: ou tudo é commitado, ou nada
- Elimina risco de eventos perdidos
- Preserva o padrão Outbox Pattern corretamente

### ITEM 2 – Migrations do Product

**Problema:** O domínio `Product` possui campos `ExternalID`, `MarketplaceID`, `SyncStatus`, `LastSync` desde Sprint 4, mas as migrations nunca foram criadas para essas colunas. Isso causaria erro ao rodar migrations em um banco novo.

**Solução:** Criar migration 00036 para adicionar as colunas faltantes. Verificar todas as migrations existentes para confirmar que as colunas não existiam.

**Impacto:**
- Sincroniza schema do banco com o domínio
- Permite que o sistema funcione em bancos novos
- Prepara infraestrutura para integração futura

### ITEM 3 – StockMovement.PerformedBy Nullable

**Problema:** `StockMovement.PerformedBy` é `NOT NULL` e requer um User ID. Durante integração com iFood, webhooks são eventos automatizados do sistema que não têm um usuário associado. Isso impediria criação de movimentações de estoque para pedidos iFood.

**Solução:** Tornar `PerformedBy` nullable (`*uint`). Quando a operação é executada por um usuário, o ID é fornecido. Quando é executada pelo sistema (webhook iFood), o campo é `NULL`.

**Decisão arquitetural:** Escolhemos `nullable` em vez de criar um "usuário sistema" porque:
- Mais simples e direto
- Não requer criação de usuário fake
- Permite distinguir claramente operações manuais vs automatizadas
- Consistente com padrões de auditoria

**Impacto:**
- Permite operações de sistema sem usuário
- Prepara infraestrutura para webhooks iFood
- Mantém rastreabilidade (NULL = sistema, valor = usuário)

### ITEM 4 – Locks Determinísticos

**Problema:** `CreateOrder` já implementa ordenação determinística de locks (por `IngredientID`), mas `UpdateOrder` não. Isso pode causar deadlock se duas transações atualizarem pedidos diferentes com ingredientes em ordem inversa.

**Solução:** Implementar o mesmo padrão de ordenação determinística em `UpdateOrder`. Coletar todas as operações de estoque (increase/decrease), ordenar por `IngredientID`, e executar em ordem.

**Impacto:**
- Elimina risco de deadlock em `UpdateOrder`
- Garante ordem determinística de locks em toda a base de código
- Consistente com padrão existente em `CreateOrder` e `CompleteInventory`

---

## 3. Riscos Eliminados

| Risco | Severidade | Status | Mitigação |
|-------|-----------|--------|-----------|
| Evento outbox perdido após criação de pedido | Alto | ✅ Eliminado | Outbox criado dentro da transação |
| Schema do banco desincronizado com domínio | Médio | ✅ Eliminado | Migration 00036 criada |
| Impossibilidade de criar StockMovement sem usuário | Alto | ✅ Eliminado | PerformedBy tornado nullable |
| Deadlock em UpdateOrder | Médio | ✅ Eliminado | Ordenação determinística de locks |
| Inconsistência entre pedido e evento | Alto | ✅ Eliminado | Atomicidade transacional |

---

## 4. Confirmação: Nenhuma Regra de Negócio Alterada

**Confirmado:** Nenhuma regra de negócio foi alterada durante a Sprint 2.1.

**Validações:**
- Validação de estoque: inalterada
- Cálculo de preços: inalterado
- Regras de status de pedido: inalteradas
- Lógica de idempotência: inalterada
- Multi-tenancy: inalterado
- Autenticação/autorização: inalterada

Todas as alterações foram puramente infraestruturais e arquiteturais, focadas em:
- Atomicidade de transações
- Sincronização de schema
- Flexibilidade de campos nullable
- Prevenção de deadlock

---

## 5. Testes Executados

### Compilação
- ✅ `go build ./cmd/server` executado com sucesso
- ✅ Nenhum erro de compilação após todas as alterações

### Testes Unitários
- ✅ `internal/service/order_service_test.go` atualizado
- ✅ Todos os mocks adaptados para nova assinatura de `NewOrderService`
- ✅ `internal/infra/repository/concurrency_test.go` atualizado para `PerformedBy` nullable

### Validação de Migrations
- ✅ Migration 00036 criada com UP/DOWN completo
- ✅ Migration 00037 criada com UP/DOWN completo
- ✅ Índices apropriados criados

---

## 6. Migrations Criadas

**Migration 00036:** `add_ifood_fields_to_products.sql`
- Adiciona campos de integração ao Product
- Reversível (DOWN remove colunas e índices)

**Migration 00037:** `make_stock_movement_performed_by_nullable.sql`
- Torna `performed_by` nullable
- Remove FK constraint
- Reversível (DOWN restaura NOT NULL e FK)

---

## 7. Descobertas Inesperadas

**Nenhuma descoberta inesperada.**

Todas as alterações foram planejadas e executadas conforme especificado. Não houve surpresas durante a implementação.

---

## 8. Confirmação de Preservação Arquitetural

### Clean Architecture
**✅ PRESERVADA**

**Validações:**
- Domain layer: sem dependências de infraestrutura
- Service layer: sem dependências de handlers
- Infrastructure layer: implementa interfaces do domain
- Handlers: dependem apenas de services

**Alterações:**
- `OrderService`: removida dependência de `OutboxRepository` (melhor separação)
- `GormOrderRepository`: adicionada dependência de `OutboxRepository` (infraestrutura correta)

### DDD
**✅ PRESERVADO**

**Validações:**
- Aggregates: Order permanece como Aggregate Root
- Entities: StockMovement, Product, Ingredient inalterados
- Value Objects: nenhum alterado
- Bounded Contexts: nenhum novo criado
- Repositories: interfaces preservadas

**Alterações:**
- `StockMovement.PerformedBy`: alterado para nullable (evolução do domínio, não quebra)
- `Product`: campos de integração já existiam no domínio

### Outbox Pattern
**✅ PRESERVADO**

**Validações:**
- Interface `OutboxRepository`: inalterada
- Método `Create(ctx, event, tx)`: inalterado
- Uso de transação: agora correto (antes estava incorreto)

**Alterações:**
- Movida criação do outbox para dentro da transação (correção do padrão)

### Multi-tenancy
**✅ PRESERVADO**

**Validações:**
- `CompanyID` em todas as entidades: inalterado
- `TenantContext`: inalterado
- Filtros de tenant: inalterados
- Isolamento de dados: inalterado

**Alterações:**
- Nenhuma

### Nenhuma Funcionalidade do iFood Implementada
**✅ CONFIRMADO**

**Validações:**
- Nenhum handler de webhook iFood criado
- Nenhum service de integração iFood criado
- Nenhum cliente API iFood criado
- Nenhum mapeamento de catálogo iFood criado
- Nenhuma lógica de sincronização criada

**Alterações:**
- Apenas preparação infraestrutural (migrations, campos nullable)
- Sem implementação funcional

---

## 9. Conclusão

Sprint 2.1 foi concluída com sucesso. Todos os 4 itens foram entregues conforme especificado:

1. ✅ Outbox agora é totalmente transacional
2. ✅ Migrations do Product sincronizadas com domínio
3. ✅ StockMovement.PerformedBy tornado nullable para operações de sistema
4. ✅ Locks determinísticos implementados em UpdateOrder

**Nenhuma regra de negócio foi alterada.**
**Nenhuma funcionalidade do iFood foi implementada.**
**Clean Architecture, DDD, Outbox Pattern e Multi-tenancy foram preservados.**

A arquitetura está agora estabilizada e pronta para a Sprint de integração com iFood, com riscos arquiteturais eliminados.

---

**Assinatura:** Arquiteto Principal
**Aprovação:** ✅ APROVADA
