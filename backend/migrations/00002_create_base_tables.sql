-- +goose Up
-- +goose StatementBegin

-- ============================================================================
-- COMPATIBILITY MIGRATION
--
-- Esta migration foi originalmente utilizada durante a fase SQLite do projeto.
--
-- Durante a migração definitiva para PostgreSQL, sua lógica foi consolidada
-- nas migrations posteriores.
--
-- Ela permanece apenas para preservar a sequência histórica das migrations.
--
-- Não remover.
-- Não reutilizar este número.
-- ============================================================================

SELECT 1;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

SELECT 1;

-- +goose StatementEnd
