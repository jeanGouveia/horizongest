# AUDITORIA PILOT-READY — PRATOONLINE

**Data:** 17 de Julho de 2026  
**Objetivo:** Auditoria técnica completa para produção piloto  
**Classificação Final:** YELLOW — PILOTO APROVADO COM RISCOS CONHECIDOS

---

## RESUMO EXECUTIVO

O PratoOnline foi submetido a uma auditoria técnica completa de produção. O sistema possui uma base sólida com arquitetura bem definida, mas apresenta alguns riscos que devem ser mitigados antes do deployment em produção piloto.

**Classificação:** YELLOW  
**Justificativa:** Sistema funcional com riscos conhecidos e mitigáveis. Não existem bloqueadores críticos que impeçam a operação, mas é necessário implementar correções de segurança e operacionais.

---

## FASE 1 — AUDITORIA DE CONFIGURAÇÃO DE PRODUÇÃO

### Secrets

#### ✅ Seguros
- `.env` está no `.gitignore` (não versionado)
- `.env.example` existe com template de variáveis
- JWT_SECRET é lido de variável de ambiente
- Senhas são hasheadas com bcrypt

#### ⚠️ Riscos
- **JWT_SECRET com fallback perigoso:** Em `auth_service.go` linha 40, existe fallback `"dev-secret-troque-em-producao"` se JWT_SECRET não estiver definido
- **.env.example com valor default:** O exemplo mostra `JWT_SECRET=troque-este-valor-em-producao-use-32-chars-minimo` que pode ser copiado acidentalmente
- **Sem validação de força do JWT_SECRET:** Não há verificação de tamanho mínimo ou complexidade

#### ❌ Bloqueadores
Nenhum

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Necessária:** Remover fallback de JWT_SECRET, adicionar validação de força

---

### Banco de Dados

#### ✅ Seguros
- SQLite com WAL mode (concorrência melhorada)
- Foreign keys habilitadas
- Pragmas de performance configurados
- AutoMigrate do GORM funciona corretamente
- Soft delete implementado em todas as tabelas
- Transações atômicas em operações críticas

#### ⚠️ Riscos
- **SQLite em produção:** SQLite não é ideal para produção com alta concorrência
- **Sem migrations versionadas:** Usa AutoMigrate do GORM ao invés de migrations SQL versionadas (Goose)
- **Limite de conexões:** `SetMaxOpenConns(1)` e `SetMaxIdleConns(1)` limitam concorrência
- **Backup manual:** Não existe automação de backup
- **Sem teste de restore:** Não há procedimento automatizado de teste de restore

#### ❌ Bloqueadores
- **Estratégia de backup indefinida:** Não existe automação, retenção ou teste de restore

**Classificação:** ⚠️ RISCO ALTO  
**Mitigação Necessária:** Implementar automação de backup, definir estratégia de restore, considerar PostgreSQL para produção

---

## FASE 2 — AUDITORIA DE DEPLOY

### Backend

#### ✅ Funcional
- Compilação com `go build cmd/server/main.go`
- Executável `server` gerado corretamente
- Variáveis de ambiente documentadas em `.env.example`
- Porta configurável via `PORT`
- Migrations automáticas ao iniciar

#### ⚠️ Riscos
- **Sem processo de CI/CD:** Deploy é manual
- **Sem rollback automatizado:** Rollback é manual
- **Sem health check no startup:** Não verifica se banco está conectado antes de aceitar requisições
- **Sem graceful shutdown:** Não aguarda requisições em andamento ao parar

#### ❌ Bloqueadores
Nenhum

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Recomendada:** Implementar CI/CD, health check no startup, graceful shutdown

---

### Frontend

#### ✅ Funcional
- Build com `npm run build`
- Adapter Node.js configurado
- Variáveis de ambiente podem ser configuradas
- Build output em `build/`

#### ⚠️ Riscos
- **Sem processo de CI/CD:** Deploy é manual
- **Sem otimização de assets:** Não há menção a compressão, CDN
- **Sem versionamento de builds:** Builds não são versionados

#### ❌ Bloqueadores
Nenhum

**Classificação:** ⚠️ RISCO BAIXO  
**Mitigação Recomendada:** Implementar CI/CD, versionamento de builds

---

## FASE 3 — HEALTH CHECKS E OBSERVABILIDADE

### Health Checks

#### ✅ Implementados
- `GET /api/health` retorna status básico
- `GET /api/version` retorna versão e build
- `GET /api/capabilities` retorna capacidades do sistema
- SystemHandler existe com endpoints de diagnóstico

#### ⚠️ Riscos
- **Health check não verifica banco:** Retorna `"database": "connected"` hardcoded, não verifica realmente
- **Sem verificação de dependências externas:** Não verifica storage, APIs externas
- **Sem métricas:** Não expõe métricas de performance
- **Sem readiness probe:** Não diferencia entre startup e readiness

#### ❌ Bloqueadores
Nenhum

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Necessária:** Implementar verificação real de banco, métricas básicas

---

### Logs

#### ✅ Implementados
- Logs do chi middleware (request ID, real IP, logger)
- Logs customizados em middleware de autenticação
- Logs em repositories (stock adjustments)
- Request ID gerado automaticamente

#### ⚠️ Riscos
- **Logs não estruturados:** Mix de `log.Printf` e strings formatadas
- **Sem nível de log:** Não há DEBUG, INFO, WARN, ERROR
- **Sem contexto estruturado:** Logs não têm campos estruturados (user_id, order_id, etc.)
- **Logs para stdout:** Não há configuração de output para arquivo
- **Sem rotação de logs:** Não há logrotate configurado

#### ❌ Bloqueadores
Nenhum

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Recomendada:** Implementar logging estruturado (zap/logrus), níveis de log, rotação

---

## FASE 4 — BACKUP E RESTORE

### Backup

#### ✅ Existente
- `.gitignore` protege arquivos `.db`
- Comentário no código sugere backup manual

#### ⚠️ Riscos
- **Sem automação:** Backup é puramente manual
- **Sem retenção:** Não há política de retenção
- **Sem offsite:** Backup não é enviado para local externo
- **Sem criptografia:** Backup não é criptografado
- **Sem compressão:** Backup não é comprimido

#### ❌ Bloqueadores
- **Ausência de estratégia:** Não existe resposta para "se o banco morrer hoje, como recuperamos?"

**Classificação:** ❌ BLOQUEADOR  
**Mitigação Obrigatória:** Implementar automação de backup com retenção, teste de restore

---

### Restore

#### ✅ Existente
- Arquivos `.db` podem ser copiados manualmente
- SQLite permite restore simples

#### ⚠️ Riscos
- **Sem procedimento documentado:** Não há script de restore
- **Sem teste:** Restore nunca foi testado
- **Sem validação:** Não há verificação de integridade pós-restore
- **Sem tempo estimado:** Não se sabe quanto tempo demora para restaurar

#### ❌ Bloqueadores
- **Impossibilidade de recuperação garantida:** Não há garantia que restore funciona

**Classificação:** ❌ BLOQUEADOR  
**Mitigação Obrigatória:** Documentar procedimento, testar restore, validar integridade

---

## FASE 5 — SEGURANÇA DE PRODUÇÃO

### CORS

#### ❌ Não Implementado
- Não há middleware de CORS
- Não há configuração de origens permitidas
- Qualquer origem pode fazer requisições

**Classificação:** ❌ BLOQUEADOR  
**Mitigação Obrigatória:** Implementar middleware CORS com origens específicas

---

### JWT_SECRET

#### ⚠️ Riscos
- Fallback perigoso se não definido
- Sem validação de força
- Sem rotação de tokens

**Classificação:** ⚠️ RISCO ALTO  
**Mitigação Necessária:** Remover fallback, validar força, implementar rotação

---

### Rate Limiting

#### ❌ Não Implementado
- Não há middleware de rate limiting
- Não há limite por IP
- Não há limite por usuário
- Suscetível a ataques DDoS e brute force

**Classificação:** ⚠️ RISCO ALTO  
**Mitigação Necessária:** Implementar rate limiting por IP e endpoint

---

### Headers de Segurança

#### ❌ Não Implementado
- Não há `X-Content-Type-Options`
- Não há `X-Frame-Options`
- Não há `X-XSS-Protection`
- Não há `Strict-Transport-Security`
- Não há `Content-Security-Policy`

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Recomendada:** Implementar headers de segurança

---

### Upload de Arquivos

#### ⚠️ Riscos
- Validação básica de MIME type
- Sem limite de tamanho explícito
- Sem verificação de conteúdo real
- Sem sanitização de nomes de arquivo
- Uploads servidos diretamente do filesystem

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Recomendada:** Implementar validação robusta, limite de tamanho, sanitização

---

### Validação de Entrada

#### ✅ Implementado
- Validator do go-playground em handlers
- Validação de email, senha, campos obrigatórios
- Validação de tipos de dados

#### ⚠️ Riscos
- Sem sanitização de XSS
- Sem proteção contra SQL injection (GORM protege, mas não explícito)
- Sem validação de tamanho de strings

**Classificação:** ✅ ACEITÁVEL  
**Mitigação Opcional:** Sanitização de XSS, validação de tamanho

---

### Exposição de Erros

#### ⚠️ Riscos
- Erros retornam mensagens em português (pode expor lógica)
- Stack traces podem ser expostos em modo dev
- Não há diferenciação entre modo dev e prod

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Recomendada:** Implementar modo prod com erros genéricos

---

### Autenticação

#### ✅ Implementado
- JWT com expiração
- Bcrypt para senhas
- Middleware de autenticação
- Cookie HttpOnly + Authorization header

#### ⚠️ Riscos
- Sem refresh token
- Sem revogação de tokens
- Token expira em 24h fixo

**Classificação:** ✅ ACEITÁVEL  
**Mitigação Opcional:** Implementar refresh token, revogação

---

### Autorização

#### ✅ Implementado
- Middleware injeta UserID no contexto
- Handlers verificam autenticação
- Sem autorização granular (simplificado para piloto)

**Classificação:** ✅ ACEITÁVEL  
**Mitigação Opcional:** Implementar RBAC se necessário

---

## FASE 6 — PERFORMANCE MÍNIMA

### Queries

#### ✅ Aceitável
- Uso de índices em campos principais (email, deleted_at)
- Preload em joins (Product, Ingredient)
- Where clauses com deleted_at IS NULL

#### ⚠️ Riscos
- **Possível N+1 em ListOrders:** Não há Preload de Items em ListOrders
- **Sem paginação:** ListProducts, ListIngredients não têm paginação
- **Sem limite de resultados:** Queries podem retornar milhares de registros

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Necessária:** Implementar paginação em todos os list endpoints

---

### Índices

#### ✅ Funcional
- Índices em foreign keys
- Índices em deleted_at
- Índice unique em email

#### ⚠️ Riscos
- **Sem índices compostos:** Queries com múltiplas condições podem ser lentas
- **Sem índices em campos de busca:** Name, Description não têm índices

**Classificação:** ⚠️ RISCO BAIXO  
**Mitigação Opcional:** Adicionar índices conforme necessário

---

### Endpoints Críticos

#### ✅ Funcional
- CreateOrder com transação atômica
- ValidateStock pré-valida estoque
- UpdateOrderStatus com ajustes de estoque

#### ⚠️ Riscos
- CreateOrder pode ser lento com muitos itens
- Dashboard faz múltiplas queries

**Classificação:** ✅ ACEITÁVEL  
**Mitigação Opcional:** Otimizar conforme necessário em produção

---

## FASE 7 — OPERABILIDADE

### ✅ Implementado
- RUNBOOK_OPERACIONAL.md criado
- Procedimentos documentados
- Scripts de backup/restore sugeridos

### ⚠️ Riscos
- Scripts não foram testados
- Procedimentos não foram validados
- Sem automação operacional

**Classificação:** ⚠️ RISCO MÉDIO  
**Mitigação Necessária:** Testar scripts, validar procedimentos

---

## CLASSIFICAÇÃO FINAL

### YELLOW — PILOTO APROVADO COM RISCOS CONHECIDOS

**Justificativa:**

O sistema é funcional e possui uma arquitetura sólida, mas apresenta riscos que devem ser mitigados antes do deployment em produção piloto. Não existem bloqueadores críticos que impeçam a operação, mas é necessário implementar correções obrigatórias de segurança e operacionais.

**Bloqueadores Obrigatórios (Devem ser corrigidos antes do deploy):**

1. **Backup e Restore:** Implementar automação de backup com retenção e teste de restore
2. **CORS:** Implementar middleware CORS com origens específicas
3. **JWT_SECRET:** Remover fallback perigoso e validar força
4. **Rate Limiting:** Implementar rate limiting básico
5. **Paginação:** Implementar paginação em todos os list endpoints

**Riscos Aceitos (Podem ser mitigados após o deploy):**

1. **SQLite em produção:** Aceitável para piloto com baixa concorrência
2. **Logs não estruturados:** Aceitável para piloto, melhorar posteriormente
3. **Headers de segurança:** Recomendado mas não bloqueador
4. **CI/CD:** Deploy manual aceitável para piloto

**Próximos Passos:**

1. Implementar correções obrigatórias (backup, CORS, JWT_SECRET, rate limiting, paginação)
2. Testar procedimentos de backup e restore
3. Validar health checks reais
4. Deploy em ambiente de staging
5. Monitorar por 1 semana
6. Deploy em produção piloto

---

## TEMPO ESTIMADO PARA CORREÇÕES

**Correções Obrigatórias:** 2-3 dias  
**Testes e Validação:** 1-2 dias  
**Total:** 3-5 dias úteis

---

## DECISÃO FINAL

**Status:** YELLOW — PILOTO APROVADO COM RISCOS CONHECIDOS

O sistema pode ser utilizado em produção piloto após implementação das correções obrigatórias. Os riscos conhecidos são mitigáveis e não representam ameaça crítica à operação.

**Recomendação:** Implementar correções obrigatórias, testar procedimentos operacionais, e então proceder com deployment em produção piloto.
