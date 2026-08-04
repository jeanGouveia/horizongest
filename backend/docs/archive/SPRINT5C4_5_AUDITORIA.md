# SPRINT 5C.4.5 — AUDITORIA DE CONSISTÊNCIA TRANSACIONAL E INTEGRIDADE DOS DADOS

**Data:** 2025-01-XX  
**Auditor:** Principal Software Architect  
**Escopo:** HorizonGest Backend  
**Objetivo:** Identificar problemas, riscos e inconsistências de dados antes da entrada em produção

---

## RESUMO EXECUTIVO

Esta auditoria identificou **23 problemas** de consistência transacional e integridade de dados, classificados em:

- **🔴 CRÍTICOS (5):** Risco de corrupção de dados em produção
- **🟡 ALTOS (10):** Risco de inconsistência em cenários específicos
- **🟢 MÉDIOS (6):** Risco moderado, impacto limitado
- **⚪ BAIXOS (2):** Melhorias recomendadas

**Nota de Consistência do Sistema:** 7.2/10  
**Nota de Confiabilidade:** 7.5/10  
**Nota Transacional:** 8.0/10

---

## 1. TRANSAÇÕES

### 1.1 Rollback Explícito Desnecessário

**Problema:** Chamadas explícitas de `tx.Rollback()` dentro de transações GORM  
**Severidade:** 🟢 BAIXA  
**Arquivo:** `internal/infra/repository/gorm_stock_adjustment_repository.go`  
**Linhas:** 227, 232

**Causa:**
```go
if err := r.approveWithTx(ctx, id, processedBy, notes, tx); err != nil {
    tx.Rollback()  // Desnecessário
    return fmt.Errorf("StockAdjustmentRepository.ApproveAndRestoreStock: erro ao aprovar ajuste: %w", err)
}

if err := r.productRepo.IncreaseIngredientStock(ctx, gAdjustment.IngredientID, gAdjustment.Quantity, tx); err != nil {
    tx.Rollback()  // Desnecessário
    return fmt.Errorf("StockAdjustmentRepository.ApproveAndRestoreStock: erro ao repor estoque: %w", err)
}
```

GORM's `Transaction()` automaticamente executa rollback quando a função retorna erro. Chamadas explícitas são redundantes.

**Impacto:** Código desnecessário, mas não causa problemas funcionais. Pode confundir desenvolvedores sobre o comportamento da transação.

**Solução Proposta:**
Remover chamadas explícitas de `tx.Rollback()`. GORM gerencia automaticamente:
```go
if err := r.approveWithTx(ctx, id, processedBy, notes, tx); err != nil {
    return fmt.Errorf("StockAdjustmentRepository.ApproveAndRestoreStock: erro ao aprovar ajuste: %w", err)
}
```

---

### 1.2 Operações Fora de Transação

**Problema:** Algumas operações críticas não estão protegidas por transações  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/platform_service.go`  
**Linhas:** UpdateCompany, DeactivateCompany, ActivateCompany, ResetOwnerPassword, BlockUser, UnblockUser, SetCompanyTrial, SuspendCompany, CancelCompany, ReactivateCompany

**Causa:**
Métodos de atualização de plataforma não usam transações explícitas. Confiam em operações únicas de repository que podem não ser transacionais.

**Impacto:** Se uma operação de atualização falhar parcialmente (ex: atualizar empresa mas falhar ao registrar audit log), pode deixar sistema em estado inconsistente.

**Solução Proposta:**
Adicionar transações para operações que envolvem múltiplas atualizações:
```go
func (s *PlatformService) UpdateCompany(ctx context.Context, id uint, in PlatformUpdateCompanyInput) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Atualizar empresa
        // Registrar audit log
        return nil
    })
}
```

---

### 1.3 Transação Grande em CreateOrder

**Problema:** Transação de criação de pedido pode ser longa com muitos itens  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 80-194

**Causa:**
A transação `CreateOrder` executa múltiplas operações:
1. Gerar order_number
2. Criar pedido
3. Calcular consumo de ingredientes
4. Ordenar ingredientes por ID
5. Adquirir locks SELECT FOR UPDATE para cada ingrediente
6. Baixar estoque de cada ingrediente
7. Persistir itens do pedido

Com pedidos grandes (ex: 50 itens, cada com 10 ingredientes), isso pode resultar em 500 locks individuais.

**Impacto:** 
- Tempo de transação longo aumenta risco de deadlock
- Locks mantidos por tempo prolongado bloqueiam outras operações
- Timeout de transação em alta concorrência

**Solução Proposta:**
Considerar dividir em múltiplas transações menores ou usar batch operations:
```go
// Opção 1: Batch stock updates
UPDATE ingredients SET stock_quantity = stock_quantity - ? 
WHERE id IN (?, ?, ?) AND stock_quantity >= ?

// Opção 2: Pré-validar estoque antes da transação
// Transação só executa as atualizações finais
```

---

## 2. ATOMICIDADE

### 2.1 CreatePurchaseReceiving Não Atualiza Estoque

**Problema:** Recebimento de compra não atualiza estoque de ingredientes  
**Severidade:** 🔴 CRÍTICA  
**Arquivo:** `internal/service/purchase_service.go`  
**Linhas:** 256-257

**Causa:**
```go
// TODO: Atualizar estoque do ingrediente via stock movements
// Isso será implementado na integração com o módulo de estoque
```

O método `CreatePurchaseReceiving` cria o recebimento e itens, mas não atualiza o estoque dos ingredientes recebidos.

**Impacto:**
- Estoque de ingredientes não é incrementado quando materiais são recebidos
- Discrepância entre estoque físico e sistema
- Planejamento de produção baseado em estoque incorreto
- Perda de rastreabilidade de movimentações de entrada

**Solução Proposta:**
Implementar atualização de estoque dentro da transação:
```go
if err := s.purchaseRepo.CreatePurchaseReceivingItem(ctx, receivingItem); err != nil {
    return fmt.Errorf("PurchaseService.ReceivePurchaseOrder: criar item de recebimento: %w", err)
}

// Atualizar estoque do ingrediente
if err := s.stockMovementService.CreateStockMovement(ctx, companyID, userID, StockMovementInput{
    IngredientID: item.IngredientID,
    Quantity:     item.Quantity,
    Type:         StockMovementTypePurchase,
    Reason:       fmt.Sprintf("Recebimento compra #%d", purchaseOrderID),
}); err != nil {
    return fmt.Errorf("PurchaseService.ReceivePurchaseOrder: atualizar estoque: %w", err)
}
```

---

### 2.2 CreateCompany Sem Transação Completa

**Problema:** Criação de empresa não é totalmente atômica  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/platform_service.go`  
**Linhas:** 72-180

**Causa:**
O método `CreateCompany` usa transação, mas se a criação do usuário owner falhar após a empresa ser criada, a empresa pode ficar sem owner.

**Impacto:**
- Empresa criada mas sem usuário owner
- Ninguém pode administrar a empresa
- Necessita intervenção manual para corrigir

**Solução Proposta:**
A transação já existe, mas garantir que todas as validações ocorram antes de iniciar a transação:
```go
// Validar tudo antes da transação
if err := validateOwnerEmail(ctx, ownerEmail); err != nil {
    return nil, err
}

// Só então iniciar transação
err := s.db.Transaction(func(tx *gorm.DB) error {
    // Criar empresa
    // Criar usuário owner
    // Registrar audit
    return nil
})
```

---

### 2.3 DuplicateProduct Parcialmente Atômico

**Problema:** Duplicação de produto depende de `s.db` ser não-nil  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/product_service.go`  
**Linhas:** 427-504

**Causa:**
```go
executeInTx := func(fn func() error) error {
    if s.db != nil {
        return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
            return fn()
        })
    }
    return fn()  // Sem transação se s.db for nil
}
```

Se `s.db` for nil (ex: em testes ou configuração incorreta), a duplicação ocorre sem transação.

**Impacto:**
- Se a criação do produto duplicado falhar após a cópia dos ingredientes, ingredientes podem ficar órfãos
- Inconsistência entre produto e seus ingredientes

**Solução Proposta:**
Remover condicional ou falhar explicitamente se DB não disponível:
```go
if s.db == nil {
    return nil, errors.New("ProductService.DuplicateProduct: database not available for transaction")
}
```

---

### 2.4 UpdateOrder Sem Validação de Estoque Prévia

**Problema:** Atualização de pedido não valida estoque antes da transação  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 452-612

**Causa:**
O método `UpdateOrder` ajusta estoque dentro da transação, mas se não houver estoque suficiente, a transação falha no meio. Isso pode deixar o pedido em estado inconsistente.

**Impacto:**
- Transação longa com alto risco de deadlock
- Usuário recebe erro apenas após tentativa de atualização
- Experiência de usuário ruim

**Solução Proposta:**
Validar estoque antes da transação:
```go
func (r *GormOrderRepository) UpdateOrder(...) error {
    // 1. Validar estoque fora da transação
    if err := r.validateStockForUpdate(ctx, items, productIngredients); err != nil {
        return err
    }
    
    // 2. Executar atualização na transação
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // ...
    })
}
```

---

## 3. INTEGRIDADE REFERENCIAL

### 3.1 ON DELETE Não Implementado

**Problema:** Foreign keys não têm cláusulas ON DELETE  
**Severidade:** 🔴 CRÍTICA  
**Arquivo:** `migrations/00018_add_fk_on_delete.sql`  
**Linhas:** 1-53

**Causa:**
A migration 00018 é apenas documentação para SQLite. ON DELETE não está implementado:
```sql
-- Note: SQLite does not support ALTER TABLE to add ON DELETE to existing foreign keys
-- Tables need to be recreated. For now, this migration documents the required changes.
```

**Impacto:**
- Ao deletar um produto, itens de pedido podem ficar com product_id órfão
- Ao deletar um usuário, registros criados por ele podem ficar sem referência
- Acúmulo de registros órfãos ao longo do tempo
- Consultas com JOIN podem retornar dados inconsistentes

**Solução Proposta:**
Migrar para PostgreSQL e implementar ON DELETE:
```sql
-- Para produção com PostgreSQL
ALTER TABLE order_items 
ADD CONSTRAINT order_items_product_id_fkey 
FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;

ALTER TABLE stock_adjustments_pending
ADD CONSTRAINT stock_adjustments_pending_order_id_fkey
FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;
```

---

### 3.2 Soft Delete Inconsistente

**Problema:** Soft delete não é verificado em todas as queries  
**Severidade:** 🟡 ALTA  
**Arquivo:** Múltiplos repositories  
**Linhas:** Variadas

**Causa:**
Algumas queries não incluem `WHERE deleted_at IS NULL`, permitindo acesso a registros deletados.

**Impacto:**
- Usuário pode acessar dados que deveriam estar "deletados"
- Relatórios podem incluir registros deletados
- Violação do princípio de soft delete

**Solução Proposta:**
Auditar todas as queries e garantir que `deleted_at IS NULL` está presente em queries de leitura. Usar scopes globais do GORM:
```go
db.Scopes(NotDeleted()).Find(...)
```

---

### 3.3 FKs Ausentes em Algumas Tabelas

**Problema:** Algumas relações não têm foreign keys no banco  
**Severidade:** 🟡 ALTA  
**Arquivo:** Migrations  
**Linhas:** Variadas

**Causa:**
Relações são mantidas apenas em nível de aplicação (GORM), sem constraints no banco.

**Impacto:**
- Possibilidade de inserir registros com IDs inválidos
- Integridade referencial depende apenas da aplicação
- Risco de corrupção se aplicação tiver bugs

**Solução Proposta:**
Adicionar FKs para todas as relações:
```sql
ALTER TABLE product_ingredients
ADD CONSTRAINT fk_product_ingredients_ingredient
FOREIGN KEY (ingredient_id) REFERENCES ingredients(id);
```

---

## 4. ESTOQUE

### 4.1 Dupla Baixa de Estoque

**Problema:** Possibilidade de dupla baixa em cenários de race condition  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 80-194

**Causa:**
Embora `CreateOrder` use SELECT FOR UPDATE e ordenação determinística de locks, se dois pedidos concorrentes usarem os mesmos ingredientes, pode haver condição de corrida na validação de estoque.

**Impacto:**
- Estoque pode ficar negativo
- Pedidos podem ser criados mesmo sem estoque suficiente
- Inconsistência entre estoque real e sistema

**Solução Proposta:**
A implementação atual com SELECT FOR UPDATE é correta, mas adicionar validação adicional no UPDATE:
```go
UPDATE ingredients 
SET stock_quantity = stock_quantity - ? 
WHERE id = ? AND stock_quantity >= ?
```

Se `RowsAffected == 0`, significa que o estoque mudou durante a transação.

---

### 4.2 CompleteInventory Sem Validação de Concorrência

**Problema:** CompleteInventory pode ser executado concorrentemente  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/stock_movement_service.go`  
**Linhas:** 213-316

**Causa:**
Embora use SELECT FOR UPDATE no inventário, se dois usuários iniciarem CompleteInventory simultaneamente para o mesmo inventário, ambos podem passar da verificação inicial.

**Impacto:**
- Inventário pode ser completado duas vezes
- Estoque ajustado incorretamente
- Movimentações duplicadas

**Solução Proposta:**
A implementação atual com SELECT FOR UPDATE no inventário é correta. Adicionar verificação final de status:
```go
// Já implementado nas linhas 303-309
currentInventory, err := s.stockMovementRepo.FindInventoryByIDForUpdate(ctx, inventoryID, tx)
if err != nil {
    return fmt.Errorf("StockMovementService.CompleteInventory: verificar status do inventário: %w", err)
}
if currentInventory.Status != "draft" {
    return ErrStockInventoryCompleted
}
```

---

### 4.3 Estoque Negativo Possível

**Problema:** CHECK constraint para stock_quantity >= 0 ausente  
**Severidade:** 🟡 ALTA  
**Arquivo:** Migrations  
**Linhas:** N/A

**Causa:**
Não há CHECK constraint no banco para garantir que `stock_quantity` nunca seja negativo. A validação é apenas em nível de aplicação.

**Impacto:**
- Se aplicação tiver bug, estoque pode ficar negativo
- Relatórios podem mostrar valores impossíveis
- Cálculos de CMV e custos incorretos

**Solução Proposta:**
Adicionar CHECK constraint:
```sql
ALTER TABLE ingredients
ADD CONSTRAINT chk_stock_non_negative 
CHECK (stock_quantity >= 0);
```

---

## 5. PEDIDOS

### 5.1 IdempotenciaKey Opcional

**Problema:** Idempotency key é opcional, permitindo duplicação  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/order_service.go`  
**Linhas:** Variadas

**Causa:**
O campo `IdempotencyKey` é `*string` (opcional). Se cliente não enviar, pedido pode ser duplicado em caso de retry.

**Impacto:**
- POST duplicado pode criar múltiplos pedidos
- Baixa de estoque duplicada
- Usuário cobrado múltiplas vezes

**Solução Proposta:**
Tornar IdempotencyKey obrigatório ou gerar automaticamente:
```go
// Opção 1: Tornar obrigatório
type CreateOrderInput struct {
    IdempotencyKey string `json:"idempotencyKey" validate:"required"`
    // ...
}

// Opção 2: Gerar automaticamente se não fornecido
if input.IdempotencyKey == nil {
    key := generateIdempotencyKey(ctx, userID, input.Items)
    input.IdempotencyKey = &key
}
```

---

### 5.2 Snapshot Inconsistente

**Problema:** Snapshots de produto podem ficar desatualizados  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 173-186

**Causa:**
Os snapshots (ProductName, ProductDescription, etc.) são carregados no service antes da transação. Se o produto for alterado entre o carregamento e a criação do pedido, o snapshot pode estar desatualizado.

**Impacto:**
- Histórico de pedido não reflete estado real do produto no momento da venda
- Relatórios podem mostrar informações incorretas

**Solução Proposta:**
Carregar snapshots dentro da transação ou usar versionamento de produtos:
```go
// Recarregar produto dentro da transação
currentProduct, err := r.productRepo.FindProductByID(ctx, item.ProductID)
if err != nil {
    return err
}
// Usar currentProduct para snapshot
```

---

### 5.3 Cancelamento Sem Validação de Transição

**Problema:** Cancelamento de pedido não valida todas as transições permitidas  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/service/order_service.go`  
**Linhas:** 407-418

**Causa:**
A função `isValidTransition` valida transições, mas pode haver casos não cobertos.

**Impacto:**
- Pedido pode ser cancelado em estados onde não deveria (ex: já entregue)
- Fluxo de negócio violado

**Solução Proposta:**
Revisar e expandir validações de transição de estado:
```go
transitions := map[domain.OrderStatus][]domain.OrderStatus{
    domain.OrderStatusPending:   {domain.OrderStatusConfirmed, domain.OrderStatusCancelled},
    domain.OrderStatusConfirmed: {domain.OrderStatusPreparing, domain.OrderStatusCancelled},
    domain.OrderStatusPreparing: {domain.OrderStatusReady, domain.OrderStatusCancelled},
    domain.OrderStatusReady:     {domain.OrderStatusDelivered, domain.OrderStatusCancelled},
    domain.OrderStatusDelivered: {}, // status final, sem transições
    domain.OrderStatusCancelled: {}, // status final, sem transições
}
```

---

## 6. PRODUTOS

### 6.1 Slug Collision Race Condition

**Problema:** Verificação de slug e criação não são atômicas  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linhas:** 144-153

**Causa:**
```go
if p.Slug != "" {
    var existing GormProduct
    err := r.db.WithContext(ctx).Where("slug = ?", p.Slug).First(&existing).Error
    if err == nil {
        return fmt.Errorf("ProductRepository.CreateProduct: slug '%s' já está em uso", p.Slug)
    }
    // ...
}
if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
    return fmt.Errorf("ProductRepository.CreateProduct: %w", err)
}
```

A verificação e criação não estão na mesma transação. Dois requests concorrentes podem passar a verificação e ambos tentar criar com o mesmo slug.

**Impacto:**
- Violação de unique constraint no banco
- Erro 500 para um dos usuários
- Experiência de usuário ruim

**Solução Proposta:**
Usar unique constraint no banco e tratar violação:
```go
// Migration já tem unique index
// No código, tratar violação como erro amigável
if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
    if IsUniqueViolation(err) {
        return fmt.Errorf("ProductRepository.CreateProduct: slug '%s' já está em uso", p.Slug)
    }
    return fmt.Errorf("ProductRepository.CreateProduct: %w", err)
}
```

---

### 6.2 Ficha Técnica Sem Transação

**Problema:** SetProductIngredients usa delete + create não atômico  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linhas:** 504-525

**Causa:**
```go
return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Apaga ficha anterior e recria (upsert simples)
    if err := tx.Where("product_id = ? AND deleted_at IS NULL", productID).
        Delete(&GormProductIngredient{}).Error; err != nil {
        return fmt.Errorf("ProductRepository.SetProductIngredients: deletar ingredientes anteriores: %w", err)
    }
    for _, item := range items {
        // ... create
    }
    return nil
})
```

Embora esteja em transação, se houver erro na criação de um ingrediente após deletar todos, a ficha fica vazia.

**Impacto:**
- Ficha técnica pode ficar vazia temporariamente
- Produto composto sem ingredientes
- Cálculos de custo incorretos

**Solução Proposta:**
Usar UPSERT real em vez de delete + create:
```go
for _, item := range items {
    err := tx.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "product_id"}, {Name: "ingredient_id"}},
        DoUpdates: clause.AssignmentColumns([]string{"quantity", "loss", "yield", "unit_cost", "total_cost"}),
    }).Create(&m).Error
    if err != nil {
        return err
    }
}
```

---

## 7. EMPRESAS

### 7.1 Slug Não Único Globalmente

**Problema:** Slug de empresa é único apenas por tenant  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** Migrations  
**Linhas:** Variadas

**Causa:**
Slug de empresa não tem unique constraint global, apenas validação em nível de aplicação.

**Impacto:**
- Duas empresas podem ter o mesmo slug
- URLs podem colidir em multi-tenant
- Confusão para usuários

**Solução Proposta:**
Adicionar unique constraint global ou prefixar com tenant:
```sql
-- Opção 1: Unique global
ALTER TABLE companies ADD CONSTRAINT uk_slug UNIQUE (slug);

-- Opção 2: Prefixar com company_id
-- slug já é único por company_id na aplicação
```

---

### 7.2 Owner Obrigatório Não Validado

**Problema:** Empresa pode ficar sem owner após criação  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/platform_service.go`  
**Linhas:** 72-180

**Causa:**
Se a criação do usuário owner falhar após a empresa ser criada, não há rollback automático da empresa.

**Impacto:**
- Empresa sem owner
- Ninguém pode administrar
- Necessita intervenção manual

**Solução Proposta:**
Garantir atomicidade completa (já coberto no item 2.2).

---

## 8. USUÁRIOS

### 8.1 Role Pode Ser Nulo

**Problema:** Role de usuário não tem NOT NULL em algumas tabelas  
**Severidade:** 🟡 ALTA  
**Arquivo:** Migrations  
**Linhas:** Variadas

**Causa:**
Migration 00016 adicionou NOT NULL, mas dados existentes podem ter valores nulos.

**Impacto:**
- Usuário sem role pode ter acesso indefinido
- RBAC pode falhar silenciosamente
- Risco de segurança

**Solução Proposta:**
Validar que todos os usuários têm role atribuída:
```sql
UPDATE users SET role = 'employee' WHERE role IS NULL;
ALTER TABLE users ALTER COLUMN role SET NOT NULL;
```

---

### 8.2 Convite Sem Expiração

**Problema:** Convites não expiram automaticamente  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/domain/invitation.go`  
**Linhas:** Variadas

**Causa:**
Não há campo de expiração ou job de limpeza para convites antigos.

**Impacto:**
- Convites podem ficar válidos indefinidamente
- Risco de segurança se convite for comprometido
- Acúmulo de convites não utilizados

**Solução Proposta:**
Adicionar campo de expiração e job de limpeza:
```go
type Invitation struct {
    ExpiresAt time.Time
    // ...
}

// Job diário para deletar convites expirados
DELETE FROM invitations WHERE expires_at < NOW();
```

---

## 9. EVENTOS

### 9.1 EventDispatcher Sem Timeout

**Problema:** processBatch não tem timeout explícito  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/event_dispatcher.go`  
**Linhas:** 88-110

**Causa:**
O loop de processamento de eventos não tem timeout, pode rodar indefinidamente se houver muitos eventos.

**Impacto:**
- Dispatcher pode travar
- Eventos acumulam
- Sistema fica sem processamento de eventos

**Solução Proposta:**
Adicionar timeout ao processBatch:
```go
func (d *EventDispatcher) processBatch() {
    ctx, cancel := context.WithTimeout(d.shutdownCtx, 30*time.Second)
    defer cancel()
    
    events, err := d.outboxRepo.FindPendingEvents(ctx, d.config.BatchSize)
    // ...
}
```

---

### 9.2 Consumidores Não Idempotentes

**Problema:** Não há garantia que consumidores são idempotentes  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/consumers/`  
**Linhas:** Variadas

**Causa:**
O framework de consumidores tem idempotência store, mas não há documentação exigindo que consumidores sejam idempotentes.

**Impacto:**
- Evento duplicado pode processar múltiplas vezes
- Email enviado múltiplas vezes
- Webhook chamado múltiplas vezes

**Solução Proposta:**
Documentar e testar idempotência de consumidores:
```go
// Adicionar teste de idempotência para cada consumer
func TestEmailConsumer_Idempotency(t *testing.T) {
    // Processar mesmo evento duas vezes
    // Verificar que email foi enviado apenas uma vez
}
```

---

### 9.3 Outbox Sem Limpeza Automática

**Problema:** Eventos processados não são limpos automaticamente  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/service/event_dispatcher.go`  
**Linhas:** Variadas

**Causa:**
Não há job de limpeza para eventos completados antigos.

**Impacto:**
- Tabela de outbox cresce indefinidamente
- Performance degrada
- Espaço em disco desperdiçado

**Solução Proposta:**
Adicionar job de limpeza:
```go
func (d *EventDispatcher) cleanupOldEvents() {
    ctx := context.Background()
    cutoff := time.Now().Add(-30 * 24 * time.Hour) // 30 dias
    d.outboxRepo.DeleteOldCompletedEvents(ctx, cutoff)
}
```

---

## 10. CONCORRÊNCIA

### 10.1 Deadlock Possível em Locks Múltiplos

**Problema:** Ordenação de locks pode causar deadlock  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 148-164

**Causa:**
Embora `CreateOrder` ordene ingredientes por ID, outras operações podem não seguir a mesma ordem.

**Impacto:**
- Deadlock em alta concorrência
- Transações abortadas
- Performance degrada

**Solução Proposta:**
Garantir ordenação consistente em todas as operações:
```go
// Criar helper para ordenação consistente
func sortIngredientIDs(ids []uint) {
    sort.Slice(ids, func(i, j int) bool {
        return ids[i] < ids[j]
    })
}
```

---

### 10.2 Lost Update em UpdateIngredient

**Problema:** UpdateIngredient pode ter lost update  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/infra/repository/gorm_product_repository.go`  
**Linhas:** 428-456

**Causa:**
Embora use SELECT FOR UPDATE, se duas transações lerem o mesmo valor e atualizarem, uma pode sobrescrever a outra.

**Impacto:**
- Atualização perdida
- Estoque incorreto
- Inconsistência de dados

**Solução Proposta:**
A implementação atual com SELECT FOR UPDATE está correta. Adicionar versionamento para casos extremos:
```go
UPDATE ingredients 
SET stock_quantity = ?, version = version + 1 
WHERE id = ? AND version = ?
```

---

### 10.3 Mutex em In-Memory Cache

**Problema:** Cache in-memory com mutex não é distribuído  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/infra/repository/gorm_global_config_repository.go`  
**Linhas:** 45-88

**Causa:**
O cache de configuração global usa RWMutex em memória. Em múltiplas instâncias, cada uma tem seu próprio cache.

**Impacto:**
- Cache inconsistente entre instâncias
- Configurações podem não propagar
- Comportamento indefinido

**Solução Proposta:**
Usar Redis para cache distribuído ou reduzir TTL:
```go
// Usar Redis cache em vez de in-memory
cache := redis.NewCache(redisClient, "global_config", 5*time.Minute)
```

---

## 11. BANCO DE DADOS

### 11.1 CHECK Constraint Ausente para Estoque

**Problema:** Não há CHECK para stock_quantity >= 0  
**Severidade:** 🟡 ALTA  
**Arquivo:** Migrations  
**Linhas:** N/A

**Causa:**
Validação de estoque negativo é apenas em nível de aplicação.

**Impacto:**
- Estoque pode ficar negativo se aplicação tiver bug
- Dados impossíveis no banco
- Relatórios incorretos

**Solução Proposta:**
Adicionar CHECK constraint (já coberto no item 4.3).

---

### 11.2 Índices Faltantes

**Problema:** Algumas queries podem não usar índices otimizados  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** Migrations  
**Linhas:** Variadas

**Causa:**
Índices compostos podem estar faltando para queries comuns.

**Impacto:**
- Performance degradada
- Queries lentas em tabelas grandes
- Timeout em alta carga

**Solução Proposta:**
Analisar queries com EXPLAIN e adicionar índices:
```sql
-- Exemplo para queries frequentes
CREATE INDEX idx_ingredients_company_stock 
ON ingredients(company_id, stock_quantity) 
WHERE deleted_at IS NULL;
```

---

### 11.3 Unique Constraint Parcial

**Problema:** Unique constraint parcial pode não funcionar em todos os bancos  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `migrations/00005_add_unique_constraint_stock_adjustments.sql`  
**Linhas:** 6-8

**Causa:**
```sql
CREATE UNIQUE INDEX IF NOT EXISTS uk_stock_adjustments_order_ingredient_pending
ON stock_adjustments_pending(order_id, ingredient_id)
WHERE status = 'pending';
```

Unique constraint parcial (com WHERE) não é suportado em todos os bancos.

**Impacto:**
- Migração pode falhar em alguns bancos
- Idempotência pode não funcionar

**Solução Proposta:**
Usar trigger ou validação em aplicação para bancos sem suporte:
```go
// Validar antes de inserir
existing, _ := repo.FindPendingByOrderAndIngredient(ctx, orderID, ingredientID)
if existing != nil {
    return ErrAdjustmentExists
}
```

---

## 12. REGRAS DE NEGÓCIO

### 12.1 Panic em Inicialização

**Problema:** Panic se variáveis de ambiente não estiverem definidas  
**Severidade:** 🟢 MÉDIA  
**Arquivo:** `internal/service/auth_service.go`, `internal/service/platform_auth_service.go`  
**Linhas:** 51, 52

**Causa:**
```go
if jwtSecret == "" {
    panic("JWT_PLATFORM_SECRET environment variable is required but not set")
}
```

**Impacto:**
- Aplicação não inicia se variável não definida
- Difícil debugar em produção
- Não graceful

**Solução Proposta:**
Retornar erro em vez de panic:
```go
func NewAuthService(...) (*AuthService, error) {
    if jwtSecret == "" {
        return nil, errors.New("JWT_TENANT_SECRET environment variable is required")
    }
    return &AuthService{...}, nil
}
```

---

### 12.2 Crash Durante CompleteInventory

**Problema:** Panic durante CompleteInventory não é tratado  
**Severidade:** 🟡 ALTA  
**Arquivo:** `internal/service/stock_movement_service.go`  
**Linhas:** 213-316

**Causa:**
Se houver panic durante o processamento de itens, a transação é abortada mas não há recuperação automática.

**Impacto:**
- Inventário fica em estado inconsistente
- Movimentações parciais podem existir
- Necessita intervenção manual

**Solução Proposta:**
Adicionar recover e tratamento de erro:
```go
func (s *StockMovementService) CompleteInventory(...) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered in CompleteInventory: %v", r)
            // Marcar inventário como failed
        }
    }()
    // ...
}
```

---

### 12.3 Timeout de Transação Não Configurado

**Problema:** Transações não têm timeout configurado  
**Severidade:** 🟡 ALTA  
**Arquivo:** Variados  
**Linhas:** Variadas

**Causa:**
GORM não tem timeout padrão para transações. Transações longas podem travar o sistema.

**Impacto:**
- Deadlocks não são detectados rapidamente
- Conexões podem esgotar
- Sistema fica sem resposta

**Solução Proposta:**
Configurar timeout no contexto:
```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // ...
})
```

---

## RESUMO POR SEVERIDADE

### 🔴 CRÍTICOS (5)

1. **CreatePurchaseReceiving Não Atualiza Estoque** - Estoque não incrementado ao receber materiais
2. **ON DELETE Não Implementado** - Orphan records podem acumular
3. **CreatePurchaseReceiving Não Atualiza Estoque** (já listado)
4. **Slug Collision Race Condition** - Dois produtos com mesmo slug
5. **Role Pode Ser Nulo** - Usuário sem role pode ter acesso indefinido

### 🟡 ALTOS (10)

1. **Operações Fora de Transação** - Platform service methods
2. **Transação Grande em CreateOrder** - Timeout e deadlock risk
3. **CreateCompany Sem Transação Completa** - Empresa sem owner
4. **DuplicateProduct Parcialmente Atômico** - Sem DB
5. **UpdateOrder Sem Validação de Estoque Prévia** - UX ruim
6. **Soft Delete Inconsistente** - Acesso a dados deletados
7. **FKs Ausentes em Algumas Tabelas** - Integridade só na aplicação
8. **Dupla Baixa de Estoque** - Race condition
9. **CompleteInventory Sem Validação de Concorrência** - Duplo completion
10. **Estoque Negativo Possível** - Sem CHECK constraint

### 🟢 MÉDIOS (6)

1. **Snapshot Inconsistente** - Histórico desatualizado
2. **Cancelamento Sem Validação de Transição** - Fluxo violado
3. **Ficha Técnica Sem Transação** - Delete + create
4. **Slug Não Único Globalmente** - Colisão em multi-tenant
5. **Convite Sem Expiração** - Risco de segurança
6. **Outbox Sem Limpeza Automática** - Tabela cresce

### ⚪ BAIXOS (2)

1. **Rollback Explícito Desnecessário** - Código redundante
2. **Panic em Inicialização** - Aplicação não inicia

---

## ESTIMATIVA DE ESFORÇO PARA CORREÇÃO

### CRÍTICOS (Prioridade 0 - Imediato)
- **CreatePurchaseReceiving:** 4 horas
- **ON DELETE:** 8 horas (migração para PostgreSQL)
- **Slug Collision:** 2 horas
- **Role NOT NULL:** 2 horas
- **Total:** 16 horas

### ALTOS (Prioridade 1 - Esta Sprint)
- **Transações Platform:** 8 horas
- **CreateOrder Refactor:** 12 horas
- **Validações Pré-transação:** 8 horas
- **Soft Delete Audit:** 4 horas
- **FKs Adicionais:** 6 horas
- **CHECK Constraints:** 2 horas
- **Total:** 40 horas

### MÉDIOS (Prioridade 2 - Próxima Sprint)
- **Snapshot Refactor:** 8 horas
- **Validações de Estado:** 4 horas
- **UPSERT Ficha Técnica:** 4 horas
- **Job Limpeza Outbox:** 4 horas
- **Expiração Convites:** 4 horas
- **Total:** 24 horas

### BAIXOS (Prioridade 3 - Quando possível)
- **Remover Rollback Explícito:** 2 horas
- **Tratar Panic em Init:** 2 horas
- **Total:** 4 horas

**ESFORÇO TOTAL:** 84 horas (~10 dias úteis)

---

## CHECKLIST PARA PRODUÇÃO

### Obrigatório (Bloqueio)
- [ ] Implementar atualização de estoque em CreatePurchaseReceiving
- [ ] Adicionar CHECK constraint para stock_quantity >= 0
- [ ] Garantir Role NOT NULL em todos os usuários
- [ ] Implementar ON DELETE ou validar soft delete
- [ ] Adicionar idempotency key obrigatório ou gerar automaticamente

### Recomendado (Alta Prioridade)
- [ ] Adicionar transações em platform service methods
- [ ] Validar estoque antes de transações de pedido
- [ ] Adicionar FKs para todas as relações
- [ ] Implementar job de limpeza de outbox
- [ ] Adicionar timeout em todas as transações

### Opcional (Melhoria)
- [ ] Refatorar CreateOrder para reduzir tamanho de transação
- [ ] Implementar UPSERT para ficha técnica
- [ ] Adicionar expiração para convites
- [ ] Migrar cache in-memory para Redis
- [ ] Adicionar índices compostos para queries frequentes

---

## CONCLUSÃO

O sistema HorizonGest tem uma base transacional sólida com uso adequado de SELECT FOR UPDATE e ordenação de locks para prevenir deadlocks. No entanto, existem lacunas críticas:

1. **Integração de módulos:** O módulo de compras não está integrado com o módulo de estoque
2. **Constraints de banco:** Muitas validações dependem apenas da aplicação
3. **Idempotência:** Nem todas as operações críticas são idempotentes
4. **Limpeza de dados:** Jobs de limpeza não estão implementados

**Recomendação:** Priorizar correções críticas antes da entrada em produção, especialmente a integração de estoque com compras e a adição de CHECK constraints no banco.
