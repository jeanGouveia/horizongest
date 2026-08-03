# FASE A — BLOCKERS DE PRODUÇÃO
## Roadmap Hardening - HorizonGest

**Data**: 01/08/2026  
**Arquiteto**: Senior Software Architect (Go, Clean Architecture, DDD, OWASP ASVS)  
**Objetivo**: Transformar o HorizonGest em sistema pronto para produção

---

## EXECUTIVO

**Status**: ❌ REPROVADA

**Percentual de Conclusão**: 45%

**Bloqueadores Críticos Restantes**: 12

---

## A1 — SEGURANÇA

### Status: ❌ REPROVADA (50% completo)

#### ✅ IMPLEMENTADOS

1. **JWT Validation**
   - Algoritmo validado corretamente (HMAC-SHA256)
   - Expiração verificada com `jwt.WithExpirationRequired()`
   - Token blacklist implementada via banco de dados
   - Arquivo: `internal/service/auth_service.go` (linhas 195-214)

2. **Session Revocation (Tenant)**
   - Logout funcional com token blacklist
   - Tokens revogados persistidos no banco
   - Arquivo: `internal/service/auth_service.go` (linhas 225-244)

3. **Login Rate Limiting**
   - Rate limiter implementado (5 req/min IP, 30 req/hour usuário)
   - Aplicado em endpoints de login
   - Arquivo: `internal/middleware/rate_limiter.go`

4. **Security Headers**
   - X-Frame-Options: DENY
   - X-Content-Type-Options: nosniff
   - Referrer-Policy: strict-origin-when-cross-origin
   - Permissions-Policy implementado
   - Arquivo: `internal/middleware/security_headers.go`

5. **CORS**
   - Whitelist de origens por ambiente
   - Não reflete Origin (security best practice)
   - Arquivo: `internal/middleware/cors.go`

6. **HSTS**
   - Ativado apenas em produção
   - max-age=31536000; includeSubDomains; preload
   - Arquivo: `internal/middleware/security_headers.go` (linhas 44-47)

#### ❌ BLOQUEADORES CRÍTICOS

##### B1 - JWT Secret Fraco (CRÍTICO)
**Problema**: Secrets JWT são fracos/placeholder no `.env.example`

**Risco**: 
- Secrets padrão podem ser usados em produção
- Força bruta possível com secrets curtos
- Comprometimento total do sistema de autenticação

**Arquivo**: `backend/.env.example` (linhas 9-10)
```env
JWT_PLATFORM_SECRET=troque-este-valor-em-producao-platform-use-32-chars-minimo
JWT_TENANT_SECRET=troque-este-valor-em-producao-tenant-use-32-chars-minimo
```

**Implementação Necessária**:
1. Remover secrets placeholder do `.env.example`
2. Documentar requisito de mínimo 32 caracteres aleatórios
3. Adicionar validação de entropia mínima
4. Gerar secrets via script seguro em produção

**Impacto**: Alto - Autenticação comprometida

---

##### B2 - JWT Rotation Ausente (CRÍTICO)
**Problema**: Não existe mecanismo de rotação de chaves JWT

**Risco**:
- Impossível rotacionar chaves comprometidas sem derrubar todos os usuários
- Viola OWASP ASVS 2.10.2
- Não atende requisitos de compliance

**Arquivos**: `internal/service/auth_service.go`, `internal/service/platform_auth_service.go`

**Implementação Necessária**:
1. Implementar JWT kid (key ID) header
2. Suportar múltiplas chaves ativas simultaneamente
3. Mecanismo de rotação gradual (grace period)
4. Key store seguro (Vault/KMS)

**Impacto**: Alto - Impossível responder a incidentes de segurança

---

##### B3 - JWT kid Ausente (CRÍTICO)
**Problema**: Tokens JWT não possuem header `kid` (Key ID)

**Risco**:
- Impossível identificar qual chave assinou o token
- Bloqueia rotação de chaves
- Viola OWASP ASVS 2.10.2

**Arquivos**: `internal/service/auth_service.go` (linha 304), `internal/service/platform_auth_service.go` (linha 188)

**Implementação Necessária**:
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
token.Header["kid"] = "key-2024-08-01" // Key identifier
signed, err := token.SignedString(s.secret)
```

**Impacto**: Alto - Bloqueia rotação de chaves

---

##### B4 - Platform Session Validation Ausente (CRÍTICO)
**Problema**: Platform auth service não valida sessão no banco

**Risco**:
- Tokens platform JWT permanecem válidos mesmo após logout
- Session hijacking possível
- Account takeover via token reuso

**Arquivo**: `internal/service/platform_auth_service.go` (linhas 169-171)
```go
// Check if session exists and is valid
// Note: This would require passing context, but for simplicity we skip here
// In production, you should validate the session in the database
```

**Implementação Necessária**:
1. Validar sessão no banco em cada request
2. Implementar session blacklist similar ao tenant
3. Adicionar session expiration check

**Impacto**: Crítico - Account Takeover possível

---

##### B5 - Password Reset Race Condition (ALTO)
**Problema**: Password reset não usa transação nem SELECT FOR UPDATE

**Risco**:
- Race condition permite reuso do mesmo token múltiplas vezes
- Token pode ser usado concorrentemente
- Viola integridade do fluxo de recuperação

**Arquivo**: `internal/service/auth_service.go` (linhas 360-408)

**Implementação Necessária**:
1. Envolver operação em transação
2. Usar SELECT FOR UPDATE no token
3. Validar status do token dentro da transação
4. Marcar como usado atomicamente

**Impacto**: Alto - Violação de integridade de segurança

---

##### B6 - Accept Invitation Race Condition (ALTO)
**Problema**: Accept invitation não usa transação nem SELECT FOR UPDATE

**Risco**:
- Mesmo convite pode ser aceito múltiplas vezes
- Usuário pode ser associado a múltiplas empresas
- Viola regras de negócio de multi-tenancy

**Arquivo**: `internal/service/invitation_service.go` (linhas 216-277)

**Implementação Necessária**:
1. Transação atômica em AcceptInvitation
2. SELECT FOR UPDATE no invitation
3. Validar status do usuário dentro da transação
4. Atualizar invitation e usuário atomicamente

**Impacto**: Alto - Tenant escape possível

---

##### B7 - Media Upload Sem Tenant Validation (MÉDIO)
**Problema**: Upload de mídia não valida tenant context no handler

**Risco**:
- Upload possível sem autenticação adequada
- IDOR em endpoints de mídia
- Data leakage entre tenants

**Arquivo**: `internal/handler/media_handler.go` (linhas 25-103)

**Implementação Necessária**:
1. Adicionar tenant middleware nas rotas de upload
2. Validar CompanyID no handler
3. Garantir que entity_id pertence ao tenant

**Impacto**: Médio - Data leakage

---

##### B8 - Media Serve Directory Traversal (MÉDIO)
**Problema**: Validação de path insuficiente em ServeFile

**Risco**:
- Directory traversal via path manipulation
- Acesso a arquivos fora do uploads/
- Information disclosure

**Arquivo**: `internal/handler/media_handler.go` (linhas 167-187)
```go
if strings.Contains(filePath, "..") {
    http.Error(w, "caminho inválido", http.StatusBadRequest)
    return
}
```

**Implementação Necessária**:
1. Usar filepath.Clean() e validar
2. Verificar se path está dentro do uploads/
3. Usar http.ServeFile com validação rigorosa
4. Considerar usar storage service (S3/GCS)

**Impacto**: Médio - Information disclosure

---

##### B9 - CSP unsafe-inline (MÉDIO)
**Problema**: Content-Security-Policy permite unsafe-inline e unsafe-eval

**Risco**:
- XSS via scripts inline
- Ataques via eval()
- Viola OWASP ASVS 5.5.3

**Arquivo**: `internal/middleware/security_headers.go` (linhas 23-30)
```go
w.Header().Set("Content-Security-Policy",
    "default-src 'self'; "+
        "script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
        "style-src 'self' 'unsafe-inline'; "+
```

**Implementação Necessária**:
1. Migrar para hash-based CSP ou nonce-based
2. Remover unsafe-inline e unsafe-eval
3. Implementar CSP report-only mode primeiro
4. Testar frontend extensivamente

**Impacto**: Médio - XSS vulnerability

---

##### B10 - CSRF Protection Ausente (MÉDIO)
**Problema**: Não existe proteção CSRF

**Risco**:
- CSRF attacks em endpoints state-changing
- Ações não autorizadas em nome do usuário
- Viola OWASP ASVS 8.3.1

**Arquivos**: Todos os handlers POST/PUT/DELETE

**Implementação Necessária**:
1. Implementar CSRF tokens para state-changing operations
2. Validar CSRF token em middleware
3. Usar SameSite cookies
4. Considerar double-submit cookie pattern

**Impacto**: Médio - CSRF vulnerability

---

##### B11 - JWT Logado em Plaintext (CRÍTICO)
**Problema**: JWT token é logado em plaintext em impersonation

**Risco**:
- Tokens expostos em logs
- Logs podem ser comprometidos
- Account takeover via log leakage

**Arquivo**: `internal/service/impersonation_service.go` (linha 113)
```go
// FORENSIC: Log JWT generated
log.Printf("[FORENSIC] StartImpersonation - JWT gerado: %s", token)
```

**Implementação Necessária**:
1. Remover log de token completo
2. Log apenas hash/token prefixo
3. Garantir que nenhum token seja logado em plaintext

**Impacto**: Crítico - Credential leakage via logs

---

##### B12 - Secrets em .env.example (ALTO)
**Problema**: Secrets placeholder em arquivo versionado

**Risco**:
- Desenvolvedores podem usar secrets padrão
- Secrets expostos em version control
- Falsa sensação de segurança

**Arquivo**: `backend/.env.example`

**Implementação Necessária**:
1. Remover todos os secrets do .env.example
2. Usar placeholders como `JWT_PLATFORM_SECRET=<generate-secure-secret>`
3. Documentar processo de geração de secrets
4. Adicionar script de geração de secrets

**Impacto**: Alto - Credential exposure

---

## A2 — MULTI-TENANCY

### Status: ✅ APROVADA (95% completo)

#### ✅ IMPLEMENTADOS

1. **Tenant Context**
   - TenantContext populado no middleware
   - CompanyID sempre requerido (Sprint 3)
   - Arquivo: `internal/middleware/tenant_middleware.go`

2. **Tenant Filtering**
   - ApplyTenantFilter em todos os repositórios
   - ApplyTenantFilterWithID para queries por ID
   - Arquivo: `internal/infra/repository/tenant_helper.go`

3. **CompanyID Validation**
   - CompanyID NOT NULL em todos os models
   - Auto-fill do tenant context em creates
   - Validação em handlers

4. **IDOR Protection**
   - Helper ValidateCompanyOwnership
   - Validação em handlers críticos
   - Arquivo: `internal/middleware/resource_ownership.go`

#### ⚠️ PENDENTE (NÃO CRÍTICO)

1. **Media Upload Tenant Validation**: Ver B7 acima

---

## A3 — BANCO

### Status: ✅ APROVADA (90% completo)

#### ✅ IMPLEMENTADOS

1. **Transações em Operações Críticas**
   - Stock movements usam transação
   - SELECT FOR UPDATE implementado
   - Arquivo: `internal/service/stock_movement_service.go` (linhas 61-100)

2. **ACID**
   - Transações GORM usadas corretamente
   - Rollback automático em erro
   - Commit explícito em sucesso

3. **Locks**
   - SELECT FOR UPDATE em stock movements
   - Previne race conditions em estoque

#### ⚠️ PENDENTE (NÃO CRÍTICO)

1. **Password Reset Transaction**: Ver B5 acima
2. **Accept Invitation Transaction**: Ver B6 acima

---

## A4 — INFRA

### Status: ❌ REPROVADA (40% completo)

#### ✅ IMPLEMENTADOS

1. **Environment Variables**
   - Validação de secrets em produção
   - Fallback seguro em desenvolvimento
   - Arquivo: `cmd/server/main.go` (linhas 112-134)

2. **Docker**
   - Docker-compose básico implementado
   - Redis com healthcheck
   - Arquivo: `docker-compose.yml`

#### ❌ BLOQUEADORES CRÍTICOS

##### B13 - TLS/HTTPS Não Enforcado (CRÍTICO)
**Problema**: Servidor HTTP não enforça TLS/HTTPS

**Risco**:
- Tráfego não criptografado
- Credentials em plaintext
- Viola OWASP ASVS 9.2.1

**Arquivo**: `cmd/server/main.go` (linha 621)
```go
if err := http.ListenAndServe(":"+port, r); err != nil {
```

**Implementação Necessária**:
1. Usar http.ListenAndServeTLS em produção
2. Configurar certificados TLS
3. Enforçar HTTPS via middleware
4. Redirecionar HTTP para HTTPS

**Impacto**: Crítico - Credentials em trânsito

---

##### B14 - Docker Security Insuficiente (ALTO)
**Problema**: Docker-compose sem hardening de segurança

**Risco**:
- Containers rodando como root
- Sem resource limits
- Sem security options

**Arquivo**: `docker-compose.yml`

**Implementação Necessária**:
1. Adicionar user não-root nos containers
2. Configurar resource limits (CPU, memory)
3. Adicionar security options (no-new-privileges)
4. Usar networks isoladas
5. Readonly filesystem onde possível

**Impacto**: Alto - Container escape

---

##### B15 - CORS Configuration (MÉDIO)
**Problema**: CORS permite localhost em staging/production

**Risco**:
- Development configs podem vazar para produção
- Origens não autorizadas podem acessar

**Arquivo**: `internal/middleware/cors.go` (linhas 69-78)

**Implementação Necessária**:
1. Remover fallback localhost em production
2. Exigir CORS_ALLOWED_ORIGINS em produção
3. Validar formato das origens

**Impacto**: Médio - Data leakage

---

## A5 — OBSERVABILIDADE

### Status: ❌ REPROVADA (30% completo)

#### ✅ IMPLEMENTADOS

1. **Request ID**
   - chimiddleware.RequestID implementado
   - Arquivo: `cmd/server/main.go` (linha 310)

2. **Health Check**
   - Endpoint /api/health implementado
   - Arquivo: `cmd/server/main.go` (linhas 319-323)

#### ❌ BLOQUEADORES CRÍTICOS

##### B16 - Logs Estruturados Ausentes (ALTO)
**Problema**: Logs não são estruturados (JSON)

**Risco**:
- Difícil parsear e analisar logs
- Impossível integrar com SIEM
- Difícil troubleshooting em produção

**Arquivos**: Todos os arquivos com log.Printf

**Implementação Necessária**:
1. Implementar logger estruturado (zap/logrus)
2. Logs em formato JSON
3. Campos padronizados (timestamp, level, request_id, user_id, company_id)
4. Context propagation

**Impacto**: Alto - Operability em produção

---

##### B17 - Audit Logs Incompletos (ALTO)
**Problema**: Audit logs limitados a algumas operações

**Risco**:
- Rastro de auditoria incompleto
- Impossível investigar incidentes
- Não atende compliance

**Arquivos**: `internal/service/impersonation_service.go` (partial)

**Implementação Necessária**:
1. Audit trail para todas as operações críticas
2. Log de: who, what, when, where, why
3. Audit logs imutáveis
4. Integração com SIEM

**Impacto**: Alto - Compliance e forensics

---

##### B18 - Sensitive Data em Logs (CRÍTICO)
**Problema**: Sensitive data logado (JWT, passwords, tokens)

**Risco**:
- Credentials expostas em logs
- Logs podem ser comprometidos
- Viola OWASP ASVS 7.3.1

**Arquivos**: 
- `internal/service/impersonation_service.go` (linha 113)
- `cmd/server/main.go` (linha 80)

**Implementação Necessária**:
1. Remover todos os logs de sensitive data
2. Implementar data masking automático
3. Lista negra de campos sensíveis
4. Sanitização de logs

**Impacto**: Crítico - Credential leakage

---

##### B19 - Tracing Ausente (MÉDIO)
**Problema**: Distributed tracing não implementado

**Risco**:
- Difícil debuggar requests multi-serviço
- Performance issues difíceis de diagnosticar
- Difícil entender latência

**Implementação Necessária**:
1. Implementar OpenTelemetry
2. Tracing spans para operações críticas
3. Propagação de trace context
4. Integração com Jaeger/Tempo

**Impacto**: Médio - Operability

---

##### B20 - Métricas Ausentes (MÉDIO)
**Problema**: Métricas de aplicação não implementadas

**Risco**:
- Difícil monitorar saúde do sistema
- Impossível alertar proativamente
- Difícil entender performance

**Implementação Necessária**:
1. Implementar Prometheus metrics
2. Métricas de: request rate, error rate, latency
3. Business metrics (stock, orders, users)
4. Exposição em /metrics

**Impacto**: Médio - Operability

---

##### B21 - Correlation ID Incompleto (MÉDIO)
**Problema**: Correlation ID não propagado consistentemente

**Risco**:
- Difícil rastrear requests através do sistema
- Troubleshooting dificultado

**Implementação Necessária**:
1. Propagar correlation ID em todos os serviços
2. Incluir em todos os logs
3. Passar para serviços externos (RabbitMQ, Redis)

**Impacto**: Médio - Operability

---

## RESUMO DOS BLOQUEADORES

### CRÍTICOS (Must Fix Before Production)
1. B1 - JWT Secret Fraco
2. B2 - JWT Rotation Ausente
3. B3 - JWT kid Ausente
4. B4 - Platform Session Validation Ausente
5. B11 - JWT Logado em Plaintext
6. B13 - TLS/HTTPS Não Enforcado
7. B18 - Sensitive Data em Logs

### ALTOS (Should Fix Before Production)
8. B5 - Password Reset Race Condition
9. B6 - Accept Invitation Race Condition
10. B12 - Secrets em .env.example
11. B14 - Docker Security Insuficiente
12. B16 - Logs Estruturados Ausentes
13. B17 - Audit Logs Incompletos

### MÉDIOS (Fix Before Production)
14. B7 - Media Upload Sem Tenant Validation
15. B8 - Media Serve Directory Traversal
16. B9 - CSP unsafe-inline
17. B10 - CSRF Protection Ausente
18. B15 - CORS Configuration
19. B19 - Tracing Ausente
20. B20 - Métricas Ausentes
21. B21 - Correlation ID Incompleto

---

## PARECER FINAL

### FASE A: ❌ REPROVADA

**Motivo**: Existem 7 bloqueadores críticos que impedem a produção segura do sistema:

1. **B1, B2, B3, B4, B11**: Vulnerabilidades críticas em JWT que podem levar a account takeover e session hijacking
2. **B13**: Ausência de TLS/HTTPS expõe credentials em trânsito
3. **B18**: Sensitive data em logs expõe credentials via logging

**Recomendação**: Corrigir todos os bloqueadores críticos antes de considerar produção. Os bloqueadores altos e médios devem ser priorizados subsequentemente.

---

## PRÓXIMOS PASSOS

1. **Imediato (Críticos)**:
   - Implementar TLS/HTTPS
   - Corrigir logging de sensitive data
   - Implementar platform session validation
   - Fortalecer JWT secrets e implementar rotation

2. **Curto Prazo (Altos)**:
   - Corrigir race conditions em password reset e invitation
   - Implementar logs estruturados
   - Completar audit logs
   - Hardening Docker

3. **Médio Prazo (Médios)**:
   - Implementar CSRF protection
   - Melhorar CSP
   - Implementar tracing e métricas
   - Completar correlation ID propagation

---

**Assinatura**: Arquiteto de Software Sênior  
**Data**: 01/08/2026
