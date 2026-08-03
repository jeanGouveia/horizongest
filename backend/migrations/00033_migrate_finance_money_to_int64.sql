-- +goose Up
-- +goose StatementBegin

-- Migration: Migrar campos monetários de Finance de NUMERIC(10,2) para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: transactions
ALTER TABLE transactions ADD COLUMN amount_cents BIGINT DEFAULT 0;

UPDATE transactions SET amount_cents = ROUND(amount * 100);

ALTER TABLE transactions DROP COLUMN amount;

ALTER TABLE transactions RENAME COLUMN amount_cents TO amount;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert migration: migrar campos monetários de Finance de BIGINT para NUMERIC(10,2)
-- Tabela: transactions
ALTER TABLE transactions ADD COLUMN amount NUMERIC(10,2) DEFAULT 0;

UPDATE transactions SET amount = amount_cents / 100.0;

ALTER TABLE transactions DROP COLUMN amount_cents;

ALTER TABLE transactions RENAME COLUMN amount TO amount_cents;

-- +goose StatementEnd
