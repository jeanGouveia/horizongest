# Session Management

## Overview

O HorizonGest possui uma arquitetura profissional de gerenciamento de sessões que garante segurança, consistência e previsibilidade em todo o ciclo de vida da autenticação. Esta arquitetura é equivalente a sistemas profissionais como Stripe, Notion, GitHub, Slack e Linear.

## Arquitetura

### SessionManager

**Arquivo:** `frontend/src/lib/managers/sessionManager.ts`

**Responsabilidades:**
- Validar sessão na inicialização
- Gerenciar Platform Session
- Gerenciar Tenant Session (via TenantSessionManager)
- Executar logout completo
- Destruir todas as sessões

### TenantSessionManager

**Arquivo:** `frontend/src/lib/managers/tenantSessionManager.ts`

**Responsabilidades:**
- Gerenciar entrada em empresa (`enterCompany`)
- Gerenciar saída de empresa (`leaveCompany`)
- Gerenciar troca de empresa (`switchCompany`)
- Destruir contexto completo (`destroy`)

## Política Oficial do HorizonGest

Esta é uma regra arquitetural permanente.

### Platform Session

Responsável pelo acesso à Plataforma.

**Características:**
- Exige login
- Possui expiração
- Nunca pode ser restaurada sem validação do backend
- Perde validade após reinício do backend
- Perde validade após logout
- Perde validade após expiração do JWT

**Armazenamento:**
- Cookie: `platform_auth_token`

### Tenant Session

Responsável pelo acesso à empresa.

**Características:**
- Só existe enquanto existir uma Platform Session válida
- Nunca pode sobreviver sozinha
- Nunca pode sobreviver ao logout
- Nunca pode sobreviver ao reinício do backend
- Nunca pode sobreviver à troca de empresa

**Armazenamento:**
- Cookie: `auth_token`
- LocalStorage: `impersonation`

## Fluxo de Inicialização

Quando a aplicação iniciar:

```
App
↓
Existe token?
↓
SIM
↓
SessionManager.validateSession()
↓
Validar Platform Session no backend
↓
Sessão válida?
↓
SIM
↓
Validar Tenant Session (se existir)
↓
Tenant válida?
↓
SIM
↓
Hidratar stores
↓
Dashboard
↓
NÃO
↓
Destroy Tenant Session
↓
Dashboard (Platform)
↓
NÃO
↓
Destroy All Sessions
↓
Limpar LocalStorage
↓
Limpar SessionStorage
↓
Limpar Cookies
↓
Limpar todas as Stores
↓
Limpar caches
↓
Tela Login
```

**Regra:** Nunca abrir Dashboard apenas porque existe token salvo.

## Entrar na Empresa

Fluxo obrigatório:

```
Clique "Entrar na Empresa"
↓
TenantSessionManager.enterCompany(companyId)
↓
Valida se já está entrando (previne race conditions)
↓
Valida se já está na empresa selecionada
↓
Destruir contexto anterior
↓
Solicitar novo Tenant JWT
↓
Limpar Stores Tenant
↓
Limpar Cache (localStorage, sessionStorage, navigation cache)
↓
Hidratar contexto
↓
Carregar Empresa
↓
Carregar Branding
↓
Carregar Permissões
↓
Carregar Dashboard
↓
Navegar para Dashboard
```

**Regra:** Nunca navegar antes da hidratação terminar.

## Trocar de Empresa

Trocar empresa nunca poderá apenas trocar JWT.

Fluxo:

```
TenantSessionManager.switchCompany(newCompanyId)
↓
leaveCompany()
↓
destroy()
↓
encerrar impersonation
↓
limpar stores tenant
↓
limpar cache tenant
↓
voltar plataforma
↓
enterCompany(newCompanyId)
↓
destruir contexto anterior
↓
solicitar novo Tenant JWT (backend encerra impersonação anterior automaticamente)
↓
limpar stores
↓
limpar cache
↓
hidratar contexto
↓
carregar empresa
↓
carregar branding
↓
carregar permissões
↓
navegar para dashboard
```

## Sair da Empresa

Fluxo oficial:

```
TenantSessionManager.leaveCompany()
↓
encerrar impersonation no backend
↓
destruir contexto tenant
↓
limpar stores tenant
↓
limpar cache tenant
↓
voltar plataforma
```

**Regra:** A Platform Session permanece.

## Logout

Logout significa destruir absolutamente tudo.

Fluxo:

```
SessionManager.logout()
↓
Encerrar impersonation (se existir)
↓
Encerrar Platform Session
↓
Limpar cookies (platform_auth_token, auth_token)
↓
Limpar LocalStorage
↓
Limpar SessionStorage
↓
Limpar todas as Stores
↓
Limpar todos os caches
↓
Tela Login
```

## Backend

### Impersonation Service

**Arquivo:** `backend/internal/service/impersonation_service.go`

**Regra Arquitetural:** A troca de empresa é responsabilidade exclusiva do backend.

**Comportamento:**
1. **Endpoint idempotente:**
   - POST `/api/platform/impersonation/start` encerra automaticamente qualquer impersonation ativa
   - Não retorna mais HTTP 409 "já existe uma sessão de impersonation ativa"
   - Permite troca livre de empresas sem necessidade de logout

2. **Gerenciamento automático:**
   - Se existir impersonation ativa (independente da idade)
   - Backend encerra automaticamente (persiste EndedAt)
   - Cria nova impersonation
   - Gera novo Tenant JWT
   - Retorna 200 OK

3. **Frontend simplificado:**
   - Não precisa chamar `endPreviousImpersonation()` antes de `enterCompany()`
   - Apenas solicita `requestTenantJWT()`
   - Backend resolve toda a troca

4. **Auditoria preservada:**
   - Impersonation anterior recebe EndedAt corretamente
   - Nova impersonation é criada
   - Nunca existem duas impersonations ativas para o mesmo PlatformUser
   - Histórico permanece íntegro

### Tratamento de 401 Unauthorized

**Arquivo:** `frontend/src/lib/api/client.ts`

**Implementação:**
- Qualquer endpoint que retorna 401
- Automaticamente destrói todas as sessões
- Redireciona para login
- Previne acesso com sessão inválida

## Stores

Todas as stores possuem métodos de limpeza:

### userStore
- **Arquivo:** `frontend/src/lib/stores/userStore.svelte.ts`
- **Método:** `logout()`
- **Responsabilidade:** Armazenar usuário logado

### companyStore
- **Arquivo:** `frontend/src/lib/stores/companyStore.svelte.ts`
- **Método:** `clear()`
- **Responsabilidade:** Armazenar configurações da empresa

### rbacStore
- **Arquivo:** `frontend/src/lib/stores/rbacStore.svelte.ts`
- **Método:** `reset()`
- **Responsabilidade:** Armazenar role e permissões

### themeStore
- **Arquivo:** `frontend/src/lib/stores/themeStore.svelte.ts`
- **Método:** `clear()`
- **Responsabilidade:** Armazenar tema visual

### brandStore
- **Arquivo:** `frontend/src/lib/stores/brandStore.ts`
- **Método:** `clear()`
- **Responsabilidade:** Armazenar configurações de branding

### toast
- **Arquivo:** `frontend/src/lib/stores/toast.ts`
- **Método:** `clear()`
- **Responsabilidade:** Armazenar notificações

**Regra:** Não recriar singleton. Reinicializar conteúdo.

## Cache

Todo cache é invalidado na troca de empresa:

### LocalStorage
- **Chave:** `impersonation`
- **Limpeza:** Removido em `clearTenantLocalStorage()`

### SessionStorage
- **Limpeza:** Limpeza granular em `clearTenantSessionStorage()`
- **Chaves limpas:** `tenant_navigation`, `tenant_forms`, `tenant_filters`
- **Regra:** Nunca usar `sessionStorage.clear()` - limpar apenas chaves específicas do tenant

### Navigation Cache
- **Limpeza:** NÃO utiliza Navigation API Experimental
- **Motivo:** `window.navigation.entries` não é suportada por todos os navegadores
- **Solução:** Limpeza granular de localStorage e sessionStorage é suficiente

### Cookies
- **Chaves:** `auth_token`, `platform_auth_token`
- **Limpeza:** Removidos via `document.cookie` com `max-age=0`

### Storage Keys Centralizadas
- **Arquivo:** `frontend/src/lib/constants/storage-keys.ts`
- **Propósito:** Centralizar todas as chaves de storage
- **Benefícios:**
  - Evitar typos
  - Facilitar refatoração
  - Documentar propósito de cada chave
- **Regra:** NUNCA usar strings literais para storage keys

## Integração com Componentes

### Botão "Entrar na Empresa"

**Arquivo:** `frontend/src/routes/platform/companies/[id]/+page.svelte`

**Implementação:**
```typescript
async function loginAsCompany() {
	const result = await tenantSessionManager.enterCompany(company.ID);
	if (result.success) {
		showSuccess('Sucesso', `Entrando na empresa ${company.Name}`);
	} else {
		showError('Erro', result.error);
	}
}
```

### ImpersonationBanner

**Arquivo:** `frontend/src/lib/components/ImpersonationBanner.svelte`

**Implementação:**
```typescript
async function endImpersonation() {
	const result = await tenantSessionManager.leaveCompany();
	if (result.success) {
		impersonationInfo = null;
		visible = false;
	}
}
```

### Sidebar Logout

**Arquivo:** `frontend/src/lib/components/layout/Sidebar.svelte`

**Implementação:**
```typescript
async function handleLogout() {
	const result = await sessionManager.logout();
	if (!result.success) {
		console.error('Logout error:', result.error);
	}
}
```

## Regras Importantes

### 1. Único Ponto de Troca
Nenhum componente deve trocar empresa diretamente. Toda troca deve passar pelo `TenantSessionManager`.

### 2. Navegação Apenas Após Contexto Consistente
Nunca navegar antes do contexto estar consistente. Os managers garantem que:
- Contexto anterior foi destruído
- Novo token foi obtido
- Stores foram limpas
- Caches foram limpos
- Novo contexto foi hidratado
- Branding foi carregado
- Permissões foram carregadas
- Empresa foi carregada

### 3. Stores Continuam Singleton
As stores NÃO deixaram de ser singleton. Singleton NÃO é o problema. O problema é o ciclo de vida.

### 4. Sem Hacks
- Não usar `setTimeout`
- Não usar `location.reload()`
- Não usar `window.reload()`
- Não usar `invalidateAll()` como solução principal
- Não duplicar código
- Não espalhar limpeza de stores em vários componentes
- Não usar APIs experimentais (ex: Navigation API)
- Não usar `@ts-ignore` para silenciar erros
- Não usar `sessionStorage.clear()` ou `localStorage.clear()` globalmente

### 5. APIs Compatíveis
- Usar apenas APIs suportadas por todos os navegadores
- Evitar APIs experimentais ou não padronizadas
- Testar em múltiplos navegadores (Chrome, Firefox, Safari, Edge)
- Verificar suporte em caniuse.com antes de usar novas APIs

### 6. Tratamento de Erros Diferenciado
- Separar erros por tipo: infrastructure, session, backend, ui
- Fornecer mensagens específicas para cada tipo de erro
- Tratar erros de rede (fetch failed) separadamente
- Tratar erros de autenticação (401) com redirecionamento
- Logar erros com contexto adequado

## Testes

### Testes Unitários

**Arquivo:** `frontend/src/lib/managers/__tests__/tenantSessionManager.test.ts`

**Cenários:**
1. Entrar na empresa com sucesso
2. Entrar na empresa com erro
3. Entrar na empresa já na empresa atual
4. Sair da empresa com sucesso
5. Sair da empresa com erro
6. Trocar de empresa com sucesso
7. Trocar de empresa com erro na saída
8. Trocar de empresa com erro na entrada

### Testes de Estresse

**Arquivo:** `frontend/src/lib/managers/__tests__/tenantSessionManager.stress.test.ts`

**Cenários:**
1. 100 trocas de empresa consecutivas
2. Garantir nenhuma store contaminada
3. Garantir nenhum cache persistente
4. Garantir empresa correta
5. Garantir permissões corretas
6. Garantir branding correto
7. Garantir dashboard correto
8. Manter performance estável

### Testes de Integração (Planejados)

**Cenários:**
1. Login → Dashboard
2. Logout → Login
3. Entrar Empresa → Sair Empresa → Entrar Novamente
4. Trocar Empresa 100 vezes consecutivas
5. Backend reiniciado → Sessão inválida
6. JWT expirado → 401 → Logout automático

## Benefícios Arquiteturais

### 1. Centralização
- Lógica de sessão em um único lugar
- Fácil de manter e evoluir
- Fácil de testar
- Responsabilidades claras

### 2. Consistência
- Contexto sempre consistente antes da navegação
- Nenhum estado fantasma
- Nenhum vazamento de contexto
- Comportamento previsível

### 3. Prevenção de Erros
- Flags previnem race conditions
- Validações previnem operações inválidas
- Tratamento de erros em cada etapa
- Recuperação robusta

### 4. Manutenibilidade
- Código mais limpo
- Sem lógica espalhada
- Fácil de adicionar novas funcionalidades
- Fácil de debugar

### 5. Testabilidade
- Fácil de testar isoladamente
- Fácil de mockar dependências
- Testes unitários e de estresse
- Cobertura de cenários críticos

### 6. Segurança
- Validação de sessão no backend
- Tratamento de 401 automático
- Prevenção de sessões órfãs
- Auto-limpeza de sessões stale

## Riscos Eliminados

### 1. Vazamento de Contexto
- **Antes:** Dados da Empresa A podiam aparecer na Empresa B
- **Depois:** Contexto é completamente destruído antes da troca

### 2. Estados Fantasmas
- **Antes:** Stores mantinham estado de empresas anteriores
- **Depois:** Stores são limpas em cada troca

### 3. Caches Persistentes
- **Antes:** localStorage e cookies não eram limpos
- **Depois:** Todos os caches são limpos

### 4. Race Conditions
- **Antes:** Múltiplos cliques podiam causar problemas
- **Depois:** Flags previnem operações simultâneas

### 5. Comportamento Aleatório
- **Antes:** Usuário com múltiplas empresas experimentava inconsistências
- **Depois:** Comportamento consistente e previsível

### 6. Sessões Órfãs
- **Antes:** Impersonation ativa sobrevivia ao reinício do backend
- **Depois:** Auto-limpeza de sessões stale e validação rigorosa

### 7. Erro 409 Eliminado
- **Antes:** "Já existe uma sessão de impersonation ativa" bloqueava usuários
- **Depois:** Backend encerra automaticamente qualquer impersonation ativa, permitindo troca livre

### 8. Acesso com Sessão Inválida
- **Antes:** Dashboard podia abrir com JWT expirado
- **Depois:** 401 automático destrói sessão e redireciona para login

## Nota da Arquitetura Antes

**Nota:** 3/10

**Problemas:**
- Sem ciclo de vida definido
- Lógica espalhada em múltiplos componentes
- Sem limpeza de contexto
- Sem prevenção de race conditions
- Difícil de testar
- Difícil de manter
- Estados fantasmas
- Vazamento de contexto
- Sessões órfãs
- Erro 409 frequente
- Acesso com sessão inválida

## Nota da Arquitetura Depois

**Nota:** 9/10

**Melhorias:**
- Ciclo de vida bem definido
- Lógica centralizada
- Limpeza completa de contexto
- Prevenção de race conditions
- Fácil de testar
- Fácil de manter
- Nenhum estado fantasma
- Nenhum vazamento de contexto
- Prevenção de sessões órfãs
- Auto-limpeza de sessões stale
- Tratamento de 401 automático

**Pontos a melhorar:**
- Adicionar testes de integração
- Adicionar monitoramento de erros
- Adicionar logging detalhado
- Adicionar métricas de performance

## Conclusão

O HorizonGest agora possui um gerenciamento de sessões profissional, centralizado e definitivo. A sessão nunca permanece inconsistente. O usuário nunca pode abrir um Dashboard sem possuir uma sessão validada pelo backend. O sistema transmite segurança, previsibilidade e robustez em todo o ciclo de vida da autenticação.

No futuro, qualquer funcionalidade relacionada à troca de empresa ou gerenciamento de sessão deve utilizar exclusivamente estes Managers.
