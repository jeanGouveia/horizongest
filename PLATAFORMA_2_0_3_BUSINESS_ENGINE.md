# Plataforma PratoOnline 2.0 - Sprint 3: Business Engine

**Data:** 18 de Julho de 2026  
**Versão:** 2.0.3  
**Status:** Concluído  
**Objetivo:** Transformar a entidade Company na identidade funcional da plataforma, permitindo que cada empresa defina seu segmento de negócio

---

## Resumo Executivo

Esta sprint estabelece o Business Engine da Plataforma PratoOnline 2.0. O objetivo principal foi transformar a entidade `Company` na identidade funcional da plataforma, permitindo que cada empresa defina seu segmento de negócio, localização, moeda e fuso horário.

A implementação focou em:
- Criação do enum `BusinessType` com 11 tipos de negócio suportados
- Adição de campos `business_type`, `locale`, `currency`, e `timezone` à entidade Company
- Criação do Business Engine desacoplado (domain, service, handler)
- Implementação de endpoint `GET /api/business/profile`
- Sistema de fallback automático para empresas antigas
- Garantia de retrocompatibilidade completa com o Core V1

**Resultado:** A identidade funcional da plataforma está estabelecida, permitindo customização por tenant sem impactar a funcionalidade existente do Core V1.

---

## Arquitetura Proposta

### Estratégia de Business Engine

A arquitetura Business Engine adotada utiliza o padrão **Profile Aggregation** com fallback automático. Esta estratégia foi escolhida por:

- **Desacoplamento:** Business Engine completamente separado do Core V1
- **Flexibilidade:** Permite expansão futura com mais propriedades de negócio
- **Retrocompatibilidade:** Fallback automático para perfil padrão se Company não configurada
- **Simplicidade:** Não requer alterações em regras de negócio existentes
- **Escalabilidade:** Suporta múltiplos tipos de negócio com configurações específicas

### Modelo de Dados

```
BusinessType (Enum)
├── restaurant
├── bakery
├── confectionery
├── coffee_shop
├── pizzeria
├── burger
├── ice_cream
├── acai
├── food_truck
├── dark_kitchen
└── generic

Company (Domain Entity - Sprint 3 updates)
├── ID, Name, Slug, Description, Active
├── LogoURL, PrimaryColor, SecondaryColor (White Label)
├── BusinessType (Business Engine - NEW)
├── Locale (Business Engine - NEW)
├── Currency (Business Engine - NEW)
├── Timezone (Business Engine - NEW)
└── DeletedAt, CreatedAt, UpdatedAt

BusinessProfile (New Domain Entity)
├── CompanyID, CompanyName, CompanySlug, Active
├── BusinessType, Locale, Currency, Timezone
├── LogoURL, PrimaryColor, SecondaryColor
├── LoadedAt
└── IsDefault (flag se usando perfil padrão)
```

### Camadas da Arquitetura

```
┌─────────────────────────────────────────┐
│         API Layer (Handlers)            │
│  - BusinessHandler                      │
│  - GET /api/business/profile            │
│  - GET /api/business/profile/default    │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│        Service Layer                    │
│  - BusinessService                      │
│  - GetBusinessProfile(userID)          │
│  - GetDefaultBusinessProfile()          │
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
│  - BusinessProfile (nova entidade)      │
│  - BusinessType (novo enum)             │
│  - Company (entidade atualizada)         │
│  - User (entidade existente)            │
└─────────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│         Database (SQLite)               │
│  - companies (tabela atualizada)        │
│  - users (tabela existente)             │
└─────────────────────────────────────────┘
```

---

## Diagrama de Fluxo do Business Engine

```
┌──────────────┐
│ App Inicia   │
└──────┬───────┘
       │
       ↓
┌──────────────────────────────────────┐
│ Frontend chama GET /api/business/profile│
│ (com auth token)                      │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ BusinessHandler.GetBusinessProfile() │
│ - Extrai userID do contexto          │
│ - businessService.GetBusinessProfile()│
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ BusinessService.GetBusinessProfile() │
│ - userRepo.FindByID(userID)          │
│ - Se user.CompanyID == null:         │
│   return DefaultBusinessProfile()    │
│ - companyRepo.FindByID(companyID)    │
│ - Se company == null:                │
│   return DefaultBusinessProfile()    │
│ - return BusinessProfileFromCompany() │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ BusinessProfileFromCompany()         │
│ - Valida BusinessType                │
│ - Aplica fallbacks para Locale,      │
│   Currency, Timezone                 │
│ - Cria BusinessProfile com dados     │
│   da Company e Theme                 │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ Response JSON: BusinessProfile object│
│ - CompanyID, CompanyName, etc.       │
│ - BusinessType, Locale, Currency     │
│ - Timezone, LogoURL, Colors          │
│ - IsDefault flag                     │
└──────┬───────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────┐
│ Frontend recebe BusinessProfile       │
│ - Pode usar para customizar UI       │
│ - Formatar moedas, datas, etc.       │
│ - Adaptar funcionalidades por tipo  │
└──────────────────────────────────────┘
```

---

## Impacto em Cada Camada

### 1. Backend - Domain Layer

**Status:** Novos arquivos criados + atualização

**Novos arquivos:**
- `internal/domain/business_type.go` - Enum BusinessType com 11 tipos
- `internal/domain/business_profile.go` - Entidade BusinessProfile

**Arquivo atualizado:**
- `internal/domain/company.go` - Adicionados campos BusinessType, Locale, Currency, Timezone

**BusinessType Enum:**
```go
type BusinessType string

const (
    BusinessTypeRestaurant      BusinessType = "restaurant"
    BusinessTypeBakery          BusinessType = "bakery"
    BusinessTypeConfectionery   BusinessType = "confectionery"
    BusinessTypeCoffeeShop      BusinessType = "coffee_shop"
    BusinessTypePizzeria        BusinessType = "pizzeria"
    BusinessTypeBurger          BusinessType = "burger"
    BusinessTypeIceCream       BusinessType = "ice_cream"
    BusinessTypeAcai            BusinessType = "acai"
    BusinessTypeFoodTruck       BusinessType = "food_truck"
    BusinessTypeDarkKitchen     BusinessType = "dark_kitchen"
    BusinessTypeGeneric         BusinessType = "generic"
)
```

**BusinessProfile Entity:**
```go
type BusinessProfile struct {
    CompanyID     uint
    CompanyName   string
    CompanySlug   string
    Active        bool
    BusinessType  BusinessType
    Locale        string
    Currency      string
    Timezone      string
    LogoURL       string
    PrimaryColor  string
    SecondaryColor string
    LoadedAt      time.Time
    IsDefault     bool
}
```

**Impacto:**
- Zero impacto em entidades existentes além de Company
- Entidades desacopladas do Core V1
- Pode ser expandida futuramente com mais propriedades

---

### 2. Backend - Service Layer

**Status:** Novo arquivo criado + atualização

**Novo arquivo:**
- `internal/service/business_service.go` - BusinessService

**Arquivo atualizado:**
- `internal/service/company_service.go` - Adicionados campos aos inputs

**BusinessService:**
```go
type BusinessService struct {
    companyRepo ports.CompanyRepository
    userRepo    ports.UserRepository
}

func (s *BusinessService) GetBusinessProfile(ctx, userID) (*domain.BusinessProfile, error)
func (s *BusinessService) GetDefaultBusinessProfile() *domain.BusinessProfile
```

**Lógica:**
1. Busca usuário por ID
2. Se usuário não tem CompanyID → retorna perfil padrão
3. Busca empresa por CompanyID
4. Se empresa não encontrada → retorna perfil padrão
5. Cria perfil a partir da empresa
6. Retorna perfil

**CompanyService Updates:**
- `CreateCompanyInput` - Adicionados business_type, locale, currency, timezone
- `UpdateCompanyInput` - Adicionados business_type, locale, currency, timezone
- `CreateCompany` - Valida BusinessType e aplica fallbacks
- `UpdateCompany` - Atualiza campos se fornecidos

**Impacto:**
- Zero impacto em serviços existentes
- Lógica de fallback robusta
- Retrocompatibilidade garantida

---

### 3. Backend - Handler Layer

**Status:** Novo arquivo criado

**Arquivo:** `internal/handler/business_handler.go`

**BusinessHandler:**
```go
type BusinessHandler struct {
    businessService *service.BusinessService
}

func (h *BusinessHandler) GetBusinessProfile(w, r)
func (h *BusinessHandler) GetDefaultBusinessProfile(w, r)
func (h *BusinessHandler) RegisterRoutes(r)
```

**Endpoints:**
- `GET /api/business/profile` - Retorna perfil do usuário autenticado
- `GET /api/business/profile/default` - Retorna perfil padrão

**Impacto:**
- Novos endpoints públicos (mas protegidos por auth)
- Zero impacto em handlers existentes
- Segurança mantida via AuthMiddleware

---

### 4. Backend - Repository Layer

**Status:** Arquivo atualizado

**Arquivo:** `internal/infra/repository/gorm_company_repository.go`

**GormCompanyModel Updates:**
```go
type GormCompanyModel struct {
    // ... campos existentes ...
    BusinessType string `gorm:"type:text;default:'generic'"`
    Locale       string `gorm:"type:text;default:'pt-BR'"`
    Currency     string `gorm:"type:text;default:'BRL'"`
    Timezone     string `gorm:"type:text;default:'America/Sao_Paulo'"`
    // ... campos existentes ...
}
```

**Métodos atualizados:**
- `Create` - Mapeia novos campos
- `Update` - Mapeia novos campos
- `toDomainCompany` - Converte BusinessType string para enum

**Impacto:**
- GORM AutoMigrate adiciona colunas automaticamente
- Valores padrão garantem retrocompatibilidade
- Zero impacto em queries existentes

---

### 5. Backend - Database Layer

**Status:** Arquivos criados + atualização

**Novo arquivo SQL:**
- `migrations/00010_add_business_fields_to_companies.sql`

**Arquivo atualizado:**
- `internal/infra/database/migrate.go` - Adicionado GormCompanyModel ao AutoMigrate

**Migration SQL:**
```sql
ALTER TABLE companies ADD COLUMN business_type TEXT;
ALTER TABLE companies ADD COLUMN locale TEXT DEFAULT 'pt-BR';
ALTER TABLE companies ADD COLUMN currency TEXT DEFAULT 'BRL';
ALTER TABLE companies ADD COLUMN timezone TEXT DEFAULT 'America/Sao_Paulo';

UPDATE companies SET 
    business_type = 'generic',
    locale = 'pt-BR',
    currency = 'BRL',
    timezone = 'America/Sao_Paulo'
WHERE business_type IS NULL OR locale IS NULL OR currency IS NULL OR timezone IS NULL;
```

**Impacto:**
- Migração automática via GORM AutoMigrate
- Valores padrão aplicados a empresas existentes
- Zero impacto em tabelas existentes

---

### 6. Backend - Main.go

**Status:** Modificado

**Alterações:**
- Adicionado `businessSvc := service.NewBusinessService(companyRepo, userRepo)`
- Adicionado `businessHandler := handler.NewBusinessHandler(businessSvc)`
- Adicionadas rotas:
  - `r.Get("/api/business/profile", businessHandler.GetBusinessProfile)`
  - `r.Get("/api/business/profile/default", businessHandler.GetDefaultBusinessProfile)`

**Impacto:**
- Injeção de dependência manual mantida
- Rotas registradas em grupo autenticado
- Zero impacto em rotas existentes

---

### 7. Frontend - Business Engine

**Status:** Não implementado nesta sprint

**Racional:**
- Foco desta sprint foi na infraestrutura backend
- Frontend pode consumir a API quando necessário
- Padrão estabelecido para futura implementação

**Impacto:**
- Zero impacto no frontend
- API disponível para consumo futuro
- Retrocompatibilidade mantida

---

## Compatibilidade com Core V1

### Garantias de Retrocompatibilidade

1. **Perfil Padrão como Fallback**
   - Se usuário não tem Company → retorna perfil padrão
   - Se Company não encontrada → retorna perfil padrão
   - Se API falha → frontend pode usar perfil padrão local

2. **Valores Padrão em Database**
   - Colunas novas têm valores padrão
   - Migração atualiza empresas existentes
   - Zero impacto em empresas antigas

3. **Zero Alterações em Regras de Negócio**
   - Nenhuma modificação em services do Core V1
   - Nenhuma modificação em handlers do Core V1
   - Nenhuma modificação em repositories do Core V1

4. **API Existente Inalterada**
   - Endpoints existentes não modificados
   - Payloads de request/response mantidos
   - Apenas novos endpoints adicionados

5. **Company Continua Opcional**
   - Usuários podem existir sem Company
   - Funcionalidades do Core V1 funcionam normalmente
   - Business Engine é completamente opcional

### Testes de Compatibilidade Realizados

✅ **Backend - Business Profile API**
- `GET /api/business/profile` com usuário sem Company: OK (retorna perfil padrão)
- `GET /api/business/profile/default`: OK (retorna perfil padrão)
- `POST /api/companies` com business_type: OK (cria empresa com campos de negócio)

✅ **Backend - Core V1 APIs**
- `GET /api/categories`: OK (retorna dados com CompanyID: null)
- `POST /api/companies`: OK (cria empresa com novos campos)
- `GET /api/companies`: OK (retorna empresas com novos campos)

✅ **Database - Migrações**
- GORM AutoMigrate adicionou colunas: OK
- Valores padrão aplicados: OK
- Empresas existentes atualizadas: OK

✅ **BusinessType Validation**
- BusinessType inválido → fallback para generic: OK
- BusinessType válido → armazenado corretamente: OK
- DisplayName funciona corretamente: OK

---

## Riscos e Mitigações

### Riscos Identificados

1. **Validação de BusinessType**
   - **Risco:** Valores inválidos de business_type podem causar inconsistência
   - **Mitigação:** Validação no service com fallback para generic

2. **Timezone Handling**
   - **Risco:** Timezones inválidos podem causar problemas de data/hora
   - **Mitigação:** Validação básica de string; futuro: usar biblioteca de timezone

3. **Currency Validation**
   - **Risco:** Moedas inválidas podem causar problemas de formatação
   - **Mitigação:** Validação básica de string; futuro: usar ISO 4217

4. **Locale Validation**
   - **Risco:** Locales inválidos podem causar problemas de internacionalização
   - **Mitigação:** Validação básica de string; futuro: usar padrão IETF BCP 47

### Riscos Futuros (Próximas Sprints)

1. **Expansão do Business Engine**
   - Implementar validação mais rigorosa de timezone, currency, locale
   - Criar configurações específicas por tipo de negócio
   - Adicionar capabilities baseadas em business_type

2. **Frontend Integration**
   - Criar store para BusinessProfile
   - Implementar formatação de moedas baseada em currency
   - Implementar formatação de datas baseada em timezone
   - Adaptar UI baseada em business_type

3. **Multi-tenant Filtering**
   - Implementar filtros automáticos por company_id
   - Garantir isolamento de dados por tenant
   - Adicionar testes de segurança

---

## Próximos Passos

### Sprint 4: Filtros Automáticos por Tenant (Planejado)

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

### Sprint 5: Internacionalização (i18n) (Planejado)

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

### Sprint 6: RBAC por Tenant (Planejado)

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

A Sprint 3 da Plataforma PratoOnline 2.0 foi concluída com sucesso. A identidade funcional da plataforma está estabelecida, garantindo:

✅ **Business Engine desacoplado** do Core V1  
✅ **11 tipos de negócio** suportados  
✅ **Retrocompatibilidade total** com Core V1  
✅ **Arquitetura escalável** para futuras expansões  
✅ **Código limpo** seguindo padrões existentes  
✅ **Sistema robusto** com fallbacks automáticos  
✅ **API funcional** para consumo futuro  

O sistema está pronto para evoluir gradualmente em direção a um Business Engine completo, mantendo a estabilidade e funcionalidades do Core V1.

---

**Próximos Passos Imediatos:**
1. Planejamento detalhado da Sprint 4 (Filtros Automáticos)
2. Definição de requisitos de isolamento de dados
3. Design de middleware de filtro por tenant
4. Preparação de ambiente de desenvolvimento para multi-tenant

---

**Assinaturas:**

Desenvolvido por: Jean Gouveia  
Data: 18 de Julho de 2026  
Versão: 2.0.3  
Status: Concluído
