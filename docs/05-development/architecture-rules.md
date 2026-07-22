# ARCHITECTURE_RULES.md

**Constituição do Projeto HorizonGest**

Este documento define as regras arquiteturais obrigatórias do projeto. Qualquer violação deve ser justificada e aprovada por arquiteto sênior.

---

## 1. Separação de Camadas

### 1.1 Camadas Obrigatórias

O projeto segue arquitetura em camadas com separação estrita:

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

### 1.2 Regras de Direção

- **Handler NUNCA chama Repository**
- **Handler NUNCA chama outro Handler**
- **Handler NUNCA acessa banco de dados diretamente**
- **Service NUNCA chama Handler**
- **Service NUNCA acessa banco de dados diretamente**
- **Repository NUNCA chama Service**
- **Repository NUNCA chama Handler**

### 1.3 Dependências Permitidas

- Handler → Service
- Handler → Domain
- Service → Repository
- Service → Domain
- Repository → Domain
- Repository → Database (via GORM)

---

## 2. Regras de Negócio

### 2.1 Localização de Regras de Negócio

- **TODA regra de negócio deve estar no Service Layer**
- **NENHUMA regra de negócio pode existir no Frontend**
- **NENHUMA regra de negócio pode existir no Handler**
- **NENHUMA regra de negócio pode existir no Repository**

### 2.2 Validação

- Validação de input: Handler (via validator)
- Validação de negócio: Service
- Validação de dados: Repository (constraints de banco)

### 2.3 Exemplos de Regras de Negócio

- Cálculos financeiros
- Validação de permissões
- Regras de estoque
- Lógica de pedidos
- Validação de relacionamentos

---

## 3. Regras de Dados

### 3.1 CompanyID

- **TODA entidade deve possuir CompanyID** (exceto entidades globais)
- CompanyID deve ser obrigatório em tabelas de tenant
- CompanyID deve ser usado em todas as queries de tenant

### 3.2 Entidades Globais

Entidades sem CompanyID (globais):
- PlatformUser
- PlatformSession
- PlatformAudit
- Plan
- PlatformBrandConfig
- GlobalConfig

### 3.3 Soft Delete

- Toda entidade deve ter DeletedAt (soft delete)
- Queries devem filtrar DeletedAt IS NULL por padrão
- Hard delete apenas em casos excepcionais

---

## 4. Regras de Autenticação e Autorização

### 4.1 Autenticação

- Toda rota protegida deve passar pelo middleware de autenticação
- JWT deve ser validado em cada requisição
- Token deve conter UserID e CompanyID

### 4.2 Autorização

- Toda rota protegida deve passar pelo middleware de autorização
- RBAC deve ser usado para controle de permissões
- Permissões devem ser verificadas no Service Layer

### 4.3 Platform vs Tenant

- Platform routes: `/api/platform/*` (autenticação platform)
- Tenant routes: `/api/*` (autenticação tenant)
- Platform users não podem acessar rotas de tenant
- Tenant users não podem acessar rotas de platform

---

## 5. Regras de Branding

### 5.1 Platform Branding

- Platform Branding deve vir de `PlatformBrandConfig` (banco)
- NENHUMA referência hardcoded de branding no código
- Nome da plataforma, logo, cores devem ser dinâmicos
- Frontend deve consumir `/api/public/brand` (endpoint público)

### 5.2 Tenant Branding

- Tenant Branding deve vir de `Theme` e `BusinessProfile` (banco)
- Cada tenant pode ter seu próprio branding
- Tenant branding é separado de platform branding

### 5.3 Separação

- Platform Branding: identidade institucional da plataforma
- Tenant Branding: identidade visual da empresa cliente
- NÃO misturar os dois conceitos

---

## 6. Regras de Configuração

### 6.1 PlatformBrandConfig vs GlobalConfig

- **PlatformBrandConfig**: Branding/institucional (nome, logo, cores, e-mail, copyright)
- **GlobalConfig**: Configurações técnicas (timezone, locale, upload limits, feature flags)
- NÃO misturar branding com configurações técnicas

### 6.2 Environment Variables

- Segredos devem vir de environment variables (JWT secrets, DB password)
- Configurações de infraestrutura devem vir de environment variables (DB host, port)
- Configurações de negócio devem vir do banco (branding, feature flags)

### 6.3 APP_VERSION

- Versão da aplicação deve vir de `APP_VERSION` environment variable
- Helper `util.PlatformVersion()` deve ser usado para obter versão
- Versão NÃO deve ser armazenada no banco

---

## 7. Regras de Feature Flags

### 7.1 Implementação

- Feature flags devem ser armazenadas em `GlobalConfig`
- Feature flags devem ser usadas para habilitar/desabilitar módulos inteiros
- Service deve verificar feature flag antes de executar lógica de módulo

### 7.2 Módulos

- Todo novo módulo deve ser registrado em `ModuleRegistry`
- Todo novo módulo deve ter feature flag correspondente
- Todo novo módulo deve declarar dependências no registry

### 7.3 Uso

- Handler deve verificar feature flag antes de expor rotas
- Service deve verificar feature flag antes de executar lógica
- Frontend deve verificar feature flag antes de exibir UI

---

## 8. Regras de Cache

### 8.1 Implementação

- Cache deve ser implementado no Repository Layer
- Cache deve usar `sync.RWMutex` para thread-safety
- Cache deve ser invalidado automaticamente após Update

### 8.2 Cache-First Logic

- Repository deve usar cache-first logic
- Cache miss deve carregar do banco e atualizar cache
- Service Layer não deve conhecer detalhes do cache

### 8.3 Invalidação

- Update deve invalidar cache automaticamente
- Métodos explícitos: `InvalidateCache()`, `ReloadCache()`
- Service pode chamar `ReloadCache()` se necessário

---

## 9. Regras de Migrations

### 9.1 Idempotência

- Toda migration deve ser idempotente
- Use `INSERT OR IGNORE` para inserts
- Use `CREATE TABLE IF NOT EXISTS` para creates

### 9.2 Versionamento

- Toda migration deve ter nome com timestamp: `YYYYMMDD_descrição.sql`
- Toda migration deve ter `+goose Up` e `+goose Down`
- Down migration deve reverter completamente Up migration

### 9.3 Schema Evolution

- Campos novos devem ser opcionais quando possível
- Valores padrão razoáveis devem ser fornecidos
- NÃO remover campos sem migration de rollback

---

## 10. Regras de Frontend

### 10.1 Separação de Responsabilidades

- Frontend NÃO deve conter regra de negócio
- Frontend deve focar em UI/UX
- Frontend deve consumir APIs do backend

### 10.2 Branding Dinâmico

- Frontend deve consumir `/api/public/brand` para branding
- Frontend NÃO deve ter branding hardcoded
- Cores, logos, nomes devem vir do backend

### 10.3 Validação

- Frontend pode fazer validação de input (UX)
- Validação de negócio deve ser feita no backend
- Frontend NÃO deve confiar apenas em validação frontend

---

## 11. Regras de Nomenclatura

### 11.1 Go Code

- Structs: PascalCase (ex: `User`, `Product`)
- Métodos: PascalCase (ex: `GetUser`, `CreateProduct`)
- Variáveis: camelCase (ex: `userID`, `productName`)
- Constantes: PascalCase ou UPPER_CASE (ex: `MaxRetries`, `MAX_RETRIES`)
- Pacotes: lowercase (ex: `service`, `repository`)

### 11.2 Database

- Tabelas: snake_case (ex: `users`, `platform_brand_config`)
- Colunas: snake_case (ex: `user_id`, `platform_name`)
- Índices: `idx_nome_tabela_colunas` (ex: `idx_users_email`)

### 11.3 API

- Rotas: kebab-case (ex: `/api/users`, `/api/platform-brand`)
- JSON: camelCase (ex: `userId`, `platformName`)
- HTTP methods: GET, POST, PUT, DELETE (semântica correta)

---

## 12. Regras de Error Handling

### 12.1 Errors

- Errors devem ser retornados, não panic
- Errors devem ser específicos (ex: `ErrUserNotFound`, `ErrInvalidEmail`)
- Errors devem ser propagados até o Handler

### 12.2 Handler Response

- Errors de validação: 400 Bad Request
- Errors de autenticação: 401 Unauthorized
- Errors de autorização: 403 Forbidden
- Errors de não encontrado: 404 Not Found
- Errors de servidor: 500 Internal Server Error

### 12.3 Logging

- Errors devem ser logados no Service Layer
- Logs devem incluir contexto (userID, companyID)
- Logs NÃO devem conter informações sensíveis (senhas, tokens)

---

## 13. Regras de Testes

### 13.1 Testes Unitários

- Todo Service deve ter testes unitários
- Todo Repository deve ter testes unitários
- Testes devem usar mocks para dependências externas

### 13.2 Testes de Integração

- Handlers devem ter testes de integração
- Migrations devem ser testadas
- API endpoints devem ser testados

### 13.3 Coverage

- Cobertura mínima: 70%
- Camada crítica (Service): 90%+
- Camada de negócio: 80%+

---

## 14. Regras de Performance

### 14.1 Database

- Queries devem usar índices apropriados
- N+1 queries devem ser evitadas (use eager loading)
- Queries grandes devem ser paginadas

### 14.2 Cache

- Dados frequentemente lidos devem ser cacheados
- Cache deve ter TTL apropriado
- Cache deve ser invalidado corretamente

### 14.3 API

- Respostas devem ser paginadas quando apropriado
- Campos desnecessários não devem ser retornados
- Compression deve ser usada para respostas grandes

---

## 15. Regras de Segurança

### 15.1 Autenticação

- Senhas devem ser hashadas com bcrypt
- JWT deve ter expiração apropriada
- Tokens devem ser revogados no logout

### 15.2 Autorização

- RBAC deve ser implementado corretamente
- Permissões devem ser verificadas em cada operação
- Platform users devem ter permissões separadas

### 15.3 Input Validation

- TODA input deve ser validada
- SQL injection deve ser prevenida (use GORM)
- XSS deve ser prevenida (escape HTML)

### 15.4 Secrets

- Secrets nunca devem estar no código
- Secrets devem estar em environment variables
- Secrets nunca devem ser logadas

---

## 16. Regras de White Label

### 16.1 Preparação

- Arquitetura deve permitir múltiplas marcas
- Branding deve ser completamente dinâmico
- NENHUMA referência hardcoded de marca

### 16.2 Implementação Futura

- Adicionar `brand_key` para distinguir marcas
- Repository deve suportar múltiplas configurações
- Cache deve ser map-based (keyed por brand)
- Middleware deve selecionar brand baseado em domínio

### 16.3 Separação

- Platform Branding vs Tenant Branding
- Platform Branding: marca da plataforma
- Tenant Branding: marca do cliente
- NÃO misturar os dois

---

## 17. Regras de Versionamento

### 17.1 SemVer

- Versões devem seguir Semantic Versioning (MAJOR.MINOR.PATCH)
- MAJOR: mudanças incompatíveis
- MINOR: funcionalidades novas compatíveis
- PATCH: correções de bugs compatíveis

### 17.2 API Versioning

- API deve ser versionada via URL (ex: `/api/v1/users`)
- Versões antigas devem ser mantidas por período de deprecation
- Breaking changes requerem nova versão major

---

## 18. Regras de Documentação

### 18.1 Code Comments

- Funções públicas devem ter godoc
- Regras de negócio complexas devem ser comentadas
- TODOs devem ter contexto e responsável

### 18.2 README

- README deve ter instruções de setup
- README deve ter instruções de execução
- README deve ter arquitetura overview

### 18.3 API Documentation

- Endpoints devem ser documentados
- Request/Response schemas devem ser documentados
- Exemplos de uso devem ser fornecidos

---

## 19. Regras de Code Review

### 19.1 Processo

- Todo código deve ser reviewado
- Reviewer deve verificar violações de ARCHITECTURE_RULES
- Violações devem ser justificadas e aprovadas

### 19.2 Checklist

- [ ] Separação de camadas respeitada
- [ ] Regras de negócio no Service
- [ ] CompanyID presente em entidades de tenant
- [ ] Autenticação/Autorização implementada
- [ ] Branding dinâmico (sem hardcoded)
- [ ] Feature flags usadas apropriadamente
- [ ] Cache implementado se necessário
- [ ] Migration idempotente
- [ ] Testes escritos
- [ ] Documentação atualizada

---

## 20. Regras de Violação

### 20.1 Processo

- Violação deve ser documentada em issue
- Violação deve ter justificativa técnica
- Violação deve ter plano de correção
- Violação deve ser aprovada por arquiteto sênior

### 20.2 Exceções

- Exceções temporárias são permitidas com timeline de correção
- Exceções permanentes requerem aprovação de arquiteto sênior
- Exceções devem ser documentadas no código

---

## Conclusão

Estas regras são a "Constituição" do projeto HorizonGest. Qualquer violação deve ser tratada com seriedade e justificada apropriadamente. A arquitetura deve evoluir, mas as regras fundamentais devem permanecer estáveis.

**Última atualização:** Sprint 3.6 - Foundation Final
**Responsável:** Arquiteto do Projeto
