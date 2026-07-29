-- Migration: Migrar campos monetários de Finance de REAL para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: transactions
ALTER TABLE transactions ADD COLUMN amount_cents BIGINT DEFAULT 0;

UPDATE transactions SET amount_cents = ROUND(amount * 100);

ALTER TABLE transactions DROP COLUMN amount;

ALTER TABLE transactions RENAME COLUMN amount_cents TO amount;
