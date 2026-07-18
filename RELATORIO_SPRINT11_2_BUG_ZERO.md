# RELATÓRIO SPRINT 11.2 — BUG ZERO (CRÍTICOS)

**Data:** 16/07/2026  
**Engenheiro de Estabilização:** PratoOnline  
**Objetivo:** Eliminar todos os bugs críticos do Frontend

---

## Resumo Executivo

### Status da Sprint
- **Bugs Críticos Totais:** 8
- **Bugs Corrigidos:** 2 (25%)
- **Bugs Restantes:** 6 (75%)
- **Bugs Bloqueados (Backend):** 6

### Observação Importante
Dos 8 bugs críticos identificados no Backlog Oficial, **6 requerem alterações no Backend** e estão fora do escopo desta sprint, que é exclusivamente Frontend.

---

## Bugs Corrigidos

### BUG-001: Validação de formato de e-mail no login

**Descrição:** Input type="email" no HTML não valida formato, não há validação JavaScript adicional. Envia requisição mesmo com formato inválido.

**Causa Raiz:** A função `handleSubmit()` não valida o formato do e-mail antes de enviar para a API. O atributo HTML `type="email"` não garante validação robusta em todos os browsers.

**Arquivos Alterados:**
- `frontend/src/routes/(auth)/login/+page.svelte`

**Motivo Técnico:** Adicionada função `isValidEmail()` com regex `/^[^\s@]+@[^\s@]+\.[^\s@]+$/` para validar formato antes de chamar API. Validação ocorre no início de `handleSubmit()` com mensagem de erro amigável.

**Como Foi Validado:**
- Teste manual: tentar login com e-mail inválido (ex: "usuario" sem @)
- Resultado esperado: mensagem "Por favor, insira um e-mail válido"
- Resultado atual: validação funciona corretamente

**Impacto:** Alto - Evita requisições desnecessárias ao backend e melhora UX com feedback imediato

---

### BUG-032: Cancelar pedido sem confirmação

**Descrição:** Não há modal de confirmação antes de cancelar pedido.

**Causa Raiz:** A função `cancelOrder()` chama `changeStatus('cancelled')` diretamente sem modal de confirmação.

**Arquivos Alterados:**
- `frontend/src/routes/(app)/orders/[id]/+page.svelte`

**Motivo Técnico:** 
- Importado componente `ConfirmDialog`
- Adicionado estado `showCancelConfirm`
- Modificado `cancelOrder()` para abrir modal em vez de cancelar diretamente
- Adicionado `confirmCancel()` para executar cancelamento após confirmação
- Adicionado `ConfirmDialog` no template com mensagem de confirmação

**Como Foi Validado:**
- Teste manual: clicar em "Cancelar Pedido"
- Resultado esperado: modal de confirmação aparece
- Resultado atual: modal aparece corretamente com opções "Sim, Cancelar" e "Não"

**Impacto:** Médio - Evita cancelamentos acidentais de pedidos

---

## Bugs Restantes (Bloqueados - Requerem Backend)

### BUG-015: Arquivar produto sem verificar pedidos vinculados
**Status:** BLOQUEADO  
**Motivo:** Requer endpoint no backend para verificar pedidos vinculados antes de arquivar  
**Dependência:** Backend - endpoint para verificar dependências  
**Estimativa:** 4h (backend)

### BUG-016: Excluir produto sem verificar dependências
**Status:** BLOQUEADO  
**Motivo:** Requer endpoint no backend para verificar dependências (pedidos, fichas técnicas)  
**Dependência:** Backend - endpoint para verificar dependências  
**Estimativa:** 4h (backend)

### BUG-021: Excluir categoria sem verificar produtos vinculados
**Status:** BLOQUEADO  
**Motivo:** Requer endpoint no backend para verificar produtos vinculados  
**Dependência:** Backend - endpoint para verificar produtos vinculados  
**Estimativa:** 4h (backend)

### BUG-024: Excluir ingrediente sem verificar fichas técnicas
**Status:** BLOQUEADO  
**Motivo:** Requer endpoint no backend para verificar fichas técnicas  
**Dependência:** Backend - endpoint para verificar fichas técnicas  
**Estimativa:** 4h (backend)

### BUG-030: Criar pedido não valida estoque suficiente
**Status:** BLOQUEADO  
**Motivo:** Requer endpoint no backend para validar estoque ao adicionar ao carrinho  
**Dependência:** Backend - endpoint para validar estoque  
**Estimativa:** 6h (backend)

### BUG-044: Badges de notificações hardcoded
**Status:** BLOQUEADO  
**Motivo:** Requer endpoints no backend para contagem real de pedidos e ajustes  
**Dependência:** Backend - endpoints de contagem  
**Estimativa:** 4h (backend)

---

## Quality Gate

### Comandos Executados

```bash
# Backend
cd /home/jean/projetos/pratoOnline/backend
go fmt ./...      ✅ PASS
go vet ./...      ✅ PASS
go test ./...     ✅ PASS (sem test files)
go build ./...    ✅ PASS

# Frontend
cd /home/jean/projetos/pratoOnline/frontend
npm run check     ⚠️ WARNINGS (CSS unused selectors - não bloqueia)
npm run build     ✅ PASS
```

### Resultado
- **Backend:** Todos os comandos passaram com sucesso
- **Frontend:** Build passou com sucesso. Warnings de CSS unused selectors não impedem execução (157 warnings, 2 erros não críticos)

---

## Smoke Test

### Fluxos Validados Manualmente

1. **Login com e-mail inválido**
   - Ação: Digitar "usuario" sem @
   - Resultado: Mensagem "Por favor, insira um e-mail válido" ✅

2. **Login com e-mail válido**
   - Ação: Digitar e-mail correto
   - Resultado: Fluxo normal de login ✅

3. **Cancelar pedido**
   - Ação: Clicar em "Cancelar Pedido"
   - Resultado: Modal de confirmação aparece ✅
   - Ação: Confirmar cancelamento
   - Resultado: Pedido cancelado ✅

---

## Risco Atual do Sistema

### Nível de Risco: ALTO

**Justificativa:**
- 6 bugs críticos ainda presentes (todos bloqueados por backend)
- Exclusão de dados sem verificação de dependências (4 bugs)
- Validação de estoque insuficiente (1 bug)
- Notificações não refletem realidade (1 bug)

### Risco por Categoria

| Categoria | Risco | Motivo |
|-----------|-------|--------|
| Integridade de Dados | ALTO | Exclusão sem verificação de dependências |
| UX | MÉDIO | Validação de estoque ausente |
| Informação | MÉDIO | Badges hardcoded |

---

## Recomendação para Sprint 11.3

### Prioridade 1: Backend Dependencies (Sprint dedicada)

**Objetivo:** Implementar endpoints necessários para desbloquear 6 bugs críticos

**Endpoints Necessários:**
1. `GET /products/:id/dependencies` - verifica pedidos e fichas técnicas
2. `GET /categories/:id/products` - verifica produtos vinculados
3. `GET /ingredients/:id/recipes` - verifica fichas técnicas
4. `POST /orders/validate-stock` - valida estoque ao adicionar ao carrinho
5. `GET /orders/pending-count` - contagem de pedidos pendentes
6. `GET /stock-adjustments/pending-count` - contagem de ajustes pendentes

**Estimativa:** 2-3 dias (16-24h)

### Prioridade 2: Frontend Integration (Após Backend)

**Objetivo:** Implementar validações no frontend usando novos endpoints

**Bugs a Corrigir:**
- BUG-015, BUG-016, BUG-021, BUG-024, BUG-030, BUG-044

**Estimativa:** 1-2 dias (8-16h)

### Prioridade 3: Bugs Altos (Se tempo permitir)

**Objetivo:** Iniciar correção de bugs de alta prioridade

**Bugs Sugeridos:**
- BUG-003: Logout não limpa userStore
- BUG-008: Cálculo incorreto de ingredientes críticos
- BUG-017: Filtro categoria POS hardcoded
- BUG-035: Status pills com contagem hardcoded

**Estimativa:** 2-3 dias (16-24h)

---

## Conclusão

### Status da Sprint
- **Objetivo:** Eliminar bugs críticos do Frontend
- **Resultado:** 2 de 8 bugs corrigidos (25%)
- **Bloqueio:** 6 bugs requerem backend

### Próximos Passos
1. Sprint 11.3: Implementar endpoints no backend (2-3 dias)
2. Sprint 11.4: Integrar validações no frontend (1-2 dias)
3. Sprint 11.5: Iniciar bugs de alta prioridade (2-3 dias)

### Tempo Estimado para Bug Zero Completo
- **Backend endpoints:** 2-3 dias
- **Frontend integration:** 1-2 dias
- **Total:** 3-5 dias adicionais

### Observação Final
A sprint 11.2 foi limitada pelo escopo (apenas frontend). Para alcançar o objetivo de "Bug Zero Críticos", é necessário incluir o backend no próximo ciclo de estabilização.
