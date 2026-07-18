# RELATÓRIO SPRINT 10E — AUDITORIA FUNCIONAL FRONTEND

**Data:** 16/07/2026  
**Auditor:** Arquiteto Executor PratoOnline  
**Escopo:** Frontend completo (SvelteKit 5 + TypeScript)  
**Objetivo:** Identificar todos os bugs funcionais antes do próximo épico

---

## RESUMO EXECUTIVO

### Visão Geral
Auditoria funcional completa do frontend PratoOnline, cobrindo todas as telas e fluxos de usuário. Foram identificados **62 bugs** distribuídos entre **Críticos (8)**, **Alta (24)**, **Média (20)** e **Baixa (10)**.

### Impacto no MVP
- **Bloqueadores:** 8 bugs críticos impedem fluxos essenciais
- **Degradação de UX:** 30 bugs afetam experiência do usuário
- **Inconsistências:** 24 bugs causam comportamentos inesperados
- **Melhorias:** 10 bugs são melhorias de UX/polishing

### Estimativa de Esforço
- **Correção Críticos:** 2-3 dias (8 bugs)
- **Correção Alta:** 3-4 dias (24 bugs)
- **Correção Média:** 2-3 dias (20 bugs)
- **Correção Baixa:** 1-2 dias (10 bugs)
- **Total estimado:** 8-12 dias (1.5 a 2 sprints)

### Recomendação
Priorizar correção dos **8 bugs críticos** antes de iniciar novo épico. Os bugs de alta prioridade devem ser corrigidos na próxima sprint de estabilização.

---

## BUGS POR SEVERIDADE

### Crítica (8) - Bloqueadores de fluxo essencial
1. BUG 01: Login sem validação de formato de e-mail
2. BUG 15: Arquivar produto sem verificar pedidos vinculados
3. BUG 16: Excluir produto sem verificar dependências
4. BUG 21: Excluir categoria sem verificar produtos vinculados
5. BUG 24: Excluir ingrediente sem verificar fichas técnicas
6. BUG 30: Criar pedido não valida estoque suficiente
7. BUG 32: Cancelar pedido sem confirmação
8. BUG 44: Sidebar badges hardcoded impedem notificações reais

### Alta (24) - Impacta significativamente UX
9. BUG 02: Não há "remember me" para persistir sessão
10. BUG 03: Logout não limpa userStore completamente
11. BUG 06: KPIs do Dashboard com dados hardcoded
12. BUG 07: Atividades recentes são dados estáticos
13. BUG 08: Cálculo incorreto de ingredientes críticos
14. BUG 11: Modal produto sem campos SEO
15. BUG 12: Modal edição sem campos SEO
16. BUG 13: Modal sem campos iFood
17. BUG 14: Duplicar não copia SEO/iFood
18. BUG 17: Filtro categoria POS hardcoded
19. BUG 28: POS não permite selecionar mesa
20. BUG 29: POS não mostra estoque disponível
21. BUG 31: Tela editar pedido não existe
22. BUG 33: Fechar pedido não existe
23. BUG 35: Status pills com contagem hardcoded
24. BUG 36: Ajustes lista apenas pendentes
25. BUG 37: Ajustes mostram IDs em vez de nomes
26. BUG 40: Alterar e-mail sem confirmação senha
27. BUG 45: Links /users e /settings não existem
28. BUG 46: Header search não funciona
29. BUG 49: Erro upload usa alert()
30. BUG 58: Sidebar some em mobile sem menu
31. BUG 59: POS layout não responsivo

### Média (20) - Comportamentos inesperados
32. BUG 04: Sem tratamento de erro de rede
33. BUG 05: Loading não bloqueia múltiplos cliques
34. BUG 09: Dashboard sem refresh automático
35. BUG 10: Erro loading sem retry em todos cards
36. BUG 18: Validação promoção sem verificar datas
37. BUG 19: Validação disponibilidade sem verificar horários
38. BUG 20: Link voltar usa href em vez de goto
39. BUG 22: Ordem categoria não impede duplicatas
40. BUG 23: Sem validação de nome único categoria
41. BUG 25: Ajuste estoque sem audit trail
42. BUG 26: Unidade sem validação formato
43. BUG 34: Sem cálculo tempo preparo estimado
44. BUG 38: Ajustes sem busca
45. BUG 39: Ajustes sem validação permissões
46. BUG 41: Alterar senha sem validação atual
47. BUG 43: Preferências vazio (placeholder)
48. BUG 47: Breadcrumb inconsistente
49. BUG 48: Upload sem validação dimensões
50. BUG 50: Upload sem progresso
51. BUG 56: Modais não fecham com ESC
52. BUG 60: Tabelas sem scroll horizontal mobile

### Baixa (10) - Melhorias de UX
53. BUG 27: Sem alerta visual estoque zero
54. BUG 42: Sem avatar upload
55. BUG 51: Drag drop não previne padrão
56. BUG 52: Loading sem skeleton em alguns lugares
57. BUG 53: Estados vazios sem CTA consistente
58. BUG 54: Erros sem retry em todos lugares
59. BUG 55: Sem loading global navegação
60. BUG 57: Modais sem backdrop click
61. BUG 61: Grid produtos mobile cards pequenos
62. BUG 62: Filtros mobile ocupam espaço vertical

---

## DETALHAMENTO COMPLETO DOS BUGS

---

### BUG 01
**Tela:** Login  
**Como reproduzir:** Tentar fazer login com e-mail em formato inválido (ex: "usuario" sem @)  
**Resultado esperado:** Validação de formato de e-mail antes de enviar requisição  
**Resultado atual:** Envia requisição mesmo com formato inválido, erro vem do backend  
**Severidade:** Crítica  
**Causa Raiz:** Input type="email" no HTML não valida formato, não há validação JavaScript adicional  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`  
**Solução proposta:** Adicionar validação regex de e-mail no handleSubmit antes de chamar API  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 02
**Tela:** Login  
**Como reproduzir:** Fazer login, fechar browser, abrir novamente  
**Resultado esperado:** Sessão persistida com "remember me"  
**Resultado atual:** Usuário precisa fazer login novamente sempre  
**Severidade:** Alta  
**Causa Raiz:** Não há checkbox "remember me" nem implementação de persistência de token  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`, `frontend/src/lib/stores/userStore.svelte`  
**Solução proposta:** Adicionar checkbox "Lembrar-me" e salvar token em localStorage  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 03
**Tela:** Perfil/Logout  
**Como reproduzir:** Fazer logout, verificar userStore  
**Resultado esperado:** userStore completamente limpo  
**Resultado atual:** userStore pode conter dados residuais  
**Severidade:** Alta  
**Causa Raiz:** Função logout apenas chama API, não limpa store local  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`, `frontend/src/lib/stores/userStore.svelte`  
**Solução proposta:** Adicionar userStore.clear() após logout bem-sucedido  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 04
**Tela:** Todas as telas  
**Como reproduzir:** Desconectar internet, tentar qualquer ação  
**Resultado esperado:** Mensagem de erro amigável "Verifique sua conexão"  
**Resultado atual:** Erro genérico ou timeout sem tratamento específico  
**Severidade:** Média  
**Causa Raiz:** API client não trata erros de rede (fetch failed, timeout)  
**Arquivos envolvidos:** `frontend/src/lib/api/client.ts`  
**Solução proposta:** Adicionar interceptor para detectar erros de rede e mostrar mensagem específica  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 05
**Tela:** Login, Register, formulários em geral  
**Como reproduzir:** Clicar múltiplas vezes no botão de submit  
**Resultado esperado:** Botão desabilitado após primeiro clique  
**Resultado atual:** Múltiplas requisições enviadas  
**Severidade:** Média  
**Causa Raiz:** Loading state não desabilita botão adequadamente  
**Arquivos envolvidos:** `frontend/src/routes/(auth)/login/+page.svelte`, `frontend/src/routes/(auth)/register/+page.svelte`  
**Solução proposta:** Garantir que disabled={loading} esteja em todos os botões de submit  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 06
**Tela:** Dashboard  
**Como reproduzir:** Abrir Dashboard  
**Resultado esperado:** KPIs com dados reais calculados vs ontem  
**Resultado atual:** KPIs mostram "+12% vs ontem" e "+8% vs ontem" hardcoded  
**Severidade:** Alta  
**Causa Raiz:** Dados de comparação são estáticos no template  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte` (linhas 189, 206)  
**Solução proposta:** Calcular variação real comparando com dados de ontem (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 07
**Tela:** Dashboard  
**Como reproduzir:** Abrir Dashboard, ver "Atividades Recentes"  
**Resultado esperado:** Atividades reais do sistema  
**Resultado atual:** Dados estáticos/hardcoded no código  
**Severidade:** Alta  
**Causa Raiz:** Array recentActivities é estático, não vem da API  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte` (linhas 44-49)  
**Solução proposta:** Criar endpoint de atividades ou remover seção temporariamente  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 08
**Tela:** Dashboard  
**Como reproduzir:** Ter produtos simples sem ingredientes, ver "Estoque Baixo"  
**Resultado esperado:** Mostrar apenas produtos compostos com estoque baixo  
**Resultado atual:** Calcula estoque baixo para todos produtos, mesmo simples  
**Severidade:** Alta  
**Causa Raiz:** Lógica filtra produtos com ingredients.stock < 10 sem verificar IsComposto  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte` (linha 82)  
**Solução proposta:** Adicionar filtro p.IsComposto antes de verificar estoque  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 09
**Tela:** Dashboard  
**Como reproduzir:** Abrir Dashboard e deixar aberto  
**Resultado esperado:** Dados atualizam automaticamente (polling ou WebSocket)  
**Resultado atual:** Dados estáticos, requer refresh manual  
**Severidade:** Média  
**Causa Raiz:** Não há mecanismo de refresh automático  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Solução proposta:** Adicionar setInterval para recarregar dados a cada X segundos  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 10
**Tela:** Dashboard  
**Como reproduzir:** Simular erro em uma das APIs (products, orders, adjustments)  
**Resultado esperado:** Botão "Tentar novamente" em todos os cards  
**Resultado atual:** Apenas alert geral, retry não disponível em cada card  
**Severidade:** Média  
**Causa Raiz:** Tratamento de erro é global, não por seção  
**Arquivos envolvidos:** `frontend/src/routes/(app)/dashboard/+page.svelte`  
**Solução proposta:** Separar loading/error por seção (metrics, orders, ingredients)  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 11
**Tela:** Produtos - Novo Produto  
**Como reproduzir:** Criar novo produto  
**Resultado esperado:** Ver campos SEO (Slug, MetaTitle, MetaDescription, AltImage, Canonical)  
**Resultado atual:** Campos SEO não aparecem no formulário  
**Severidade:** Alta  
**Causa Raiz:** Form não inclui campos SEO apesar de backend suportar  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`  
**Solução proposta:** Adicionar aba "SEO" com campos Slug, MetaTitle, MetaDescription, AltImage, Canonical  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 12
**Tela:** Produtos - Editar Produto  
**Como reproduzir:** Editar produto existente  
**Resultado esperado:** Ver e editar campos SEO  
**Resultado atual:** Campos SEO não aparecem no formulário  
**Severidade:** Alta  
**Causa Raiz:** Form de edição não inclui campos SEO  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Solução proposta:** Adicionar aba "SEO" com campos SEO (igual BUG 11)  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 13
**Tela:** Produtos - Novo/Editar  
**Como reproduzir:** Criar ou editar produto  
**Resultado esperado:** Ver campos iFood (ExternalID, MarketplaceID, SyncStatus, LastSync)  
**Resultado atual:** Campos iFood não aparecem no formulário  
**Severidade:** Alta  
**Causa Raiz:** Form não inclui campos iFood apesar de backend suportar  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Solução proposta:** Adicionar aba "Integrações" com campos iFood  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 14
**Tela:** Produtos - Listagem  
**Como reproduzir:** Clicar em "Duplicar" em um produto  
**Resultado esperado:** Produto duplicado com todos os campos (SEO, iFood)  
**Resultado atual:** Produto duplicado sem campos SEO e iFood  
**Severidade:** Alta  
**Causa Raiz:** Função duplicateProduct não inclui campos SEO e iFood no payload  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte` (linhas 139-171)  
**Solução proposta:** Adicionar campos SEO e iFood ao payload de duplicação  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 15
**Tela:** Produtos - Listagem  
**Como reproduzir:** Clicar em "Arquivar" em produto com pedidos  
**Resultado esperado:** Alerta "Este produto possui X pedidos vinculados"  
**Resultado atual:** Arquiva sem verificar dependências  
**Severidade:** Crítica  
**Causa Raiz:** Função archiveProduct não verifica se há pedidos vinculados  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte` (linhas 173-205)  
**Solução proposta:** Verificar se produto tem pedidos antes de arquivar (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 16
**Tela:** Produtos - Listagem  
**Como reproduzir:** Clicar em "Excluir" em produto  
**Resultado esperado:** Verificação de dependências (pedidos, fichas técnicas)  
**Resultado atual:** Exclui sem verificar dependências  
**Severidade:** Crítica  
**Causa Raiz:** Função deleteProductById não verifica dependências  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte` (linhas 129-137)  
**Solução proposta:** Verificar dependências antes de excluir (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 17
**Tela:** Pedidos - Novo Pedido (POS)  
**Como reproduzir:** Selecionar categoria "Bebidas"  
**Resultado esperado:** Filtrar por categoria real do produto  
**Resultado atual:** Filtra por palavras no nome (suco, refrigerante, café)  
**Severidade:** Alta  
**Causa Raiz:** Lógica de filtro usa strings hardcoded no nome do produto  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte` (linhas 47-50)  
**Solução proposta:** Usar CategoryID real do produto para filtrar  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 18
**Tela:** Produtos - Novo/Editar  
**Como reproduzir:** Definir promoção com data fim anterior à data início  
**Resultado esperado:** Validação "Data fim deve ser posterior à data início"  
**Resultado atual:** Salva sem validação  
**Severidade:** Média  
**Causa Raiz:** Função validate não verifica datas de promoção  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Solução proposta:** Adicionar validação de datas na função validate  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 19
**Tela:** Produtos - Novo/Editar  
**Como reproduzir:** Definir disponibilidade com horário fim anterior ao horário início  
**Resultado esperado:** Validação "Horário fim deve ser posterior ao horário início"  
**Resultado atual:** Salva sem validação  
**Severidade:** Média  
**Causa Raiz:** Função validate não verifica horários de disponibilidade  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/new/+page.svelte`, `frontend/src/routes/(app)/products/[id]/edit/+page.svelte`  
**Solução proposta:** Adicionar validação de horários na função validate  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 20
**Tela:** Produtos - Detalhes ([id])  
**Como reproduzir:** Clicar em "Voltar para Cardápio"  
**Resultado esperado:** Navegação SvelteKit (goto)  
**Resultado atual:** Link HTML (<a href>)  
**Severidade:** Média  
**Causa Raiz:** Usa <a href> em vez de navegação programática  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/[id]/+page.svelte` (linha 89)  
**Solução proposta:** Usar goto('/products) ou Link component do SvelteKit  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 21
**Tela:** Categorias - Listagem  
**Como reproduzir:** Tentar excluir categoria com produtos vinculados  
**Resultado esperado:** Alerta "Esta categoria possui X produtos"  
**Resultado atual:** Exclui sem verificar  
**Severidade:** Crítica  
**Causa Raiz:** Função deleteCategoryById não verifica produtos vinculados  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte` (linhas 87-95)  
**Solução proposta:** Verificar produtos vinculados antes de excluir (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 22
**Tela:** Categorias - Novo/Editar  
**Como reproduzir:** Criar duas categorias com mesma ordem de exibição  
**Resultado esperado:** Validação ou aviso de duplicata  
**Resultado atual:** Permite duplicatas  
**Severidade:** Média  
**Causa Raiz:** Não há validação de DisplayOrder único  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte`  
**Solução proposta:** Adicionar validação de DisplayOrder único  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 23
**Tela:** Categorias - Novo/Editar  
**Como reproduzir:** Criar duas categorias com mesmo nome  
**Resultado esperado:** Validação de nome único  
**Resultado atual:** Permite duplicatas  
**Severidade:** Média  
**Causa Raiz:** Não há validação de nome único  
**Arquivos envolvidos:** `frontend/src/routes/(app)/categories/+page.svelte`  
**Solução proposta:** Adicionar validação de nome único  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 24
**Tela:** Ingredientes - Listagem  
**Como reproduzir:** Tentar excluir ingrediente usado em fichas técnicas  
**Resultado esperado:** Alerta "Este ingrediente é usado em X produtos"  
**Resultado atual:** Exclui sem verificar  
**Severidade:** Crítica  
**Causa Raiz:** Função deleteIngredientById não verifica fichas técnicas  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte` (linhas 99-107)  
**Solução proposta:** Verificar fichas técnicas antes de excluir (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 25
**Tela:** Ingredientes - Ajuste de Estoque  
**Como reproduzir:** Ajustar estoque de ingrediente  
**Resultado esperado:** Registro em histórico/audit trail  
**Resultado atual:** Ajuste sem rastro  
**Severidade:** Média  
**Causa Raiz:** Não há implementação de histórico de ajustes  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte` (linhas 109-131)  
**Solução proposta:** Criar tabela de histórico de ajustes (requer backend)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 26
**Tela:** Ingredientes - Novo/Editar  
**Como reproduzir:** Inserir unidade inválida (ex: "xyz")  
**Resultado esperado:** Validação de formato (kg, g, L, ml, un, etc.)  
**Resultado atual:** Aceita qualquer texto  
**Severidade:** Média  
**Causa Raiz:** Input de unidade é texto livre sem validação  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Solução proposta:** Usar Select com opções pré-definidas de unidades  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 27
**Tela:** Ingredientes - Listagem  
**Como reproduzir:** Estoque de ingrediente chega a zero  
**Resultado esperado:** Alerta visual proeminente (badge vermelho)  
**Resultado atual:** Badge "Estoque Zerado" pouco visível  
**Severidade:** Baixa  
**Causa Raiz:** Badge não é suficientemente destacado  
**Arquivos envolvidos:** `frontend/src/routes/(app)/ingredients/+page.svelte`  
**Solução proposta:** Melhorar destaque visual de estoque zero  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 28
**Tela:** Pedidos - Novo Pedido (POS)  
**Como reproduzir:** Criar novo pedido  
**Resultado esperado:** Campo para selecionar número da mesa  
**Resultado atual:** Não há seleção de mesa  
**Severidade:** Alta  
**Causa Raiz:** Form não inclui campo de mesa  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Solução proposta:** Adicionar campo "Número da Mesa" no formulário  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 29
**Tela:** Pedidos - Novo Pedido (POS)  
**Como reproduzir:** Adicionar produto ao carrinho  
**Resultado esperado:** Ver estoque disponível de ingredientes  
**Resultado atual:** Não mostra estoque  
**Severidade:** Alta  
**Causa Raiz:** Não há exibição de estoque de ingredientes  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Solução proposta:** Mostrar estoque disponível ao passar mouse no produto  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 30
**Tela:** Pedidos - Novo Pedido (POS)  
**Como reproduzir:** Adicionar produto com estoque insuficiente ao carrinho  
**Resultado esperado:** Alerta "Estoque insuficiente para X produtos"  
**Resultado atual:** Permite adicionar, erro só ao confirmar pedido  
**Severidade:** Crítica  
**Causa Raiz:** Validação de estoque só no backend ao criar pedido  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Solução proposta:** Validar estoque ao adicionar ao carrinho (requer backend endpoint)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 31
**Tela:** Pedidos - Detalhes  
**Como reproduzir:** Tentar editar pedido existente  
**Resultado esperado:** Tela de edição de pedido  
**Resultado atual:** Rota `/orders/[id]/edit` não existe  
**Severidade:** Alta  
**Causa Raiz:** Tela de edição não foi implementada  
**Arquivos envolvidos:** N/A (arquivo não existe)  
**Solução proposta:** Criar tela `frontend/src/routes/(app)/orders/[id]/edit/+page.svelte`  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 32
**Tela:** Pedidos - Detalhes  
**Como reproduzir:** Clicar em cancelar pedido  
**Resultado esperado:** Modal de confirmação "Tem certeza?"  
**Resultado atual:** Cancela sem confirmação  
**Severidade:** Crítica  
**Causa Raiz:** Não há modal de confirmação  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/[id]/+page.svelte`  
**Solução proposta:** Adicionar ConfirmDialog antes de cancelar  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 33
**Tela:** Pedidos - Detalhes  
**Como reproduzir:** Tentar fechar/concluir pedido  
**Resultado esperado:** Botão/ação para fechar pedido  
**Resultado atual:** Não há ação de fechar pedido  
**Severidade:** Alta  
**Causa Raiz:** Fluxo de fechamento não implementado  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/[id]/+page.svelte`  
**Solução proposta:** Adicionar botão "Fechar Pedido" com mudança de status  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 34
**Tela:** Pedidos - Novo Pedido (POS)  
**Como reproduzir:** Adicionar produtos ao carrinho  
**Resultado esperado:** Ver tempo de preparo estimado total  
**Resultado atual:** Não mostra tempo estimado  
**Severidade:** Média  
**Causa Raiz:** Não há cálculo de tempo de preparo  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte`  
**Solução proposta:** Somar PreparationTimeMinutes dos produtos no carrinho  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 35
**Tela:** Pedidos - Listagem  
**Como reproduzir:** Ver status pills  
**Resultado esperado:** Contagem real de pedidos por status  
**Resultado atual:** Contagem baseada em dados filtrados atuais, não total  
**Severidade:** Alta  
**Causa Raiz:** countByStatus calculado sobre array filtrado  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/+page.svelte` (linhas 116-121)  
**Solução proposta:** Calcular countByStatus sobre array orders (não filtered)  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 36
**Tela:** Ajustes de Estoque - Listagem  
**Como reproduzir:** Abrir tela de ajustes  
**Resultado esperado:** Ver todos os ajustes (pending, approved, rejected)  
**Resultado atual:** API getPendingAdjustments retorna apenas pendentes  
**Severidade:** Alta  
**Causa Raiz:** Função loadAdjustments chama getPendingAdjustments  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte` (linha 25)  
**Solução proposta:** Criar endpoint getAllAdjustments ou usar existente  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 37
**Tela:** Ajustes de Estoque - Detalhes  
**Como reproduzir:** Ver card de ajuste  
**Resultado esperado:** Ver nome do pedido e nome do ingrediente  
**Resultado atual:** Mostra apenas IDs (#123, #456)  
**Severidade:** Alta  
**Causa Raiz:** API retorna apenas IDs, não nomes  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte` (linhas 249-255)  
**Solução proposta:** Backend deve incluir nomos ou frontend deve buscar  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 38
**Tela:** Ajustes de Estoque - Listagem  
**Como reproduzir:** Tentar buscar ajuste por pedido ou ingrediente  
**Resultado esperado:** Campo de busca  
**Resultado atual:** Não há busca  
**Severidade:** Média  
**Causa Raiz:** Não implementado  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Solução proposta:** Adicionar campo de busca por ID ou nome  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 39
**Tela:** Ajustes de Estoque - Ação  
**Como reproduzir:** Tentar aprovar/rejeitar ajuste  
**Resultado esperado:** Validação de permissão  
**Resultado atual:** Qualquer usuário pode aprovar/rejeitar  
**Severidade:** Média  
**Causa Raiz:** Não há validação de permissões no frontend  
**Arquivos envolvidos:** `frontend/src/routes/(app)/stock-adjustments/+page.svelte`  
**Solução proposta:** Verificar role/permissão do usuário antes de permitir ação  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 40
**Tela:** Perfil - Informações Pessoais  
**Como reproduzir:** Alterar e-mail  
**Resultado esperado:** Pedir senha atual para confirmar  
**Resultado atual:** Altera sem confirmação  
**Severidade:** Alta  
**Causa Raiz:** Form não pede confirmação de senha  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte` (linhas 39-58)  
**Solução proposta:** Adicionar campo "Senha atual" ao alterar e-mail  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 41
**Tela:** Perfil - Segurança  
**Como reproduzir:** Alterar senha  
**Resultado esperado:** Sempre pedir senha atual  
**Resultado atual:** Já pede, mas validação pode ser melhorada  
**Severidade:** Média  
**Causa Raiz:** Validação básica implementada  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte` (linhas 60-95)  
**Solução proposta:** Melhorar validação (complexidade, caracteres especiais)  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 42
**Tela:** Perfil - Preferências  
**Como reproduzir:** Abrir aba Preferências  
**Resultado esperado:** Configurações de tema, idioma, notificações  
**Resultado atual:** Placeholder "Configurações adicionais estarão disponíveis em breve"  
**Severidade:** Baixa  
**Causa Raiz:** Funcionalidade não implementada  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte` (linhas 250-263)  
**Solução proposta:** Implementar preferências ou remover aba  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 43
**Tela:** Perfil - Informações Pessoais  
**Como reproduzir:** Ver seção de avatar  
**Resultado esperado:** Upload de avatar  
**Resultado atual:** Não há upload de avatar  
**Severidade:** Baixa  
**Causa Raiz:** Funcionalidade não implementada  
**Arquivos envolvidos:** `frontend/src/routes/(app)/profile/+page.svelte`  
**Solução proposta:** Adicionar PhotoUpload para avatar  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 44
**Tela:** Sidebar - Navegação  
**Como reproduzir:** Ver badges de notificações  
**Resultado esperado:** Contagem real de pedidos pendentes e ajustes  
**Resultado atual:** Badges hardcoded (Pedidos: 3, Ajustes: 2)  
**Severidade:** Crítica  
**Causa Raiz:** Badges são estáticos no código  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte` (linhas 38, 51)  
**Solução proposta:** Buscar contagens reais da API  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 45
**Tela:** Sidebar - Navegação  
**Como reproduzir:** Clicar em "Usuários" ou "Configurações"  
**Resultado esperado:** Navegar para telas correspondentes  
**Resultado atual:** Rotas não existem (404)  
**Severidade:** Alta  
**Causa Raiz:** Links apontam para rotas não implementadas  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte` (linhas 57, 59)  
**Solução proposta:** Remover links ou implementar telas  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 46
**Tela:** Header - Busca  
**Como reproduzir:** Digitar no campo de busca e pressionar Enter  
**Resultado esperado:** Redirecionar para busca global  
**Resultado atual:** console.log apenas  
**Severidade:** Alta  
**Causa Raiz:** Função handleSearch não implementada  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Header.svelte` (linhas 14-19)  
**Solução proposta:** Implementar busca global ou remover campo  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 47
**Tela:** Todas as telas com Workspace  
**Como reproduzir:** Comparar breadcrumb com sidebar  
**Resultado esperado:** Breadcrumb consistente com navegação atual  
**Resultado atual:** Breadcrumb pode ser inconsistente  
**Severidade:** Média  
**Causa Raiz:** Breadcrumb é manual, não automático  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Workspace.svelte`  
**Solução proposta:** Tornar breadcrumb automático baseado em rota  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 48
**Tela:** Upload de Foto (Produtos)  
**Como reproduzir:** Fazer upload de imagem muito pequena ou muito grande  
**Resultado esperado:** Validação de dimensões (ex: mínimo 200x200, máximo 4000x4000)  
**Resultado atual:** Valida apenas tamanho do arquivo (5MB)  
**Severidade:** Média  
**Causa Raiz:** PhotoUpload não valida dimensões  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte` (linhas 52-61)  
**Solução proposta:** Adicionar validação de dimensões após carregar imagem  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 49
**Tela:** Upload de Foto (Produtos)  
**Como reproduzir:** Simular erro de upload  
**Resultado esperado:** Toast/Alert com mensagem de erro  
**Resultado atual:** alert() nativo do browser  
**Severidade:** Alta  
**Causa Raiz:** Usa alert() em vez de componente de UI  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte` (linha 78)  
**Solução proposta:** Usar Alert ou Toast component  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 50
**Tela:** Upload de Foto (Produtos)  
**Como reproduzir:** Fazer upload de imagem grande  
**Resultado esperado:** Barra de progresso de upload  
**Resultado atual:** Spinner sem progresso  
**Severidade:** Média  
**Causa Raiz:** API não retorna progresso de upload  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte`, `frontend/src/lib/api/media.ts`  
**Solução proposta:** Implementar upload com progresso (requer backend suporte)  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** NÃO (requer backend)

---

### BUG 51
**Tela:** Upload de Foto (Produtos)  
**Como reproduzir:** Arrastar arquivo e soltar  
**Resultado esperado:** Comportamento consistente em todos browsers  
**Resultado atual:** Drag and drop pode não funcionar em alguns browsers  
**Severidade:** Baixa  
**Causa Raiz:** Event handlers podem não prevenir comportamento padrão corretamente  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/PhotoUpload.svelte` (linhas 33-50)  
**Solução proposta:** Melhorar handlers de drag and drop  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 52
**Tela:** Variadas  
**Como reproduzir:** Carregar dados em várias telas  
**Resultado esperado:** Skeleton loading em todos os lugares  
**Resultado atual:** Alguns lugares têm spinner, outros nada  
**Severidade:** Baixa  
**Causa Raiz:** Implementação inconsistente de loading states  
**Arquivos envolvidos:** Múltiplos arquivos  
**Solução proposta:** Padronizar uso de Skeleton component  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 53
**Tela:** Variadas  
**Como reproduzir:** Ver estados vazios (sem dados)  
**Resultado esperado:** Call-to-action consistente ("Criar X", "Adicionar Y")  
**Resultado atual:** Alguns estados vazios sem CTA  
**Severidade:** Baixa  
**Causa Raiz:** EmptyState component não usado consistentemente  
**Arquivos envolvidos:** Múltiplos arquivos  
**Solução proposta:** Padronizar uso de EmptyState component  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 54
**Tela:** Variadas  
**Como reproduzir:** Ver estados de erro  
**Resultado esperado:** Botão "Tentar novamente" em todos os erros  
**Resultado atual:** Alguns erros sem retry  
**Severidade:** Baixa  
**Causa Raiz:** Tratamento de erro inconsistente  
**Arquivos envolvidos:** Múltiplos arquivos  
**Solução proposta:** Padronizar tratamento de erro com retry  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 55
**Tela:** Navegação entre páginas  
**Como reproduzir:** Navegar entre telas  
**Resultado esperado:** Loading global durante transição  
**Resultado atual:** Sem loading, pode parecer travado  
**Severidade:** Baixa  
**Causa Raiz:** Não há loading global de navegação  
**Arquivos envolvidos:** `frontend/src/routes/+layout.svelte`  
**Solução proposta:** Adicionar loading global baseado em navigating do SvelteKit  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 56
**Tela:** Modais (Produtos, Ajustes)  
**Como reproduzir:** Abrir modal e pressionar ESC  
**Resultado esperado:** Modal fecha  
**Resultado atual:** Modal não fecha  
**Severidade:** Média  
**Causa Raiz:** Modal component não implementa keydown ESC  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Modal.svelte`  
**Solução proposta:** Adicionar handler para tecla ESC  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 57
**Tela:** Modais (Produtos, Ajustes)  
**Como reproduzir:** Abrir modal e clicar fora (backdrop)  
**Resultado esperado:** Modal fecha  
**Resultado atual:** Modal não fecha  
**Severidade:** Baixa  
**Causa Raiz:** Modal component não implementa backdrop click  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Modal.svelte`  
**Solução proposta:** Adicionar handler para backdrop click  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 58
**Tela:** Mobile (todas as telas)  
**Como reproduzir:** Abrir em mobile (< 768px)  
**Resultado esperado:** Menu hamburguer para abrir sidebar  
**Resultado atual:** Sidebar some completamente, sem acesso  
**Severidade:** Alta  
**Causa Raiz:** Sidebar tem display: none em mobile sem alternativa  
**Arquivos envolvidos:** `frontend/src/lib/components/layout/Sidebar.svelte` (linhas 509-513)  
**Solução proposta:** Implementar menu hamburguer/drawer para mobile  
**Impacto:** Alto  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 59
**Tela:** Pedidos - POS (Mobile)  
**Como reproduzir:** Abrir POS em mobile  
**Resultado esperado:** Layout responsivo com carrinho acessível  
**Resultado atual:** Carrinho some em mobile  
**Severidade:** Alta  
**Causa Raiz:** Grid layout não responsivo  
**Arquivos envolvidos:** `frontend/src/routes/(app)/orders/new/+page.svelte` (linhas 645-653)  
**Solução proposta:** Mudar layout para mobile (carrinho em drawer ou abaixo)  
**Impacto:** Médio  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 60
**Tela:** Tabelas (Pedidos, Ajustes)  
**Como reproduzir:** Abrir em mobile  
**Resultado esperado:** Scroll horizontal em tabelas  
**Resultado atual:** Tabelas podem quebrar layout  
**Severidade:** Média  
**Causa Raiz:** Tabelas não têm overflow-x  
**Arquivos envolvidos:** `frontend/src/lib/components/ui/Table.svelte`  
**Solução proposta:** Adicionar overflow-x: auto em tabelas  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 61
**Tela:** Produtos - Listagem (Mobile)  
**Como reproduzir:** Abrir em mobile  
**Resultado esperado:** Cards legíveis com tamanho adequado  
**Resultado atual:** Grid minmax(160px) pode deixar cards muito pequenos  
**Severidade:** Baixa  
**Causa Raiz:** Grid responsivo pode ser agressivo  
**Arquivos envolvidos:** `frontend/src/routes/(app)/products/+page.svelte` (linhas 684-689)  
**Solução proposta:** Ajustar minmax para mobile ou usar 1 coluna  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

### BUG 62
**Tela:** Filtros (Mobile)  
**Como reproduzir:** Abrir filtros em mobile  
**Resultado esperado:** Filtros em drawer ou accordion  
**Resultado atual:** Filtros ocupam muito espaço vertical  
**Severidade:** Baixa  
**Causa Raiz:** Filtros sempre expandidos  
**Arquivos envolvidos:** Múltiplas telas com filtros  
**Solução proposta:** Colapsar filtros em mobile ou usar drawer  
**Impacto:** Baixo  
**Pode ser corrigido sem alterar arquitetura?** SIM

---

## PLANOS DE CORREÇÃO SUGERIDOS

### Plano 1: Sprint de Estabilização (Prioridade Crítica)

**Objetivo:** Corrigir 8 bugs críticos em 2-3 dias

**Bugs a corrigir:**
- BUG 01: Validação e-mail login
- BUG 15: Verificar pedidos ao arquivar (backend)
- BUG 16: Verificar dependências ao excluir (backend)
- BUG 21: Verificar produtos ao excluir categoria (backend)
- BUG 24: Verificar fichas ao excluir ingrediente (backend)
- BUG 30: Validar estoque ao adicionar carrinho (backend)
- BUG 32: Confirmação ao cancelar pedido
- BUG 44: Badges dinâmicos (backend)

**Esforço:** 2-3 dias  
**Dependências:** 5 bugs requerem backend

### Plano 2: Sprint de Melhorias UX (Prioridade Alta)

**Objetivo:** Corrigir 24 bugs de alta prioridade em 3-4 dias

**Bugs a corrigir:**
- BUG 02: Remember me
- BUG 03: Logout limpa store
- BUG 06-08: Dashboard dados reais (backend)
- BUG 11-14: Campos SEO/iFood
- BUG 17: Filtro categoria real
- BUG 28-29: POS mesa e estoque
- BUG 31: Editar pedido
- BUG 33: Fechar pedido
- BUG 35: Contagem status pills
- BUG 36-37: Ajustes lista e nomes (backend)
- BUG 40: Alterar e-mail com senha
- BUG 45-46: Links e busca
- BUG 49: Erro upload com Alert
- BUG 58-59: Responsividade mobile

**Esforço:** 3-4 dias  
**Dependências:** 6 bugs requerem backend

### Plano 3: Sprint de Polishing (Prioridade Média/Baixa)

**Objetivo:** Corrigir 30 bugs de média/baixa prioridade em 2-3 dias

**Bugs a corrigir:**
- BUG 04-05: Erro rede e loading
- BUG 09-10: Dashboard refresh e retry
- BUG 18-20: Validações produto
- BUG 22-23: Categoria validações
- BUG 25-27: Ingrediente melhorias
- BUG 34: Tempo preparo
- BUG 38-39: Ajustes busca e permissões
- BUG 41-43: Perfil melhorias
- BUG 47: Breadcrumb automático
- BUG 48, 50-51: Upload melhorias
- BUG 52-57: UX consistência
- BUG 60-62: Responsividade tabelas e filtros

**Esforço:** 2-3 dias  
**Dependências:** 1 bug requer backend (BUG 25)

---

## CONCLUSÃO

A auditoria identificou **62 bugs** no frontend, sendo:
- **8 críticos** (bloqueadores)
- **24 alta** (impacto significativo)
- **20 média** (comportamentos inesperados)
- **10 baixa** (melhorias de UX)

**Recomendação final:**
1. Corrigir 8 bugs críticos antes do próximo épico
2. Implementar melhorias de alta prioridade na próxima sprint
3. Deixar melhorias de média/baixa para polishing contínuo

**Impacto no MVP:** Os bugs críticos e de alta prioridade afetam diretamente a usabilidade e confiabilidade do sistema. A correção é essencial antes de avançar com novos recursos.

**Tempo total estimado:** 8-12 dias (1.5 a 2 sprints) para correção completa.
