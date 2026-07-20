# Sprint 3 - Refatoração Arquitetural Multi-Tenant (Modelo SaaS Empresarial)

**Data:** 19/07/2026  
**Versão:** 3.0  
**Status:** Planejamento

---

## Resumo Executivo

Esta refatoração transforma o PratoOnline de um modelo self-service (onde usuários se registram e automaticamente recebem uma empresa) para um modelo SaaS empresarial gerenciado pela plataforma (onde apenas administradores da plataforma criam empresas e atribuem usuários).

**Mudança Fundamental:** Eliminação completa do cadastro público e implementação de arquitetura de dois níveis (Platform vs Company).

---

## Arquitetura Antiga (Pré-Sprint 3)

### Modelo Self-Service

```
┌─────────────────────────────────────────────────────────────┐
│                    Público                                  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  POST /auth/register   │
              └────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Auto-cria Empresa     │
              │  Auto-atribui Owner    │
              └────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Empresa                                   │
│  - Usuário pertence a 1 empresa (ou nenhuma)               │
│  - CompanyID pode ser NULL                                  │
│  - Roles: owner, admin, manager, cashier, kitchen, waiter   │
└─────────────────────────────────────────────────────────────┘
```

### Problemas Identificados

1. **Bloqueador Crítico:** Usuários recém-cadastrados automaticamente recebem empresa, impedindo aceitação de convites
2. **Modelo de Negócio Incorreto:** Plataforma não tem controle sobre criação de empresas
3. **Isolamento Insuficiente:** CompanyID NULL permite usuários "órfãos"
4. **Sem Gestão Plataforma:** Não existe nível administrativo para gerenciar empresas
5. **Convites Quebrados:** Fluxo de convites inoperacional devido ao auto-creation

### Componentes Afetados

**Backend:**
- `internal/domain/role.go` - Roles atuais (owner, admin, manager, cashier, kitchen, waiter)
- `internal/domain/user.go` - CompanyID nullable, Role nullable
- `internal/domain/tenant_context.go` - Suporte a CompanyID NULL
- `internal/service/auth_service.go` - Register() cria empresa automaticamente (linhas 83-119)
- `internal/service/invitation_service.go` - Validação "usuário já pertence a outra empresa"
- `internal/service/user_management_service.go` - AddExistingUser() requer CompanyID NULL
- `internal/middleware/tenant_middleware.go` - Aceita CompanyID NULL
- `cmd/server/main.go` - Rota POST /api/auth/register pública

**Frontend:**
- `frontend/src/routes/(auth)/register/+page.svelte` - Tela de cadastro público
- `frontend/src/routes/(auth)/login/+page.svelte` - Login sem distinção de nível

---

## Arquitetura Nova (Pós-Sprint 3)

### Modelo SaaS Empresarial de Dois Níveis

```
┌─────────────────────────────────────────────────────────────┐
│              NÍVEL 1: PLATAFORMA                            │
│  Acessível apenas por: PlatformAdmin, PlatformSupport       │
│  Rotas: /platform/*                                         │
│  Responsabilidades:                                          │
│  - Criar empresas                                            │
│  - Editar empresas                                           │
│  - Suspender empresas                                        │
│  - Excluir empresas                                          │
│  - Criar Owner inicial                                       │
│  - Gerenciar usuários da plataforma                         │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  POST /platform/      │
              │  companies           │
              └────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Empresa criada       │
              │  (sem usuário ainda)   │
              └────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  POST /platform/      │
              │  companies/:id/owner │
              └────────────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  Owner criado         │
              │  (CompanyID definido) │
              └────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│              NÍVEL 2: EMPRESA                                │
│  Acessível por: Owner, Admin, Manager, Employee            │
│  Rotas: /api/*                                              │
│  Responsabilidades:                                          │
│  - Gerenciar produtos                                        │
│  - Gerenciar pedidos                                         │
│  - Gerenciar ingredientes                                    │
│  - Gerenciar usuários da empresa                            │
│  - Enviar convites (opcional)                                │
└─────────────────────────────────────────────────────────────┘
```

### Novos Roles

**Nível Platform:**
- `PlatformAdmin` - Acesso total à plataforma, pode gerenciar empresas
- `PlatformSupport` - Acesso de suporte, visualização limitada

**Nível Company:**
- `Owner` - Acesso total à empresa, pode gerenciar usuários
- `Admin` - Acesso quase total, não pode alterar Owner
- `Manager` - Pode gerenciar pedidos, produtos, ver relatórios
- `Employee` - Acesso operacional limitado

### Regras de Negócio

**REGRA 1 - Eliminar Cadastro Público:**
- `POST /api/auth/register` removido ou retorna 403
- Mensagem: "Cadastro indisponível. Solicite acesso ao administrador da plataforma."

**REGRA 2 - Criar Módulo Super Admin:**
- Novo domínio: `Platform` ou `Administration`
- Entidades: `PlatformUser`, `Company` (platform-managed)
- Rotas: `/platform/*` com autenticação própria

**REGRA 3 - Criar Company Manualmente:**
- `POST /platform/companies`
- Payload: Nome, CNPJ, Email, Telefone, Plano, Status
- Não cria usuário automaticamente

**REGRA 4 - Criar Owner:**
- `POST /platform/companies/:id/owner`
- Payload: Nome, Email, Senha
- Cria usuário com RoleOwner e CompanyID correto
- Nunca cria empresa automaticamente

**REGRA 5 - Usuários Internos:**
- `POST /api/company/users` continua existindo
- CompanyID vem SEMPRE de `ctx.User.CompanyID`
- Nunca do frontend, payload ou banco

**REGRA 6 - Convites Opcionais:**
- Fluxo A: Owner cria usuário diretamente (Nome, Email, Senha, Role)
- Fluxo B: Owner envia convite (Usuário define senha, ativa conta)
- Ambos pertencem automaticamente à mesma empresa

**REGRA 7 - Usuário Sem Empresa:**
- Validação: Todo User deve possuir `CompanyID != NULL`
- Caso contrário: erro de validação

**REGRA 8 - Roles Distintos:**
- `PlatformAdmin` acessa `/platform/*`
- `Owner` acessa `/api/*` somente da própria empresa
- São níveis completamente distintos

**REGRA 9 - Remover Lógica Antiga:**
- Eliminar completamente:
  - Criar empresa automaticamente
  - CompanyID NULL
  - Owner sem empresa
  - Registro público
- Todo código morto deve ser removido

**REGRA 10 - Dashboard Plataforma:**
- Tela `/platform` com:
  - Listagem de empresas
  - Criar empresa
  - Editar empresa
  - Suspender empresa
  - Excluir empresa
  - Entrar como empresa (impersonation futuramente)

**REGRA 11 - Refatorar Middleware:**
- `PlatformAdmin` ≠ `Owner`
- Níveis completamente distintos
- Middleware deve distinguir entre os dois

---

## Decisão de Schema de Banco (REGRA 12)

### Opção Escolhida: Tabela Separada

**Decisão:** Criar tabela `platform_users` separada da tabela `users`.

**Justificativa:**

1. **Separação de Concerns:** Platform users e Company users têm propósitos completamente diferentes
2. **Campos Diferentes:** Platform users não precisam de CompanyID, Company users sempre precisam
3. **Segurança:** Separação física reduz risco de confusão entre níveis
4. **Performance:** Queries de plataforma não precisam JOIN com users
5. **Escalabilidade:** Permite diferentes índices e otimizações para cada tipo
6. **Manutenibilidade:** Código mais limpo, menos condicionais

### Schema Proposto

```sql
-- Tabela existente (modificada)
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT 1,
  company_id INTEGER NOT NULL,  -- AGORA OBRIGATÓRIO (NOT NULL)
  role TEXT NOT NULL,            -- AGORA OBRIGATÓRIO (NOT NULL)
  deleted_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  FOREIGN KEY (company_id) REFERENCES companies(id)
);

-- Nova tabela
CREATE TABLE platform_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT 1,
  role TEXT NOT NULL,              -- 'PlatformAdmin' ou 'PlatformSupport'
  deleted_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

-- Tabela existente (sem mudanças estruturais)
CREATE TABLE companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  description TEXT,
  active BOOLEAN NOT NULL DEFAULT 1,
  logo_url TEXT,
  primary_color TEXT,
  secondary_color TEXT,
  business_type TEXT,
  locale TEXT,
  currency TEXT,
  timezone TEXT,
  deleted_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### Migration

```sql
-- Step 1: Criar tabela platform_users
CREATE TABLE platform_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT 1,
  role TEXT NOT NULL,
  deleted_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

-- Step 2: Alterar tabela users para CompanyID NOT NULL
-- Primeiro, migrar dados existentes com CompanyID NULL para uma empresa temporária
-- (ou eliminar se não houver dados em produção)
-- Depois, alterar a constraint:
ALTER TABLE users ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

-- Step 3: Criar platform admin inicial
INSERT INTO platform_users (name, email, password_hash, role, created_at, updated_at)
VALUES ('Platform Admin', 'admin@pratoonline.com', '$2a$10$...', 'PlatformAdmin', datetime('now'), datetime('now'));
```

---

## Impactos por Componente

### Backend

#### Domain Layer

**Arquivos a Criar:**
- `internal/domain/platform_user.go` - Nova entidade
- `internal/domain/platform_company.go` - Company com campos de plataforma

**Arquivos a Modificar:**
- `internal/domain/role.go` - Adicionar PlatformAdmin, PlatformSupport
- `internal/domain/user.go` - Remover nullable de CompanyID e Role
- `internal/domain/tenant_context.go` - Remover suporte a CompanyID NULL

#### Service Layer

**Arquivos a Criar:**
- `internal/service/platform_service.go` - Gestão de empresas pela plataforma
- `internal/service/platform_auth_service.go` - Autenticação de platform users
- `internal/service/platform_user_service.go` - Gestão de platform users

**Arquivos a Modificar:**
- `internal/service/auth_service.go` - Remover Register() ou limitar a platform users
- `internal/service/invitation_service.go` - Remover validação "usuário já pertence a outra empresa"
- `internal/service/user_management_service.go` - Remover AddExistingUser() (CompanyID NULL não existe mais)

#### Handler Layer

**Arquivos a Criar:**
- `internal/handler/platform_handler.go` - Endpoints de plataforma
- `internal/handler/platform_auth_handler.go` - Auth de plataforma
- `internal/handler/platform_company_handler.go` - Gestão de empresas

**Arquivos a Modificar:**
- `internal/handler/auth_handler.go` - Remover Register() ou retornar 403
- `internal/handler/invitation_handler.go` - Simplificar validação

#### Middleware Layer

**Arquivos a Criar:**
- `internal/middleware/platform_auth_middleware.go` - Auth específico para plataforma
- `internal/middleware/platform_role_middleware.go` - RBAC para plataforma

**Arquivos a Modificar:**
- `internal/middleware/auth_middleware.go` - Suportar ambos os tipos de usuários
- `internal/middleware/tenant_middleware.go` - Remover suporte a CompanyID NULL
- `internal/middleware/role_middleware.go` - Distinguir PlatformAdmin de Owner

#### Repository Layer

**Arquivos a Criar:**
- `internal/infra/repository/gorm_platform_user_repository.go`
- `internal/ports/platform_user_repository.go`

#### Router

**Arquivos a Modificar:**
- `cmd/server/main.go` - Adicionar rotas `/platform/*`

### Frontend

#### Arquivos a Remover
- `frontend/src/routes/(auth)/register/+page.svelte` - Tela de cadastro público
- `frontend/src/routes/(auth)/register/+page.server.ts` - Server actions de cadastro

#### Arquivos a Criar
- `frontend/src/routes/(platform)/+page.svelte` - Dashboard plataforma
- `frontend/src/routes/(platform)/companies/+page.svelte` - Listagem empresas
- `frontend/src/routes/(platform)/companies/create/+page.svelte` - Criar empresa
- `frontend/src/routes/(platform)/companies/[id]/+page.svelte` - Detalhes empresa
- `frontend/src/routes/(platform)/companies/[id]/owner/+page.svelte` - Criar owner
- `frontend/src/routes/(platform)/auth/+page.svelte` - Login plataforma
- `frontend/src/lib/api/platform_client.ts` - Client API plataforma

#### Arquivos a Modificar
- `frontend/src/routes/(auth)/login/+page.svelte` - Adicionar opção de login plataforma
- `frontend/src/lib/api/client.ts` - Adicionar endpoints de plataforma

---

## Fluxos de Usuário

### Fluxo 1: Plataforma Cria Empresa

```
1. PlatformAdmin faz login em /platform/auth
2. PlatformAdmin acessa /platform/companies
3. PlatformAdmin clica em "Criar Empresa"
4. PlatformAdmin preenche: Nome, CNPJ, Email, Telefone, Plano, Status
5. POST /platform/companies
6. Empresa criada (sem usuário)
7. PlatformAdmin clica em "Criar Owner"
8. PlatformAdmin preenche: Nome, Email, Senha
9. POST /platform/companies/:id/owner
10. Owner criado com CompanyID definido
11. Owner recebe email com credenciais
12. Owner faz login em /api/auth/login
13. Owner acessa dashboard da empresa
```

### Fluxo 2: Owner Cria Usuário Diretamente

```
1. Owner faz login em /api/auth/login
2. Owner acessa /api/company/users
3. Owner clica em "Criar Usuário"
4. Owner preenche: Nome, Email, Senha, Role
5. POST /api/company/users/add
6. Usuário criado com CompanyID do Owner
7. Usuário recebe email com credenciais
8. Usuário faz login em /api/auth/login
9. Usuário acessa dashboard da empresa
```

### Fluxo 3: Owner Envia Convite (Opcional)

```
1. Owner faz login em /api/auth/login
2. Owner acessa /api/company/invitations
3. Owner clica em "Enviar Convite"
4. Owner preenche: Email, Role
5. POST /api/company/invitations
6. Convite criado com token
7. Usuário recebe email com link
8. Usuário clica no link (não precisa estar cadastrado)
9. Usuário define senha
10. POST /api/invitations/accept
11. Usuário criado com CompanyID da empresa
12. Usuário faz login em /api/auth/login
```

---

## Validações e Constraints

### Validações de Negócio

1. **CompanyID Obrigatório:** Todo User deve ter CompanyID != NULL
2. **Role Obrigatório:** Todo User deve ter Role != NULL
3. **Platform Users Sem Empresa:** Platform users nunca têm CompanyID
4. **Company Users Sem Plataforma:** Company users nunca têm acesso a /platform/*
5. **Um Usuário Uma Empresa:** Um usuário pertence a exatamente uma empresa
6. **Owner Único:** Cada empresa tem exatamente um Owner
7. **Platform Admin Único:** Plataforma tem pelo menos um PlatformAdmin

### Validações Técnicas

1. **Email Único:** Email único em cada tabela (users, platform_users)
2. **Slug Único:** Slug único em companies
3. **Foreign Keys:** CompanyID em users referencia companies.id
4. **Soft Delete:** DeletedAt em users, platform_users, companies
5. **Active Status:** Active boolean em users, platform_users, companies

---

## Migração de Dados

### Estratégia de Migração

**Cenário 1: Desenvolvimento (Sem dados de produção)**
1. Drop tabela users
2. Recriar users com CompanyID NOT NULL
3. Criar platform_users
4. Criar platform admin inicial
5. Criar empresa de teste
6. Criar owner inicial

**Cenário 2: Produção (Com dados existentes)**
1. Backup completo do banco
2. Criar tabela platform_users
3. Migrar platform users existentes (se houver)
4. Para users com CompanyID NULL:
   - Criar empresas temporárias
   - Atribuir CompanyID
   - Notificar usuários
5. Alterar constraint para NOT NULL
6. Validar integridade
7. Rollback se houver problemas

---

## Testes Necessários

### Testes de Integração

1. **Platform Admin:**
   - Login plataforma
   - Criar empresa
   - Editar empresa
   - Suspender empresa
   - Excluir empresa
   - Criar owner

2. **Owner:**
   - Login empresa
   - Criar usuário diretamente
   - Enviar convite
   - Alterar role
   - Desativar usuário
   - Remover usuário

3. **Admin/Manager/Employee:**
   - Login empresa
   - Acessar recursos permitidos
   - Bloqueio de recursos não permitidos

4. **Convites:**
   - Criar convite
   - Aceitar convite (usuário não cadastrado)
   - Aceitar convite (usuário cadastrado sem empresa)
   - Revogar convite

### Testes de Regressão

1. **Produtos:** CRUD, SEO, Slug, Soft Delete
2. **Pedidos:** Criação, Status, Validação de estoque
3. **Ingredientes:** CRUD, Estoque
4. **Tema:** Cores, Persistência
5. **RBAC:** Permissões por role
6. **JWT:** Geração, Validação, Blacklist
7. **Soft Delete:** Todos os recursos

---

## Riscos e Mitigações

### Risco 1: Perda de Dados na Migração

**Mitigação:**
- Backup completo antes da migração
- Teste de migração em ambiente de staging
- Script de rollback preparado
- Validação de integridade pós-migração

### Risco 2: Quebra de Funcionalidades Existentes

**Mitigação:**
- Testes de regressão completos
- Testes manuais de todos os fluxos
- Monitoramento pós-deploy
- Hotfix rápido se necessário

### Risco 3: Usuários Existentes Sem Empresa

**Mitigação:**
- Identificar todos os users com CompanyID NULL
- Criar empresas temporárias ou notificar usuários
- Comunicação clara sobre mudança
- Período de transição

### Risco 4: Complexidade de Autenticação

**Mitigação:**
- Middleware bem testado
- Documentação clara de níveis
- Logs detalhados de autenticação
- Monitoramento de erros

---

## Cronograma Sugerido

### Fase 1: Planejamento (1 dia)
- [x] Analisar arquitetura atual
- [x] Desenhar arquitetura nova
- [x] Decidir schema de banco
- [x] Documentar mudanças

### Fase 2: Backend (5-7 dias)
- [ ] Criar domain entities
- [ ] Criar services de plataforma
- [ ] Criar handlers de plataforma
- [ ] Criar middlewares de plataforma
- [ ] Refatorar auth service
- [ ] Refatorar tenant middleware
- [ ] Criar migration
- [ ] Atualizar router

### Fase 3: Frontend (3-5 dias)
- [ ] Remover telas de cadastro
- [ ] Criar dashboard plataforma
- [ ] Criar telas de gestão de empresas
- [ ] Atualizar login
- [ ] Atualizar API client

### Fase 4: Testes (2-3 dias)
- [ ] Testes de integração
- [ ] Testes de regressão
- [ ] Testes manuais
- [ ] Correção de bugs

### Fase 5: Deploy (1 dia)
- [ ] Backup produção
- [ ] Executar migration
- [ ] Deploy backend
- [ ] Deploy frontend
- [ ] Monitoramento
- [ ] Validação

**Total Estimado:** 12-17 dias

---

## Conclusão

Esta refatoração é fundamental para transformar o PratoOnline em um SaaS empresarial robusto e escalável. A separação clara entre nível de plataforma e nível de empresa elimina os problemas atuais de convites e provides controle total sobre a criação de empresas.

A decisão de usar tabelas separadas para platform users e company users garante separação de concerns, segurança e performance, facilitando manutenção futura.

A remoção completa do cadastro público e a obrigatoriedade de CompanyID eliminam ambiguidades e simplificam o modelo de dados.
