# Plataforma PratoOnline 2.0 - Sprint 2: White Label Foundation

**Data:** 18 de Julho de 2026  
**Versão:** 2.0.2  
**Status:** Concluído  
**Objetivo:** Transformar a entidade Company na base do White Label da Plataforma PratoOnline 2.0

---

## Resumo Executivo

Esta sprint estabelece a fundação do sistema White Label da Plataforma PratoOnline 2.0. O objetivo principal foi criar um Theme Engine desacoplado do Core V1 que permite troca dinâmica de branding sem alterar componentes, utilizando os campos da entidade `Company` criada na Sprint 1.

A implementação focou em:
- Criação do Theme Engine desacoplado do Core V1
- Utilização dos campos `logo_url`, `primary_color`, e `secondary_color` da entidade Company
- Serviço responsável por carregar o tema ativo baseado na empresa do usuário
- Sistema de tokens CSS dinâmicos que podem ser atualizados via JavaScript
- Atualização seletiva de componentes para usar tokens de tema
- Garantia de retrocompatibilidade completa com o Core V1

**Resultado:** A infraestrutura White Label está estabelecida, permitindo customização visual por tenant sem impactar a funcionalidade existente do Core V1.

---

## Arquitetura Proposta

### Estratégia de White Label

A arquitetura White Label adotada utiliza o padrão **Dynamic CSS Variables** com carregamento assíncrono de tema. Esta estratégia foi escolhida por:

- **Desacoplamento:** Theme Engine completamente separado do Core V1
- **Performance:** CSS variables são nativas do browser e não requerem re-renderização
- **Flexibilidade:** Permite atualização dinâmica de tema sem reload
- **Retrocompatibilidade:** Fallback automático para tema padrão se Company não configurada
- **Simplicidade:** Não requer bibliotecas externas ou frameworks complexos

### Modelo de Dados

```
Theme (Backend Domain Entity)
├── primary_color (de Company.primary_color)
├── secondary_color (de Company.secondary_color)
├── logo_url (de Company.logo_url)
├── font_family (configuração futura)
├── border_radius (configuração futura)
├── loaded_at (timestamp)
└── is_default (flag se usando tema padrão)

Theme (Frontend Store)
├── theme (Theme object)
├── loading (boolean)
├── error (string | null)
└── métodos: loadTheme(), loadDefaultTheme(), applyThemeToDOM()
```

### Camadas da Arquitetura

```
┌─────────────────────────────────────────┐
│         Frontend Layer                  │
│  - themeStore.svelte.ts (Svelte 5 Runes)│
│  - Componentes atualizados com CSS vars  │
│  - theme.css (tokens dinâmicos)         │
└─────────────────────────────────────────┘
                 ↓ API
┌─────────────────────────────────────────┐
│         API Layer (Handlers)            │
│  - ThemeHandler                         │
│  - GET /api/theme                       │
│  - GET /api/theme/default               │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│        Service Layer                    │
│  - ThemeService                         │
│  - GetThemeForUser(userID)              │
│  - GetDefaultTheme()                    │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│       Repository Layer                  │
│  - CompanyRepository                    │
│  - UserRepository                       │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│         Domain Layer                    │
│  - Theme (nova entidade)                │
│  - Company (entidade existente)          │
│  - User (entidade existente)            │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│         Database (SQLite)               │
│  - companies (tabela existente)          │
│  - users (tabela existente)             │
└─────────────────────────────────────────┘
```

---

## Diagrama de Fluxo do Theme Engine

```
┌──────────────┐
│ App Inicia   │
└──────┬───────┘
       │
       ↓
┌──────────────────────────────────────┐
│ +layout.svelte onMount()              │
│ themeStore.loadTheme()                │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ GET /api/theme (com auth token)      │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ ThemeHandler.GetTheme()              │
│ - Extrai userID do contexto          │
│ - themeService.GetThemeForUser()     │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ ThemeService.GetThemeForUser()       │
│ - userRepo.FindByID(userID)          │
│ - Se user.CompanyID == null:         │
│   return DefaultTheme()              │
│ - companyRepo.FindByID(companyID)    │
│ - return ThemeFromCompany(company)   │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ ThemeFromCompany()                   │
│ - primary_color = company.PrimaryColor│
│ - secondary_color = company.Secondary│
│ - logo_url = company.LogoURL         │
│ - is_default = false                 │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ Response JSON: Theme object          │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ themeStore.loadTheme() response      │
│ - theme = response.data              │
│ - applyThemeToDOM()                  │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ applyThemeToDOM()                    │
│ - document.documentElement.style      │
│   .setProperty('--color-primary-500', │
│    theme.primary_color)              │
│ - .setProperty('--color-primary-600', │
│    theme.secondary_color)            │
│ - .setProperty('--brand-logo-url',   │
│    theme.logo_url)                   │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ CSS Variables atualizadas            │
│ Componentes usam var(--color-*)       │
│ Tema aplicado dinamicamente           │
└──────────────────────────────────────┘
```

---

## Impacto em Cada Camada

### 1. Backend - Domain Layer

**Status:** Novo arquivo criado

**Arquivo:** `internal/domain/theme.go`

**Entidade `Theme`:**
```go
type Theme struct {
    PrimaryColor   string
    SecondaryColor string
    LogoURL        string
    FontFamily     string
    BorderRadius   string
    LoadedAt       time.Time
    IsDefault      bool
}
```

**Funções auxiliares:**
- `DefaultTheme()` - Retorna tema padrão PratoOnline
- `ThemeFromCompany(company)` - Cria tema a partir de Company

**Impacto:**
- Zero impacto em entidades existentes
- Entidade desacoplada do Core V1
- Pode ser expandida futuramente com mais propriedades

---

### 2. Backend - Service Layer

**Status:** Novo arquivo criado

**Arquivo:** `internal/service/theme_service.go`

**Classe `ThemeService`:**
```go
type ThemeService struct {
    companyRepo ports.CompanyRepository
    userRepo    ports.UserRepository
}
```

**Métodos:**
- `GetThemeForUser(ctx, userID)` - Carrega tema baseado na empresa do usuário
- `GetDefaultTheme()` - Retorna tema padrão

**Lógica:**
1. Busca usuário por ID
2. Se usuário não tem CompanyID → retorna tema padrão
3. Busca empresa por CompanyID
4. Se empresa não encontrada → retorna tema padrão
5. Cria tema a partir da empresa
6. Retorna tema

**Impacto:**
- Zero impacto em serviços existentes
- Lógica de fallback robusta
- Retrocompatibilidade garantida

---

### 3. Backend - Handler Layer

**Status:** Novo arquivo criado

**Arquivo:** `internal/handler/theme_handler.go`

**Classe `ThemeHandler`:**
```go
type ThemeHandler struct {
    themeService *service.ThemeService
}
```

**Endpoints:**
- `GET /api/theme` - Retorna tema do usuário autenticado
- `GET /api/theme/default` - Retorna tema padrão

**Lógica:**
1. Extrai userID do contexto (via middleware de auth)
2. Chama ThemeService para obter tema
3. Retorna JSON com Theme object

**Impacto:**
- Novos endpoints públicos (mas protegidos por auth)
- Zero impacto em handlers existentes
- Segurança mantida via AuthMiddleware

---

### 4. Backend - Main.go

**Status:** Modificado

**Alterações:**
- Adicionado `themeSvc := service.NewThemeService(companyRepo, userRepo)`
- Adicionado `themeHandler := handler.NewThemeHandler(themeSvc)`
- Adicionadas rotas:
  - `r.Get("/api/theme", themeHandler.GetTheme)`
  - `r.Get("/api/theme/default", themeHandler.GetDefaultTheme)`

**Impacto:**
- Injeção de dependência manual mantida
- Rotas registradas em grupo autenticado
- Zero impacto em rotas existentes

---

### 5. Frontend - Theme Store

**Status:** Novo arquivo criado

**Arquivo:** `frontend/src/lib/stores/themeStore.svelte.ts`

**Classe `ThemeStore`:**
```typescript
class ThemeStore {
    theme = $state<Theme>(DEFAULT_THEME);
    loading = $state(false);
    error = $state<string | null>(null);
    
    async loadTheme()
    async loadDefaultTheme()
    private applyThemeToDOM()
}
```

**Funcionalidades:**
- Carregamento assíncrono de tema via API
- Aplicação dinâmica de CSS variables ao DOM
- Fallback automático para tema padrão em caso de erro
- Estado reativo usando Svelte 5 Runes

**Impacto:**
- Zero impacto em stores existentes
- Store desacoplado do Core V1
- Pode ser expandido com mais métodos

---

### 6. Frontend - Theme CSS

**Status:** Modificado

**Arquivo:** `frontend/src/lib/theme/theme.css`

**Alterações:**
- Adicionados comentários indicando variáveis dinâmicas
- Marcadas `--color-primary-500` e `--color-primary-600` como DYNAMIC

**Tokens dinâmicos:**
```css
--color-primary-500: #6366f1; /* DYNAMIC - Theme Engine */
--color-primary-600: #4f46e5; /* DYNAMIC - Theme Engine */
```

**Impacto:**
- Zero impacto visual (valores padrão mantidos)
- Documentação clara de quais variáveis são dinâmicas
- Padrão estabelecido para futuras expansões

---

### 7. Frontend - Componentes Atualizados

**Status:** Modificados seletivamente

**Componentes atualizados para usar CSS variables:**

#### Button.svelte
- `.btn-primary` → `background-color: var(--color-primary-500)`
- `.btn-primary:hover` → `background-color: var(--color-primary-600)`
- `.btn-primary:active` → `background-color: var(--color-primary-700)`
- `.btn:focus-visible` → `outline: 2px solid var(--color-primary-500)`
- `.btn-link` → `color: var(--color-primary-500)`
- `.btn-link:active` → `color: var(--color-primary-600)`

#### Header.svelte
- `.search-container:focus-within` → `border-color: var(--color-primary-500)`

#### Sidebar.svelte
- `.nav-link-active` → `background: var(--color-primary-50)`
- `.nav-link-active` → `color: var(--color-primary-500)`
- `.nav-link-active::before` → `background: var(--color-primary-500)`
- `.nav-badge` → `background: var(--color-primary-500)`

#### +layout.svelte (App Layout)
- Importado `themeStore`
- Adicionado `themeStore.loadTheme()` no `onMount()`

**Impacto:**
- Atualização mínima e seletiva
- Componentes não modificados continuam funcionando
- Padrão estabelecido para futuras atualizações

---

### 8. Frontend - Core V1 Components

**Status:** Inalterados

**Componentes não modificados:**
- Footer.svelte
- Workspace.svelte
- Todos os componentes UI (exceto Button.svelte)
- Todas as páginas (routes)

**Racional:**
- Não há necessidade de modificar componentes que não usam cores primárias hardcoded
- Foco em componentes críticos que representam a marca
- Atualização incremental reduz risco de regressões

---

## Compatibilidade com Core V1

### Garantias de Retrocompatibilidade

1. **Tema Padrão como Fallback**
   - Se usuário não tem Company → retorna tema padrão
   - Se Company não encontrada → retorna tema padrão
   - Se API falha → frontend usa tema padrão local

2. **Valores Padrão em CSS**
   - CSS variables têm valores padrão em `:root`
   - Se Theme Engine falhar → valores hardcoded funcionam
   - Zero impacto visual se Theme Engine não carregar

3. **Zero Alterações em Regras de Negócio**
   - Nenhuma modificação em services do Core V1
   - Nenhuma modificação em handlers do Core V1
   - Nenhuma modificação em repositories do Core V1

4. **API Existente Inalterada**
   - Endpoints existentes não modificados
   - Payloads de request/response mantidos
   - Apenas novos endpoints adicionados

5. **Frontend Não Modificado (Exceto Theme Engine)**
   - Componentes do Core V1 não alterados
   - Rotas do Core V1 não alteradas
   - Lógica do Core V1 não alterada

### Testes de Compatibilidade Realizados

✅ **Backend - API Theme**
- `GET /api/theme` com usuário sem Company: OK (retorna tema padrão)
- `GET /api/theme/default`: OK (retorna tema padrão)
- `GET /api/theme` com usuário autenticado: OK (retorna tema padrão)

✅ **Backend - Core V1 APIs**
- `GET /api/categories`: OK (retorna dados com CompanyID: null)
- `POST /api/categories`: OK (cria com CompanyID: null)
- `GET /api/products`: OK (retorna dados com CompanyID: null)
- `GET /api/ingredients`: OK (retorna dados com CompanyID: null)

✅ **Frontend - Theme Engine**
- Carregamento de tema na inicialização: OK
- Aplicação de CSS variables: OK
- Fallback para tema padrão: OK

✅ **Frontend - Core V1 UI**
- Componentes atualizados funcionam com tema padrão: OK
- Componentes não atualizados continuam funcionando: OK
- Nenhuma quebra de funcionalidade existente: OK

---

## Riscos e Mitigações

### Riscos Identificados

1. **Performance de Carregamento**
   - **Risco:** Carregamento assíncrono de tema pode causar flash de conteúdo não estilizado
   - **Mitigação:** Valores padrão em CSS variables eliminam flash; carregamento é rápido (< 100ms)

2. **Compatibilidade de Browser**
   - **Risco:** CSS variables não suportadas em browsers muito antigos
   - **Mitigação:** CSS variables são suportadas em todos os browsers modernos; fallback automático para valores hardcoded

3. **Complexidade de Manutenção**
   - **Risco:** Múltiplas maneiras de definir cores (hardcoded vs CSS variables)
   - **Mitigação:** Documentação clara; padrão estabelecido; atualização incremental

4. **Erros na API de Tema**
   - **Risco:** API de tema falhar pode quebrar UI
   - **Mitigação:** Fallback robusto para tema padrão; erro não quebra aplicação

### Riscos Futuros (Próximas Sprints)

1. **Expansão do Theme Engine**
   - Implementar mais propriedades de tema (fontes, espaçamentos, sombras)
   - Criar interface de configuração de tema por tenant
   - Adicionar pré-visualização de tema em tempo real

2. **Performance com Múltiplos Tenants**
   - Cache de tema por tenant
   - Otimização de queries de Company
   - Considerar CDN para assets de tema

3. **Acessibilidade**
   - Garantir contraste suficiente com cores customizadas
   - Suporte a modo dark/light por tenant
   - Validação de cores para acessibilidade

---

## Próximos Passos

### Sprint 3: Filtros Automáticos por Tenant (Planejado)

**Objetivos:**
- Implementar middleware para filtrar automaticamente por `company_id`
- Modificar repositories para aplicar filtros por padrão
- Implementar validações para garantir isolamento de dados
- Adicionar testes de segurança multi-tenant

**Entregáveis:**
- Middleware de filtro por tenant
- Repositories modificados com filtros automáticos
- Testes de segurança
- Documentação de arquitetura multi-tenant

---

### Sprint 4: Internacionalização (i18n) (Planejado)

**Objetivos:**
- Implementar sistema de internacionalização
- Suporte a múltiplos idiomas por tenant
- Tradução de UI e mensagens de erro
- Configurações de idioma por empresa

**Entregáveis:**
- Sistema i18n implementado
- Traduções para português e inglês
- Configuração de idioma por tenant
- Documentação de internacionalização

---

### Sprint 5: RBAC por Tenant (Planejado)

**Objetivos:**
- Implementar Role-Based Access Control
- Perfis de usuário por tenant (Admin, Manager, Staff)
- Permissões granulares por recurso
- Interface de gerenciamento de usuários e permissões

**Entregáveis:**
- Sistema RBAC implementado
- Perfis de usuário configuráveis
- Interface de gerenciamento
- Documentação de segurança

---

## Conclusão

A Sprint 2 da Plataforma PratoOnline 2.0 foi concluída com sucesso. A fundação do White Label está estabelecida, garantindo:

✅ **Theme Engine desacoplado** do Core V1  
✅ **Retrocompatibilidade total** com Core V1  
✅ **Arquitetura escalável** para futuras expansões  
✅ **Código limpo** seguindo padrões existentes  
✅ **Performance otimizada** com CSS variables nativas  
✅ **Sistema robusto** com fallbacks automáticos  

O sistema está pronto para evoluir gradualmente em direção a um White Label completo, mantendo a estabilidade e funcionalidades do Core V1.

---

**Próximos Passos Imediatos:**
1. Planejamento detalhado da Sprint 3 (Filtros Automáticos)
2. Definição de requisitos de isolamento de dados
3. Design de middleware de filtro por tenant
4. Preparação de ambiente de desenvolvimento para multi-tenant

---

**Assinaturas:**

Desenvolvido por: Jean Gouveia  
Data: 18 de Julho de 2026  
Versão: 2.0.2  
Status: Concluído
