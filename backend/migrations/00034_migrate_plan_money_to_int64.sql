-- Migration: Migrar campos monetários de Plan de REAL para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: plans
ALTER TABLE plans ADD COLUMN price_cents BIGINT DEFAULT 0;

UPDATE plans SET price_cents = ROUND(price * 100);

ALTER TABLE plans DROP COLUMN price;

ALTER TABLE plans RENAME COLUMN price_cents TO price;
