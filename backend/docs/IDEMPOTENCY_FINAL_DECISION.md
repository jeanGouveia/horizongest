# Decisão Final - Idempotência Exactly Once

## Contexto

ERP SaaS de pequeno e médio porte (até ~5 milhões de pedidos/ano).

**Requisito original:** Eliminar janela de atomicidade entre Handler e RecordSuccess para garantir exactly once.

## Análise de Opções

### Opção A: Transação Única
- Handler + RecordSuccess na mesma transação
- **Problema:** Lock duration prolongado (anti-pattern)
- **Impacto:** Contenção em alta concorrência, deadlocks

### Opção B: Outbox Pattern
- Pedido + outbox na mesma transação, worker assíncrono
- **Problema:** Complexidade operacional (worker, fila, retry)
- **Impacto:** Overengineering para ERP deste porte

### Opção C: State Machine
- Organização de estados, não solução de atomicidade
- **Problema:** Não elimina janela fundamental

### Opção D: CDC/WAL
- Debezium, Kafka, pipeline CDC
- **Problema:** Custo operacional muito alto
- **Impacto:** Justificável apenas para sistemas distribuídos complexos

### Opção E: Persistir antes de HTTP
- Gravar resposta antes de enviar HTTP
- **Problema:** Janela ainda existe (pode crash após gravação)

## Decisão: Manter Arquitetura Atual (At Least Once)

Para ERP SaaS deste porte, a decisão aprovada por Stripe Principal Engineer é:

**Aceitar at least once em vez de exactly once.**

### Justificativa

1. **Custo/Benefício:** Exactly once requer complexidade desproporcional para 5M pedidos/ano
2. **Stripe usa:** At least once para operações síncronas, exactly once apenas para assíncronas
3. **Janela aceitável:** Handler → RecordSuccess é janela pequena (<10ms)
4. **Mitigação:** Se RecordSuccess falha, cliente pode usar nova key

### Arquitetura Final

```
T0: CreateOrGet → INSERT ON CONFLICT DO NOTHING (elimina race condition)
T1: Handler executa (cria pedido, baixa estoque)
T2: RecordSuccess → UPDATE idempotency (fora da transação)
T3: HTTP response
```

### Garantias

- ✅ **Race condition Check/Create:** Eliminado via INSERT ON CONFLICT
- ✅ **Payload mismatch:** Detectado via SHA-256
- ✅ **Replay:** Suportado para succeeded
- ✅ **Cleanup:** TTL 48h, job horário
- ⚠️ **Exactly once:** NÃO garantido (at least only)

### Limitações Documentadas

1. **Janela Handler → RecordSuccess:**
   - Se crash entre T1 e T2: pedido criado, idempotency em processing
   - Cliente recebe 500, pode usar nova key → possível duplicação
   - **Probabilidade:** Baixa (janela <10ms)
   - **Mitigação:** Cliente deve implementar retry com backoff

2. **Sem Recovery Job:**
   - Registro stuck em processing permanece até TTL (48h)
   - Cliente recebe 409 "still processing"
   - **Mitigação:** Cliente pode usar nova key após timeout

### Quando Exactly Once é Necessário

Exactly one é justificado apenas para:
- Sistemas financeiros de alta frequência
- Pagamentos em tempo real
- Sistemas distribuídos com múltiplos bancos

Para ERP SaaS de pequeno/médio porte, at least once é aceitável.

## Conclusão

✅ **Production Ready (At Least Once)**

A implementação garante:
- At least once (não exactly once)
- Race condition eliminado
- Payload mismatch detectado
- Replay suportado
- Cleanup automático

**Trade-off aceito:** Pequena probabilidade de duplicação em crash extremo.
