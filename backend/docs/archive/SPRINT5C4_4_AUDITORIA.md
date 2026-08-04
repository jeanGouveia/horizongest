# SPRINT 5C.4.4 — AUDITORIA DE INFRAESTRUTURA ASSÍNCRONA

**Data:** 2026-07-31  
**Objetivo:** Auditoria extremamente profunda da infraestrutura assíncrona do HorizonGest  
**Escopo:** Outbox Pattern, RabbitMQ, Consumer Framework, Redis, Performance, Observabilidade, Startup/Shutdown, Concorrência, Clean Architecture  
**Status:** 🔴 CRÍTICO - Infraestrutura assíncrona NÃO está operacional

---

## Sumário Executivo

A auditoria identificou **24 problemas** distribuídos em 5 níveis de severidade:

- 🔴 **Crítico:** 12 problemas (50%)
- 🟠 **Alto:** 6 problemas (25%)
- 🟡 **Médio:** 4 problemas (17%)
- 🔵 **Baixo:** 2 problemas (8%)

**Conclusão:** A infraestrutura assíncrona está **COMPLETAMENTE NÃO FUNCIONAL** em produção. Todo o código existe mas NÃO está inicializado no main.go.

---

## Problemas Detalhados

### 🔴 CRÍTICO - Outbox Pattern Não Funcional

#### O1 - EventDispatcher NÃO é iniciado no main.go
- **Arquivo:** `cmd/server/main.go` (linhas 1-553)
- **Causa:** EventDispatcher existe, é criado em testes, mas NÃO é inicializado no main.go
- **Impacto:** Outbox pattern está completamente inoperacional. Eventos são escritos na tabela outbox_events mas NUNCA são publicados no RabbitMQ
- **Risco:** Perda de eventos, inconsistência de dados, sistema assíncrono não funciona
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // No main.go, após inicialização de serviços:
  outboxRepo := repository.NewGormOutboxRepository(db)
  
  // Inicializar RabbitMQ publisher
  rabbitMQConfig := rabbitmq.Config{
      URL:                    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
      Exchange:               "horizongest.events",
      ExchangeType:           "topic",
      QueuePrefix:            "horizongest",
      RetryCount:             3,
      PublisherTimeout:       10 * time.Second,
      ReconnectDelay:         5 * time.Second,
      EnablePublisherConfirm: true,
  }
  eventPublisher, err := rabbitmq.NewRabbitMQPublisher(rabbitMQConfig)
  if err != nil {
      log.Fatalf("FATAL: falha ao criar RabbitMQ publisher: %v", err)
  }
  defer eventPublisher.Close()
  
  // Criar e iniciar dispatcher
  dispatcherConfig := service.DefaultDispatcherConfig()
  dispatcher := service.NewEventDispatcher(outboxRepo, eventPublisher, dispatcherConfig)
  dispatcher.Start(context.Background())
  defer dispatcher.Shutdown()
  ```

#### O2 - RabbitMQ Publisher NÃO é inicializado no main.go
- **Arquivo:** `cmd/server/main.go` (linhas 1-553)
- **Causa:** Código do RabbitMQ publisher existe mas não é instanciado no main.go
- **Impacto:** Não há conexão com RabbitMQ, eventos não podem ser publicados
- **Risco:** Sistema assíncrono completamente não funcional
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:** Ver O1

#### O3 - Redis Client NÃO é inicializado no main.go
- **Arquivo:** `cmd/server/main.go` (linhas 1-553)
- **Causa:** Código do Redis client existe mas não é instanciado no main.go
- **Impacto:** Cache, sessões, locks, rate limiter, idempotência não funcionam
- **Risco:** Perda de performance, sem cache, sem rate limiting, sem locks distribuídos
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // No main.go, após inicialização de DB:
  redisConfig := redis.Config{
      Host:         getEnv("REDIS_HOST", "localhost"),
      Port:         getEnvInt("REDIS_PORT", 6379),
      Password:     getEnv("REDIS_PASSWORD", ""),
      DB:           getEnvInt("REDIS_DB", 0),
      PoolSize:     10,
      MinIdleConns: 2,
      MaxIdleConns: 5,
      MaxRetries:   3,
      DialTimeout:  5 * time.Second,
      ReadTimeout:  3 * time.Second,
      WriteTimeout: 3 * time.Second,
      PoolTimeout:  4 * time.Second,
  }
  redisClient, err := redis.NewClient(redisConfig)
  if err != nil {
      log.Fatalf("FATAL: falha ao conectar Redis: %v", err)
  }
  defer redisClient.Close()
  ```

#### O4 - Consumers NÃO são iniciados no main.go
- **Arquivo:** `cmd/server/main.go` (linhas 1-553)
- **Causa:** Código dos consumers existe mas não são instanciados nem iniciados no main.go
- **Impacto:** Email consumer, webhook consumer não funcionam
- **Risco:** Emails não são enviados, webhooks não são chamados
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // No main.go, após inicialização de RabbitMQ:
  // Criar conexão RabbitMQ para consumers
  rabbitConn, err := amqp.Dial(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
  if err != nil {
      log.Fatalf("FATAL: falha ao conectar RabbitMQ para consumers: %v", err)
  }
  defer rabbitConn.Close()
  
  // Inicializar email consumer
  emailProvider := email.NewMockEmailProvider() // ou implementação real
  emailConfig := framework.Config{
      Queue:              "horizongest.email",
      ConsumerName:        "EmailConsumer",
      MaxRetries:         3,
      InitialRetryDelay:  1 * time.Second,
      MaxRetryDelay:      30 * time.Second,
      RetryMultiplier:    2.0,
      OperationTimeout:   30 * time.Second,
      CircuitBreakerThreshold: 5,
      CircuitBreakerTimeout:   1 * time.Minute,
      DeadLetterQueue:    "horizongest.email.dlq",
      MaxRetryAttempts:   5,
      EnableMetrics:      true,
  }
  emailConsumer := email.NewEmailConsumer(rabbitConn, emailConfig.Queue, emailProvider, emailConfig)
  go emailConsumer.Start(context.Background())
  defer emailConsumer.Close()
  
  // Inicializar webhook consumer (similar)
  ```

#### O5 - Idempotência é In-Memory por Padrão
- **Arquivo:** `internal/consumers/framework/consumer.go` (linha 48)
- **Causa:** `NewIdempotencyStore()` cria idempotência in-memory em vez de Redis
- **Impacto:** Em caso de restart do consumer, eventos podem ser processados duplicadamente
- **Risco:** Duplicação de emails, webhooks, inconsistência de dados
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // No consumer.go, mudar de:
  idempotencyStore := NewIdempotencyStore()
  
  // Para:
  redisCache := redis.NewCache(redisClient)
  idempotencyStore := NewRedisIdempotencyStore(redisCache, 24*time.Hour)
  ```

#### O6 - Dead Letter Handler NÃO usa DLQ do RabbitMQ
- **Arquivo:** `internal/consumers/framework/dead_letter.go` (linhas 34-89)
- **Causa:** DeadLetterHandler declara fila manualmente em vez de usar DLQ nativa do RabbitMQ
- **Impacto:** Não há integração com DLQ nativa do RabbitMQ, eventos mortos ficam em fila separada
- **Risco:** Gerenciamento manual de DLQ, perda de eventos se fila não for monitorada
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:** Configurar DLQ nativa do RabbitMQ com x-dead-letter-exchange e x-dead-letter-routing-key

#### O7 - EventDispatcher Processa Todos os Tenants (tenantID=0)
- **Arquivo:** `internal/service/event_dispatcher.go` (linha 108)
- **Causa:** `FindPendingEvents(ctx, 0, d.config.BatchSize)` usa tenantID=0, que processa todos os tenants
- **Impacto:** Dispatcher processa eventos de todos os tenants em batch, pode causar cross-tenant contamination
- **Risco:** Eventos de um tenant podem atrasar processamento de outro, violação de isolamento
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Implementar processamento por tenant:
  // 1. Buscar lista de tenants ativos
  // 2. Para cada tenant, processar eventos pendentes
  // 3. Usar tenantID específico em FindPendingEvents
  ```

#### O8 - EventDispatcher Não Tem Timeout no Loop Principal
- **Arquivo:** `internal/service/event_dispatcher.go` (linhas 88-96)
- **Causa:** Loop principal não tem timeout, pode bloquear indefinidamente
- **Impacto:** Se operação de banco travar, dispatcher pode ficar bloqueado
- **Risco:** Deadlock, dispatcher não responde, eventos não são processados
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Adicionar timeout no processBatch:
  ctx, cancel := context.WithTimeout(d.shutdownCtx, 30*time.Second)
  defer cancel()
  d.processBatch(ctx)
  ```

#### O9 - Publisher Confirm Pode Bloquear Indefinidamente
- **Arquivo:** `internal/infra/messaging/rabbitmq/rabbitmq_publisher.go` (linhas 129-139)
- **Causa:** Publisher confirm usa select sem timeout, pode bloquear se RabbitMQ não responder
- **Impacto:** Publisher pode ficar bloqueado esperando confirm
- **Risco:** Deadlock, eventos não são publicados, dispatcher trava
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Adicionar timeout no publisher confirm:
  confirmChan := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
  select {
  case confirm := <-confirmChan:
      if !confirm.Ack {
          return fmt.Errorf("publisher confirm: message not acknowledged")
      }
  case <-publishCtx.Done():
      return fmt.Errorf("publisher confirm timeout")
  case <-time.After(5 * time.Second):
      return fmt.Errorf("publisher confirm timeout after 5s")
  }
  ```

#### O10 - RabbitMQ Reconnect Loop Sem Backoff
- **Arquivo:** `internal/infra/messaging/rabbitmq/rabbitmq_connection.go` (linhas 89-109)
- **Causa:** Reconnect usa delay fixo sem backoff exponencial
- **Impacto:** Em caso de RabbitMQ down, reconexão é agressiva e pode sobrecarregar
- **Risco:** Thundering herd, sobrecarga de rede, RabbitMQ não recupera
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Implementar backoff exponencial:
  reconnectDelay := c.config.ReconnectDelay
  for {
      time.Sleep(reconnectDelay)
      // ... tentar reconectar ...
      if err == nil {
          break
      }
      reconnectDelay = time.Duration(float64(reconnectDelay) * 1.5)
      if reconnectDelay > 60*time.Second {
          reconnectDelay = 60*time.Second
      }
  }
  ```

#### O11 - Não Há Prefetch Configurado no Consumer
- **Arquivo:** `internal/consumers/framework/consumer.go` (linhas 116-127)
- **Causa:** `channel.Consume` não configura prefetch (QoS)
- **Impacto:** RabbitMQ pode enviar muitas mensagens de uma vez, sobrecarregando consumer
- **Risco:** OOM, consumer trava, perda de mensagens
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Antes de Consume, configurar QoS:
  err := channel.Qos(
      10,    // prefetch count
      0,     // prefetch size
      false, // global
  )
  if err != nil {
      return fmt.Errorf("failed to set QoS: %w", err)
  }
  ```

#### O12 - Não Há Health Check de Async Components
- **Arquivo:** `cmd/server/main.go` (linha 240-244)
- **Causa:** Health check `/api/health` apenas retorna "ok", não verifica RabbitMQ, Redis, Dispatcher, Consumers
- **Impacto:** Não há visibilidade da saúde da infraestrutura assíncrona
- **Risco:** Problemas não detectados, downtime não identificado
- **Severidade:** 🔴 CRÍTICO
- **Solução Proposta:**
  ```go
  // Implementar health check detalhado:
  r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
      health := map[string]interface{}{
          "status": "ok",
          "service": "horizongest",
          "components": map[string]interface{}{
              "database": checkDBHealth(db),
              "redis": checkRedisHealth(redisClient),
              "rabbitmq": checkRabbitMQHealth(rabbitConn),
              "dispatcher": dispatcher.IsRunning(),
              "consumers": checkConsumersHealth(consumers),
          },
      }
      json.NewEncoder(w).Encode(health)
  })
  ```

---

### 🟠 ALTO - Performance e Confiabilidade

#### A1 - IdempotencyStore In-Memory Sem TTL
- **Arquivo:** `internal/consumers/framework/idempotency.go` (linhas 17-51)
- **Causa:** IdempotencyStore in-memory não tem TTL, IDs acumulam indefinidamente
- **Impacto:** Memory leak em consumers de longa duração
- **Risco:** OOM, consumer crash após tempo prolongado
- **Severidade:** 🟠 ALTO
- **Solução Proposta:** Implementar TTL ou usar Redis idempotency store

#### A2 - DeadLetterMiddleware Usa Map Para Tracking de Attempts
- **Arquivo:** `internal/consumers/framework/middleware.go` (linhas 122-158)
- **Causa:** `attemptMap` é map in-memory, não persiste entre restarts
- **Impacto:** Após restart, contador de attempts é resetado, eventos podem exceder limite real
- **Risco:** Loop infinito de retry, eventos nunca vão para DLQ
- **Severidade:** 🟠 ALTO
- **Solução Proposta:** Usar Redis para tracking de attempts ou usar campo `attempts` do evento

#### A3 - Retry Sem Jitter
- **Arquivo:** `internal/consumers/framework/retry.go` (linhas 18-47)
- **Causa:** Retry usa backoff exponencial sem jitter
- **Impacto:** Múltiplos consumers podem retryar sincronizados, causando thundering herd
- **Risco:** Sobrecarga de recursos, race conditions
- **Severidade:** 🟠 ALTO
- **Solução Proposta:**
  ```go
  // Adicionar jitter ao delay:
  jitter := time.Duration(rand.Int63n(int64(delay) / 10)) // 10% jitter
  delay = delay + jitter
  ```

#### A4 - Circuit Breaker Sem Half-Open Success Threshold
- **Arquivo:** `internal/consumers/framework/circuit_breaker.go` (linhas 72-78)
- **Causa:** Circuit breaker transiciona para closed após 1 sucesso em half-open
- **Impacto:** Circuit breaker pode fechar em falso positivo
- **Risco:** Oscilação de estado, instabilidade
- **Severidade:** 🟠 ALTO
- **Solução Proposta:** Implementar threshold de sucessos (ex: 3 sucessos consecutivos)

#### A5 - Redis Rate Limiter Usa Token Bucket Simples
- **Arquivo:** `internal/infra/redis/ratelimiter.go` (linhas 42-59)
- **Causa:** Rate limiter usa INCR + EXPIRE, que não é sliding window verdadeiro
- **Impacto:** Rate limiting não é preciso, pode permitir mais requests do que o limite
- **Risco:** Abuso de API, sobrecarga
- **Severidade:** 🟠 ALTO
- **Solução Proposta:** Implementar sliding window com Redis sorted sets ou usar token bucket com Lua script

#### A6 - Metrics são In-Memory
- **Arquivo:** `internal/consumers/framework/metrics.go` (não existe, mas InMemoryMetrics é usado)
- **Causa:** Metrics são armazenadas in-memory, perdidas em restart
- **Impacto:** Não há histórico de métricas, impossível analisar tendências
- **Risco:** Falta de observabilidade histórica
- **Severidade:** 🟠 ALTO
- **Solução Proposta:** Integrar com Prometheus ou OpenTelemetry

---

### 🟡 MÉDIO - Observabilidade e Monitoramento

#### M1 - Logs São Printf em vez de Structured Logging
- **Arquivo:** Múltiplos arquivos (consumer.go, event_dispatcher.go, etc.)
- **Causa:** Logs usam `log.Printf` em vez de structured logging (ex: logrus, zap)
- **Impacto:** Difícil de parsear logs, falta de contexto, não há trace ID
- **Risco:** Debugging difícil, não há correlação entre logs
- **Severidade:** 🟡 MÉDIO
- **Solução Proposta:** Implementar structured logging com trace ID propagation

#### M2 - Não Há Distributed Tracing
- **Arquivo:** N/A (não implementado)
- **Causa:** Não há integração com OpenTelemetry ou similar
- **Impacto:** Não há visibilidade de fluxo de requisições através de componentes assíncronos
- **Risco:** Difícil debugar problemas distribuídos
- **Severidade:** 🟡 MÉDIO
- **Solução Proposta:** Implementar OpenTelemetry com trace propagation

#### M3 - Não Há Alerting
- **Arquivo:** N/A (não implementado)
- **Causa:** Não há integração com sistema de alerting (Prometheus Alertmanager, PagerDuty, etc.)
- **Impacto:** Problemas não geram alertas automáticos
- **Risco:** Downtime não detectado em tempo hábil
- **Severidade:** 🟡 MÉDIO
- **Solução Proposta:** Configurar alertas baseados em métricas

#### M4 - Não Há Dashboard de Monitoramento
- **Arquivo:** N/A (não implementado)
- **Causa:** Não há dashboard Grafana ou similar
- **Impacto:** Não há visibilidade visual da saúde do sistema
- **Risco:** Operação manual, difícil identificar problemas
- **Severidade:** 🟡 MÉDIO
- **Solução Proposta:** Criar dashboards no Grafana

---

### 🔵 BAIXO - Melhorias de Código

#### B1 - DeadLetterHandler Declara Fila a Cada Envio
- **Arquivo:** `internal/consumers/framework/dead_letter.go` (linhas 42-53)
- **Causa:** `QueueDeclare` é chamado a cada envio para DLQ
- **Impacto:** Overhead desnecessário de rede
- **Risco:** Performance degradation
- **Severidade:** 🔵 BAIXO
- **Solução Proposta:** Declarar fila uma vez na inicialização

#### B2 - IdempotencyStore.Clear Não Implementado para Redis
- **Arquivo:** `internal/consumers/framework/idempotency_redis.go` (linhas 63-69)
- **Causa:** Método Clear retorna erro "not implemented"
- **Impacto:** Não é possível limpar idempotência em Redis (útil para testes)
- **Risco:** Dificuldade em testes
- **Severidade:** 🔵 BAIXO
- **Solução Proposta:** Implementar scan + delete ou usar prefixo separado por teste

---

### 🟢 CLEAN ARCHITECTURE - Avaliação

#### CA1 - Infra NÃO Depende de Application ✅
- **Status:** APROVADO
- **Evidência:** `internal/infra/messaging/rabbitmq` e `internal/infra/redis` não dependem de `internal/service`
- **Conclusão:** Camada de infraestrutura está isolada

#### CA2 - Application NÃO Depende de Infra ✅
- **Status:** APROVADO
- **Evidência:** `internal/service/event_dispatcher` depende de interfaces (`ports.OutboxRepository`, `ports.EventPublisher`)
- **Conclusão:** Service layer depende de abstrações, não de implementações concretas

#### CA3 - Consumers Respeitam DDD ✅
- **Status:** APROVADO
- **Evidência:** Consumers usam framework que delega para `Processor` interface, que contém lógica de negócio
- **Conclusão:** Separação entre infra (framework) e domínio (processor)

#### CA4 - Redis Está Isolado ✅
- **Status:** APROVADO
- **Evidência:** `internal/infra/redis` é pacote isolado com interfaces bem definidas
- **Conclusão:** Redis pode ser substituído por outra implementação

#### CA5 - RabbitMQ Está Isolado ✅
- **Status:** APROVADO
- **Evidência:** `internal/infra/messaging/rabbitmq` é pacote isolado com interfaces bem definidas
- **Conclusão:** RabbitMQ pode ser substituído por outro broker

#### CA6 - Framework É Reutilizável ✅
- **Status:** APROVADO
- **Evidência:** `internal/consumers/framework` é genérico e reutilizável para diferentes consumers
- **Conclusão:** Framework bem desenhado

---

## Conclusão da Clean Architecture

**Status:** ✅ **APROVADO**

A arquitetura limpa está bem implementada. As camadas estão corretamente separadas, dependências apontam na direção correta, e o framework é reutilizável. Os problemas são de **integração** (não inicializar componentes no main.go) e não de **arquitetura**.

---

## Resumo por Categoria

| Categoria | Crítico | Alto | Médio | Baixo | Total |
|-----------|---------|------|-------|-------|-------|
| Outbox Pattern | 3 | 0 | 0 | 0 | 3 |
| RabbitMQ | 3 | 2 | 0 | 1 | 6 |
| Consumer Framework | 2 | 2 | 0 | 1 | 5 |
| Idempotência | 2 | 1 | 0 | 1 | 4 |
| Redis | 1 | 1 | 0 | 0 | 2 |
| Performance | 0 | 2 | 0 | 0 | 2 |
| Observabilidade | 1 | 1 | 4 | 0 | 6 |
| Startup/Shutdown | 4 | 0 | 0 | 0 | 4 |
| Clean Architecture | 0 | 0 | 0 | 0 | 0 |
| **TOTAL** | **12** | **6** | **4** | **2** | **24** |

---

## Priorização de Correção

### Fase 1: Crítico - Tornar Infraestrutura Funcional (8-12 horas)
1. Inicializar Redis client no main.go
2. Inicializar RabbitMQ publisher no main.go
3. Inicializar EventDispatcher no main.go
4. Inicializar consumers no main.go
5. Implementar health check de async components
6. Mudar idempotência para Redis
7. Configurar prefetch no consumer
8. Implementar timeout no loop do dispatcher

### Fase 2: Alto - Melhorar Confiabilidade (4-6 horas)
1. Implementar backoff exponencial no reconnect RabbitMQ
2. Adicionar jitter no retry
3. Implementar threshold no circuit breaker
4. Melhorar rate limiter com sliding window
5. Integrar metrics com Prometheus
6. Usar Redis para tracking de attempts

### Fase 3: Médio - Melhorar Observabilidade (4-6 horas)
1. Implementar structured logging
2. Implementar distributed tracing
3. Configurar alerting
4. Criar dashboards

### Fase 4: Baixo - Otimizações (1-2 horas)
1. Declarar DLQ uma vez na inicialização
2. Implementar Clear para Redis idempotency

---

**Fim do Relatório de Auditoria**
