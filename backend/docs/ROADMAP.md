# HorizonGest - Roadmap

**Version:** 1.0
**Last Updated:** 31/07/2026

---

## Sprint 5C - Event Publishing Infrastructure

### Sprint 5C.1 - Outbox Pattern Infrastructure ✅ COMPLETED

**Status:** ✅ COMPLETED
**Data:** 30/07/2026

**Objetivo:** Implementar infraestrutura base do Outbox Pattern.

**Entregáveis:**
- ✅ Tabela `outbox_events` com schema otimizado
- ✅ Interface `OutboxRepository`
- ✅ Implementação `GormOutboxRepository`
- ✅ Índices otimizados para workload do Outbox
- ✅ Suporte a multi-tenancy via `tenant_id`
- ✅ Unique constraint para idempotência
- ✅ Testes unitários do repository

**Documentação:**
- `docs/OUTBOX_SPRINT5C1_AUDIT.md`

---

### Sprint 5C.2 - Event Publisher Implementation ✅ COMPLETED

**Status:** ✅ COMPLETED
**Data:** 30/07/2026

**Objetivo:** Implementar infraestrutura de publicação de eventos usando RabbitMQ.

**Entregáveis:**
- ✅ Interface `EventPublisher` em `internal/ports`
- ✅ RabbitMQ adapter completo (conexão, reconexão, exchanges, filas, bindings)
- ✅ `RabbitMQPublisher` implementando `EventPublisher`
- ✅ `EventDispatcher` com batch processing, retry, dead letter
- ✅ Configuração centralizada via environment variables
- ✅ Logs estruturados para observabilidade
- ✅ Testes unitários (Publisher e Dispatcher)

**Documentação:**
- `docs/SPRINT5C2_AUDITORIA.md`
- `docs/SPRINT5C2_ARQUITETURA.md`
- `docs/SPRINT5C2_VALIDACAO_ARQUITETURAL.md`
- `docs/SPRINT5C2_RELATORIO_FINAL.md`

---

### Sprint 5C.2.1 - Hardening da Infraestrutura ✅ COMPLETED

**Status:** ✅ HARDENED ✅ READY FOR CONSUMERS
**Data:** 31/07/2026

**Objetivo:** Endurecer a infraestrutura de eventos antes do desenvolvimento de consumidores.

**Entregáveis:**
- ✅ Lock otimista implementado (previne race condition)
- ✅ Retry refinado (Publisher: conexão, Dispatcher: negócio)
- ✅ `available_at` implementado corretamente para agendamento de retry
- ✅ Documentação de idempotência (`EVENT_IDEMPOTENCY.md`)
- ✅ Testes abrangentes (concorrência, retry, available_at, lock, shutdown)
- ✅ Auditoria completa (6 perguntas críticas respondidas)

**Bugs Corrigidos:**
- ✅ BUG-1: Race condition no Dispatcher
- ✅ BUG-2: Retry duplicado
- ✅ BUG-3: Available_at não atualizado

**Documentação:**
- `docs/EVENT_IDEMPOTENCY.md`
- `docs/SPRINT5C2_1_AUDITORIA_IMPLEMENTACAO.md`
- `docs/SPRINT5C2_1_AUDITORIA_FINAL.md`
- `docs/SPRINT5C2_1_RELATORIO_FINAL.md`

**Nota Arquitetural:** 9.5/10 (ALTA)

**Status:** ✅ **INFRAESTRUTURA PRONTA PARA CONSUMIDORES**

---

## Sprint 5C.3 - Email Consumer ✅ COMPLETED

**Status:** ✅ COMPLETED
**Data:** 31/07/2026
**Prioridade:** ALTA

**Objetivo:** Implementar o primeiro consumidor oficial do HorizonGest e validar a arquitetura Event-Driven.

**Entregáveis:**
- ✅ EmailConsumer independente que consome do RabbitMQ
- ✅ EmailProvider interface (substituível: Log, SMTP, SendGrid, SES, Mailgun)
- ✅ LogEmailProvider (implementação para desenvolvimento/teste)
- ✅ Template engine (InvitationTemplate, OrderCreatedTemplate, CompanyCreatedTemplate)
- ✅ IdempotencyStore (proteção contra processamento duplicado)
- ✅ Logs estruturados (recebido, ignorado, processado, tempo, falha)
- ✅ 11 testes unitários (100% pass)
- ✅ Auditoria completa (8 perguntas respondidas)
- ✅ Documentação completa

**Eventos Suportados:**
- ✅ invitation.created
- ✅ order.created
- ✅ company.created

**Documentação:**
- `docs/EMAIL_CONSUMER_ARCHITECTURE.md`
- `docs/SPRINT5C3_AUDITORIA.md`
- `docs/SPRINT5C3_RELATORIO_FINAL.md`
- `docs/ADR-006.md` (Arquitetura dos Consumers)

**Nota Arquitetural:** 9/10 (ALTA)

**Status:** ✅ **REFERÊNCIA VÁLIDA PARA PRÓXIMOS CONSUMIDORES**

---

## Sprint 5C.4 - Webhook Consumer (PLANNED)

**Status:** 📋 PLANNED
**Prioridade:** ALTA

**Objetivo:** Implementar Webhook Consumer seguindo o padrão do Email Consumer.

**Entregáveis Planejados:**
- 📋 WebhookConsumer independente que consome do RabbitMQ
- 📋 WebhookProvider interface (substituível: Log, HTTP)
- 📋 LogWebhookProvider (implementação para desenvolvimento/teste)
- 📋 Template engine para webhooks
- 📋 IdempotencyStore (reutilizar ou específico)
- 📋 Logs estruturados
- 📋 Testes unitários
- 📋 Documentação

**Pré-requisitos:**
- ✅ Infraestrutura de eventos hardened (Sprint 5C.2.1)
- ✅ Email Consumer como referência (Sprint 5C.3)
- ✅ ADR-006 (Arquitetura dos Consumers)

---

## Sprint 5C.5 - iFood Consumer (PLANNED)

**Status:** 📋 PLANNED
**Prioridade:** ALTA

**Objetivo:** Implementar iFood Consumer seguindo o padrão do Email Consumer.

**Entregáveis Planejados:**
- 📋 iFoodConsumer independente que consome do RabbitMQ
- 📋 iFoodProvider interface (substituível: Log, iFood API)
- 📋 LogiFoodProvider (implementação para desenvolvimento/teste)
- 📋 Template engine para iFood
- 📋 IdempotencyStore (reutilizar ou específico)
- 📋 Logs estruturados
- 📋 Testes unitários
- 📋 Documentação

**Pré-requisitos:**
- ✅ Infraestrutura de eventos hardened (Sprint 5C.2.1)
- ✅ Email Consumer como referência (Sprint 5C.3)
- ✅ ADR-006 (Arquitetura dos Consumers)

---

## Sprint 5C.4 - Operational Improvements (PLANNED)

**Status:** 📋 PLANNED
**Prioridade:** MÉDIA

**Objetivo:** Melhorias operacionais na infraestrutura de eventos.

**Entregáveis Planejados:**
- 📋 Job de recovery para eventos presos em `processing`
- 📋 Dead letter queue dedicada (tabela `outbox_dead_letters`)
- 📋 Metrics Prometheus (Dispatcher, Publisher, RabbitMQ)
- 📋 Dashboard de monitoramento (Grafana)
- 📋 Alertas de DLQ e latência

**Pré-requisitos:**
- ✅ Infraestrutura de eventos hardened (Sprint 5C.2.1)
- 📋 Consumidores implementados (Sprint 5C.3)

---

## Sprint 5C.5 - Optimizations (FUTURE)

**Status:** 📋 FUTURE
**Prioridade:** BAIXA

**Objetivo:** Otimizações de performance e arquitetura.

**Entregáveis Planejados:**
- 📋 LISTEN/NOTIFY ou CDC para reduzir latência (substituir polling)
- 📋 Multi-tenant dispatcher (processamento por tenant específico)
- 📋 Schema evolution para eventos (versionamento)
- 📋 Event Sourcing (se aplicável)
- 📋 CQRS (se aplicável)

**Pré-requisitos:**
- ✅ Infraestrutura de eventos hardened (Sprint 5C.2.1)
- 📋 Consumidores implementados (Sprint 5C.3)
- 📋 Melhorias operacionais (Sprint 5C.4)

---

## Status Geral

### Event Publishing Infrastructure

| Componente | Status | Nota |
|------------|--------|------|
| Outbox Pattern | ✅ COMPLETED | 9/10 |
| Event Publisher | ✅ COMPLETED | 9/10 |
| RabbitMQ Adapter | ✅ COMPLETED | 9/10 |
| Event Dispatcher | ✅ HARDENED | 9.5/10 |
| Idempotência | ✅ DOCUMENTED | 10/10 |
| Testes | ✅ COMPLETED | 10/10 |
| **INFRAESTRUTURA GERAL** | ✅ **READY FOR CONSUMERS** | **9.5/10** |

---

## Próximos Passos Imediatos

1. ✅ **COMPLETO:** Sprint 5C.2.1 - Hardening
2. 📋 **PRÓXIMO:** Sprint 5C.3 - Event Consumers
3. 📋 **DEPOIS:** Sprint 5C.4 - Operational Improvements

---

**END OF ROADMAP**
