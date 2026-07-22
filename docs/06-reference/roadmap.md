# BACKLOG OFICIAL PRATOONLINE

**Fase 11 — Hardening do Produto**  
**Data:** 16/07/2026  
**Arquiteto-Chefe:** PratoOnline

---

## Bugs Críticos

### BUG-001: Validação de formato de e-mail no login
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Input type="email" no HTML não valida formato, não há validação JavaScript adicional. Envia requisição mesmo com formato inválido.  
**Impacto:** Alto - Usuário pode tentar login com e-mail inválido  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-015: Arquivar produto sem verificar pedidos vinculados
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Função archiveProduct não verifica se há pedidos vinculados antes de arquivar.  
**Impacto:** Alto - Pode arquivar produto com pedidos ativos  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte`  
**Dependências:** Backend - endpoint para verificar pedidos vinculados  
**Pode esperar?** NÃO

### BUG-016: Excluir produto sem verificar dependências
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Função deleteProductById não verifica dependências (pedidos, fichas técnicas).  
**Impacto:** Alto - Pode excluir produto com dados vinculados  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte`  
**Dependências:** Backend - endpoint para verificar dependências  
**Pode esperar?** NÃO

### BUG-021: Excluir categoria sem verificar produtos vinculados
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Função deleteCategoryById não verifica produtos vinculados.  
**Impacto:** Alto - Pode excluir categoria com produtos  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte`  
**Dependências:** Backend - endpoint para verificar produtos vinculados  
**Pode esperar?** NÃO

### BUG-024: Excluir ingrediente sem verificar fichas técnicas
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Função deleteIngredientById não verifica fichas técnicas.  
**Impacto:** Alto - Pode excluir ingrediente usado em produtos  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Dependências:** Backend - endpoint para verificar fichas técnicas  
**Pode esperar?** NÃO

### BUG-030: Criar pedido não valida estoque suficiente
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Validação de estoque só no backend ao criar pedido, não ao adicionar ao carrinho.  
**Impacto:** Alto - Usuário pode adicionar produtos sem estoque  
**Esforço:** 6h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Backend - endpoint para validar estoque  
**Pode esperar?** NÃO

### BUG-032: Cancelar pedido sem confirmação
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Não há modal de confirmação antes de cancelar pedido.  
**Impacto:** Médio - Usuário pode cancelar por engano  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/[id]/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-044: Badges de notificações hardcoded
**Categoria:** BUG  
**Prioridade:** Crítica  
**Descrição:** Badges de pedidos (3) e ajustes (2) são estáticos no código.  
**Impacto:** Médio - Notificações não refletem realidade  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte`  
**Dependências:** Backend - endpoints de contagem  
**Pode esperar?** NÃO

---

## Bugs Altos

### BUG-003: Logout não limpa userStore completamente
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Função logout apenas chama API, não limpa store local.  
**Impacto:** Médio - Dados residuais podem permanecer  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`, `frontend/src/lib/stores/userStore.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-004: Sem tratamento de erro de rede
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** API client não trata erros de rede (fetch failed, timeout).  
**Impacto:** Alto - Erros genéricos sem contexto  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/lib/api/client.ts`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-006: KPIs do Dashboard com dados hardcoded
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** KPIs mostram "+12% vs ontem" e "+8% vs ontem" hardcoded.  
**Impacto:** Alto - Dados incorretos  
**Esforço:** 8h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Backend - endpoint de comparação  
**Pode esperar?** NÃO

### BUG-007: Atividades recentes são dados estáticos
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Array recentActivities é estático, não vem da API.  
**Impacto:** Alto - Dados falsos  
**Esforço:** 8h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Backend - endpoint de atividades  
**Pode esperar?** NÃO

### BUG-008: Cálculo incorreto de ingredientes críticos
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Lógica filtra produtos com ingredients.stock < 10 sem verificar IsComposto.  
**Impacto:** Médio - Mostra estoque baixo para produtos simples  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-011: Modal produto sem campos SEO
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Form de novo produto não inclui campos SEO apesar de backend suportar.  
**Impacto:** Médio - Campos não disponíveis para usuário  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-012: Modal edição sem campos SEO
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Form de edição não inclui campos SEO.  
**Impacto:** Médio - Campos não disponíveis para usuário  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-013: Modal sem campos iFood
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Form não inclui campos iFood apesar de backend suportar.  
**Impacto:** Médio - Campos não disponíveis para usuário  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-014: Duplicar não copia SEO/iFood
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Função duplicateProduct não inclui campos SEO e iFood no payload.  
**Impacto:** Baixo - Produto duplicado incompleto  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-017: Filtro categoria POS hardcoded
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Lógica de filtro usa strings hardcoded no nome do produto.  
**Impacto:** Médio - Filtragem incorreta  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-028: POS não permite selecionar mesa
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Form não inclui campo de mesa.  
**Impacto:** Médio - Não é possível associar pedido a mesa  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-029: POS não mostra estoque disponível
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Não há exibição de estoque de ingredientes.  
**Impacto:** Médio - Usuário não sabe estoque disponível  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-031: Tela editar pedido não existe
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Rota `/orders/[id]/edit` não existe.  
**Impacto:** Alto - Não é possível editar pedidos  
**Esforço:** 12h  
**Arquivos envolvidos:** N/A (criar arquivo)  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-033: Fechar pedido não existe
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Não há ação de fechar pedido.  
**Impacto:** Médio - Não é possível concluir pedido  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/[id]/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-035: Status pills com contagem hardcoded
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** countByStatus calculado sobre array filtrado, não total.  
**Impacto:** Médio - Contagem incorreta  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-036: Ajustes lista apenas pendentes
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** API getPendingAdjustments retorna apenas pendentes.  
**Impacto:** Médio - Não vê histórico completo  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Dependências:** Backend - endpoint getAllAdjustments  
**Pode esperar?** NÃO

### BUG-037: Ajustes mostram IDs em vez de nomes
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** API retorna apenas IDs, não nomes de pedido e ingrediente.  
**Impacto:** Médio - Difícil identificar ajustes  
**Esforço:** 4h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Dependências:** Backend - incluir nomes na resposta  
**Pode esperar?** NÃO

### BUG-040: Alterar e-mail sem confirmação senha
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Form não pede confirmação de senha ao alterar e-mail.  
**Impacto:** Médio - Risco de segurança  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-041: Alterar senha sem validação atual
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Validação básica implementada, pode ser melhorada.  
**Impacto:** Médio - Risco de segurança  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-045: Links /users e /settings não existem
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Links apontam para rotas não implementadas (404).  
**Impacto:** Baixo - Links quebrados  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-046: Header search não funciona
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Função handleSearch apenas faz console.log.  
**Impacto:** Médio - Busca não funcional  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Header.svelte`  
**Dependências:** Backend - endpoint de busca global  
**Pode esperar?** SIM

### BUG-049: Erro upload usa alert()
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Usa alert() nativo em vez de componente de UI.  
**Impacto:** Baixo - UX inconsistente  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-058: Sidebar some em mobile sem menu
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Sidebar tem display: none em mobile sem alternativa.  
**Impacto:** Alto - Navegação inacessível em mobile  
**Esforço:** 8h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-059: POS layout não responsivo
**Categoria:** BUG  
**Prioridade:** Alta  
**Descrição:** Carrinho some em mobile.  
**Impacto:** Médio - POS inutilizável em mobile  
**Esforço:** 6h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

---

## Bugs Médios

### BUG-005: Loading não bloqueia múltiplos cliques
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Loading state não desabilita botão adequadamente.  
**Impacto:** Baixo - Múltiplas requisições  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`, `frontend/src/routes/(auth)/register/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-009: Dashboard sem refresh automático
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há mecanismo de refresh automático.  
**Impacto:** Médio - Dados desatualizados  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-010: Erro loading sem retry em todos cards
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Tratamento de erro é global, não por seção.  
**Impacto:** Médio - Retry não disponível por seção  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-018: Validação promoção sem verificar datas
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Função validate não verifica datas de promoção.  
**Impacto:** Baixo - Salva datas inválidas  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-019: Validação disponibilidade sem verificar horários
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Função validate não verifica horários de disponibilidade.  
**Impacto:** Baixo - Salva horários inválidos  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-020: Link voltar usa href em vez de goto
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Usa <a href> em vez de navegação programática.  
**Impacto:** Baixo - Navegação não otimizada  
**Esforço:** 0.5h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/[id]/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-022: Ordem categoria não impede duplicatas
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há validação de DisplayOrder único.  
**Impacto:** Baixo - Pode haver duplicatas  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-023: Sem validação de nome único categoria
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há validação de nome único.  
**Impacto:** Baixo - Pode haver duplicatas  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-025: Ajuste estoque sem audit trail
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há implementação de histórico de ajustes.  
**Impacto:** Alto - Sem rastro de alterações  
**Esforço:** 8h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Dependências:** Backend - tabela de histórico  
**Pode esperar?** SIM

### BUG-026: Unidade sem validação formato
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Input de unidade é texto livre sem validação.  
**Impacto:** Baixo - Unidades inconsistentes  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-034: Sem cálculo tempo preparo estimado
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há cálculo de tempo de preparo no carrinho.  
**Impacto:** Baixo - Usuário não sabe tempo estimado  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-038: Ajustes sem busca
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Não há campo de busca.  
**Impacto:** Baixo - Difícil encontrar ajustes  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-039: Ajustes sem validação permissões
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Qualquer usuário pode aprovar/rejeitar.  
**Impacto:** Médio - Risco de segurança  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### BUG-047: Breadcrumb inconsistente
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Breadcrumb é manual, não automático.  
**Impacto:** Baixo - Pode ser inconsistente  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Workspace.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-048: Upload sem validação dimensões
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Valida apenas tamanho do arquivo (5MB), não dimensões.  
**Impacto:** Baixo - Pode aceitar imagens inadequadas  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-051: Drag drop não previne padrão
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Event handlers podem não prevenir comportamento padrão corretamente.  
**Impacto:** Baixo - Drag and drop pode falhar  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-056: Modais não fecham com ESC
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Modal component não implementa keydown ESC.  
**Impacto:** Baixo - UX não padrão  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Modal.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-060: Tabelas sem scroll horizontal mobile
**Categoria:** BUG  
**Prioridade:** Média  
**Descrição:** Tabelas não têm overflow-x.  
**Impacto:** Médio - Layout quebra em mobile  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Table.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

---

## Bugs Baixos

### BUG-027: Sem alerta visual estoque zero
**Categoria:** BUG  
**Prioridade:** Baixa  
**Descrição:** Badge "Estoque Zerado" pouco visível.  
**Impacto:** Baixo - Alerta sutil  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-061: Grid produtos mobile cards pequenos
**Categoria:** BUG  
**Prioridade:** Baixa  
**Descrição:** Grid minmax(160px) pode deixar cards muito pequenos.  
**Impacto:** Baixo - Cards pequenos em mobile  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### BUG-062: Filtros mobile ocupam espaço vertical
**Categoria:** BUG  
**Prioridade:** Baixa  
**Descrição:** Filtros sempre expandidos ocupam muito espaço.  
**Impacto:** Baixo - Layout vertical excessivo  
**Esforço:** 2h  
**Arquivos envolvidos:** Múltiplas telas com filtros  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

---

## UX

### UX-002: Remember me para persistir sessão
**Categoria:** UX  
**Prioridade:** Alta  
**Descrição:** Não há checkbox "remember me" nem implementação de persistência de token.  
**Impacto:** Médio - Usuário precisa login sempre  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`, `frontend/src/lib/stores/userStore.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-010: Retry por seção no Dashboard
**Categoria:** UX  
**Prioridade:** Média  
**Descrição:** Tratamento de erro é global, não por seção.  
**Impacto:** Médio - Retry não disponível por seção  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-042: Preferências vazias (placeholder)
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Placeholder "Configurações adicionais estarão disponíveis em breve".  
**Impacto:** Baixo - Funcionalidade não implementada  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-043: Sem avatar upload
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Não há upload de avatar no perfil.  
**Impacto:** Baixo - Personalização limitada  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-049: Erro upload usa alert()
**Categoria:** UX  
**Prioridade:** Alta  
**Descrição:** Usa alert() nativo em vez de componente de UI.  
**Impacto:** Baixo - UX inconsistente  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-052: Loading sem skeleton em alguns lugares
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Implementação inconsistente de loading states.  
**Impacto:** Baixo - Experiência inconsistente  
**Esforço:** 4h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-053: Estados vazios sem CTA consistente
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Alguns estados vazios sem call-to-action.  
**Impacto:** Baixo - UX inconsistente  
**Esforço:** 3h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-054: Erros sem retry em todos lugares
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Tratamento de erro inconsistente.  
**Impacto:** Baixo - UX inconsistente  
**Esforço:** 3h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-055: Sem loading global navegação
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Não há loading global durante transição.  
**Impacto:** Baixo - Pode parecer travado  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/+layout.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-057: Modais sem backdrop click
**Categoria:** UX  
**Prioridade:** Baixa  
**Descrição:** Modal component não implementa backdrop click.  
**Impacto:** Baixo - UX não padrão  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Modal.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### UX-059: POS layout não responsivo
**Categoria:** UX  
**Prioridade:** Alta  
**Descrição:** Carrinho some em mobile.  
**Impacto:** Médio - POS inutilizável em mobile  
**Esforço:** 6h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

---

## Features

### FEAT-002: Remember me para persistir sessão
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Implementar checkbox "Lembrar-me" e salvar token em localStorage.  
**Impacto:** Médio - Melhora conveniência  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`, `frontend/src/lib/stores/userStore.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-009: Dashboard refresh automático
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Implementar polling ou WebSocket para atualizar dados automaticamente.  
**Impacto:** Médio - Dados sempre atualizados  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-025: Histórico de ajustes de estoque
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Criar tabela de histórico/audit trail de ajustes.  
**Impacto:** Alto - Rastro completo de alterações  
**Esforço:** 8h (requer backend)  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Dependências:** Backend - tabela de histórico  
**Pode esperar?** SIM

### FEAT-028: Seleção de mesa no POS
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Adicionar campo "Número da Mesa" no formulário.  
**Impacto:** Médio - Associação pedido-mesa  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### FEAT-029: Exibir estoque disponível no POS
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Mostrar estoque disponível ao passar mouse no produto.  
**Impacto:** Médio - Usuário vê estoque  
**Esforço:** 3h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-031: Tela editar pedido
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Criar tela para editar pedidos existentes.  
**Impacto:** Alto - Permite correções  
**Esforço:** 12h  
**Arquivos envolvidos:** Criar `frontend/src/routes/(app)/orders/[id]/edit/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### FEAT-033: Fechar pedido
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Adicionar botão/ação para fechar pedido.  
**Impacto:** Médio - Fluxo completo  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/[id]/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

### FEAT-034: Cálculo tempo preparo estimado
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Somar PreparationTimeMinutes dos produtos no carrinho.  
**Impacto:** Baixo - Usuário sabe tempo  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-038: Busca em ajustes de estoque
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Adicionar campo de busca por ID ou nome.  
**Impacto:** Baixo - Facilita encontrar ajustes  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-042: Preferências de usuário
**Categoria:** FEATURE  
**Prioridade:** Baixa  
**Descrição:** Implementar configurações de tema, idioma, notificações.  
**Impacto:** Baixo - Personalização  
**Esforço:** 8h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-043: Upload de avatar
**Categoria:** FEATURE  
**Prioridade:** Baixa  
**Descrição:** Adicionar PhotoUpload para avatar do usuário.  
**Impacto:** Baixo - Personalização  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-046: Busca global
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Implementar busca global em todo o sistema.  
**Impacto:** Médio - Facilita navegação  
**Esforço:** 8h (requer backend)  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Header.svelte`  
**Dependências:** Backend - endpoint de busca global  
**Pode esperar?** SIM

### FEAT-048: Validação de dimensões de upload
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Validar dimensões (mínimo 200x200, máximo 4000x4000).  
**Impacto:** Baixo - Qualidade de imagens  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-050: Progresso de upload
**Categoria:** FEATURE  
**Prioridade:** Média  
**Descrição:** Implementar barra de progresso de upload.  
**Impacto:** Baixo - Feedback visual  
**Esforço:** 6h (requer backend)  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`, `frontend/src/lib/api/media.ts`  
**Dependências:** Backend - suporte a progresso  
**Pode esperar?** SIM

### FEAT-055: Loading global de navegação
**Categoria:** FEATURE  
**Prioridade:** Baixa  
**Descrição:** Adicionar loading global baseado em navigating do SvelteKit.  
**Impacto:** Baixo - Feedback visual  
**Esforço:** 2h  
**Arquivos envolvidos:** `frontend/src/routes/+layout.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### FEAT-058: Menu hamburguer mobile
**Categoria:** FEATURE  
**Prioridade:** Alta  
**Descrição:** Implementar menu hamburguer/drawer para mobile.  
**Impacto:** Alto - Navegação acessível em mobile  
**Esforço:** 8h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

---

## Tech Debt

### DEBT-005: Padronizar loading states
**Categoria:** TECH DEBT  
**Prioridade:** Baixa  
**Descrição:** Garantir que disabled={loading} esteja em todos os botões de submit.  
**Impacto:** Baixo - Consistência  
**Esforço:** 2h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### DEBT-047: Breadcrumb automático
**Categoria:** TECH DEBT  
**Prioridade:** Média  
**Descrição:** Tornar breadcrumb automático baseado em rota.  
**Impacto:** Médio - Manutenibilidade  
**Esforço:** 4h  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Workspace.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### DEBT-052: Padronizar Skeleton loading
**Categoria:** TECH DEBT  
**Prioridade:** Baixa  
**Descrição:** Padronizar uso de Skeleton component em todos os lugares.  
**Impacto:** Baixo - Consistência  
**Esforço:** 4h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### DEBT-053: Padronizar EmptyState
**Categoria:** TECH DEBT  
**Prioridade:** Baixa  
**Descrição:** Padronizar uso de EmptyState component.  
**Impacto:** Baixo - Consistência  
**Esforço:** 3h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### DEBT-054: Padronizar tratamento de erro
**Categoria:** TECH DEBT  
**Prioridade:** Baixa  
**Descrição:** Padronizar tratamento de erro com retry.  
**Impacto:** Baixo - Consistência  
**Esforço:** 3h  
**Arquivos envolvidos:** Múltiplos arquivos  
**Dependências:** Nenhuma  
**Pode esperar?** SIM

### DEBT-060: Responsividade de tabelas
**Categoria:** TECH DEBT  
**Prioridade:** Média  
**Descrição:** Adicionar overflow-x: auto em tabelas.  
**Impacto:** Médio - Responsividade  
**Esforço:** 1h  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Table.svelte`  
**Dependências:** Nenhuma  
**Pode esperar?** NÃO

---

## RESUMO EXECUTIVO

### Quantitativo por Categoria

| Categoria | Quantidade | % do Total |
|-----------|-----------|------------|
| Bugs Críticos | 8 | 12.9% |
| Bugs Altos | 24 | 38.7% |
| Bugs Médios | 17 | 27.4% |
| Bugs Baixos | 3 | 4.8% |
| UX | 11 | 17.7% |
| Features | 14 | 22.6% |
| Tech Debt | 6 | 9.7% |
| **TOTAL** | **62** | **100%** |

### Resumo Consolidado

- **Total de Bugs:** 52 (83.9%)
  - Críticos: 8
  - Altos: 24
  - Médios: 17
  - Baixos: 3
- **Total de UX:** 11 (17.7%)
- **Total de Features:** 14 (22.6%)
- **Total de Tech Debt:** 6 (9.7%)

### Estimativa de Esforço

| Categoria | Horas | Dias (8h) |
|-----------|-------|-----------|
| Bugs Críticos | 30h | 3.75 dias |
| Bugs Altos | 68h | 8.5 dias |
| Bugs Médios | 35.5h | 4.44 dias |
| Bugs Baixos | 4h | 0.5 dias |
| UX | 29h | 3.63 dias |
| Features | 61h | 7.63 dias |
| Tech Debt | 17h | 2.13 dias |
| **TOTAL** | **244.5h** | **30.56 dias** |

### Risco do Sistema

**Nível de Risco:** ALTO

**Justificativa:**
- 8 bugs críticos que podem causar perda de dados ou bloquear fluxos essenciais
- 24 bugs altos que impactam significativamente a experiência
- 12 bugs requerem alterações no backend (dependência externa)
- Responsividade em mobile comprometida (2 bugs críticos)
- Validações de segurança ausentes (alterar e-mail/senha sem confirmação adequada)
- Exclusão de dados sem verificação de dependências (4 bugs críticos)

### Índice de Maturidade

**Índice:** 3.2 / 10

**Cálculo:**
- Funcionalidade básica implementada: 4/10
- Validações e segurança: 2/10
- UX e consistência: 3/10
- Responsividade: 2/10
- Código limpo e manutenível: 5/10

**Interpretação:** Sistema em fase inicial de desenvolvimento. Funcionalidades básicas presentes, mas com muitas lacunas em validações, segurança e UX.

### Bugs que Impedem Produção

**Quantidade:** 8 bugs críticos

**Lista:**
1. BUG-001: Validação de formato de e-mail no login
2. BUG-015: Arquivar produto sem verificar pedidos vinculados
3. BUG-016: Excluir produto sem verificar dependências
4. BUG-021: Excluir categoria sem verificar produtos vinculados
5. BUG-024: Excluir ingrediente sem verificar fichas técnicas
6. BUG-030: Criar pedido não valida estoque suficiente
7. BUG-032: Cancelar pedido sem confirmação
8. BUG-044: Badges de notificações hardcoded

**Justificativa:** Estes bugs podem causar perda de dados, corrupção de integridade referencial ou bloquear fluxos essenciais do negócio.

### Recomendação

**Fase 1 - Estabilização Crítica (1 semana):**
- Corrigir 8 bugs críticos
- Foco: segurança e integridade de dados

**Fase 2 - Estabilização Alta (2 semanas):**
- Corrigir 24 bugs altos
- Foco: UX e funcionalidades essenciais

**Fase 3 - Polishing (2 semanas):**
- Corrigir bugs médios/baixos
- Implementar features prioritárias
- Reduzir tech debt

**Total estimado:** 5 semanas para sistema em nível de produção aceitável.
