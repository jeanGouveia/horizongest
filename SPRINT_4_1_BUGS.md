# SPRINT 4.1 - Bugs Found

**Data:** 2026-07-20  
**Auditor:** Cascade AI

---

## Bugs Encontrados

### BUG-001: Login via API não funciona
**Severidade:** Alta  
**Arquivo provável:** `internal/service/auth_service.go`  
**Reprodução:**
1. Iniciar backend
2. Executar `curl -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"email":"jwtfinal@test.com","password":"admin123"}'`
3. Obter erro: `{"error":"e-mail ou senha incorretos. Verifique suas credenciais."}`

**Impacto:** Usuários não conseguem fazer login via API. Sistema inutilizável.

**Correção sugerida:** Verificar hash de senhas no banco e garantir que senhas de teste estejam corretamente configuradas.

---

### BUG-002: Endpoint forgot-password não existe
**Severidade:** Média  
**Arquivo provável:** `cmd/server/main.go` (rotas)  
**Reprodução:**
1. Executar `curl -X POST http://localhost:8080/api/auth/forgot-password -H "Content-Type: application/json" -d '{"email":"qa@test.com"}'`
2. Obter erro: `404 page not found`

**Impacto:** Usuários não conseguem recuperar senha.

**Correção sugerida:** Implementar endpoint `/api/auth/forgot-password` ou verificar se rota está registrada corretamente.

---

### BUG-003: Endpoint register não existe
**Severidade:** Baixa  
**Arquivo provável:** `cmd/server/main.go` (rotas)  
**Reprodução:**
1. Executar `curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"email":"qa@pratoonline.com","password":"qa123456","name":"QA Tester"}'`
2. Obter erro: `404 page not found`

**Impacto:** Registro público não funciona (intencional - removido no Sprint 3).

**Correção sugerida:** Nenhuma (comportamento esperado - registro removido).

---

### BUG-004: Tabela platform_users não existe
**Severidade:** Alta  
**Arquivo provável:** `migrations/`  
**Reprodução:**
1. Executar `sqlite3 app.db "SELECT * FROM platform_users;"`
2. Obter erro: `Error: in prepare, no such table: platform_users`

**Impacto:** Plataforma de administração não funcional.

**Correção sugerida:** Verificar se migration para tabela platform_users foi executada ou se tabela foi renomeada.

---

## Resumo

| Bug ID | Severidade | Status |
|--------|------------|--------|
| BUG-001 | Alta | Aberto |
| BUG-002 | Média | Aberto |
| BUG-003 | Baixa | Fechado (comportamento esperado) |
| BUG-004 | Alta | Aberto |

---

## Assinatura

**Auditor:** Cascade AI  
**Data:** 2026-07-20
