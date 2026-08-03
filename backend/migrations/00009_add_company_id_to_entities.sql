-- +goose Up
-- Migration: Add company_id foreign key to tenant-scoped entities
-- This migration adds company_id columns to support multi-tenant isolation
-- Core V1 compatibility: All company_id columns are nullable (default NULL)
-- This ensures existing Core V1 data continues to work without modification

-- Add company_id to users table
ALTER TABLE users ADD COLUMN company_id BIGINT REFERENCES companies(id);

-- Note: categories table doesn't exist in this codebase

-- Add company_id to ingredients table
ALTER TABLE ingredients ADD COLUMN company_id BIGINT REFERENCES companies(id);

-- Add company_id to products table
ALTER TABLE products ADD COLUMN company_id BIGINT REFERENCES companies(id);

-- Note: orders table will be created in migration 00027 with company_id already included
-- This avoids the error since orders table doesn't exist yet

-- Add company_id to stock_adjustments_pending table
ALTER TABLE stock_adjustments_pending ADD COLUMN company_id BIGINT REFERENCES companies(id);

-- Add company_id to media table
ALTER TABLE media ADD COLUMN company_id BIGINT REFERENCES companies(id);

-- Create indexes for efficient tenant filtering
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
-- Note: categories table doesn't exist in this codebase
CREATE INDEX IF NOT EXISTS idx_ingredients_company_id ON ingredients(company_id);
CREATE INDEX IF NOT EXISTS idx_products_company_id ON products(company_id);
-- Note: orders table will be created in migration 00027 with indexes already included
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_pending_company_id ON stock_adjustments_pending(company_id);
CREATE INDEX IF NOT EXISTS idx_media_company_id ON media(company_id);

-- +goose Down
-- Drop indexes
DROP INDEX IF EXISTS idx_media_company_id;
DROP INDEX IF EXISTS idx_stock_adjustments_pending_company_id;
-- Note: orders table will be dropped in migration 00027
DROP INDEX IF EXISTS idx_products_company_id;
DROP INDEX IF EXISTS idx_ingredients_company_id;
-- Note: categories table doesn't exist in this codebase
DROP INDEX IF EXISTS idx_users_company_id;

-- Drop company_id columns
ALTER TABLE media DROP COLUMN company_id;
ALTER TABLE stock_adjustments_pending DROP COLUMN company_id;
-- Note: orders table will be dropped in migration 00027
ALTER TABLE products DROP COLUMN company_id;
ALTER TABLE ingredients DROP COLUMN company_id;
-- Note: categories table doesn't exist in this codebase
ALTER TABLE users DROP COLUMN company_id;
