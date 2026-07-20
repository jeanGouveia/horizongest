# SPRINT 4 - Ficha Técnica

**Data:** 2025-01-XX  
**Implementador:** Cascade AI  
**Escopo:** Completar ficha técnica com ingredientes, quantidades, perdas, rendimento, custo  
**Objetivo:** Transformar ficha técnica básica em sistema inteligente de cálculo de custos

---

## Resumo Executivo

Ficha técnica expandida com cálculo automático de custos, CMV, margem e lucro. Sistema calcula custo total do produto baseado em ingredientes, considerando perdas e rendimento. Validação automática de fichas técnicas com alertas para produtos sem ficha ou ingredientes inativos.

**Status:** ✅ **IMPLEMENTADO**

---

## 1. Funcionalidades Implementadas

### 1.1 Campos de Ficha Técnica
- ✅ Perdas (loss) - % de desperdício
- ✅ Rendimento (yield) - % de rendimento
- ✅ Custo unitário (unit_cost) - custo do ingrediente
- ✅ Custo total (total_cost) - custo no produto
- ✅ Custo do produto (cost) - soma dos ingredientes
- ✅ CMV (cost of goods sold) - custo / preço
- ✅ Margem (margin) - (preço - custo) / preço
- ✅ Lucro (profit) - preço - custo
- ✅ Preço sugerido (suggested_price) - baseado em custo e margem desejada

### 1.2 Métodos de Cálculo
- ✅ `CalculateCost()` - Calcula custo total do ingrediente no produto
- ✅ `GetEffectiveQuantity()` - Calcula quantidade efetiva considerando perdas
- ✅ `CalculateCost()` - Calcula custo total do produto
- ✅ `CalculateCMV()` - Calcula custo merca da venda
- ✅ `CalculateMargin()` - Calcula margem de lucro
- ✅ `CalculateProfit()` - Calcula lucro por unidade
- ✅ `CalculateSuggestedPrice()` - Calcula preço sugerido

### 1.3 Validação
- ✅ `HasRecipe()` - Verifica se produto tem ficha técnica
- ✅ `ValidateRecipe()` - Valida ficha técnica do produto
- ✅ Alerta quando produto sem ficha
- ✅ Alerta quando ingrediente não existe
- ✅ Alerta quando ingrediente inativo

---

## 2. Arquivos Modificados

### 2.1 Domain
- `internal/domain/product_ingredient.go` - Expandido com campos de ficha técnica avançada

**Mudanças:**
- Adicionado `Loss float64` - perda em %
- Adicionado `Yield float64` - rendimento em %
- Adicionado `UnitCost float64` - custo unitário
- Adicionado `TotalCost float64` - custo total
- Adicionado método `CalculateCost()` - calcula custo total
- Adicionado método `GetEffectiveQuantity()` - calcula quantidade efetiva

- `internal/domain/product.go` - Expandido com campos de custo e métodos de cálculo

**Mudanças:**
- Adicionado `Cost float64` - custo total do produto
- Adicionado `CMV float64` - custo merca da venda
- Adicionado `Margin float64` - margem de lucro
- Adicionado `Profit float64` - lucro por unidade
- Adicionado `SuggestedPrice float64` - preço sugerido
- Adicionado método `CalculateCost()` - calcula custo do produto
- Adicionado método `CalculateCMV()` - calcula CMV
- Adicionado método `CalculateMargin()` - calcula margem
- Adicionado método `CalculateProfit()` - calcula lucro
- Adicionado método `CalculateSuggestedPrice()` - calcula preço sugerido
- Adicionado método `HasRecipe()` - verifica se tem ficha técnica
- Adicionado método `ValidateRecipe()` - valida ficha técnica
- Adicionado import `errors` para validação

### 2.2 Repository
- `internal/infra/repository/gorm_product_repository.go` - Expandido com novos campos

**Mudanças:**
- Adicionado campos ao `GormProductIngredient`: Loss, Yield, UnitCost, TotalCost
- Adicionado campos ao `GormProduct`: Cost, CMV, Margin, Profit, SuggestedPrice
- Atualizado `SetProductIngredients()` para incluir novos campos
- Atualizado `GetProductIngredients()` para retornar novos campos

### 2.3 Migration
- `migrations/00020_add_recipe_fields.sql` - Adição de campos nas tabelas

**Mudanças:**
- Adicionado `loss`, `yield`, `unit_cost`, `total_cost` à tabela `product_ingredients`
- Adicionado `cost`, `cmv`, `margin`, `profit`, `suggested_price` à tabela `products`

---

## 3. Fórmulas Implementadas

### 3.1 Custo do Ingrediente no Produto
```
Custo = (Quantidade × CustoUnitário) / Rendimento
```

### 3.2 Quantidade Efetiva
```
QuantidadeEfetiva = Quantidade / (1 - Perda)
```

### 3.3 Custo do Produto
```
Custo = Σ(Custo de cada ingrediente)
```

### 3.4 CMV (Custo Merca da Venda)
```
CMV = Custo / Preço
```

### 3.5 Margem de Lucro
```
Margem = (Preço - Custo) / Preço
```

### 3.6 Lucro por Unidade
```
Lucro = Preço - Custo
```

### 3.7 Preço Sugerido
```
PreçoSugerido = Custo / (1 - MargemDesejada)
```

---

## 4. Integrações Futuras

### 4.1 Estoque
- Integrar custo unitário com preço de compra do ingrediente
- Atualizar custo automaticamente ao atualizar preço de compra

### 4.2 Dashboard
- Calcular CMV real baseado em fichas técnicas
- Substituir CMV simplificado (30% fixo) por CMV real

### 4.3 Service Layer
- Implementar cálculo automático ao salvar ficha técnica
- Implementar recálculo ao atualizar preço de ingrediente
- Implementar recálculo ao atualizar quantidade de ingrediente

### 4.4 Handler
- Aceitar novos campos na criação/edição de ficha técnica
- Retornar campos calculados nas respostas

---

## 5. Limitações

### 5.1 Cálculo Automático
- Cálculos são manuais (métodos disponíveis mas não chamados automaticamente)
- Necessário implementar cálculo automático no service layer
- Necessário implementar recálculo ao atualizar ingredientes

### 5.2 Custo Unitário
- Custo unitário é manual (não integrado com compras)
- Necessário integrar com módulo de compras (ETAPA 5)
- Necessário atualizar custo automaticamente ao comprar ingrediente

### 5.3 Validação Automática
- Validação de ficha técnica disponível mas não aplicada automaticamente
- Necessário implementar validação ao criar/editar produto
- Necessário implementar alertas no dashboard

---

## 6. Testes

### 6.1 Testes Manuais Requeridos
- [ ] Criar ficha técnica com perdas e rendimento
- [ ] Calcular custo total do produto
- [ ] Calcular CMV do produto
- [ ] Calcular margem de lucro
- [ ] Calcular preço sugerido
- [ ] Validar ficha técnica de produto
- [ ] Verificar alerta de produto sem ficha técnica
- [ ] Verificar alerta de ingrediente inativo

### 6.2 Testes de Integração
- [ ] Testar integração com estoque (quando implementado)
- [ ] Testar integração com compras (quando implementado)
- [ ] Testar integração com dashboard (quando implementado)

---

## 7. Próximos Passos

1. **Service Layer:** Implementar cálculo automático ao salvar ficha técnica
2. **Handler:** Aceitar novos campos na criação/edição
3. **Compras (ETAPA 5):** Integrar custo unitário com preço de compra
4. **Dashboard:** Substituir CMV simplificado por CMV real
5. **Validação:** Implementar validação automática ao criar/editar produto

---

## 8. Assinatura

**Implementador:** Cascade AI  
**Data:** 2025-01-XX  
**Versão:** 1.0  
**Status:** ✅ IMPLEMENTADO (Backend - Domain e Repository)
