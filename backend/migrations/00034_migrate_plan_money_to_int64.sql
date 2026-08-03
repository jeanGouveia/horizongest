-- +goose Up
-- +goose StatementBegin

-- Migration: Migrar campos monetários de Plan de NUMERIC(10,2) para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: plans
ALTER TABLE plans ADD COLUMN price_cents BIGINT DEFAULT 0;

UPDATE plans SET price_cents = ROUND(price * 100);

ALTER TABLE plans DROP COLUMN price;

ALTER TABLE plans RENAME COLUMN price_cents TO price;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert migration: migrar campos monetários de Plan de BIGINT para NUMERIC(10,2)
-- Tabela: plans
ALTER TABLE plans ADD COLUMN price NUMERIC(10,2) DEFAULT 0;

UPDATE plans SET price = price_cents / 100.0;

ALTER TABLE plans DROP COLUMN price_cents;

ALTER TABLE plans RENAME COLUMN price TO price_cents;

-- +goose StatementEnd
