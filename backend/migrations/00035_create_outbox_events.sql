-- +goose Up
-- +goose StatementBegin

-- Outbox Events Table
-- Tabela para implementação do Outbox Pattern
-- Garante consistência transacional entre operações de domínio e publicação de eventos
-- Isolamento multi-tenant via tenant_id (company_id)

CREATE TABLE IF NOT EXISTS outbox_events (
    id BIGSERIAL PRIMARY KEY,
    
    -- Identificação do agregado
    aggregate_type VARCHAR(100) NOT NULL,        -- 'order', 'product', 'ingredient', etc.
    aggregate_id BIGINT NOT NULL,                -- ID do agregado (ex: order_id)
    
    -- Metadados do evento
    event_type VARCHAR(100) NOT NULL,            -- 'order.created', 'product.updated', etc.
    event_version VARCHAR(20) NOT NULL DEFAULT '1.0',  -- Versão do schema do evento
    
    -- Payload do evento (dados completos em JSON)
    payload TEXT NOT NULL,                       -- JSON string (PostgreSQL suporta JSONB, mas TEXT é compatível)
    
    -- Isolamento multi-tenant
    tenant_id BIGINT NOT NULL,                   -- company_id para isolamento
    
    -- Controle de processamento
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    priority INTEGER NOT NULL DEFAULT 5,         -- 1=crítico, 5=normal, 10=baixa prioridade
    attempts INTEGER NOT NULL DEFAULT 0,         -- Número de tentativas de processamento
    
    -- Controle de tempo
    available_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Quando o evento fica disponível para processamento
    processed_at TIMESTAMP,                       -- Quando foi processado com sucesso
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Tratamento de erros
    last_error TEXT,                             -- Último erro (se houver)
    
    -- Unique constraint para idempotência por agregado+evento
    -- Previne duplicação de eventos para o mesmo agregado
    UNIQUE (aggregate_type, aggregate_id, event_type)
);

-- Índices otimizados para o workload do Outbox Pattern

-- Índice principal para dispatcher: buscar eventos pending por tenant
CREATE INDEX IF NOT EXISTS idx_outbox_tenant_status 
ON outbox_events(tenant_id, status);

-- Índice para disponibilidade temporal: eventos disponíveis para processamento
CREATE INDEX IF NOT EXISTS idx_outbox_available_at 
ON outbox_events(available_at);

-- Índice para lookup por agregado: auditoria e debugging
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate 
ON outbox_events(aggregate_type, aggregate_id);

-- Índice para limpeza de eventos processados antigos
CREATE INDEX IF NOT EXISTS idx_outbox_processed_at 
ON outbox_events(processed_at);

-- Índice para prioridade: processar eventos críticos primeiro
CREATE INDEX IF NOT EXISTS idx_outbox_priority 
ON outbox_events(priority, available_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_outbox_priority;
DROP INDEX IF EXISTS idx_outbox_processed_at;
DROP INDEX IF EXISTS idx_outbox_aggregate;
DROP INDEX IF EXISTS idx_outbox_available_at;
DROP INDEX IF EXISTS idx_outbox_tenant_status;
DROP TABLE IF EXISTS outbox_events;

-- +goose StatementEnd
