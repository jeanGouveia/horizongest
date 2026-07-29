# Arquitetura de Idempotência - HorizonGest ERP

## Visão Geral

Solução de idempotência pragmática para ERP SaaS de pequeno e médio porte (~5M pedidos/ano), adaptada dos padrões da Stripe mas sem overengineering.

## Objetivos Atendidos

- ✅ Manter exatamente uma execução por `idempotency_key`
- ✅ Detectar reutilização da mesma chave com payload diferente
- ✅ Permitir response replay
- ✅ Expirar automaticamente chaves antigas
- ✅ Manter simplicidade operacional
- ✅ Evitar tabelas ou estruturas desnecessárias

## Componentes da Solução

### 1. Schema SQL (`migrations/00029_create_idempotency_table.sql`)

```sql
CREATE TABLE idempotency_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(255) NOT NULL,
    company_id INTEGER NOT NULL,
    request_hash VARCHAR(64) NOT NULL,      -- SHA-256 do payload normalizado
    request_params TEXT NOT NULL,           -- method, path, headers, query
    status VARCHAR(20) NOT NULL DEFAULT 'processing',
    response_status_code INTEGER,
    response_headers TEXT,
    response_body TEXT,
    error_message TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (company_id, key)
);

-- Índices otimizados
CREATE INDEX idx_idempotency_lookup ON idempotency_keys(company_id, key);
CREATE INDEX idx_idempotency_created_at ON idempotency_keys(created_at);
CREATE INDEX idx_idempotency_status ON idempotency_keys(status);
```

**Por que este schema?**
- Tabela dedicada: desacoplada de entities específicas (orders, payments, etc.)
- `request_hash`: detecta reutilização com payload diferente
- `response_*`: permite replay sem reexecutar lógica de negócio
- Índices minimalistas: apenas os necessários para o workload real

### 2. TTL e Limpeza

**TTL recomendado**: 48 horas
- Cobertura para todos os cenários de retry de rede
- Janela para debugging de problemas recentes
- Suficiente para timeouts de pagamento externo

**Job de limpeza**: Executar a cada 1 hora
```go
DELETE FROM idempotency_keys WHERE created_at < NOW() - INTERVAL '48 hours'
```

**Por que job periódico em vez de TTL automático do banco?**
- SQLite não suporta TTL automático (PostgreSQL suporta mas é complexo)
- Job periódico é mais simples e portável
- Impacto mínimo no banco (poucos milhares de registros por execução)

### 3. Fluxo Completo: Handler → Service → Repository

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ POST /api/orders
       │ Idempotency-Key: uuid-v4
       ▼
┌─────────────────────────────────┐
│  IdempotencyMiddleware          │
│  1. Extrai key do header        │
│  2. Lê body da requisição       │
│  3. Calcula hash do payload     │
│  4. Chama CheckAndCreate()      │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│  IdempotencyService             │
│  CheckAndCreate():              │
│  - Se replayable → retorna resp │
│  - Se payload mismatch → erro   │
│  - Se processing → erro/wait    │
│  - Se não existe → cria record  │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│  GormIdempotencyRepository     │
│  - Check(): SELECT por key      │
│  - Create(): INSERT com status  │
│  - UpdateSuccess(): UPDATE resp │
│  - UpdateFailure(): UPDATE erro │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│  Business Logic Handler        │
│  - Executa operação real        │
│  - ResponseWriter captura resp   │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│  IdempotencyService             │
│  - RecordSuccess() ou            │
│    RecordFailure()               │
└─────────────────────────────────┘
```

### 4. Tratamento de Retries

**Cenário 1: Retry com mesma chave e mesmo payload**
```
Request 1: POST /api/orders, Idempotency-Key=abc, body={...}
→ Criado record (status=processing)
→ Executa handler
→ RecordSuccess (status=succeeded, response=...)

Request 2: POST /api/orders, Idempotency-Key=abc, body={...}
→ Check: encontrado, status=succeeded, hash match
→ Replay: retorna response armazenada instantaneamente
```

**Cenário 2: Retry com mesma chave e payload diferente**
```
Request 1: POST /api/orders, Idempotency-Key=abc, body={items: [...]}
→ Criado record (status=processing)
→ Executa handler
→ RecordSuccess

Request 2: POST /api/orders, Idempotency-Key=abc, body={items: [...DIFERENTE...]}
→ Check: encontrado, status=succeeded, hash MISMATCH
→ Erro: 400 Bad Request
→ Mensagem: "idempotency_key já foi utilizada com payload diferente"
```

**Cenário 3: Retry enquanto ainda processing**
```
Request 1: POST /api/orders, Idempotency-Key=abc
→ Criado record (status=processing)
→ Executa handler (demora 500ms)

Request 2: POST /api/orders, Idempotency-Key=abc (50ms depois)
→ Check: encontrado, status=processing
→ Erro: 409 Conflict
→ Mensagem: "operação com esta chave ainda está em processamento"
→ Cliente deve aguardar e retry
```

### 5. Detecção de Payload Diferente

```go
// Hash normalizado: JSON ordenado, sem whitespace
func ComputeRequestHash(body []byte) string {
    var normalized map[string]interface{}
    if err := json.Unmarshal(body, &normalized); err == nil {
        normalizedBytes, _ := json.Marshal(normalized)
        hash := sha256.Sum256(normalizedBytes)
        return hex.EncodeToString(hash[:])
    }
    // Fallback para não-JSON
    hash := sha256.Sum256(body)
    return hex.EncodeToString(hash[:])
}
```

**Por que hash normalizado?**
- Detecta mudanças semânticas mesmo se whitespace mudar
- Chaves ordenadas: `{a:1, b:2}` == `{b:2, a:1}`
- Eficiente: SHA-256 é rápido e colisões são improváveis

### 6. Armazenamento da Resposta

```go
type IdempotencyResponse struct {
    StatusCode int
    Headers    map[string]string  // Relevantes apenas
    Body       string             // JSON completo
}
```

**Por que armazenar a resposta completa?**
- Replay instantâneo sem reexecutar lógica de negócio
- Útil para debugging (ver o que foi retornado originalmente)
- Headers relevantes (Content-Type, etc.) preservam contexto

### 7. Exemplos de Uso em Go

**Middleware Setup:**
```go
// Inicialização
idempotencyRepo := repository.NewGormIdempotencyRepository(db)
idempotencyService := service.NewIdempotencyService(idempotencyRepo)
idempotencyMiddleware := middleware.NewIdempotencyMiddleware(idempotencyService)

// Aplicar em endpoints mutáveis
router.Post("/api/orders", 
    idempotencyMiddleware.Handler(
        orderHandler.CreateOrder,
    ),
)
```

**Job de Limpeza:**
```go
// Inicialização
cleanupJob := service.NewIdempotencyCleanupJob(db, idempotencyRepo, 48) // 48 horas

// Iniciar em background
cleanupJob.Start()

// Parar graceful shutdown
defer cleanupJob.Stop()
```

**Uso Manual (sem middleware):**
```go
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    idempotencyKey := r.Header.Get("Idempotency-Key")
    body, _ := service.ReadBody(r)
    params := service.ExtractRequestParams(r)
    
    result, recordID, err := h.idempotencyService.CheckAndCreate(
        r.Context(),
        companyID,
        idempotencyKey,
        body,
        params,
    )
    
    if err != nil {
        // Tratar erro
        return
    }
    
    if result.Replayable {
        // Retornar response armazenada
        w.WriteHeader(result.Response.StatusCode)
        w.Write([]byte(result.Response.Body))
        return
    }
    
    // Executar lógica de negócio
    order, err := h.svc.CreateOrder(...)
    
    // Registrar resultado
    if err != nil {
        h.idempotencyService.RecordFailure(r.Context(), recordID, err.Error())
    } else {
        h.idempotencyService.RecordSuccess(r.Context(), recordID, statusCode, headers, body)
    }
}
```

## Componentes Stripe Descartados (Overengineering)

### ❌ Descartado: Tabela por recurso
**Stripe approach**: `payment_intents` tem `idempotency_key` como coluna, `charges` tem outra, etc.

**Por que descartar?**
- Duplicação de lógica em múltiplas tabelas
- Dificulta mudanças na estratégia de idempotência
- Não escala bem para novos recursos

**Nossa solução**: Tabela única `idempotency_keys` genérica

### ❌ Descartado: Lock distribuído para status "processing"
**Stripe approach**: Usa Redis lock para garantir que apenas uma requisição processe uma chave.

**Por que descartar?**
- Adiciona complexidade operacional (Redis)
- Race condition é rara em ERP (não é gateway de pagamento)
- Unique constraint no banco já garante unicidade

**Nossa solução**: Unique constraint + erro 409 Conflict se processing

### ❌ Descartado: Idempotency key gerada pelo servidor
**Stripe approach**: Se cliente não enviar key, servidor gera uma.

**Por que descartar?**
- Cliente não sabe qual key usar para retry
- Viola princípio de responsabilidade (cliente deve controlar retries)

**Nossa solução**: Key obrigatória, gerada pelo cliente (UUID v4)

### ❌ Descartado: Versionamento de chave
**Stripe approach**: `idempotency_key` pode ter sufixo de versão para permitir reuso.

**Por que descartar?**
- Complexidade desnecessária para ERP
- Cliente pode simplesmente gerar nova UUID

**Nossa solução**: Nova UUID para nova operação

### ❌ Descartado: Webhook replay com idempotency
**Stripe approach**: Webhooks podem ser reprocessados com idempotency.

**Por que descartar?**
- ERP não tem webhooks de alta frequência
- Webhooks podem ser idempotentes por natureza (UPSERT)

**Nossa solução**: Não implementado (pode ser adicionado se necessário)

### ❌ Descartado: Idempotency key com prefixo de recurso
**Stripe approach**: `order_abc123`, `payment_xyz789` para evitar colisões entre recursos.

**Por que descartar?**
- UUID v4 já garante unicidade global
- Prefixo adiciona complexidade sem benefício real

**Nossa solução**: UUID v4 puro

### ❌ Descartado: TTL por recurso
**Stripe approach**: Diferentes recursos têm diferentes TTLs.

**Por que descartar?**
- Complexidade desnecessária
- ERP tem padrão uniforme de operações

**Nossa solução**: TTL único de 48 horas para tudo

### ❌ Descartado: Idempotency key em query param
**Stripe approach**: Aceita key em header ou query param.

**Por que descartar?**
- Query param pode vazar em logs
- Header é mais seguro e padrão

**Nossa solução**: Apenas header `Idempotency-Key`

## Métricas e Monitoramento

Sugestão de métricas para monitorar a saúde do sistema:

- `idempotency_check_duration`: Tempo de verificação (p50, p95, p99)
- `idempotency_replay_rate`: % de requisições que são replay
- `idempotency_payload_mismatch_rate`: % de erros por payload diferente
- `idempotency_processing_conflict_rate`: % de erros por processing
- `idempotency_cleanup_duration`: Tempo do job de limpeza
- `idempotency_records_count`: Total de registros ativos

## Considerações de Escala

Para ~5M pedidos/ano (~13.7K/dia):

- **Volume de registros**: ~27K registros ativos (48h TTL)
- **Tamanho médio**: ~1KB por registro (response body comprimido se necessário)
- **Storage total**: ~27MB (trivial)
- **Lookup performance**: <1ms com índice `(company_id, key)`
- **Cleanup impact**: <100ms por execução (DELETE de poucos milhares de registros)

## Conclusão

Esta solução atinge o equilíbrio ideal entre:
- **Simplicidade**: Tabela única, job periódico, sem Redis
- **Funcionalidade**: Replay, detecção de payload diff, TTL
- **Escalabilidade**: Suporta facilmente 5M+ operações/ano
- **Operacionalidade**: Fácil de debugar, monitorar e manter

Componentes da Stripe foram propositalmente descartados porque adicionam complexidade sem benefício real para um ERP deste porte.
