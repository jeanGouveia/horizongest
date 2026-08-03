-- +goose Up
-- +goose StatementBegin

-- Idempotency Table
-- Tabela dedicada para idempotência, desacoplada de entities específicas
-- Suporta qualquer operação mutável (POST, PUT, PATCH, DELETE)
-- Permite response replay e detecção de payload diferente
-- TTL implementado via job de limpeza (DELETE WHERE created_at < expiry)

CREATE TABLE IF NOT EXISTS idempotency_keys (
    id BIGSERIAL PRIMARY KEY,
    
    -- Chave de idempotência (gerada pelo cliente, UUID v4 recomendado)
    key VARCHAR(255) NOT NULL,
    
    -- Escopo da chave (company_id para isolamento multi-tenant)
    company_id BIGINT NOT NULL,
    
    -- Hash do payload request para detectar reutilização com payload diferente
    -- SHA-256 do body JSON normalizado (chaves ordenadas, sem whitespace)
    request_hash VARCHAR(64) NOT NULL,
    
    -- Parâmetros da requisição (method, path, query params)
    -- Armazenado como JSON para contexto e debugging
    request_params TEXT NOT NULL,
    
    -- Status da operação original
    -- 'processing' = em execução, 'succeeded' = concluída, 'failed' = falhou
    status VARCHAR(20) NOT NULL DEFAULT 'processing',
    
    -- Resposta HTTP armazenada para replay
    -- Inclui: status_code, headers, body
    response_status_code INTEGER,
    response_headers TEXT,      -- JSON headers
    response_body TEXT,         -- JSON body
    
    -- Erro se a operação falhou
    error_message TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint: uma chave só pode ser usada uma vez por company
    UNIQUE (company_id, key)
);

-- Índices otimizados para o workload de idempotência
-- Lookup principal: buscar por (company_id, key)
CREATE INDEX IF NOT EXISTS idx_idempotency_lookup 
ON idempotency_keys(company_id, key);

-- Índice para TTL/limpeza: buscar registros antigos
CREATE INDEX IF NOT EXISTS idx_idempotency_created_at 
ON idempotency_keys(created_at);

-- Índice para status: monitorar operações em processing (stuck)
CREATE INDEX IF NOT EXISTS idx_idempotency_status 
ON idempotency_keys(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_idempotency_status;
DROP INDEX IF EXISTS idx_idempotency_created_at;
DROP INDEX IF EXISTS idx_idempotency_lookup;
DROP TABLE IF EXISTS idempotency_keys;

-- +goose StatementEnd
