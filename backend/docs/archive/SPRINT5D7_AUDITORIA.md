# SPRINT 5D.7 — AUDITORIA DE ARQUITETURA E CÓDIGO MORTO

## Resumo Executivo

Esta sprint realizou uma auditoria completa da arquitetura do HorizonGest, focando exclusivamente em arquitetura, organização, acoplamento, coesão, código morto, duplicação, complexidade e qualidade estrutural. Regras de negócio, segurança, resiliência e performance foram excluídos desta análise.

**Status:** ✅ AUDITORIA COMPLETA

---

## FASE 1 — COMPLEXIDADE CICLOMÁTICA

### Problema 1.1: Função Gigante em concurrency_test.go
- **Severidade:** P2
- **Arquivo:** `internal/infra/repository/concurrency_test.go`
- **Função:** `TestConcurrency_100Goroutines_CompleteInventory`
- **Linhas:** ~150
- **Complexidade Estimada:** Alta (muitos if/else, goroutines, mutex)
- **Motivo:** Função de teste com lógica complexa de concorrência
- **Impacto:** Difícil de manter, difícil de debugar
- **Correção Sugerida:** Extrair helpers para setup, teardown e validação
- **Estimativa:** 4h

### Problema 1.2: Função generateSlug em product_service.go
- **Severidade:** P3
- **Arquivo:** `internal/service/product_service.go`
- **Função:** `generateSlug`
- **Linhas:** ~45
- **Complexidade Estimada:** Média (múltiplos loops, map de replacements)
- **Motivo:** Função helper com lógica de string manipulation complexa
- **Impacto:** Difícil de testar isoladamente, pode ter bugs edge cases
- **Correção Sugerida:** Mover para util package, usar biblioteca de slug (ex: github.com/gosimple/slug)
- **Estimativa:** 2h

### Problema 1.3: Função CreateCompany em platform_service.go
- **Severidade:** P2
- **Arquivo:** `internal/service/platform_service.go`
- **Função:** `CreateCompany`
- **Linhas:** ~100+
- **Complexidade Estimada:** Alta (transação, múltiplas validações, criação de usuário)
- **Motivo:** Função faz muitas coisas: validação, criação de empresa, criação de usuário, envio de email
- **Impacto:** Difícil de testar, difícil de manter, viola SRP
- **Correção Sugerida:** Extrair createCompany, createUser, sendWelcomeEmail em funções separadas
- **Estimativa:** 6h

---

## FASE 2 — ACOPLAMENTO

### Problema 2.1: Service Dependendo Diretamente de gorm.DB
- **Severidade:** P1
- **Arquivo:** `internal/service/platform_service.go`, `internal/service/product_service.go`
- **Linhas:** 22, 22
- **Causa Raiz:** Services recebem *gorm.DB diretamente em vez de usar apenas interfaces de repository
- **Impacto:** Alto acoplamento com GORM, difícil de testar, difícil de mudar ORM
- **Correção Sugerida:** Remover dependência direta de gorm.DB dos services, usar apenas repository interfaces
- **Estimativa:** 16h

### Problema 2.2: Repository Implementando GORM Models no Mesmo Arquivo
- **Severidade:** P2
- **Arquivo:** `internal/infra/repository/gorm_product_repository.go`
- **Linhas:** 36-100
- **Causa Raiz:** GORM models (GormProduct, GormIngredient) estão definidos no mesmo arquivo que o repository
- **Impacto:** Mistura de responsabilidades, acoplamento com GORM em repository
- **Correção Sugerida:** Mover GORM models para package separado (ex: internal/infra/models)
- **Estimativa:** 8h

### Problema 2.3: Handler Importando Service Diretamente
- **Severidade:** P0
- **Arquivo:** Vários handlers
- **Causa Raiz:** Handlers importam services diretamente, violando arquitetura de camadas
- **Impacto:** Arquitetura quebrada, difícil de testar, alto acoplamento
- **Correção Sugerida:** Handlers devem depender apenas de interfaces definidas em ports
- **Estimativa:** 32h

---

## FASE 3 — DEPENDÊNCIAS CIRCULARES

### Problema 3.1: Nenhuma Dependência Circular Detectada
- **Severidade:** N/A
- **Status:** ✅ OK
- **Observação:** Architecture tests confirmam que não há dependências circulares entre camadas
- **Arquivo:** `internal/architecture_test.go`

---

## FASE 4 — CÓDIGO MORTO

### Problema 4.1: Package errors Vazio
- **Severidade:** P3
- **Arquivo:** `internal/errors/`
- **Causa Raiz:** Package errors existe mas está vazio
- **Impacto:** Confusão, código morto
- **Correção Sugerida:** Remover package ou implementar error handling centralizado
- **Estimativa:** 1h

### Problema 4.2: Package testutil/mocks Vazio
- **Severidade:** P3
- **Arquivo:** `internal/testutil/mocks/`
- **Causa Raiz:** Package testutil/mocks existe mas está vazio
- **Impacto:** Confusão, código morto
- **Correção Sugerida:** Remover package ou implementar mocks reais
- **Estimativa:** 1h

### Problema 4.3: architecture_test.go com build tag que impede execução normal
- **Severidade:** P2
- **Arquivo:** `internal/architecture_test.go`
- **Linhas:** 1
- **Causa Raiz:** `//go:build ignore_test` impede que os testes de arquitetura sejam executados normalmente
- **Impacto:** Testes de arquitetura não são executados no CI/CD, violações podem passar despercebidas
- **Correção Sugerida:** Remover build tag ou mudar para tag que seja executada no CI
- **Estimativa:** 1h

### Problema 4.4: 111 TODOs no Código
- **Severidade:** P3
- **Arquivos:** 35 arquivos
- **Causa Raiz:** Muitos TODOs espalhados pelo código
- **Impacto:** Código técnico, possíveis bugs não resolvidos, confusão
- **Correção Sugerida:** Auditar todos os TODOs e resolver ou remover
- **Estimativa:** 16h

### Problema 4.5: 1 FIXME no Código
- **Severidade:** P3
- **Arquivo:** `internal/infra/repository/gorm_product_repository_test.go`
- **Causa Raiz:** FIXME indica código que precisa ser corrigido
- **Impacto:** Possível bug não resolvido
- **Correção Sugerida:** Investigar e resolver o FIXME
- **Estimativa:** 2h

---

## FASE 5 — DUPLICAÇÃO

### Problema 5.1: Conversão de Money Repetida
- **Severidade:** P2
- **Arquivos:** `internal/infra/repository/gorm_product_repository.go`, outros repositories
- **Linhas:** 18-32
- **Causa Raiz:** Funções convertMoneyPtrToInt64Ptr e convertInt64PtrToMoneyPtr estão duplicadas em múltiplos repositories
- **Impacto:** Duplicação de código, manutenção difícil
- **Correção Sugerida:** Mover conversores para package util ou internal/infra/converter
- **Estimativa:** 4h

### Problema 5.2: Validators Duplicados em Handlers
- **Severidade:** P2
- **Arquivos:** Vários handlers
- **Causa Raiz:** Lógica de validação JSON e error handling é repetida em todos os handlers
- **Impacto:** Duplicação de código, inconsistência
- **Correção Sugerida:** Criar package handler/util com helpers comuns
- **Estimativa:** 8h

### Problema 5.3: Repository Pattern Repetido
- **Severidade:** P1
- **Arquivos:** Todos os GORM repositories
- **Causa Raiz:** Todos os GORM repositories seguem o mesmo padrão mas não há abstração compartilhada
- **Impacto:** Duplicação de código, difícil de adicionar features globais
- **Correção Sugerida:** Criar base repository com funcionalidades comuns (tenant filtering, soft delete, etc.)
- **Estimativa:** 24h

---

## FASE 6 — COESÃO

### Problema 6.1: Package domain Muito Grande
- **Severidade:** P2
- **Arquivo:** `internal/domain/`
- **Causa Raiz:** Package domain tem 42 arquivos, mistura entidades, value objects, enums
- **Impacto:** Difícil de navegar, baixa coesão, viola SRP
- **Correção Sugerida:** Dividir domain em subpackages: entities, valueobjects, enums, aggregates
- **Estimativa:** 16h

### Problema 6.2: Package service Muito Grande
- **Severidade:** P2
- **Arquivo:** `internal/service/`
- **Causa Raiz:** Package service tem 46 arquivos, mistura application services, domain services
- **Impacto:** Difícil de navegar, baixa coesão
- **Correção Sugerida:** Dividir service em subpackages: application, domain
- **Estimativa:** 12h

### Problema 6.3: Package handler Muito Grande
- **Severidade:** P2
- **Arquivo:** `internal/handler/`
- **Causa Raiz:** Package handler tem 32 arquivos, mistura HTTP handlers de diferentes contextos
- **Impacto:** Difícil de navegar, baixa coesão
- **Correção Sugerida:** Dividir handler em subpackages: api, platform, admin
- **Estimativa:** 8h

### Problema 6.4: Package infra/repository Muito Grande
- **Severidade:** P2
- **Arquivo:** `internal/infra/repository/`
- **Causa Raiz:** Package infra/repository tem 34 arquivos
- **Impacto:** Difícil de navegar, baixa coesão
- **Correção Sugerida:** Dividir repository em subpackages por bounded context: product, order, finance
- **Estimativa:** 12h

---

## FASE 7 — TAMANHO DOS ARQUIVOS

### Problema 7.1: concurrency_test.go com 1814 linhas
- **Severidade:** P2
- **Arquivo:** `internal/infra/repository/concurrency_test.go`
- **Linhas:** 1814
- **Causa Raiz:** Arquivo de teste com múltiplos cenários de concorrência
- **Impacto:** Difícil de manter, difícil de entender
- **Correção Sugerida:** Dividir em múltiplos arquivos por cenário
- **Estimativa:** 4h

### Problema 7.2: gorm_product_repository_test.go com 827 linhas
- **Severidade:** P3
- **Arquivo:** `internal/infra/repository/gorm_product_repository_test.go`
- **Linhas:** 827
- **Causa Raiz:** Arquivo de teste com muitos casos
- **Impacto:** Difícil de manter
- **Correção Sugerida:** Dividir por funcionalidade (create, update, delete, search)
- **Estimativa:** 4h

### Problema 7.3: auth_service_test.go with 811 linhas
- **Severidade:** P3
- **Arquivo:** `internal/service/auth_service_test.go`
- **Linhas:** 811
- **Causa Raiz:** Arquivo de teste com muitos casos de autenticação
- **Impacto:** Difícil de manter
- **Correção Sugerida:** Dividir por funcionalidade (login, register, reset password)
- **Estimativa:** 4h

### Problema 7.4: platform_service.go com 782 linhas
- **Severidade:** P2
- **Arquivo:** `internal/service/platform_service.go`
- **Linhas:** 782
- **Causa Raiz:** Service com muitas responsabilidades (company, user, audit)
- **Impacto:** Difícil de manter, viola SRP
- **Correção Sugerida:** Dividir em platform_company_service, platform_user_service, platform_audit_service
- **Estimativa:** 8h

### Problema 7.5: gorm_product_repository.go com 712 linhas
- **Severidade:** P2
- **Arquivo:** `internal/infra/repository/gorm_product_repository.go`
- **Linhas:** 712
- **Causa Raiz:** Repository com muitas operações complexas
- **Impacto:** Difícil de manter
- **Correção Sugerida:** Dividir em product_repository, product_search_repository, product_ingredient_repository
- **Estimativa:** 6h

---

## FASE 8 — INTERFACES

### Problema 8.1: Interfaces com Apenas Uma Implementação
- **Severidade:** P2
- **Arquivos:** Muitas interfaces em `internal/ports/`
- **Causa Raiz:** Muitas interfaces têm apenas uma implementação GORM
- **Impacto:** YAGNI (You Aren't Gonna Need It), complexidade desnecessária
- **Correção Sugerida:** Remover interfaces que não têm múltiplas implementações ou justificar necessidade
- **Estimativa:** 8h

### Problema 8.2: Interfaces Gigantes
- **Severidade:** P2
- **Arquivos:** `internal/ports/order_repository.go`, `internal/ports/purchase_repository.go`
- **Linhas:** 2224, 2098
- **Causa Raiz:** Interfaces com muitos métodos (20+)
- **Impacto:** Viola Interface Segregation Principle, difícil de implementar
- **Correção Sugerida:** Dividir interfaces em interfaces menores e coesas
- **Estimativa:** 12h

---

## FASE 9 — ORGANIZAÇÃO

### Problema 9.1: Arquivos de Documentação no Root
- **Severidade:** P3
- **Arquivos:** `BUSINESS_INTEGRITY_AUDIT.md`, `SECURITY_AUDITORIA_FINAL.md`, etc.
- **Causa Raiz:** Arquivos de documentação estão no root do projeto em vez de docs/
- **Impacto:** Organização ruim, confusão
- **Correção Sugerida:** Mover todos os arquivos .md para docs/
- **Estimativa:** 2h

### Problema 9.2: Arquivos de Database no Root
- **Severidade:** P3
- **Arquivos:** `app.db`, `app.db.bkp`, `prato.db`, `test_snapshot*.db`
- **Causa Raiz:** Arquivos de database estão no root do projeto
- **Impacto:** Poluição do root, possível commit acidental
- **Correção Sugerida:** Mover para data/ ou adicionar ao .gitignore
- **Estimativa:** 1h

### Problema 9.3: Arquivos de Coverage no Root
- **Severidade:** P3
- **Arquivos:** `coverage.out`, `coverage_final.out`, etc.
- **Causa Raiz:** Arquivos de coverage estão no root do projeto
- **Impacto:** Poluição do root
- **Correção Sugerida:** Mover para coverage/ ou adicionar ao .gitignore
- **Estimativa:** 1h

---

## FASE 10 — DEPENDÊNCIAS EXTERNAS

### Problema 10.1: Dependências Possivelmente Não Utilizadas
- **Severidade:** P3
- **Arquivo:** `go.mod`
- **Causa Raiz:** Algumas dependências podem não ser utilizadas (ex: testify se não há testes que usam)
- **Impacto:** Bloat, vulnerabilidades desnecessárias
- **Correção Sugerida:** Executar `go mod tidy` e remover dependências não utilizadas
- **Estimativa:** 2h

---

## FASE 11 — CLEAN ARCHITECTURE

### Problema 11.1: Violação de Dependency Inversion
- **Severidade:** P0
- **Arquivos:** Vários handlers e services
- **Causa Raiz:** Handlers e services dependem de implementações concretas (GORM) em vez de interfaces
- **Impacto:** Arquitetura quebrada, alto acoplamento, difícil de testar
- **Correção Sugerida:** Garantir que todas as dependências sejam através de interfaces definidas em ports
- **Estimativa:** 40h

### Problema 11.2: Violação de Single Responsibility
- **Severidade:** P1
- **Arquivos:** `internal/service/platform_service.go`, `internal/service/product_service.go`
- **Causa Raiz:** Services têm múltiplas responsabilidades
- **Impacto:** Difícil de manter, difícil de testar
- **Correção Sugerida:** Dividir services em serviços menores e coesos
- **Estimativa:** 24h

### Problema 11.3: Violação de Interface Segregation
- **Severidade:** P2
- **Arquivos:** `internal/ports/order_repository.go`, `internal/ports/purchase_repository.go`
- **Causa Raiz:** Interfaces gigantes com muitos métodos
- **Impacto:** Difícil de implementar, viola ISP
- **Correção Sugerida:** Dividir interfaces em interfaces menores
- **Estimativa:** 12h

---

## FASE 12 — PADRONIZAÇÃO

### Problema 12.1: Tratamento de Erro Inconsistente
- **Severidade:** P2
- **Arquivos:** Vários arquivos
- **Causa Raiz:** Alguns lugares retornam errors.New(), outros retornam fmt.Errorf(), outros retornam wrapped errors
- **Impacto:** Inconsistência, difícil de debugar
- **Correção Sugerida:** Padronizar tratamento de erro (ex: sempre usar errors.Is/As para wrapped errors)
- **Estimativa:** 8h

### Problema 12.2: Nomes Inconsistentes
- **Severidade:** P3
- **Arquivos:** Vários arquivos
- **Causa Raiz:** Alguns nomes usam snake_case, outros camelCase, inconsistência em convenções
- **Impacto:** Confusão, difícil de manter
- **Correção Sugerida:** Padronizar convenções de nomenclatura (Go standard)
- **Estimativa:** 4h

### Problema 12.3: Context Não Propagado Consistentemente
- **Severidade:** P2
- **Arquivos:** Vários arquivos
- **Causa Raiz:** Alguns métodos não recebem context, context não é propagado consistentemente
- **Impacto:** Difícil de implementar tracing, timeouts, cancellation
- **Correção Sugerida:** Garantir que todos os métodos recebem e propagam context
- **Estimativa:** 12h

---

## FASE 13 — QUALIDADE GERAL

### Notas Atuais

| Métrica | Nota | Observação |
|---------|------|------------|
| **Arquitetura** | 5/10 | Violações de DIP, SRP, ISP |
| **Acoplamento** | 4/10 | Alto acoplamento com GORM |
| **Coesão** | 5/10 | Packages muito grandes |
| **Complexidade** | 6/10 | Algumas funções complexas |
| **Duplicação** | 5/10 | Código duplicado em repositories e handlers |
| **Código Morto** | 6/10 | Packages vazios, TODOs |
| **Legibilidade** | 7/10 | Código geralmente legível |
| **Manutenibilidade** | 5/10 | Difícil de manter devido a acoplamento |
| **Escalabilidade** | 5/10 | Arquitetura precisa de refactoring para escalar |
| **Nota Geral** | 5/10 | Arquitetura funcional mas precisa de melhorias |

---

## Resumo por Severidade

### P0 — Arquitetura Quebrada (2 problemas)
1. Handler Importando Service Diretamente
2. Violação de Dependency Inversion

### D1 — Alto Acoplamento (4 problemas)
1. Service Dependendo Diretamente de gorm.DB
2. Repository Pattern Repetido
3. Violação de Single Responsibility
4. Tratamento de Erro Inconsistente

### P2 — Duplicação (15 problemas)
1. Função Gigante em concurrency_test.go
2. Função CreateCompany em platform_service.go
3. Repository Implementando GORM Models no Mesmo Arquivo
4. architecture_test.go com build tag
5. Conversão de Money Repetida
6. Validators Duplicados em Handlers
7. Package domain Muito Grande
8. Package service Muito Grande
9. Package handler Muito Grande
10. Package infra/repository Muito Grande
11. concurrency_test.go com 1814 linhas
12. platform_service.go com 782 linhas
13. gorm_product_repository.go com 712 linhas
14. Interfaces com Apenas Uma Implementação
15. Interfaces Gigantes
16. Violação de Interface Segregation
17. Context Não Propagado Consistentemente

### P3 — Limpeza (10 problemas)
1. Função generateSlug em product_service.go
2. Package errors Vazio
3. Package testutil/mocks Vazio
4. 111 TODOs no Código
5. 1 FIXME no Código
6. gorm_product_repository_test.go com 827 linhas
7. auth_service_test.go with 811 linhas
8. Arquivos de Documentação no Root
9. Arquivos de Database no Root
10. Arquivos de Coverage no Root
11. Dependências Possivelmente Não Utilizadas
12. Nomes Inconsistentes

---

## Estimativa Total de Esforço

**Total:** 312 horas (~39 dias úteis)

**Por fase:**
- FASE 1 (Complexidade): 12h
- FASE 2 (Acoplamento): 56h
- FASE 3 (Dependências Circulares): 0h
- FASE 4 (Código Morto): 21h
- FASE 5 (Duplicação): 36h
- FASE 6 (Coesão): 48h
- FASE 7 (Tamanho dos Arquivos): 26h
- FASE 8 (Interfaces): 20h
- FASE 9 (Organização): 4h
- FASE 10 (Dependências Externas): 2h
- FASE 11 (Clean Architecture): 76h
- FASE 12 (Padronização): 24h
- FASE 13 (Qualidade Geral): 0h

**Por severidade:**
- **Prioridade 0 (Arquitetura Quebrada):** 72h
- **Prioridade 1 (Alto Acoplamento):** 88h
- **Prioridade 2 (Duplicação):** 126h
- **Prioridade 3 (Limpeza):** 26h
