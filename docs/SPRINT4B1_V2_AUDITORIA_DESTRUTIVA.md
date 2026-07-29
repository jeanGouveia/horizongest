# AUDITORIA DESTRUTIVA — Sprint 4B.1 v2

**Data:** 27 de Julho de 2026  
**Auditor:** Principal Software Engineer  
**Status:** ❌ **REPROVADO** — Problemas Críticos Encontrados

---

## RESUMO EXECUTIVO

A Sprint 4B.1 v2 contém **PROBLEMAS CRÍTICOS** que comprometem a integridade transacional:

1. **BUG CRÍTICO #1:** CompleteInventory faz leituras fora da transação (GetInventoryByID, ListInventoryItems)
2. **BUG CRÍTICO #2:** CreateOrder não ordena locks → DEADLOCK POSSÍVEL
3. **BUG CRÍTICO #3:** UpdateIngredient faz SELECT sem FOR UPDATE → LOST UPDATE POSSÍVEL
4. **BUG CRÍTICO #4:** CompleteInventory não valida se inventário foi modificado durante processamento

**PAREDER:** **REPROVAR** — Não aprovar para produção sem correções.

---

## 1. PROPAGAÇÃO DE TRANSAÇÕES

### 1.1 ❌ BUG CRÍTICO: CompleteInventory - Leituras Fora da Transação

**Arquivo:** `backend/internal/service/stock_movement_service.go`  
**Linhas:** 211-227

```go
func (s *StockMovementService) CompleteInventory(ctx context.Context, inventoryID, userID uint) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // ❌ BUG: GetInventoryByID usa r.db, não tx
        inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID)
        if err != nil {
            return ErrStockInventoryNotFound
        }

        // ❌ BUG: ListInventoryItems usa r.db, não tx
        items, err := s.stockMovementRepo.ListInventoryItems(ctx, inventoryID)
        if err != nil {
            return fmt.Errorf("erro ao buscar itens do inventário: %w", err)
        }
        
        // ... processamento ...
    })
}
```

**Problema:**
- `GetInventoryByID` e `ListInventoryItems` estão DENTRO da transação mas usam `r.db`
- Isso significa que estas leituras não participam da transação
- Se outra transação modificar o inventário entre a leitura e o processamento, haverá inconsistência

**Cenário de Falha:**
```
Tempo | Transação A (CompleteInventory) | Transação B (AddInventoryItem)
------|--------------------------------|----------------------------------
T1    | BEGIN                           |
T2    | GetInventoryByID (r.db)         |
T3    |                                 | AddInventoryItem (r.db)
T4    |                                 | Commit (item adicionado)
T5    | ListInventoryItems (r.db)       | 
T6    | Processa itens (inclui novo)    |
T7    | Commit                          |
```

**Resultado:** O novo item adicionado por B é processado por A, mas A não tinha lock no inventário.

**Correção Necessária:**
```go
// Repository methods precisam aceitar tx
func (r *GormStockMovementRepository) GetInventoryByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error) {
    query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
    // ...
}

func (r *GormStockMovementRepository) ListInventoryItems(ctx context.Context, inventoryID uint, tx *gorm.DB) ([]domain.StockInventoryItem, error) {
    var items []domain.StockInventoryItem
    err := r.getDB(ctx, tx).Where("inventory_id = ? AND deleted_at IS NULL", inventoryID).
        Preload("Ingredient").
        Find(&items).Error
    return items, err
}
```

---

## 2. SELECT FOR UPDATE

### 2.1 ❌ BUG CRÍTICO: UpdateIngredient - SELECT Sem FOR UPDATE

**Arquivo:** `backend/internal/infra/repository/gorm_product_repository.go`  
**Linhas:** 409-430

```go
func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient, tx *gorm.DB) error {
    // ❌ BUG: SELECT sem FOR UPDATE
    var existing GormIngredient
    query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), i.ID)
    if err := query.Where("deleted_at IS NULL").First(&existing).Error; err != nil {
        // ...
    }

    // Update sem lock
    m := GormIngredient{...}
    if err := r.getDB(ctx, tx).Save(&m).Error; err != nil {
        // ...
    }
    return nil
}
```

**Problema:**
- O SELECT inicial não tem FOR UPDATE
- Entre o SELECT e o UPDATE, outra transação pode modificar o ingrediente
- Isso pode causar lost update

**Cenário de Falha:**
```
Tempo | Transação A (UpdateIngredient) | Transação B (CreateStockMovement)
------|--------------------------------|----------------------------------
T1    | BEGIN                           |
T2    | SELECT ingredient (sem lock)    |
T3    |                                 | BEGIN
T4    |                                 | SELECT FOR UPDATE ingredient
T5    |                                 | UPDATE stock_quantity = 100
T6    |                                 | COMMIT
T7    | UPDATE ingredient (sobrescreve) |
T8    | COMMIT                          |
```

**Resultado:** O update de B é perdido (lost update).

**Correção Necessária:**
```go
func (r *GormProductRepository) UpdateIngredient(ctx context.Context, i *domain.Ingredient, tx *gorm.DB) error {
    var existing GormIngredient
    query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), i.ID)
    
    // ✅ Adicionar FOR UPDATE
    if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("deleted_at IS NULL").
        First(&existing).Error; err != nil {
        // ...
    }

    // ...
}
```

---

## 3. LOST UPDATE

### 3.1 ❌ PROVA MATEMÁTICA: Lost Update em UpdateIngredient

**Cenário:**
- Ingrediente X tem stock_quantity = 50
- Transação A: UpdateIngredient para 60
- Transação B: CreateStockMovement para 55

**Sequência sem FOR UPDATE:**
```
1. A: SELECT stock_quantity = 50
2. B: SELECT FOR UPDATE stock_quantity = 50
3. B: UPDATE stock_quantity = 55
4. B: COMMIT
5. A: UPDATE stock_quantity = 60 (sobrescreve 55)
6. A: COMMIT
```

**Resultado Final:** stock_quantity = 60  
**Resultado Esperado:** stock_quantity = 65 (55 + 5) ou 60 (dependendo da ordem)

**Prova:**
- A leu 50, calculou 60
- B leu 50, calculou 55, escreveu 55
- A sobrescreveu 55 com 60
- O update de B foi perdido

**Conclusão:** Lost update é possível sem FOR UPDATE em UpdateIngredient.

---

## 4. DEADLOCKS

### 4.1 ❌ BUG CRÍTICO: CreateOrder - Sem Ordenação de Locks

**Arquivo:** `backend/internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 143-148

```go
// ❌ BUG: Loop sem ordenação de ingredientes
for _, pi := range ingredients {
    consumo := pi.Quantity * item.Quantity
    if err := r.productRepo.DecreaseIngredientStock(ctx, pi.IngredientID, consumo, tx, ...); err != nil {
        return fmt.Errorf("CreateOrder: baixa estoque: %w", err)
    }
}
```

**Problema:**
- Os ingredientes são processados na ordem que aparecem no map
- Maps em Go não têm ordem garantida
- Se dois pedidos têm ingredientes em ordens diferentes, deadlock é possível

**Cenário de Deadlock:**
```
Tempo | Transação A (Pedido 1) | Transação B (Pedido 2)
------|-----------------------|-----------------------
T1    | BEGIN                 |
T2    | Lock ingrediente 3    |
T3    |                       | BEGIN
T4    |                       | Lock ingrediente 1
T5    | Tentar lock 1         |
T6    |                       | Tentar lock 3
T7    | ❌ DEADLOCK           | ❌ DEADLOCK
```

**Diagrama de Deadlock:**
```
Pedido 1: [Produto X → Ingrediente 3, Produto Y → Ingrediente 1]
Pedido 2: [Produto Z → Ingrediente 1, Produto W → Ingrediente 3]

A locks 3 → B locks 1 → A espera 1 → B espera 3 → DEADLOCK
```

**Correção Necessária:**
```go
// ✅ Ordenar ingredientes antes do loop
ingredientIDs := make([]uint, 0, len(ingredients))
for id := range ingredients {
    ingredientIDs = append(ingredientIDs, id)
}
sort.Slice(ingredientIDs, func(i, j int) bool {
    return ingredientIDs[i] < ingredientIDs[j]
})

for _, ingredientID := range ingredientIDs {
    pi := ingredients[ingredientID]
    consumo := pi.Quantity * item.Quantity
    if err := r.productRepo.DecreaseIngredientStock(ctx, ingredientID, consumo, tx, ...); err != nil {
        return fmt.Errorf("CreateOrder: baixa estoque: %w", err)
    }
}
```

### 4.2 ✅ CompleteInventory - Ordenação Correta

**Arquivo:** `backend/internal/service/stock_movement_service.go`  
**Linhas:** 229-233

```go
// ✅ Ordenação correta
sort.Slice(items, func(i, j int) bool {
    return items[i].IngredientID < items[j].IngredientID
})
```

**Status:** ✅ CORRETO - Ordenação implementada

---

## 5. PHANTOM READS

### 5.1 ❌ POSSÍVEL: CompleteInventory - Phantom Read em ListInventoryItems

**Cenário:**
- Transação A: CompleteInventory
- Transação B: AddInventoryItem

**Sequência:**
```
1. A: BEGIN
2. A: ListInventoryItems (retorna 3 itens)
3. B: BEGIN
4. B: AddInventoryItem (adiciona 4º item)
5. B: COMMIT
6. A: Processa 3 itens (ignora 4º)
7. A: UpdateInventoryStatus (completed)
8. A: COMMIT
```

**Resultado:** O 4º item adicionado por B nunca é processado.

**Prova:**
- A leu 3 itens
- B adicionou 1 item
- A processou apenas os 3 itens originais
- O novo item fica em estado inconsistente (inventário completed mas item não processado)

**Correção Necessária:**
- ListInventoryItems deve usar FOR UPDATE na tabela stock_inventories
- Ou validar status do inventário após o loop

---

## 6. NON REPEATABLE READ

### 6.1 ✅ CreateStockMovement - Protegido por FOR UPDATE

**Status:** ✅ CORRETO - FindIngredientByIDForUpdate usa FOR UPDATE

### 6.2 ❌ CompleteInventory - Non Repeatable Read em GetInventoryByID

**Cenário:**
- Transação A: CompleteInventory
- Transação B: DeleteInventory

**Sequência:**
```
1. A: BEGIN
2. A: GetInventoryByID (status = draft)
3. B: BEGIN
4. B: DeleteInventory (soft delete)
5. B: COMMIT
6. A: ListInventoryItems (pode falhar ou retornar vazio)
7. A: Processa itens
8. A: UpdateInventoryStatus (pode falhar)
```

**Resultado:** Inconsistência no estado do inventário.

**Correção Necessária:**
- GetInventoryByID deve usar FOR UPDATE
- Ou validar status após cada operação

---

## 7. LOCK ESCALATION

### 7.1 ⚠️ RISCO: CompleteInventory - Muitos Locks

**Cenário:**
- Inventário com 100 ingredientes
- CompleteInventory faz 100 SELECT FOR UPDATE
- Cada lock é mantido durante toda a transação

**Risco:**
- Timeout de lock wait
- Baixo throughput em alta concorrência

**Mitigação:**
- Limitar número de ingredientes por inventário
- Usar batch processing
- Considerar SERIALIZABLE isolation level

**Status:** ⚠️ RISCO ACEITÁVEL para ERP típico (< 50 ingredientes/inventário)

---

## 8. ROLLBACK

### 8.1 ✅ CreateStockMovement - Rollback Completo

**Status:** ✅ CORRETO - Todas as operações dentro da transação

### 8.2 ❌ CompleteInventory - Rollback Parcial Possível

**Problema:**
- GetInventoryByID e ListInventoryItems estão fora da transação
- Se falharem após o loop, o status do inventário pode ficar inconsistente

**Cenário:**
```
1. A: BEGIN
2. A: GetInventoryByID (fora de tx) - OK
3. A: ListInventoryItems (fora de tx) - OK
4. A: Processa ingrediente 1 - OK
5. A: Processa ingrediente 2 - ERRO
6. A: ROLLBACK
7. Resultado: Movimentação do ingrediente 1 foi revertida
8. Mas as leituras iniciais já foram feitas
```

**Status:** ⚠️ RISCO BAIXO - Rollback ainda funciona para operações dentro da transação

---

## 9. ISOLATION LEVEL

### 9.1 ❌ DEPENDÊNCIA IMPLÍCITA EM READ COMMITTED

**Problema:**
- O código assume READ COMMITTED
- Se o banco estiver em READ UNCOMMITTED, pode ler dados não commitados
- Se estiver em SERIALIZABLE, pode ter mais deadlocks

**Correção Necessária:**
- Documentar explicitamente a dependência
- Ou configurar explicitamente no GORM

```go
db.Config().PrepareStmt = false
db.Set("default isolation level", "READ COMMITTED")
```

---

## 10. GORM

### 10.1 ❌ USO INCORRETO: Save() em UpdateIngredient

**Arquivo:** `backend/internal/infra/repository/gorm_product_repository.go`  
**Linha:** 427

```go
if err := r.getDB(ctx, tx).Save(&m).Error; err != nil {
```

**Problema:**
- `Save()` atualiza TODOS os campos, incluindo timestamps
- Pode causar race condition em `updated_at`
- `Update()` ou `UpdateColumn()` é mais seguro

**Correção Necessária:**
```go
if err := r.getDB(ctx, tx).Model(&GormIngredient{}).
    Where("id = ?", i.ID).
    Updates(map[string]interface{}{
        "name": i.Name,
        "unit": i.Unit,
        "stock_quantity": i.StockQuantity,
        "min_stock": i.MinStock,
        "active": i.Active,
    }).Error; err != nil {
```

---

## 11. POSTGRESQL

### 11.1 ✅ FOR UPDATE Implementado Corretamente

**Status:** ✅ CORRETO - Clauses(clause.Locking{Strength: "UPDATE"})

### 11.2 ⚠️ Índices Necessários

**Índices recomendados:**
```sql
-- Para SELECT FOR UPDATE eficiente
CREATE INDEX CONCURRENTLY idx_ingredients_company_id_deleted_at 
ON ingredients(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Para CompleteInventory
CREATE INDEX CONCURRENTLY idx_stock_inventory_items_inventory_id 
ON stock_inventory_items(inventory_id) WHERE deleted_at IS NULL;

-- Para CreateOrder
CREATE INDEX CONCURRENTLY idx_order_items_order_id 
ON order_items(order_id) WHERE deleted_at IS NULL;
```

---

## 12. MULTI-TENANT

### 12.1 ✅ Filtro de Tenant Aplicado

**Status:** ✅ CORRETO - ApplyTenantFilter usado em todas as queries

### 12.2 ⚠️ RISCO: Lock pode atingir outra empresa?

**Análise:**
- ApplyTenantFilter adiciona `WHERE company_id = ?`
- FOR UPDATE aplica lock apenas nas linhas retornadas
- Locks são por linha, não por tabela
- **Conclusão:** ✅ CORRETO - Locks não atingem outras empresas

---

## 13. PERFORMANCE

### 13.1 ⚠️ SELECT FOR UPDATE pode ser lento

**Problema:**
- FOR UPDATE bloqueia leituras concorrentes
- Em alta concorrência, pode causar lock wait timeout

**Mitigação:**
- Manter transações curtas
- Usar índices apropriados
- Considerar optimistic locking para leituras frequentes

**Status:** ⚠️ RISCO ACEITÁVEL para ERP típico

---

## 14. CONSISTÊNCIA

### 14.1 ❌ POSSÍVEL: Movimentação sem Estoque Atualizado

**Cenário:**
- CreateStockMovement cria movimentação
- UpdateIngredient falha
- Rollback reverte ambos

**Status:** ✅ CORRETO - Rollback funciona

### 14.2 ❌ POSSÍVEL: Estoque Atualizado sem Movimentação

**Cenário:**
- UpdateIngredient é chamado diretamente (sem CreateStockMovement)
- Não há movimentação registrada

**Status:** ⚠️ RISCO ACEITÁVEL - UpdateIngredient é usado apenas por services que garantem consistência

---

## 15. PAREDER FINAL

### ❌ **REPROVAR**

**Motivos:**
1. **BUG CRÍTICO #1:** CompleteInventory faz leituras fora da transação
2. **BUG CRÍTICO #2:** CreateOrder não ordena locks → DEADLOCK POSSÍVEL
3. **BUG CRÍTICO #3:** UpdateIngredient faz SELECT sem FOR UPDATE → LOST UPDATE POSSÍVEL
4. **BUG CRÍTICO #4:** CompleteInventory não valida modificações durante processamento

**Correções Necessárias antes de aprovar:**
1. Adicionar parâmetro `tx` em GetInventoryByID e ListInventoryItems
2. Ordenar ingredientes em CreateOrder
3. Adicionar FOR UPDATE em UpdateIngredient
4. Validar status do inventário após processamento
5. Substituir Save() por Updates() em UpdateIngredient
6. Documentar dependência de READ COMMITTED
7. Adicionar índices recomendados

**Estimativa de esforço:** 4-6 horas de desenvolvimento + 2 horas de testes

---

## 16. CONCLUSÃO

A Sprint 4B.1 v2 representa uma melhoria significativa em relação à v1, mas ainda contém **problemas críticos** que comprometem a integridade transacional. Os bugs identificados podem causar:

- **Deadlocks** em alta concorrência (CreateOrder)
- **Lost updates** em operações concorrentes (UpdateIngredient)
- **Inconsistência de dados** em CompleteInventory (leituras fora de transação)

**Recomendação:** Não aprovar para produção sem correções.
