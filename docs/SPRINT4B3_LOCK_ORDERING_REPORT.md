# Sprint 4B.3 — Relatório de Prevenção de Deadlock (Lock Ordering)

**Data:** 27 de Julho de 2026  
**Versão:** 1.0  
**Status:** ✅ Implementado e Compilando

---

## Resumo Executivo

Esta Sprint 4B.3 corrigiu o **BUG CRÍTICO #2** identificado na auditoria destrutiva da Sprint 4B.1 v2. O problema era que `CreateOrder` adquiria locks pessimistas (SELECT FOR UPDATE) dos ingredientes na ordem em que apareciam nos produtos do pedido, que não era determinística. Em alta concorrência, isso podia causar deadlock quando duas transações adquiriam locks em ordens diferentes.

A correção consistiu em:
- Coletar todos os ingredientes únicos do pedido antes de adquirir locks
- Eliminar duplicados acumulando o consumo total
- Ordenar os ingredientes por ID (ordem global determinística)
- Adquirir locks na ordem ordenada
- Persistir os itens do pedido após os locks já estarem adquiridos

---

## 1. Estratégia Implementada

### 1.1 Princípio Fundamental

**Teorema:** Se todas as transações adquirirem locks em uma ordem global determinística, deadlock é impossível.

**Prova:**
- Se T1 adquire locks na ordem L1 < L2 < L3 < ... < Ln
- E T2 também adquire locks na ordem L1 < L2 < L3 < ... < Ln
- Então T1 e T2 nunca podem esperar por locks que a outra transação já possui
- Pois se T1 possui Li, T2 só pode possuir Lj onde j < i
- T2 nunca espera por Li pois já possui Lj < Li
- Portanto, deadlock é impossível

**Q.E.D.**

### 1.2 Estratégia Específica

1. **Coleta de Ingredientes:** Iterar sobre todos os itens do pedido e coletar todos os ingredientes
2. **Eliminação de Duplicados:** Usar um map para acumular o consumo total de cada ingrediente
3. **Ordenação:** Converter o map para slice e ordenar por `IngredientID` crescente
4. **Aquisição de Locks:** Iterar sobre a lista ordenada e chamar `DecreaseIngredientStock` (que faz SELECT FOR UPDATE)
5. **Persistência:** Após todos os locks adquiridos, persistir os itens do pedido

---

## 2. Fluxograma da Nova Sequência de Locks

### 2.1 Fluxo Antes da Sprint 4B.3

```
CreateOrder (Repository)
  │
  ├─ db.Transaction(func(tx *gorm.DB) error {
  │   │
  │   ├─ Criar pedido (INSERT orders)
  │   │
  │   ├─ Para cada item do pedido:
  │   │   │
  │   │   ├─ Criar item (INSERT order_items)
  │   │   │
  │   │   └─ Para cada ingrediente do produto:
  │   │       │
  │   │       └─ DecreaseIngredientStock (SELECT FOR UPDATE)
  │   │           └─ Ordem: depende da ordem dos produtos ❌
  │   │
  │   └─ COMMIT
  │
  └─ ROLLBACK (se erro)
```

**Problema:** A ordem dos locks depende da ordem dos produtos no pedido, que não é determinística.

### 2.2 Fluxo Depois da Sprint 4B.3

```
CreateOrder (Repository)
  │
  ├─ db.Transaction(func(tx *gorm.DB) error {
  │   │
  │   ├─ Criar pedido (INSERT orders)
  │   │
  │   ├─ [Fase 1: Coleta de Ingredientes]
  │   │   │
  │   │   ├─ Para cada item do pedido:
  │   │   │   │
  │   │   │   └─ Para cada ingrediente do produto:
  │   │   │       │
  │   │   │       └─ Acumular consumo no map[ingredientID]
  │   │   │
  │   │   ├─ Converter map para slice
  │   │   │
  │   │   └─ Ordenar slice por IngredientID crescente ✅
  │   │
  │   ├─ [Fase 2: Aquisição de Locks]
  │   │   │
  │   │   └─ Para cada ingrediente na lista ordenada:
  │   │       │
  │   │       └─ DecreaseIngredientStock (SELECT FOR UPDATE)
  │   │           └─ Ordem: IngredientID crescente ✅
  │   │
  │   ├─ [Fase 3: Persistência de Itens]
  │   │   │
  │   │   └─ Para cada item do pedido:
  │   │       │
  │   │       └─ Criar item (INSERT order_items)
  │   │
  │   └─ COMMIT
  │
  └─ ROLLBACK (se erro)
```

**Melhoria:** A ordem dos locks é sempre por `IngredientID` crescente, independentemente da ordem dos produtos.

---

## 3. Código Implementado

### 3.1 Arquivo Modificado

**Arquivo:** `internal/infra/repository/gorm_order_repository.go`  
**Método:** `CreateOrder`  
**Linhas:** 104-180

### 3.2 Implementação

```go
// 2. Sprint 4B.3: Coletar todos os ingredientes únicos do pedido para ordenação de locks
// Isso previne deadlock garantindo ordem global determinística
type ingredientConsumption struct {
	ingredientID uint
	totalQty     float64
	name         string
	currentStock float64
}

ingredientMap := make(map[uint]*ingredientConsumption)

for i := range order.Items {
	item := &order.Items[i]
	ingredients, ok := productIngredients[item.ProductID]
	if !ok {
		return fmt.Errorf("CreateOrder: ingredientes não pré-carregados para produto_id=%d", item.ProductID)
	}
	
	for _, pi := range ingredients {
		consumo := pi.Quantity * item.Quantity
		if existing, found := ingredientMap[pi.IngredientID]; found {
			existing.totalQty += consumo  // Acumula consumo se duplicado
		} else {
			ingredientMap[pi.IngredientID] = &ingredientConsumption{
				ingredientID: pi.IngredientID,
				totalQty:     consumo,
				name:         pi.Ingredient.Name,
				currentStock: pi.Ingredient.StockQuantity,
			}
		}
	}
}

// Converter para slice e ordenar por IngredientID (ordem global determinística)
ingredientList := make([]*ingredientConsumption, 0, len(ingredientMap))
for _, ic := range ingredientMap {
	ingredientList = append(ingredientList, ic)
}
sort.Slice(ingredientList, func(i, j int) bool {
	return ingredientList[i].ingredientID < ingredientList[j].ingredientID
})

// 3. Sprint 4B.3: Adquirir locks em ordem determinística
// Isso garante que todas as transações adquiram locks na mesma ordem
for _, ic := range ingredientList {
	consumo := ic.totalQty
	if err := r.productRepo.DecreaseIngredientStock(ctx, ic.ingredientID, consumo, tx, ic.name, ic.currentStock); err != nil {
		return fmt.Errorf("CreateOrder: baixa estoque ingrediente_id=%d: %w", ic.ingredientID, err)
	}
}

// 4. Persistir os itens do pedido (após locks já adquiridos)
for i := range order.Items {
	item := &order.Items[i]
	item.OrderID = order.ID

	// 4a. Persiste o item com snapshot pré-carregado
	gItem := GormOrderItem{
		OrderID:               item.OrderID,
		ProductID:             item.ProductID,
		Quantity:              item.Quantity,
		UnitPrice:             item.UnitPrice,
		ProductName:           item.ProductName,
		ProductDescription:    item.ProductDescription,
		ProductIsComposto:     item.ProductIsComposto,
		ProductPhotoURL:       item.ProductPhotoURL,
		ProductCategoryID:     item.ProductCategoryID,
		ProductPromotionPrice: item.ProductPromotionPrice,
		ProductFeatured:       item.ProductFeatured,
		ProductIsNew:          item.ProductIsNew,
	}
	if err := tx.Create(&gItem).Error; err != nil {
		return fmt.Errorf("CreateOrder: criar item produto_id=%d: %w", item.ProductID, err)
	}
	item.ID = gItem.ID
}
```

---

## 4. Prova Matemática de Ausência de Deadlock

### 4.1 Teorema

**Teorema:** Com a implementação da Sprint 4B.3, deadlock em `CreateOrder` é impossível.

### 4.2 Prova

**Definições:**
- Seja `I = {i1, i2, ..., in}` o conjunto de ingredientes únicos de um pedido
- Seja `order(i)` a função que retorna o ID do ingrediente i
- Seja `L(T)` a sequência de locks adquiridos pela transação T

**Premissa 1 (Implementação):**
A Sprint 4B.3 garante que para qualquer transação T:
```
L(T) = [i1, i2, ..., in] onde order(i1) < order(i2) < ... < order(in)
```

**Premissa 2 (Propriedade de Ordenação):**
A ordem `<` é uma relação de ordem total em ℕ (números naturais).

**Premissa 3 (Propriedade de Deadlock):**
Deadlock ocorre quando existem duas transações T1 e T2 tais que:
- T1 espera por um lock que T2 possui
- T2 espera por um lock que T1 possui

**Prova por Contradição:**

Assuma que deadlock é possível após a Sprint 4B.3.

Então existem T1 e T2 em deadlock.

Pela definição de deadlock, existe um ingrediente x tal que:
- T1 espera por lock(x)
- T2 possui lock(x)

E existe um ingrediente y tal que:
- T2 espera por lock(y)
- T1 possui lock(y)

Pela Premissa 1, a ordem de locks de T1 é:
```
L(T1) = [..., y, ..., x, ...] onde order(y) < order(x)
```

Pela Premissa 1, a ordem de locks de T2 é:
```
L(T2) = [..., x, ..., y, ...] onde order(x) < order(y)
```

Mas isso implica:
```
order(y) < order(x) E order(x) < order(y)
```

Isso contradiz a Premissa 2 (irreflexividade da ordem estrita).

Portanto, a suposição de que deadlock é possível é falsa.

**Q.E.D.**

### 4.3 Exemplo Prático

**Cenário Antes da Sprint 4B.3:**

```
Pedido A:
- Produto X: ingredientes [3, 7, 12]
- Produto Y: ingredientes [5, 9]

Pedido B:
- Produto Z: ingredientes [12, 7, 3]
- Produto W: ingredientes [9, 5]

Ordem de locks A: [3, 7, 12, 5, 9]
Ordem de locks B: [12, 7, 3, 9, 5]

T1 lock 3 → T2 lock 12 → T1 espera 12 → T2 espera 3 → DEADLOCK ❌
```

**Cenário Depois da Sprint 4B.3:**

```
Pedido A:
- Produto X: ingredientes [3, 7, 12]
- Produto Y: ingredientes [5, 9]

Pedido B:
- Produto Z: ingredientes [12, 7, 3]
- Produto W: ingredientes [9, 5]

Ingredientes únicos A: {3, 5, 7, 9, 12}
Ingredientes únicos B: {3, 5, 7, 9, 12}

Ordem de locks A: [3, 5, 7, 9, 12] (ordenado por ID)
Ordem de locks B: [3, 5, 7, 9, 12] (ordenado por ID)

T1 lock 3 → T2 espera 3 → T1 lock 5 → T2 espera 5 → ... → SEM DEADLOCK ✅
```

---

## 5. Arquivos Modificados

| Arquivo | Linhas Modificadas | Alteração |
|---------|-------------------|-----------|
| `internal/infra/repository/gorm_order_repository.go` | 104-180 | Implementado coleta, eliminação de duplicados e ordenação de ingredientes |

---

## 6. Impacto em Performance

### 6.1 Análise de Complexidade

**Antes da Sprint 4B.3:**
- Coleta: O(n × m) onde n = itens, m = ingredientes/item
- Locks: O(n × m)
- Total: O(n × m)

**Depois da Sprint 4B.3:**
- Coleta: O(n × m)
- Ordenação: O(k log k) onde k = ingredientes únicos (k ≤ n × m)
- Locks: O(k)
- Persistência: O(n)
- Total: O(n × m + k log k + n)

**Diferença:** Adicionado O(k log k) para ordenação.

### 6.2 Análise Prática

**Cenário Típico:**
- Pedido com 5 itens
- Cada item com 3 ingredientes
- Ingredientes únicos: ~10 (alguns duplicados)

**Custo de Ordenação:**
- k = 10
- O(10 log 10) ≈ 33 operações

**Custo Total:**
- Antes: 5 × 3 = 15 operações
- Depois: 15 + 33 + 5 = 53 operações

**Overhead:** ~3.5x (mas em operações de CPU, não I/O)

### 6.3 Impacto em I/O

**Antes:**
- n × m chamadas a `DecreaseIngredientStock` (SELECT FOR UPDATE + UPDATE)

**Depois:**
- k chamadas a `DecreaseIngredientStock` (SELECT FOR UPDATE + UPDATE)
- k ≤ n × m (devido à eliminação de duplicados)

**Melhoria:** Menos chamadas ao banco quando há ingredientes duplicados.

### 6.4 Avaliação

**Overhead de CPU:** Aceitável (ordenação de pequenos conjuntos)  
**Overhead de I/O:** Melhoria (menos chamadas quando há duplicados)  
**Lock Duration:** Similar (mesmo número de locks, apenas em ordem diferente)

**Conclusão:** ✅ Impacto positivo ou neutro em performance.

---

## 7. Efeitos Colaterais

### 7.1 Mudança na Ordem de Persistência

**Antes:** Itens eram persistidos antes da baixa de estoque.  
**Depois:** Itens são persistidos após a baixa de estoque.

**Impacto:**
- Se a baixa de estoque falhar, os itens não são persistidos
- Isso é correto pois o pedido inteiro deve ser revertido
- A transação garante atomicidade

**Conclusão:** ✅ Sem impacto funcional (transação garante consistência).

### 7.2 Acumulação de Consumo

**Antes:** Cada ingrediente era baixado individualmente por item.  
**Depois:** Consumo total por ingrediente é acumulado antes da baixa.

**Impacto:**
- Se o mesmo ingrediente aparece em múltiplos itens, apenas uma baixa é feita
- O valor baixado é a soma de todos os consumos
- Isso é funcionalmente equivalente

**Conclusão:** ✅ Sem impacto funcional (equivalente matemático).

### 7.3 Mensagens de Erro

**Antes:** Erro indicava o produto_id.  
**Depois:** Erro indica o ingrediente_id.

**Impacto:**
- Mensagens de erro ligeiramente diferentes
- Mais específicas para debugging

**Conclusão:** ✅ Melhoria em debugging.

---

## 8. Critérios de Aceitação

| Critério | Status | Evidência |
|----------|--------|-----------|
| Ordem de locks é global, determinística e crescente por IngredientID | ✅ | `sort.Slice` por `ingredientID` |
| Duplicados de ingredientes são eliminados | ✅ | Map com acumulação de consumo |
| Impossível que duas transações adquiram locks em ordens diferentes | ✅ | Prova matemática na Seção 4.2 |
| go build ./... executa sem erros | ✅ | Build exit code 0 |
| Nenhuma alteração em regras de negócio | ✅ | Apenas reordenação de operações |
| Nenhuma alteração em interfaces públicas | ✅ | Assinatura de `CreateOrder` inalterada |
| Nenhuma alteração em domínio | ✅ | Sem mudanças em structs de domínio |

---

## 9. Comparação Antes vs Depois

### 9.1 Sequência de Operações

**Antes:**
```
1. Criar pedido
2. Para cada item:
   a. Criar item
   b. Para cada ingrediente:
      - Baixar estoque (lock na ordem do produto)
```

**Depois:**
```
1. Criar pedido
2. Coletar todos os ingredientes únicos
3. Ordenar por IngredientID
4. Para cada ingrediente ordenado:
   - Baixar estoque (lock na ordem do ID)
5. Para cada item:
   - Criar item
```

### 9.2 Ordem de Locks

**Antes:**
```
Pedido A: [3, 7, 12, 5, 9] (ordem dos produtos)
Pedido B: [12, 7, 3, 9, 5] (ordem dos produtos)
→ DEADLOCK POSSÍVEL ❌
```

**Depois:**
```
Pedido A: [3, 5, 7, 9, 12] (ordenado por ID)
Pedido B: [3, 5, 7, 9, 12] (ordenado por ID)
→ DEADLOCK IMPOSSÍVEL ✅
```

---

## 10. Integração com Outras Sprints

### 10.1 Sprint 4B.1 v2

- 4B.1 v2 implementou SELECT FOR UPDATE em `DecreaseIngredientStock`
- 4B.3 usa essa implementação para adquirir locks
- Integração perfeita: 4B.3 ordena os locks que 4B.1 v2 implementou

### 10.2 Sprint 4B.2

- 4B.2 corrigiu propagação de transação em `CompleteInventory`
- 4B.3 não afeta `CompleteInventory`
- Integração independente: sprints tratam métodos diferentes

### 10.3 Próximas Sprints

A auditoria identificou 2 bugs críticos adicionais:
- BUG CRÍTICO #3: UpdateIngredient SELECT FOR UPDATE
- BUG CRÍTICO #4: CompleteInventory validação de modificações

Estes bugs não são afetados pela Sprint 4B.3.

---

## 11. Conclusão

A Sprint 4B.3 corrigiu com sucesso o BUG CRÍTICO #2 da auditoria destrutiva. A ordenação de locks em `CreateOrder` agora é global, determinística e crescente por `IngredientID`, tornando deadlock impossível.

**Status:** ✅ **APROVADO** para este critério específico.

**Nota:** Os bugs críticos #3 e #4 da auditoria ainda precisam ser corrigidos antes de aprovar o módulo de estoque para produção.

---

## 12. Próximos Passos

### 12.1 BUG CRÍTICO #3: UpdateIngredient - SELECT Sem FOR UPDATE

**Problema:** `UpdateIngredient` faz SELECT sem FOR UPDATE, podendo causar lost update.

**Recomendação:** Adicionar `Clauses(clause.Locking{Strength: "UPDATE"})` no SELECT de `UpdateIngredient`.

### 12.2 BUG CRÍTICO #4: CompleteInventory - Validação de Modificações

**Problema:** `CompleteInventory` não valida se o inventário foi modificado durante o processamento.

**Recomendação:** Adicionar validação de status após o loop de processamento.
