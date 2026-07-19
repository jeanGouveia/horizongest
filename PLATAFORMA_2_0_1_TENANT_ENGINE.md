# Plataforma PratoOnline 2.0 - Sprint 1: Tenant Engine

**Data:** 18 de Julho de 2026  
**Versão:** 2.0.1  
**Status:** Concluído  
**Objetivo:** Fundação do Tenant Engine para suporte multi-tenant futuro

---

## Resumo Executivo

Esta sprint marca o início oficial da Plataforma PratoOnline 2.0. O objetivo principal foi criar a infraestrutura necessária para suportar múltiplas empresas (tenants) no futuro, sem implementar o multi-tenant completo nesta sprint. A implementação focou em:

- Criação da entidade `Company` (Tenant)
- Adição de relacionamentos entre Company e todas as entidades do domínio
- Migrações de banco de dados para suportar a nova estrutura
- Implementação completa das camadas de Repository, Service e Handler para Company
- Garantia de retrocompatibilidade total com o Core V1

**Resultado:** A infraestrutura está pronta para futuras implementações de isolamento de dados, filtros automáticos por tenant e funcionalidades multi-tenant completas.

---

## Arquitetura Proposta

### Estratégia de Isolamento

A arquitetura multi-tenant adotada utiliza o padrão **Shared Database, Shared Schema** com isolamento por chave estrangeira (`company_id`). Esta estratégia foi escolhida por:

- **Eficiência de recursos:** Único banco de dados para todos os tenants
- **Simplicidade de migração:** Adição de colunas nullable mantém compatibilidade
- **Escalabilidade:** Índices em `company_id` garantem performance
- **Flexibilidade:** Permite evolução gradual para isolamento mais rigoroso se necessário

### Modelo de Dados

```
companies (nova tabela)
├── id (PK)
├── name
├── slug (unique)
├── description
├── active
├── logo_url
├── primary_color
├── secondary_color
├── deleted_at
├── created_at
└── updated_at

Entidades Tenant-Scoped (com company_id FK):
├── users.company_id (nullable)
├── categories.company_id (nullable)
├── ingredients.company_id (nullable)
├── products.company_id (nullable)
├── orders.company_id (nullable)
├── stock_adjustments_pending.company_id (nullable)
└── media.company_id (nullable)
```

### Camadas da Arquitetura

```
┌─────────────────────────────────────────┐
│         API Layer (Handlers)            │
│  - CompanyHandler                       │
│  - AuthHandler (modificado)             │
│  - Outros handlers (inalterados)        │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│        Service Layer                    │
│  - CompanyService                       │
│  - AuthService (inalterado)             │
│  - Outros services (inalterados)         │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│       Repository Layer                  │
│  - GormCompanyRepository                │
│  - Outros repositories (modificados)    │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│         Domain Layer                    │
│  - Company (nova entidade)              │
│  - User, Category, Product, etc.        │
│    (todos com CompanyID adicionado)      │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│         Database (SQLite)               │
│  - companies (nova tabela)              │
│  - users, categories, products, etc.     │
│    (todos com company_id adicionado)      │
└─────────────────────────────────────────┘
```

---

## Diagrama de Relacionamentos

```
┌──────────────────┐
│    Company       │
├──────────────────┤
│ id (PK)          │
│ name             │
│ slug (unique)    │
│ description      │
│ active           │
│ logo_url         │
│ primary_color    │
│ secondary_color  │
│ deleted_at       │
│ created_at       │
│ updated_at       │
└──────────────────┘
         │
         │ 1:N (nullable)
         ├────────────────────────────────────┐
         │                                     │
    ┌────▼─────┐  ┌──────────┐  ┌──────────┐
    │  User    │  │ Category │  │Ingredient│
    ├──────────┤  ├──────────┤  ├──────────┤
    │ company_id│  │company_id│  │company_id│
    └──────────┘  └──────────┘  └──────────┘
         │                                     │
         │                                     │
    ┌────▼─────┐  ┌──────────┐  ┌──────────┐
    │  Product │  │  Order   │  │  Media   │
    ├──────────┤  ├──────────┤  ├──────────┤
    │company_id│  │company_id│  │company_id│
    └──────────┘  └──────────┘  └──────────┘
                                    │
                            ┌───────▼───────┐
                            │StockAdjustment│
                            │  Pending      │
                            ├───────────────┤
                            │  company_id   │
                            └───────────────┘
```

---

## Impacto em Cada Entidade

### 1. Company (Nova Entidade)

**Status:** Criada nesta sprint

**Campos:**
- `id`: Identificador único
- `name`: Nome da empresa
- `slug`: Identificador único para URLs
- `description`: Descrição da empresa
- `active`: Status de ativação
- `logo_url`: URL do logo
- `primary_color`: Cor primária para theming
- `secondary_color`: Cor secundária para theming
- `deleted_at`: Soft delete
- `created_at`: Timestamp de criação
- `updated_at`: Timestamp de atualização

**API Endpoints:**
- `POST /api/companies` - Criar empresa
- `GET /api/companies` - Listar empresas
- `GET /api/companies/{id}` - Obter empresa por ID
- `PUT /api/companies/{id}` - Atualizar empresa
- `DELETE /api/companies/{id}` - Remover empresa (soft delete)

---

### 2. User

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: usuários existentes têm `CompanyID = null`
- Futuramente: usuários serão associados a uma empresa específica
- Autenticação: não modificada nesta sprint

---

### 3. Category

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: categorias existentes têm `CompanyID = null`
- Futuramente: categorias serão isoladas por empresa

---

### 4. Ingredient

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: ingredientes existentes têm `CompanyID = null`
- Futuramente: estoque de ingredientes será isolado por empresa

---

### 5. Product

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: produtos existentes têm `CompanyID = null`
- Futuramente: catálogo de produtos será isolado por empresa

---

### 6. Order

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: pedidos existentes têm `CompanyID = null`
- Futuramente: pedidos serão isolados por empresa

---

### 7. StockAdjustmentPending

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: ajustes pendentes existentes têm `CompanyID = null`
- Futuramente: ajustes de estoque serão isolados por empresa

---

### 8. Media

**Status:** Modificado nesta sprint

**Alterações:**
- Adicionado campo `CompanyID *uint` (nullable)

**Impacto:**
- Retrocompatível: mídias existentes têm `CompanyID = null`
- Futuramente: arquivos de mídia serão isolados por empresa

---

### 9. Dashboard

**Status:** Não modificado

**Racional:**
- Dashboard é uma view agregada, não uma entidade persistente
- Futuramente: filtros por `company_id` serão aplicados nas queries

---

## Estratégia de Migração

### Migrações de Banco de Dados

**Migration 00008:** `create_companies_table.sql`
- Cria tabela `companies`
- Adiciona índices em `slug` e `active`
- Define valores padrão para cores

**Migration 00009:** `add_company_id_to_entities.sql`
- Adiciona `company_id` (nullable) a:
  - `users`
  - `categories`
  - `ingredients`
  - `products`
  - `orders`
  - `stock_adjustments_pending`
  - `media`
- Cria índices em todos os `company_id` para performance
- Foreign keys referenciando `companies(id)`

### Processo de Migração

1. **Backup do banco de dados** (recomendado antes da migração)
2. **Executar migration 00008** - Cria tabela companies
3. **Executar migration 00009** - Adiciona company_id às entidades
4. **Verificar dados existentes** - Confirmar que company_id = null para todos
5. **Testar funcionalidades Core V1** - Garantir retrocompatibilidade

### Rollback

As migrations incluem seções `Down` para rollback:
- Migration 00009: Remove colunas company_id e índices
- Migration 00008: Remove tabela companies

---

## Riscos e Mitigações

### Riscos Identificados

1. **Performance de Queries**
   - **Risco:** Queries sem filtro por `company_id` podem retornar dados de múltiplos tenants
   - **Mitigação:** Índices em `company_id`; futuramente implementar middleware de filtro automático

2. **Dados Órfãos**
   - **Risco:** Registros com `company_id` referenciando empresas deletadas
   - **Mitigação:** Soft delete em companies; constraints de foreign key; validações em service layer

3. **Complexidade de Migração**
   - **Risco:** Erros durante migração de bancos de dados em produção
   - **Mitigação:** Migrations testadas; rollback disponível; backup obrigatório

4. **Retrocompatibilidade**
   - **Risco:** Quebra de funcionalidades existentes do Core V1
   - **Mitigação:** Campos nullable; testes extensivos; nenhum filtro automático nesta sprint

### Riscos Futuros (Próximas Sprints)

1. **Isolamento de Dados**
   - Implementar middleware para filtrar automaticamente por `company_id`
   - Validar que usuários não acessem dados de outros tenants

2. **Segurança**
   - Implementar RBAC (Role-Based Access Control) por tenant
   - Garantir que admins de uma empresa não acessem dados de outra

3. **Escalabilidade**
   - Monitorar performance com grande volume de dados
   - Considerar particionamento de dados por tenant se necessário

---

## Compatibilidade com Core V1

### Garantias de Retrocompatibilidade

1. **Campos Nullable**
   - Todos os `company_id` são nullable
   - Dados existentes têm `company_id = null`
   - Nenhuma quebra de funcionalidade existente

2. **Sem Filtros Automáticos**
   - Nenhum filtro por `company_id` foi implementado nesta sprint
   - Queries continuam retornando todos os dados
   - Comportamento do Core V1 preservado

3. **API Inalterada**
   - Endpoints existentes não foram modificados
   - Payloads de request/response mantidos
   - Apenas novos campos (`company_id`) adicionados às respostas

4. **Frontend Não Modificado**
   - Nenhuma alteração necessária no frontend
   - Novos campos podem ser ignorados
   - Funcionalidades existentes continuam funcionando

### Testes de Compatibilidade Realizados

✅ **Autenticação**
- Login com usuário existente: OK
- Registro de novo usuário: OK
- Logout: OK

✅ **Categorias**
- Listagem de categorias: OK (retorna CompanyID: null)
- Criação de categoria: OK (CompanyID: null)
- Atualização de categoria: OK
- Remoção de categoria: OK

✅ **Produtos**
- Listagem de produtos: OK (retorna CompanyID: null)
- Criação de produto: OK (CompanyID: null)
- Atualização de produto: OK
- Remoção de produto: OK

✅ **Ingredientes**
- Listagem de ingredientes: OK (retorna CompanyID: null)
- Criação de ingrediente: OK (CompanyID: null)
- Atualização de ingrediente: OK
- Remoção de ingrediente: OK

✅ **Empresas (Novo)**
- Criação de empresa: OK
- Listagem de empresas: OK
- Consulta por ID: OK

---

## Plano das Próximas Sprints

### Sprint 2: White Label (Planejado)

**Objetivos:**
- Implementar sistema de theming por tenant
- Utilizar campos `primary_color` e `secondary_color` da entidade Company
- Criar middleware para injetar tema baseado em `company_id` do usuário
- Modificar frontend para suportar temas dinâmicos

**Entregáveis:**
- Theme Engine básico
- Middleware de injeção de tema
- Frontend adaptado para temas dinâmicos
- Documentação de customização de temas

---

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

A Sprint 1 da Plataforma PratoOnline 2.0 foi concluída com sucesso. A infraestrutura do Tenant Engine está estabelecida, garantindo:

✅ **Fundação sólida** para multi-tenancy  
✅ **Retrocompatibilidade total** com Core V1  
✅ **Arquitetura escalável** para futuras implementações  
✅ **Código limpo** seguindo Clean Architecture  
✅ **Migrações seguras** com rollback disponível  

O sistema está pronto para evoluir gradualmente em direção a um multi-tenant completo, mantendo a estabilidade e funcionalidades do Core V1.

---

**Próximos Passos Imediatos:**
1. Planejamento detalhado da Sprint 2 (White Label)
2. Definição de requisitos de theming
3. Design da interface de customização de temas
4. Preparação de ambiente de desenvolvimento para temas

---

**Assinaturas:**

Desenvolvido por: Jean Gouveia  
Data: 18 de Julho de 2026  
Versão: 2.0.1  
Status: Concluído
