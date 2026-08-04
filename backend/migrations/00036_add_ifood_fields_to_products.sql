-- +goose Up
-- +goose StatementBegin

-- Sprint 2.1: Add iFood integration fields to products table
-- These fields exist in the domain model but were missing from migrations

ALTER TABLE products ADD COLUMN external_id VARCHAR(255);
ALTER TABLE products ADD COLUMN marketplace_id VARCHAR(255);
ALTER TABLE products ADD COLUMN sync_status VARCHAR(50) DEFAULT 'pending';
ALTER TABLE products ADD COLUMN last_sync TIMESTAMP;

-- Add indexes for common queries
CREATE INDEX IF NOT EXISTS idx_products_external_id ON products(external_id);
CREATE INDEX IF NOT EXISTS idx_products_marketplace_id ON products(marketplace_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_products_marketplace_id;
DROP INDEX IF EXISTS idx_products_external_id;
ALTER TABLE products DROP COLUMN last_sync;
ALTER TABLE products DROP COLUMN sync_status;
ALTER TABLE products DROP COLUMN marketplace_id;
ALTER TABLE products DROP COLUMN external_id;

-- +goose StatementEnd
