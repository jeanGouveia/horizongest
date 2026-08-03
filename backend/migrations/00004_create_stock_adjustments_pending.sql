-- +goose Up
-- +goose StatementBegin

-- Tabela para registrar ajustes de estoque pendentes de aprovação
-- Usada para estornos de estoque por cancelamento de pedidos
CREATE TABLE IF NOT EXISTS stock_adjustments_pending (
    id            BIGSERIAL PRIMARY KEY,
    order_id      BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    quantity      NUMERIC(10,2) NOT NULL,
    order_status  VARCHAR(50) NOT NULL,
    status        VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at  TIMESTAMP
);

-- Índices para consultas frequentes
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_order_id ON stock_adjustments_pending(order_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_ingredient_id ON stock_adjustments_pending(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_status ON stock_adjustments_pending(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_stock_adjustments_pending_status;
DROP INDEX IF EXISTS idx_stock_adjustments_pending_ingredient_id;
DROP INDEX IF EXISTS idx_stock_adjustments_pending_order_id;
DROP TABLE IF EXISTS stock_adjustments_pending;

-- +goose StatementEnd
