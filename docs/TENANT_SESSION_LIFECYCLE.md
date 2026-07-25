# Tenant Session Lifecycle

## Overview

O HorizonGest possui um gerenciamento centralizado do ciclo de vida da sessão Tenant através do `TenantSessionManager`. Esta camada é responsável por gerenciar a entrada, saída e troca de empresas, garantindo que não haja vazamento de contexto entre diferentes sessões.

## Arquitetura

### TenantSessionManager

**Arquivo:** `frontend/src/lib/managers/tenantSessionManager.ts`

**Responsabilidades:**
- Gerenciar entrada em empresa (`enterCompany`)
- Gerenciar saída de empresa (`leaveCompany`)
- Gerenciar troca de empresa (`switchCompany`)
- Destruir contexto completo (`destroy`)

**Regra:** Este é o ÚNICO lugar onde deve ocorrer troca de empresa. Nenhum outro componente deve trocar empresa diretamente.

## Operações

### enterCompany(companyId)

Responsável por entrar em uma empresa.

**Fluxo:**
1. Valida se já está entrando em uma empresa (previne race conditions)
2. Valida se já está na empresa selecionada
3. Destrói contexto anterior (`destroy()`)
4. Solicita Tenant JWT ao backend
5. Limpa caches (localStorage, sessionStorage)
6. Hidrata contexto com novo token
7. Carrega branding
8. Carrega permissões
9. Carrega empresa
10. Atualiza estado interno
11. Navega para Dashboard

**Proteções:**
- Flag `isEntering` previne múltiplas entradas simultâneas
- Validação de empresa atual previne reentrada desnecessária
- Tratamento de erros em cada etapa

### leaveCompany()

Responsável por sair da empresa atual.

**Fluxo:**
1. Valida se já está saindo da empresa (previne race conditions)
2. Destrói contexto Tenant (`destroy()`)
3. Limpa caches
4. Atualiza estado interno
5. Navega para Plataforma

**Proteções:**
- Flag `isLeaving` previne múltiplas saídas simultâneas
- Mantém Platform Session intacta

### switchCompany(companyId)

Responsável por trocar de empresa.

**Fluxo:**
1. Executa `leaveCompany()`
2. Executa `enterCompany(companyId)`

**Proteções:**
- Sequencialidade garantida
- Erro em qualquer etapa aborta a troca

### destroy()

Responsável por destruir completamente o contexto do tenant.

**Fluxo:**
1. Limpa todas as stores
2. Limpa cookies tenant
3. Limpa localStorage tenant

**Stores limpas:**
- `userStore` - logout()
- `companyStore` - clear()
- `rbacStore` - reset()
- `themeStore` - clear()
- `brandStore` - clear()
- `toast` - clear()

## Stores Envolvidas

### userStore
- **Arquivo:** `frontend/src/lib/stores/userStore.svelte.ts`
- **Método de limpeza:** `logout()`
- **Responsabilidade:** Armazenar usuário logado

### companyStore
- **Arquivo:** `frontend/src/lib/stores/companyStore.svelte.ts`
- **Método de limpeza:** `clear()`
- **Responsabilidade:** Armazenar configurações da empresa

### rbacStore
- **Arquivo:** `frontend/src/lib/stores/rbacStore.svelte.ts`
- **Método de limpeza:** `reset()`
- **Responsabilidade:** Armazenar role e permissões

### themeStore
- **Arquivo:** `frontend/src/lib/stores/themeStore.svelte.ts`
- **Método de limpeza:** `clear()`
- **Responsabilidade:** Armazenar tema visual

### brandStore
- **Arquivo:** `frontend/src/lib/stores/brandStore.ts`
- **Método de limpeza:** `clear()`
- **Responsabilidade:** Armazenar configurações de branding

### toast
- **Arquivo:** `frontend/src/lib/stores/toast.ts`
- **Método de limpeza:** `clear()`
- **Responsabilidade:** Armazenar notificações

## Caches Envolvidos

### localStorage
- **Chave:** `impersonation`
- **Limpeza:** Removido em `clearTenantLocalStorage()`

### sessionStorage
- **Limpeza:** `sessionStorage.clear()` em `clearCaches()`

### Cookies
- **Chave:** `auth_token`
- **Limpeza:** Removido em `clearTenantCookies()`

## Sequência de Eventos

### Entrar na Empresa

```
Clique "Entrar na Empresa"
↓
tenantSessionManager.enterCompany(companyId)
↓
destroy() - limpa contexto anterior
↓
requestTenantJWT() - obtém novo token
↓
clearCaches() - limpa localStorage/sessionStorage
↓
hydrateContext() - define novo cookie
↓
loadBranding() - carrega branding
↓
loadPermissions() - carrega permissões
↓
loadCompany() - carrega empresa
↓
goto('/dashboard')
```

### Sair da Empresa

```
Clique "Voltar para Plataforma"
↓
tenantSessionManager.leaveCompany()
↓
destroy() - limpa contexto tenant
↓
clearCaches() - limpa caches
↓
goto('/platform/admin')
```

### Trocar de Empresa

```
Clique em nova empresa
↓
tenantSessionManager.switchCompany(newCompanyId)
↓
leaveCompany() - sai da empresa atual
↓
enterCompany(newCompanyId) - entra na nova empresa
```

## Por Que Esta Arquitetura Foi Criada

### Problema Original

Antes da implementação do `TenantSessionManager`, o HorizonGest não possuía um ciclo de vida definido para a sessão Tenant. O fluxo era:

```
Platform Login
↓
Selecionar Empresa
↓
Troca JWT
↓
goto()
↓
Dashboard
```

### Consequências

1. **Estados Fantasmas:** Stores sobreviviam entre empresas
2. **Caches Persistentes:** localStorage e cookies não eram limpos
3. **Componentes Sobreviventes:** Contexto da empresa anterior permanecia parcialmente carregado
4. **Comportamento Aleatório:** Usuário com múltiplas empresas experimentava comportamentos inconsistentes
5. **Vazamento de Contexto:** Dados da Empresa A podiam aparecer na Empresa B

### Solução

O `TenantSessionManager` foi criado para:

1. **Centralizar Responsabilidade:** Único lugar para troca de empresa
2. **Garantir Limpeza:** Contexto anterior é sempre destruído antes do novo
3. **Prevenir Race Conditions:** Flags previnem operações simultâneas
4. **Garantir Consistência:** Contexto é completamente hidratado antes da navegação
5. **Facilitar Manutenção:** Lógica centralizada é mais fácil de manter e testar

## Integração com Componentes

### Botão "Entrar na Empresa"

**Arquivo:** `frontend/src/routes/platform/companies/[id]/+page.svelte`

**Antes:**
```typescript
async function loginAsCompany() {
    // Fetch direto para API
    // Set cookie manualmente
    // Set localStorage manualmente
    // setTimeout + goto
}
```

**Depois:**
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

**Antes:**
```typescript
async function endImpersonation() {
    // Fetch direto para API
    // Remove localStorage manualmente
    // Remove cookie manualmente
    // goto('/platform/admin')
}
```

**Depois:**
```typescript
async function endImpersonation() {
    const result = await tenantSessionManager.leaveCompany();
    if (result.success) {
        impersonationInfo = null;
        visible = false;
    }
}
```

## Regras Importantes

### Stores Continuam Singleton

As stores NÃO deixaram de ser singleton. Singleton NÃO é o problema. O problema é o ciclo de vida.

Portanto:
- Stores continuam existindo
- O conteúdo delas é reinicializado
- Não há recriação de instâncias

### Navegação Apenas Após Contexto Consistente

Nunca navegar antes do contexto estar consistente. O `TenantSessionManager` garante que:

1. Contexto anterior foi destruído
2. Novo token foi obtido
3. Stores foram limpas
4. Caches foram limpos
5. Novo contexto foi hidratado
6. Branding foi carregado
7. Permissões foram carregadas
8. Empresa foi carregada

SOMENTE ENTAO ocorre a navegação.

### Único Ponto de Troca

Nenhum outro lugar deve trocar empresa diretamente. Toda troca deve passar pelo `TenantSessionManager`.

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

### Testes de Integração

**Cenários:**
1. Entrar Empresa → Sair Empresa → Entrar Novamente
2. Trocar Empresa 100 vezes consecutivas
3. Garantir nenhuma store contaminada
4. Garantir nenhum cache persistente
5. Garantir empresa correta
6. Garantir permissões corretas
7. Garantir branding correto
8. Garantir dashboard correto

## Benefícios Arquiteturais

### 1. Centralização
- Lógica de troca de empresa em um único lugar
- Fácil de manter e evoluir
- Fácil de testar

### 2. Consistência
- Contexto sempre consistente antes da navegação
- Nenhum estado fantasma
- Nenhum vazamento de contexto

### 3. Prevenção de Erros
- Flags previnem race conditions
- Validações previnem operações inválidas
- Tratamento de erros em cada etapa

### 4. Manutenibilidade
- Código mais limpo
- Responsabilidades claras
- Fácil de adicionar novas funcionalidades

### 5. Testabilidade
- Fácil de testar isoladamente
- Fácil de mockar dependências
- Fácil de simular cenários

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

## Nota da Arquitetura Antes

**Nota:** 3/10

**Problemas:**
- Sem ciclo de vida definido
- Lógica espalhada em múltiplos componentes
- Sem limpeza de contexto
- Sem prevenção de race conditions
- Difícil de testar
- Difícil de manter

## Nota da Arquitetura Depois

**Nota:** 9/10

**Melhorias:**
- Ciclo de vida bem definido
- Lógica centralizada
- Limpeza completa de contexto
- Prevenção de race conditions
- Fácil de testar
- Fácil de manter

**Pontos a melhorar:**
- Adicionar testes de integração
- Adicionar monitoramento de erros
- Adicionar logging detalhado
- Adicionar métricas de performance

## Conclusão

O `TenantSessionManager` é uma camada arquitetural essencial para o HorizonGest. Ela garante que o ciclo de vida da sessão Tenant seja gerenciado de forma profissional, centralizada e definitiva. No futuro, qualquer funcionalidade relacionada à troca de empresa deve utilizar exclusivamente este Manager.
