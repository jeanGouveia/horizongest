# Sprint 4B.5 — Auditoria Destrutiva do BUG CRÍTICO #4

**Data:** 27 de Julho de 2026  
**Versão:** 1.0  
**Tipo:** Auditoria Destrutiva  
**Status:** Concluída

---

## Contexto

Esta auditoria analisa se o **BUG CRÍTICO #4** identificado na auditoria original ainda existe após as correções das Sprints 4B.2, 4B.3 e 4B.4.

**BUG CRÍTICO #4 (Original):** CompleteInventory não valida se o inventário foi modificado durante o processamento.

---

## Parte 1 — O inventário fica protegido por lock pessimista durante CompleteInventory?

**Resposta:** NÃO.

### Análise Detalhada

**Arquivo:** `internal/service/stock_movement_service.go`  
**Método:** `CompleteInventory` (linhas 212-305)

```go
func (s *StockMovementService) CompleteInventory(ctx context.Context, inventoryID, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Buscar inventário DENTRO da transação (passando tx)
		inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID, tx)
		// ...
		
		// 2. Buscar itens DENTRO da transação (passando tx)
		items, err := s.stockMovementRepo.ListInventoryItems(ctx, inventoryID, tx)
		// ...
		
		// 3. Ajustar estoque para cada item
		for _, item := range items {
			// SELECT FOR UPDATE em ingredientes (não em inventário)
			ingredient, err := s.productRepo.FindIngredientByIDForUpdate(ctx, item.IngredientID, tx)
			// ...
		}
		
		// 4. Atualizar status do inventário
		if err := s.stockMovementRepo.UpdateInventoryStatus(ctx, inventoryID, "completed", tx); err != nil {
			// ...
		}
	})
}
```

**Arquivo:** `internal/infra/repository/gorm_stock_movement_repository.go`  
**Método:** `GetInventoryByID` (linhas 71-78)

```go
func (r *GormStockMovementRepository) GetInventoryByID(ctx context.Context, id uint, tx *gorm.DB) (*domain.StockInventory, error) {
	var inventory domain.StockInventory
	query := ApplyTenantFilterWithID(ctx, r.getDB(ctx, tx), id)
	err := query.Where("deleted_at IS NULL").
		Preload("Items.Ingredient").
		First(&inventory).Error  // ❌ SEM SELECT FOR UPDATE
	return &inventory, err
}
```

**Arquivo:** `internal/infra/repository/gorm_stock_movement_repository.go`  
**Método:** `UpdateInventoryStatus` (linhas 94-98)

```go
func (r *GormStockMovementRepository) UpdateInventoryStatus(ctx context.Context, id uint, status string, tx *gorm.DB) error {
	return r.getDB(ctx, tx).Model(&domain.StockInventory{}).
		Where("id = ?", id).
		Update("status", status).Error  // ❌ SEM SELECT FOR UPDATE
}
```

### Conclusão

**NÃO há lock pessimista na tabela `stock_inventories` durante `CompleteInventory`.**

Os locks pessimistas (SELECT FOR UPDATE) são aplicados APENAS na tabela `ingredients` através de `FindIngredientByIDForUpdate`. A tabela `stock_inventories` e `stock_inventory_items` não têm nenhum lock pessimista.

---

## Parte 2 — Existe alguma possibilidade de modificação durante CompleteInventory?

**Resposta:** SIM.

### 2.1 Adicionar Item (AddInventoryItem)

**Arquivo:** `internal/service/stock_movement_service.go`  
**Método:** `AddInventoryItem` (linhas 175-207)

```go
func (s *StockMovementService) AddInventoryItem(ctx context.Context, inventoryID, ingredientID uint, expectedStock, actualStock float64, reason string) (*domain.StockInventoryItem, error) {
	// Buscar inventário
	inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID, nil)  // ❌ FORA de transação
	if err != nil {
		return nil, ErrStockInventoryNotFound
	}

	// Validar status
	if inventory.Status != "draft" {
		return nil, ErrStockInventoryCompleted
	}

	// Criar item
	if err := s.stockMovementRepo.CreateInventoryItem(ctx, item, nil); err != nil {  // ❌ FORA de transação
		return nil, fmt.Errorf("erro ao criar item de inventário: %w", err)
	}

	return item, nil
}
```

**Caminho do Código:**
1. `AddInventoryItem` chama `GetInventoryByID` com `nil` (fora de transação)
2. Valida status == "draft"
3. Chama `CreateInventoryItem` com `nil` (fora de transação)
4. INSERT em `stock_inventory_items`

**Possibilidade:** ✅ SIM - `AddInventoryItem` pode adicionar item enquanto `CompleteInventory` está executando, pois não há lock em `stock_inventories` ou `stock_inventory_items`.

### 2.2 Remover Item (DeleteInventoryItem)

**Arquivo:** `internal/infra/repository/gorm_stock_movement_repository.go`  
**Método:** `DeleteInventoryItem` (linhas 123-125)

```go
func (r *GormStockMovementRepository) DeleteInventoryItem(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.StockInventoryItem{}).Error  // ❌ FORA de transação
}
```

**Caminho do Código:**
1. DELETE em `stock_inventory_items` sem transação
2. Sem lock pessimista

**Possibilidade:** ✅ SIM - `DeleteInventoryItem` pode remover item enquanto `CompleteInventory` está executando.

### 2.3 Alterar Status (UpdateInventoryStatus)

**Arquivo:** `internal/infra/repository/gorm_stock_movement_repository.go`  
**Método:** `UpdateInventoryStatus` (linhas 94-98)

```go
func (r *GormStockMovementRepository) UpdateInventoryStatus(ctx context.Context, id uint, status string, tx *gorm.DB) error {
	return r.getDB(ctx, tx).Model(&domain.StockInventory{}).
		Where("id = ?", id).
		Update("status", status).Error  // ❌ SEM SELECT FOR UPDATE
}
```

**Caminho do Código:**
1. UPDATE em `stock_inventories` sem SELECT FOR UPDATE
2. Sem lock pessimista

**Possibilidade:** ✅ SIM - Outra transação pode chamar `UpdateInventoryStatus` (se tiver acesso ao repository) e alterar o status enquanto `CompleteInventory` está executando.

**Nota:** Não há um método público de serviço para alterar status arbitrariamente, mas o repository permite.

### 2.4 Alterar Quantidade de Item

**Análise:** Não existe método público para alterar quantidade de item de inventário após criação.

**Possibilidade:** ❌ NÃO - Não há método para essa operação.

### 2.5 Alterar Ingrediente de Item

**Análise:** Não existe método público para alterar ingrediente de item após criação.

**Possibilidade:** ❌ NÃO - Não há método para essa operação.

### Resumo da Parte 2

| Operação | Possível? | Motivo |
|----------|-----------|--------|
| Adicionar item | ✅ SIM | `AddInventoryItem` não usa transação e não há lock |
| Remover item | ✅ SIM | `DeleteInventoryItem` não usa transação e não há lock |
| Alterar quantidade | ❌ NÃO | Não existe método |
| Alterar ingrediente | ❌ NÃO | Não existe método |
| Alterar status | ✅ SIM | `UpdateInventoryStatus` não usa SELECT FOR UPDATE |

---

## Parte 3 — Existe algum cenário de race condition?

**Resposta:** SIM.

### 3.1 Cenário 1: Lost Update de Itens

**Sequência Completa:**

```
Tempo | CompleteInventory (T1) | AddInventoryItem (T2) | Estado
------|------------------------|------------------------|--------
T1    | GetInventoryByID (status=draft) | | inventory=draft
T2    | | GetInventoryByID (status=draft) | inventory=draft
T3    | ListInventoryItems (retorna 2 itens) | | items=[A,B]
T4    | | CreateInventoryItem (adiciona item C) | items=[A,B,C]
T5    | Processa item A (SELECT FOR UPDATE ingrediente) | | lock ingrediente A
T6    | Processa item B (SELECT FOR UPDATE ingrediente) | | lock ingrediente B
T7    | UpdateInventoryStatus (status=completed) | | inventory=completed
T8    | COMMIT | | transação T1 finalizada
T9    | | | item C nunca processado
```

**Linhas do Código:**

**T1 (CompleteInventory):**
- Linha 216: `GetInventoryByID(ctx, inventoryID, tx)` - lê status=draft
- Linha 227: `ListInventoryItems(ctx, inventoryID, tx)` - lê 2 itens
- Linha 239-296: Loop processando itens A e B
- Linha 299: `UpdateInventoryStatus(ctx, inventoryID, "completed", tx)` - atualiza status

**T2 (AddInventoryItem):**
- Linha 179: `GetInventoryByID(ctx, inventoryID, nil)` - lê status=draft
- Linha 185: Valida status == "draft" ✅
- Linha 202: `CreateInventoryItem(ctx, item, nil)` - adiciona item C

**Resultado:** Item C foi adicionado mas nunca foi processado no ajuste de estoque. Lost update de itens.

### 3.2 Cenário 2: Race Condition em Status

**Sequência Completa:**

```
Tempo | CompleteInventory (T1) | DeleteInventory (T2) | Estado
------|------------------------|----------------------|--------
T1    | GetInventoryByID (status=draft) | | inventory=draft
T2    | | DeleteInventory (soft delete) | inventory=deleted
T3    | ListInventoryItems (retorna itens) | | itens lidos
T4    | Processa itens (SELECT FOR UPDATE) | | locks em ingredientes
T5    | UpdateInventoryStatus (status=completed) | | ERRO: registro deletado
T6    | ROLLBACK | | transação revertida
```

**Linhas do Código:**

**T1 (CompleteInventory):**
- Linha 216: `GetInventoryByID(ctx, inventoryID, tx)` - lê inventory
- Linha 227: `ListInventoryItems(ctx, inventoryID, tx)` - lê itens
- Linha 299: `UpdateInventoryStatus(ctx, inventoryID, "completed", tx)` - tenta atualizar registro deletado

**T2 (DeleteInventory):**
- Linha 309: `DeleteInventory(ctx, id)` - soft delete
- Linha 101-102: `query.Delete(&domain.StockInventory{})` - marca deleted_at

**Resultado:** CompleteInventory falha no UPDATE pois o registro foi deletado durante o processamento.

### 3.3 Cenário 3: Double Completion

**Sequência Completa:**

```
Tempo | CompleteInventory A (T1) | CompleteInventory B (T2) | Estado
------|-------------------------|-------------------------|--------
T1    | GetInventoryByID (status=draft) | | inventory=draft
T2    | | GetInventoryByID (status=draft) | inventory=draft
T3    | ListInventoryItems (retorna itens) | | itens lidos
T4    | Processa itens (SELECT FOR UPDATE) | | locks em ingredientes
T5    | UpdateInventoryStatus (status=completed) | | inventory=completed
T6    | COMMIT | | transação T1 finalizada
T7    | | ListInventoryItems (retorna itens) | itens lidos
T8    | | Processa itens (SELECT FOR UPDATE) | locks em ingredientes
T9    | | UpdateInventoryStatus (status=completed) | inventory=completed (já estava)
T10   | | COMMIT | transação T2 finalizada
```

**Linhas do Código:**

**T1 (CompleteInventory):**
- Linha 216: `GetInventoryByID(ctx, inventoryID, tx)` - lê status=draft
- Linha 222: Valida status == "draft" ✅
- Linha 299: `UpdateInventoryStatus(ctx, inventoryID, "completed", tx)` - atualiza para completed
- Linha 303: COMMIT

**T2 (CompleteInventory):**
- Linha 216: `GetInventoryByID(ctx, inventoryID, tx)` - lê status=draft (snapshot de leitura)
- Linha 222: Valida status == "draft" ✅ (baseado no snapshot)
- Linha 299: `UpdateInventoryStatus(ctx, inventoryID, "completed", tx)` - atualiza para completed
- Linha 303: COMMIT

**Resultado:** Ambas as transações completam o mesmo inventário, processando os ajustes de estoque duas vezes. Double completion.

### 3.4 Reprodução do Cenário 1 (Lost Update de Itens)

**Passos para Reproduzir:**

1. Criar inventário em status "draft" com 2 itens
2. Iniciar CompleteInventory (T1)
3. Imediamente após T1 chamar ListInventoryItems, chamar AddInventoryItem (T2) para adicionar um terceiro item
4. T1 processa apenas os 2 itens originais
5. T1 atualiza status para "completed"
6. T2 adiciona o terceiro item com sucesso
7. Resultado: Terceiro item nunca foi processado no ajuste de estoque

**Código de Teste (Conceitual):**

```go
// T1
go func() {
    CompleteInventory(ctx, inventoryID, userID)
}()

// T2 (timing crítico)
time.Sleep(10 * time.Millisecond)
AddInventoryItem(ctx, inventoryID, ingredientID3, 10, 15, "razão")
```

---

## Parte 4 — A sugestão de validação de status após o loop é necessária?

**Resposta:** SIM, é necessária.

### Análise Técnica

**Sugestão Original da Auditoria:**
Validar novamente o status após o loop de processamento de itens.

**Por que é necessária:**

1. **Ausência de Lock Pessimista em `stock_inventories`:**
   - `GetInventoryByID` não usa SELECT FOR UPDATE
   - `UpdateInventoryStatus` não usa SELECT FOR UPDATE
   - Portanto, outra transação pode modificar o status durante o processamento

2. **Snapshot de Leitura vs Estado Atual:**
   - O status validado na linha 222 é um snapshot do início da transação
   - Entre a linha 222 e 299, o status pode ter sido alterado por outra transação
   - Sem lock pessimista, não há garantia de isolamento

3. **Cenário de Double Completion:**
   - Duas transações podem validar status == "draft" simultaneamente
   - Ambas processam os itens
   - Ambas atualizam status para "completed"
   - Resultado: Ajuste de estoque duplicado

4. **Defesa em Profundidade:**
   - A validação inicial (linha 222) é uma defesa
   - A validação final (após loop) é uma segunda defesa
   - Juntas, reduzem a janela de race condition

### Por que os Locks Pessimistas em Ingredientes Não Resolvem

**Análise:**

Os locks pessimistas em `FindIngredientByIDForUpdate` protegem APENAS a tabela `ingredients`, não a tabela `stock_inventories`.

```
CompleteInventory:
  ├─ GetInventoryByID (sem lock) ❌
  ├─ ListInventoryItems (sem lock) ❌
  ├─ FindIngredientByIDForUpdate (com lock) ✅
  ├─ UpdateIngredient (com lock ativo) ✅
  └─ UpdateInventoryStatus (sem lock) ❌
```

**Problema:** Outra transação pode modificar `stock_inventories` ou `stock_inventory_items` enquanto `CompleteInventory` está executando, pois não há lock nessas tabelas.

### Conclusão da Parte 4

**A validação de status após o loop é NECESSÁRIA** pois:

1. Não há lock pessimista em `stock_inventories`
2. O status pode ser alterado durante o processamento
3. A validação inicial é baseada em snapshot, não em estado atual
4. Os locks em ingredientes não protegem a tabela de inventários
5. É uma defesa em profundidade contra double completion

---

## Parte 5 — Conclusão Obrigatória

**Resposta:** A) O BUG #4 é real e precisa ser corrigido.

### Justificativa Técnica

**O BUG CRÍTICO #4 ainda existe após as Sprints 4B.2, 4B.3 e 4B.4.**

**Por que as Sprints Anteriores Não Corrigiram Este Bug:**

1. **Sprint 4B.2 (Propagação de Transação):**
   - Corrigiu propagação de `tx` em `GetInventoryByID` e `ListInventoryItems`
   - Isso garante que as leituras participem da transação
   - **Mas não adiciona lock pessimista em `stock_inventories`**
   - Portanto, não previne modificações concorrentes

2. **Sprint 4B.3 (Ordenação de Locks em CreateOrder):**
   - Corrigiu ordenação de locks em `CreateOrder`
   - Isso previne deadlock em pedidos
   - **Não afeta `CompleteInventory`**
   - Portanto, irrelevante para este bug

3. **Sprint 4B.4 (SELECT FOR UPDATE em UpdateIngredient):**
   - Corrigiu lost update em `UpdateIngredient`
   - Isso previne race conditions em atualizações de ingredientes
   - **Mas `CompleteInventory` não chama `UpdateIngredient` na tabela `stock_inventories`**
   - Portanto, não previne modificações concorrentes no inventário

### Resumo da Auditoria

| Aspecto | Status | Explicação |
|---------|--------|------------|
| Lock pessimista em `stock_inventories` | ❌ NÃO | `GetInventoryByID` não usa SELECT FOR UPDATE |
| Lock pessimista em `stock_inventory_items` | ❌ NÃO | `ListInventoryItems` não usa SELECT FOR UPDATE |
| Lock pessimista em `ingredients` | ✅ SIM | `FindIngredientByIDForUpdate` usa SELECT FOR UPDATE |
| Proteção contra adição de itens | ❌ NÃO | `AddInventoryItem` pode executar concorrentemente |
| Proteção contra remoção de itens | ❌ NÃO | `DeleteInventoryItem` pode executar concorrentemente |
| Proteção contra alteração de status | ❌ NÃO | `UpdateInventoryStatus` não usa SELECT FOR UPDATE |
| Race condition em itens | ✅ EXISTE | Cenário 1 demonstra lost update de itens |
| Race condition em status | ✅ EXISTE | Cenário 2 demonstra falha em registro deletado |
| Double completion | ✅ EXISTE | Cenário 3 demonstra completion duplo |
| Validação de status após loop | ❌ NÃO | Não implementada |

### Cenários de Race Condition Identificados

1. **Lost Update de Itens:** Item adicionado durante processamento não é processado
2. **Race Condition em Status:** Inventário deletado durante processamento causa falha
3. **Double Completion:** Duas transações completam o mesmo inventário

### Recomendação

**Implementar a correção do BUG CRÍTICO #4 conforme sugerido na auditoria original:**

Adicionar validação de status após o loop de processamento de itens em `CompleteInventory`.

**Alternativa mais robusta:**

Adicionar SELECT FOR UPDATE em `GetInventoryByID` para a tabela `stock_inventories`, garantindo lock pessimista no inventário durante todo o processamento.

---

## Conclusão Final

**O BUG CRÍTICO #4 é REAL e precisa ser corrigido.**

As correções das Sprints 4B.2, 4B.3 e 4B.4 não resolveram este problema específico, pois:
- 4B.2 tratou de propagação de transação (não de locks)
- 4B.3 tratou de ordenação de locks em CreateOrder (método diferente)
- 4B.4 tratou de SELECT FOR UPDATE em UpdateIngredient (tabela diferente)

O problema fundamental é a ausência de lock pessimista na tabela `stock_inventories`, que permite modificações concorrentes durante o processamento de `CompleteInventory`.
