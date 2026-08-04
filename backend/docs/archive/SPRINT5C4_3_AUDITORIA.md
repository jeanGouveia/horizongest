# SPRINT 5C.4.3 — Auditoria Completa de Bootstrap da Plataforma

**Data:** 2026-07-31  
**Objetivo:** Auditoria criteriosa do sistema após migração para PostgreSQL, introdução do Redis e criação do novo framework de Consumers  
**Status:** 🔴 CRÍTICO - Múltiplos problemas graves encontrados

---

## Sumário Executivo

A auditoria identificou **problemas CRÍTICOS** que impedem o funcionamento correto do sistema:

- **Migração PostgreSQL INCOMPLETA**: SQLite remnants ainda presentes em todo o códigobase
- **Funções SQL incompatíveis**: `FROM_UNIXTIME()`, `strftime()`, `datetime()` não funcionam em PostgreSQL
- **Tipos de dados incorretos**: `INTEGER AUTOINCREMENT`, `TEXT`, `REAL` são sintaxe SQLite
- **Filtros de tenant ausentes**: Repositórios sem ApplyTenantFilter causam vazamento de dados
- **Testes usando SQLite**: Arquivos de teste ainda usam driver SQLite
- **Dependências SQLite ainda presentes**: go.mod e go.sum contêm dependências SQLite

**Classificação dos Problemas:**
- 🔴 Crítico: 15 problemas
- 🟠 Alto: 8 problemas
- 🟡 Médio: 5 problemas
- 🔵 Baixo: 3 problemas
- ⚪ Sugestão: 2 problemas

---

## 🔴 Problemas CRÍTICOS

### P1: Migração PostgreSQL INCOMPLETA - SQLite driver ainda presente

**Arquivo:** `backend/go.mod`  
**Linha:** 15  
**Código:**
```go
gorm.io/driver/sqlite v1.6.0
```

**Causa Raiz:** A migração reportada em POSTGRES_MIGRATION_REPORT.md não removeu a dependência SQLite do go.mod

**Impacto:** 
- Dependência desnecessária aumenta tamanho do binário
- Confusão sobre qual banco está sendo usado
- Possível uso acidental de SQLite em produção

**Risco:** 🔴 CRÍTICO - Pode causar comportamento indefinido se código tentar usar SQLite

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```bash
go mod edit -droprequire=gorm.io/driver/sqlite
go mod tidy
```

---

### P2: Migração PostgreSQL INCOMPLETA - SQLite driver em go.sum

**Arquivo:** `backend/go.sum`  
**Linhas:** 76-77  
**Código:**
```
gorm.io/driver/sqlite v1.6.0 h1:WHRRrIiulaPiPFmDcod6prc4l2VGVWHz80KspNsxSfQ=
gorm.io/driver/sqlite v1.6.0/go.mod h1:AO9V1qIQddBESngQUKWL9yoH93HIeA1X6V633rBwyT8=
github.com/mattn/go-sqlite3 v1.14.22 h1:2gZY6PC6kBnID23Tichd1K+Z0oS6nE/XwU+Vz/5o4kU=
github.com/mattn/go-sqlite3 v1.14.22/go.mod h1:Uh1q+B4BYcTPb+yiD3kU8Ct7aC0hY9fxUwlHK0RXw+Y=
```

**Causa Raiz:** go mod tidy não foi executado após remoção manual

**Impacto:** Mesmo que P1 seja corrigido, go.sum ainda contém hashes SQLite

**Risco:** 🔴 CRÍTICO - go.sum inconsistente pode causar problemas de build

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```bash
go mod edit -droprequire=gorm.io/driver/sqlite
go mod tidy
```

**Status:** ✅ RESOLVIDO (SPRINT 5D.1)
- go.sum limpo após go mod tidy
- Arquivos alterados: backend/go.sum

---

### P3: Arquivo de teste usando SQLite - test_snapshot_ingredient.go

**Arquivo:** `backend/test_snapshot_ingredient.go`  
**Linhas:** 8-9, 19, 26-27  
**Código:**
```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)
db, err := gorm.Open(sqlite.Open("test_snapshot_ingredient.db"), &gorm.Config{})
sqlDB.Exec("PRAGMA journal_mode=WAL;")
sqlDB.Exec("PRAGMA foreign_keys=ON;")
```

**Causa Raiz:** Arquivo de teste não foi atualizado durante migração PostgreSQL

**Impacto:** 
- Teste não funciona com PostgreSQL
- Teste cria arquivo .db local
- PRAGMA statements são SQLite-only

**Risco:** 🔴 CRÍTICO - Testes não refletem ambiente de produção

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```go
// Corrigido para usar PostgreSQL
```

**Status:** ✅ RESOLVIDO (SPRINT 5D.1)
- gorm_outbox_repository_test.go atualizado para usar PostgreSQL
- Testes de idempotência e isolamento de tenant implementados
- Arquivos alterados: backend/internal/infra/repository/gorm_outbox_repository_test.go

---

### P4: Arquivo de teste usando SQLite - gorm_outbox_repository_test.go

**Arquivo:** `backend/internal/infra/repository/gorm_outbox_repository_test.go`  
**Linhas:** 10, 16  
**Código:**
```go
import (
    "gorm.io/driver/sqlite"
)
db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
```

**Causa Raiz:** Teste de outbox repository não foi atualizado

**Impacto:** Teste usa SQLite em memória enquanto produção usa PostgreSQL

**Risco:** 🔴 CRÍTICO - Testes não validam comportamento PostgreSQL

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Atualizar para usar PostgreSQL como outros testes

---

### P5: Migration files usando sintaxe SQLite - INTEGER PRIMARY KEY AUTOINCREMENT

**Arquivos:** TODOS os arquivos em `backend/migrations/`  
**Exemplo:** `migrations/00027_create_orders_table.sql` linha 6  
**Código:**
```sql
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
```

**Causa Raiz:** Migrations nunca foram convertidas para sintaxe PostgreSQL

**Impacto:** 
- Migrations falham em PostgreSQL
- AUTOINCREMENT não existe em PostgreSQL (usa SERIAL/BIGSERIAL)
- Tables não são criadas corretamente

**Risco:** 🔴 CRÍTICO - Sistema não inicializa em PostgreSQL

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```sql
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
```

**Arquivos afetados:**
- 00001_create_users.sql
- 00002_create_base_tables.sql
- 00004_create_stock_adjustments_pending.sql
- 00008_create_companies_table.sql
- 00012_create_invitations.sql
- 00013_create_platform_users.sql
- 00014_create_platform_sessions.sql
- 00015_create_platform_audit.sql
- 00016_make_user_companyid_role_not_null.sql
- 00019_create_stock_movements.sql
- 00021_create_purchase_tables.sql
- 00022_create_finance_tables.sql
- 00023_create_platform_brand_config.sql
- 00024_create_global_config.sql
- 00025_create_plans.sql
- 00027_create_orders_table.sql
- 00029_create_idempotency_table.sql
- 00035_create_outbox_events.sql

---

### P6: Migration files usando sintaxe SQLite - strftime('%s', 'now')

**Arquivos:** `migrations/00016_make_user_companyid_role_not_null.sql` linhas 31-32, 74-75  
**Código:**
```sql
created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
```

**Causa Raiz:** Função strftime é SQLite-only

**Impacto:** Migrations falham em PostgreSQL

**Risco:** 🔴 CRÍTICO - Sistema não inicializa

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```sql
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
```

**Arquivos afetados:**
- 00001_create_users.sql
- 00002_create_base_tables.sql
- 00004_create_stock_adjustments_pending.sql
- 00007_create_media_table.sql
- 00013_create_platform_users.sql
- 00014_create_platform_sessions.sql
- 00015_create_platform_audit.sql
- 00016_make_user_companyid_role_not_null.sql
- 00023_create_platform_brand_config.sql
- 00024_create_global_config.sql
- 00025_create_plans.sql

---

### P7: Migration files usando sintaxe SQLite - datetime('now')

**Arquivos:** `migrations/00004_create_stock_adjustments_pending.sql` linha 13  
**Código:**
```sql
created_at DATETIME NOT NULL DEFAULT (datetime('now')),
```

**Causa Raiz:** Função datetime é SQLite-only

**Impacto:** Migrations falham em PostgreSQL

**Risco:** 🔴 CRÍTICO - Sistema não inicializa

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```sql
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
```

**Arquivos afetados:**
- 00001_create_users.sql
- 00002_create_base_tables.sql
- 00004_create_stock_adjustments_pending.sql

---

### P8: Repository usando função MySQL - FROM_UNIXTIME()

**Arquivo:** `backend/internal/infra/repository/gorm_report_repository.go`  
**Linhas:** 37, 46, 62, 72, 90, 162, 198, 400  
**Código:**
```go
Where("DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
```

**Causa Raiz:** FROM_UNIXTIME() é função MySQL, não existe em PostgreSQL

**Impacto:** Queries falham em PostgreSQL com erro "function from_unixtime does not exist"

**Risco:** 🔴 CRÍTICO - Relatórios não funcionam

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```go
// PostgreSQL usa timestamp diretamente
Where("DATE(created_at) >= ? AND DATE(created_at) <= ? AND deleted_at IS NULL", startDate, endDate).
```

**Arquivos afetados:**
- gorm_report_repository.go (8 ocorrências)
- gorm_notifications_repository.go (1 ocorrência)

---

### P9: Repository usando time.Now().Unix() para soft delete

**Arquivo:** `backend/internal/infra/repository/gorm_category_repository.go`  
**Linha:** 114  
**Código:**
```go
now := time.Now().Unix()
query := ApplyTenantFilterWithID(ctx, r.db, id)
if err := query.WithContext(ctx).Model(&GormCategory{}).
    Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
```

**Causa Raiz:** Código não foi atualizado após migração para usar time.Time

**Impacto:** deleted_at espera time.Time, está recebendo int64

**Risco:** 🔴 CRÍTICO - Soft delete falha com erro de tipo

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```go
now := time.Now()
query := ApplyTenantFilterWithID(ctx, r.db, id)
if err := query.WithContext(ctx).Model(&GormCategory{}).
    Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
```

---

### P10: Repository sem ApplyTenantFilter - gorm_plan_repository.go

**Arquivo:** `backend/internal/infra/repository/gorm_plan_repository.go`  
**Linhas:** 27-33, 59-65, 68-74  
**Código:**
```go
func (r *GormPlanRepository) FindByID(ctx context.Context, id uint) (*domain.Plan, error) {
    var plan domain.Plan
    err := r.db.WithContext(ctx).First(&plan, id).Error
    // Sem ApplyTenantFilter
}

func (r *GormPlanRepository) List(ctx context.Context) ([]*domain.Plan, error) {
    var plans []*domain.Plan
    err := r.db.WithContext(ctx).Find(&plans).Error
    // Sem ApplyTenantFilter
}
```

**Causa Raiz:** Plans são entidades de plataforma, mas deveriam ter filtragem se tiverem company_id

**Impacto:** Se plans tiverem company_id no futuro, vazamento de dados

**Risco:** 🔴 CRÍTICO - Cross-tenant data leak potencial

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Verificar se plans devem ter tenant filtering. Se sim, adicionar ApplyTenantFilter

---

### P11: Repository sem ApplyTenantFilter - gorm_invitation_repository.go

**Arquivo:** `backend/internal/infra/repository/gorm_invitation_repository.go`  
**Linhas:** 51-59, 62-70  
**Código:**
```go
func (r *GormInvitationRepository) FindByID(ctx context.Context, id uint) (*domain.Invitation, error) {
    var model GormInvitationModel
    if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
        // Sem ApplyTenantFilter
    }
}

func (r *GormInvitationRepository) FindByToken(ctx context.Context, token string) (*domain.Invitation, error) {
    var model GormInvitationModel
    if err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error; err != nil {
        // Sem ApplyTenantFilter - token pode ser acessado cross-tenant
    }
}
```

**Causa Raiz:** Invitations têm company_id mas não usam ApplyTenantFilter

**Impacto:** 
- FindByID pode retornar invitation de outro tenant
- FindByToken pode acessar invitation de outro tenant via token

**Risco:** 🔴 CRÍTICO - Cross-tenant data leak confirmado

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:**
```go
func (r *GormInvitationRepository) FindByID(ctx context.Context, id uint) (*domain.Invitation, error) {
    var model GormInvitationModel
    query := ApplyTenantFilterWithID(ctx, r.db, id)
    if err := query.First(&model).Error; err != nil {
        // ...
    }
}

func (r *GormInvitationRepository) FindByToken(ctx context.Context, token string) (*domain.Invitation, error) {
    var model GormInvitationModel
    query := ApplyTenantFilter(ctx, r.db)
    if err := query.Where("token = ?", token).First(&model).Error; err != nil {
        // ...
    }
}
```

---

### P12: Repository sem ApplyTenantFilter - gorm_company_repository.go FindByID

**Arquivo:** `backend/internal/infra/repository/gorm_company_repository.go`  
**Linhas:** 81-104  
**Código:**
```go
func (r *GormCompanyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
    var model GormCompanyModel
    err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&model, id).Error
    // Sem ApplyTenantFilter - usa First(&model, id) direto
}
```

**Causa Raiz:** CompanyRepository é usado pelo platform admin, mas deveria ter validação

**Impacto:** Se este método for chamado de contexto tenant, pode acessar company de outro tenant

**Risco:** 🔴 CRÍTICO - Cross-tenant data leak potencial

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Adicionar validação de contexto ou documentar que é platform-only

---

### P13: Repository sem ApplyTenantFilter - gorm_global_config_repository.go

**Arquivo:** `backend/internal/infra/repository/gorm_global_config_repository.go`  
**Linha:** 70  
**Código:**
```go
err := r.db.WithContext(ctx).First(&gormConfig, 1).Error
```

**Causa Raiz:** Global config é singleton, mas não tem validação de contexto

**Impacto:** Se houver multi-tenancy no futuro, pode causar problemas

**Risco:** 🔴 CRÍTICO - Arquitetura não suporta multi-tenancy

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Documentar que é platform-only ou adicionar validação

---

### P14: Repository sem ApplyTenantFilter - gorm_platform_brand_repository.go

**Arquivo:** `backend/internal/infra/repository/gorm_platform_brand_repository.go`  
**Linha:** 83  
**Código:**
```go
err := r.db.WithContext(ctx).First(&gormBrand, 1).Error
```

**Causa Raiz:** Platform brand é singleton, mas não tem validação de contexto

**Impacto:** Se houver multi-tenancy no futuro, pode causar problemas

**Risco:** 🔴 CRÍTICO - Arquitetura não suporta multi-tenancy

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Documentar que é platform-only ou adicionar validação

---

### P15: Migration POSTGRES_MIGRATION_REPORT.md é FALSO

**Arquivo:** `backend/POSTGRES_MIGRATION_REPORT.md`  
**Linhas:** 235-238  
**Código:**
```
- ✅ No `sqlite.Open` references found
- ✅ No `gorm.io/driver/sqlite` references found
- ✅ No `app.db` path references found
- ✅ No SQLite-specific PRAGMA statements found
```

**Causa Raiz:** Relatório foi gerado sem verificar completamente o códigobase

**Impacto:** 
- Equipe acredita que migração está completa
- Problemas não são corrigidos
- Sistema não funciona em PostgreSQL

**Risco:** 🔴 CRÍTICO - Documentação falsa causa confusão e atrasos

**Prioridade:** 🔴 CRÍTICA

**Correção Sugerida:** Atualizar relatório com status real ou remover documento

---

## 🟠 Problemas ALTO

### A1: Migration files usando TEXT onde deveria ser VARCHAR

**Arquivos:** Todos os migrations  
**Exemplo:** `migrations/00025_create_plans.sql` linha 7  
**Código:**
```sql
name TEXT NOT NULL,
slug TEXT NOT NULL UNIQUE,
```

**Causa Raiz:** TEXT é tipo SQLite, PostgreSQL usa VARCHAR

**Impacto:** Funciona mas não é ideal para performance e indexação

**Risco:** 🟠 ALTO - Performance subótima

**Prioridade:** 🟠 ALTA

**Correção Sugerida:**
```sql
name VARCHAR(255) NOT NULL,
slug VARCHAR(255) NOT NULL UNIQUE,
```

---

### A2: Migration files usando REAL onde deveria ser NUMERIC/DECIMAL

**Arquivos:** Todos os migrations  
**Exemplo:** `migrations/00027_create_orders_table.sql` linha 9  
**Código:**
```sql
total_price REAL NOT NULL DEFAULT 0.0,
```

**Causa Raiz:** REAL é tipo SQLite, PostgreSQL usa NUMERIC ou DECIMAL para dinheiro

**Impacto:** REAL pode ter problemas de precisão com valores monetários

**Risco:** 🟠 ALTO - Perda de precisão em cálculos financeiros

**Prioridade:** 🟠 ALTA

**Correção Sugerida:**
```sql
total_price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
```

---

### A3: POSTGRES_MIGRATION_REPORT.md claims timestamp fields fixed but code still has issues

**Arquivo:** `backend/POSTGRES_MIGRATION_REPORT.md`  
**Linha:** 92  
**Código:**
```
- Removed `time.Unix()` conversions in repository methods
```

**Causa Raiz:** Relatório afirma que time.Unix() foi removido, mas gorm_category_repository.go linha 114 ainda usa

**Impacto:** Documentação inconsistente com código

**Risco:** 🟠 ALTO - Confusão sobre estado real do código

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Corrigir código ou atualizar relatório

---

### A4: gorm_password_reset_repository sem ApplyTenantFilter

**Arquivo:** `backend/internal/infra/repository/gorm_password_reset_repository.go`  
**Linhas:** 52-59, 64-69  
**Código:**
```go
func (r *GormPasswordResetRepository) FindByToken(ctx context.Context, tokenStr string) (*domain.PasswordResetToken, error) {
    var model GormPasswordResetTokenModel
    err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&model).Error
    // Sem ApplyTenantFilter
}

func (r *GormPasswordResetRepository) FindByUserID(ctx context.Context, userID uint) ([]*domain.PasswordResetToken, error) {
    var models []GormPasswordResetTokenModel
    err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error
    // Sem ApplyTenantFilter - userID pode ser de outro tenant
}
```

**Causa Raiz:** Password reset tokens têm user_id mas não usam ApplyTenantFilter

**Impacto:** Token pode ser acessado cross-tenant via user_id

**Risco:** 🟠 ALTO - Cross-tenant data leak potencial

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Adicionar ApplyTenantFilter

---

### A5: gorm_report_repository usa ApplyTenantFilter mas ignora companyID parameter

**Arquivo:** `backend/internal/infra/repository/gorm_report_repository.go`  
**Linhas:** 31, 118, 156, 192, 233, 328, 349, 391  
**Código:**
```go
// companyID parâmetro é ignorado - usamos ApplyTenantFilter do contexto
query := ApplyTenantFilter(ctx, r.db)
```

**Causa Raiz:** companyID parameter é passado mas ignorado

**Impacto:** Confusão sobre qual company_id usar

**Risco:** 🟠 ALTO - Possível bug se contexto não tiver tenant

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Remover parameter ou usar como fallback

---

### A6: gorm_purchase_repository PurchaseReceiving não tem CompanyID

**Arquivo:** `backend/internal/infra/repository/gorm_purchase_repository.go`  
**Linhas:** 189-194  
**Código:**
```go
func (r *GormPurchaseRepository) CreatePurchaseReceiving(ctx context.Context, receiving *domain.PurchaseReceiving) error {
    // PurchaseReceiving não tem CompanyID direto - o tenant é herdado através de PurchaseOrder
    // A validação de tenant deve ser feita no service layer verificando o PurchaseOrderID
    if err := r.db.WithContext(ctx).Create(receiving).Error; err != nil {
```

**Causa Raiz:** PurchaseReceiving não tem company_id, depende de PurchaseOrder

**Impacto:** Se service layer não validar, cross-tenant data leak

**Risco:** 🟠 ALTO - Cross-tenant data leak se service não validar

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Adicionar company_id a PurchaseReceiving ou validar no repository

---

### A7: gorm_outbox_repository usa tenant_id em vez de company_id

**Arquivo:** `backend/internal/infra/repository/gorm_outbox_repository.go`  
**Linhas:** 56-63  
**Código:**
```go
// applyTenantFilter applies tenant filtering for outbox events (uses tenant_id instead of company_id)
func applyTenantFilter(ctx context.Context, db *gorm.DB) *gorm.DB {
    tenantCtx, ok := domain.GetTenantContextFromContext(ctx)
    if !ok {
        return db
    }
    return db.Where("tenant_id = ?", tenantCtx.GetCompanyID())
}
```

**Causa Raiz:** Inconsistência de nomenclatura - usa tenant_id mas mapeia para company_id

**Impacto:** Confusão na arquitetura

**Risco:** 🟠 ALTO - Erros de manutenção no futuro

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Padronizar para company_id em toda a arquitetura

---

### A8: gorm_company_repository tem logs FORENSIC em produção

**Arquivo:** `backend/internal/infra/repository/gorm_company_repository.go`  
**Linhas:** 82-101  
**Código:**
```go
// FORENSIC: Log CompanyID received
log.Printf("[FORENSIC] CompanyRepository.FindByID - CompanyID recebido: %d", id)

// FORENSIC: Log SQL that will be executed
log.Printf("[FORENSIC] CompanyRepository.FindByID - SQL a ser executado: SELECT * FROM companies WHERE id = %d AND deleted_at IS NULL", id)
```

**Causa Raiz:** Logs de debug FORENSIC deixados no código

**Impacto:** Poluição de logs em produção

**Risco:** 🟠 ALTO - Logs sensíveis em produção

**Prioridade:** 🟠 ALTA

**Correção Sugerida:** Remover logs FORENSIC ou usar logger condicional

---

## 🟡 Problemas MÉDIO

### M1: Migration 00017_add_composite_indexes.sql está vazia

**Arquivo:** `backend/migrations/00017_add_composite_indexes.sql`  
**Conteúdo:** Vazio

**Causa Raiz:** Migration foi criada mas nunca preenchida

**Impacto:** Índices compostos planejados não foram criados

**Risco:** 🟡 MÉDIO - Performance subótima

**Prioridade:** 🟡 MÉDIA

**Correção Sugerida:** Preencher migration ou remover arquivo

---

### M2: Migration 00018_add_fk_on_delete.sql é documentação-only

**Arquivo:** `backend/migrations/00018_add_fk_on_delete.sql`  
**Conteúdo:** Apenas comentários

**Causa Raiz:** Documentação de ON DELETE mas sem implementação

**Impacto:** Foreign keys não têm ON DELETE CASCADE/SET NULL

**Risco:** 🟡 MÉDIO - Orphan records podem acumular

**Prioridade:** 🟡 MÉDIA

**Correção Sugerida:** Implementar ou remover migration

---

### M3: gorm_category_repository DeleteCategory usa time.Now().Unix()

**Arquivo:** `backend/internal/infra/repository/gorm_category_repository.go`  
**Linha:** 114 (já listado como P9)

**Causa Raiz:** Código não foi atualizado para time.Time

**Impacto:** Type mismatch

**Risco:** 🟡 MÉDIO - Soft delete falha

**Prioridade:** 🟡 MÉDIA

**Correção Sugerida:** (Já documentado em P9)

---

### M4: gorm_dashboard_repository pode ter N+1 queries

**Arquivo:** `backend/internal/infra/repository/gorm_dashboard_repository.go`  
**Observação:** Multiple queries sem preload otimizado

**Causa Raiz:** Queries complexas podem gerar N+1

**Impacto:** Performance subótima em dashboard

**Risco:** 🟡 MÉDIO - Performance degradation

**Prioridade:** 🟡 MÉDIA

**Correção Sugerida:** Revisar queries e adicionar preload onde necessário

---

### M5: gorm_order_repository UpdateOrder pode ter race condition

**Arquivo:** `backend/internal/infra/repository/gorm_order_repository.go`  
**Linhas:** 461-478  
**Código:**
```go
// Get current order to calculate stock adjustments
var gOrder GormOrder
query := ApplyTenantFilterWithID(ctx, tx, id)
if err := query.Where("deleted_at IS NULL").First(&gOrder).Error; err != nil {
```

**Causa Raiz:** UpdateOrder lê order sem SELECT FOR UPDATE

**Impacto:** Race condition se dois updates simultâneos

**Risco:** 🟡 MÉDIO - Inconsistência de dados

**Prioridade:** 🟡 MÉDIA

**Correção Sugerida:** Adicionar SELECT FOR UPDATE ou validar no service

---

## 🔵 Problemas BAIXO

### B1: go.mod tem comentário SQLite removido mas dependência ainda existe

**Arquivo:** `backend/go.mod`  
**Linha:** 15  
**Observação:** Comentário no POSTGRES_MIGRATION_REPORT.md diz que foi removido, mas ainda está

**Impacto:** Documentação inconsistente

**Risco:** 🔵 BAIXO - Confusão

**Prioridade:** 🔵 BAIXA

**Correção Sugerida:** (Já documentado em P1-P2)

---

### B2: .gitignore ainda tem *.sqlite

**Arquivo:** `backend/.gitignore`  
**Linhas:** 27-28, 44-45  
**Código:**
```
*.sqlite
*.sqlite3
```

**Impacto:** Não há mais arquivos .sqlite no projeto

**Risco:** 🔵 BAIXO - Configuração desnecessária

**Prioridade:** 🔵 BAIXA

**Correção Sugerida:** Remover do .gitignore

---

### B3: test_snapshot_ingredient.go não está em diretório de testes

**Arquivo:** `backend/test_snapshot_ingredient.go`  
**Localização:** Raiz do backend, não em diretório tests/

**Impacto:** Arquivo de teste fora de local padrão

**Risco:** 🔵 BAIXO - Organização de código

**Prioridade:** 🔵 BAIXA

**Correção Sugerida:** Mover para tests/ ou remover

---

## ⚪ Sugestões

### S1: Considerar usar UUID para IDs em vez de auto-increment

**Razão:** 
- Melhor para distribuição
- Evita conflitos em merge
- Mais seguro (não expõe sequência)

**Impacto:** Arquitetural

**Prioridade:** ⚪ SUGESTÃO

---

### S2: Considerar implementar cache de dashboard

**Razão:** 
- Reduz load do banco
- Melhora performance
- Dashboard é acessado frequentemente

**Impacto:** Performance

**Prioridade:** ⚪ SUGESTÃO

---

## Resumo por Categoria

| Categoria | Crítico | Alto | Médio | Baixo | Sugestão | Total |
|-----------|---------|------|-------|------|----------|-------|
| PostgreSQL Migration | 8 | 2 | 0 | 1 | 0 | 11 |
| Multi-Tenant Isolation | 5 | 2 | 0 | 0 | 0 | 7 |
| Repository Code | 2 | 4 | 2 | 0 | 0 | 8 |
| Test Files | 2 | 0 | 0 | 1 | 0 | 3 |
| Documentation | 1 | 1 | 0 | 0 | 0 | 2 |
| Performance | 0 | 0 | 2 | 0 | 1 | 3 |
| Arquitetura | 1 | 0 | 0 | 0 | 1 | 2 |
| **TOTAL** | **15** | **8** | **5** | **3** | **2** | **33** |

---

## Resumo por Prioridade

| Prioridade | Quantidade | % |
|------------|------------|---|
| 🔴 Crítico | 15 | 45% |
| 🟠 Alto | 8 | 24% |
| 🟡 Médio | 5 | 15% |
| 🔵 Baixo | 3 | 9% |
| ⚪ Sugestão | 2 | 6% |
| **TOTAL** | **33** | 100% |

---

## Arquivos Afetados

### Migrations (18 arquivos)
- 00001_create_users.sql
- 00002_create_base_tables.sql
- 00004_create_stock_adjustments_pending.sql
- 00007_create_media_table.sql
- 00008_create_companies_table.sql
- 00012_create_invitations.sql
- 00013_create_platform_users.sql
- 00014_create_platform_sessions.sql
- 00015_create_platform_audit.sql
- 00016_make_user_companyid_role_not_null.sql
- 00017_add_composite_indexes.sql
- 00018_add_fk_on_delete.sql
- 00019_create_stock_movements.sql
- 00021_create_purchase_tables.sql
- 00022_create_finance_tables.sql
- 00023_create_platform_brand_config.sql
- 00024_create_global_config.sql
- 00025_create_plans.sql
- 00027_create_orders_table.sql
- 00029_create_idempotency_table.sql
- 00035_create_outbox_events.sql

### Repositories (10 arquivos)
- gorm_plan_repository.go
- gorm_invitation_repository.go
- gorm_company_repository.go
- gorm_global_config_repository.go
- gorm_platform_brand_repository.go
- gorm_password_reset_repository.go
- gorm_report_repository.go
- gorm_purchase_repository.go
- gorm_category_repository.go
- gorm_order_repository.go
- gorm_dashboard_repository.go
- gorm_notifications_repository.go
- gorm_outbox_repository.go

### Test Files (2 arquivos)
- test_snapshot_ingredient.go
- gorm_outbox_repository_test.go

### Configuração (3 arquivos)
- go.mod
- go.sum
- .gitignore

### Documentação (2 arquivos)
- POSTGRES_MIGRATION_REPORT.md
- (Este arquivo)

---

## Plano de Correção Recomendado

### Fase 1: Crítico - Migração PostgreSQL (Estimado: 4-6 horas)
1. Remover dependências SQLite do go.mod e go.sum
2. Converter todas as migrations para sintaxe PostgreSQL
3. Atualizar test files para usar PostgreSQL
4. Atualizar repository code para remover funções MySQL
5. Corrigir time.Now().Unix() para time.Now()
6. Atualizar/remover POSTGRES_MIGRATION_REPORT.md

### Fase 2: Crítico - Multi-Tenant Isolation (Estimado: 2-3 horas)
1. Adicionar ApplyTenantFilter em gorm_invitation_repository
2. Adicionar ApplyTenantFilter em gorm_password_reset_repository
3. Validar gorm_company_repository FindByID
4. Validar gorm_global_config_repository
5. Validar gorm_platform_brand_repository
6. Validar gorm_plan_repository

### Fase 3: Alto - Consistência e Performance (Estimado: 1-2 horas)
1. Converter TEXT para VARCHAR em migrations
2. Converter REAL para NUMERIC em migrations
3. Corrigir gorm_report_repository companyID parameter
4. Adicionar company_id a PurchaseReceiving ou validar
5. Padronizar tenant_id vs company_id
6. Remover logs FORENSIC

### Fase 4: Médio - Performance e Consistência (Estimado: 1-2 horas)
1. Preencher migration 00017 ou remover
2. Implementar migration 00018 ou remover
3. Revisar gorm_dashboard_repository para N+1
4. Adicionar SELECT FOR UPDATE em UpdateOrder

### Fase 5: Baixo - Limpeza (Estimado: 30 minutos)
1. Remover *.sqlite do .gitignore
2. Mover test_snapshot_ingredient.go para tests/ ou remover

### Fase 6: Sugestões (Opcional)
1. Avaliar UUID para IDs
2. Implementar cache de dashboard

---

## Estimativa de Esforço Total

| Fase | Horas |
|------|-------|
| Fase 1: Crítico - Migração PostgreSQL | 4-6 |
| Fase 2: Crítico - Multi-Tenant Isolation | 2-3 |
| Fase 3: Alto - Consistência e Performance | 1-2 |
| Fase 4: Médio - Performance e Consistência | 1-2 |
| Fase 5: Baixo - Limpeza | 0.5 |
| Fase 6: Sugestões | 2-4 |
| **TOTAL** | **10.5-17.5 horas** |

---

## Avaliação Arquitetural Final

### Estado Atual: 🔴 CRÍTICO

O sistema **NÃO ESTÁ PRONTO PARA PRODUÇÃO** com PostgreSQL devido a:

1. **Migração PostgreSQL INCOMPLETA**: Migrations ainda usam sintaxe SQLite
2. **Dependências SQLite ainda presentes**: go.mod e go.sum não foram limpos
3. **Funções SQL incompatíveis**: FROM_UNIXTIME, strftime, datetime não funcionam em PostgreSQL
4. **Cross-tenant data leaks**: Repositórios sem ApplyTenantFilter
5. **Documentação falsa**: POSTGRES_MIGRATION_REPORT.md afirma que migração está completa

### Recomendação

**NÃO DEPLOYAR EM PRODUÇÃO** até que todos os problemas 🔴 CRÍTICOS sejam corrigidos.

A migração PostgreSQL reportada em POSTGRES_MIGRATION_REPORT.md é **INCOMPLETA** e **INCORRETA**. O sistema atual não funciona corretamente com PostgreSQL.

---

## Próximos Passos

1. **IMEDIATO**: Corrigir Fase 1 (Migração PostgreSQL)
2. **CURTO PRAZO**: Corrigir Fase 2 (Multi-Tenant Isolation)
3. **MÉDIO PRAZO**: Corrigir Fases 3-4 (Consistência e Performance)
4. **LONGO PRAZO**: Avaliar Fase 6 (Sugestões)

---

**Fim da Auditoria**
