-- +goose Up
-- +goose StatementBegin

-- Sprint 5D.3 - Performance Hardening: Add composite indexes for performance
-- These indexes optimize queries filtered by tenant (company_id) and commonly used columns

-- Orders: optimize queries by company_id, status, and created_at
CREATE INDEX IF NOT EXISTS idx_orders_company_status_created ON orders(company_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_company_deleted ON orders(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Products: optimize queries by company_id, active status, and deleted_at
CREATE INDEX IF NOT EXISTS idx_products_company_active_deleted ON products(company_id, active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_company_deleted ON products(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Ingredients: optimize queries by company_id, active status, and deleted_at
CREATE INDEX IF NOT EXISTS idx_ingredients_company_active_deleted ON ingredients(company_id, active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_ingredients_company_deleted ON ingredients(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Users: optimize queries by company_id and deleted_at
CREATE INDEX IF NOT EXISTS idx_users_company_deleted ON users(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Companies: optimize queries by deleted_at
CREATE INDEX IF NOT EXISTS idx_companies_deleted ON companies(deleted_at) WHERE deleted_at IS NULL;

-- Stock movements: optimize queries by company_id, type, and created_at
CREATE INDEX IF NOT EXISTS idx_stock_movements_company_type_created ON stock_movements(company_id, movement_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_stock_movements_company_deleted ON stock_movements(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Purchase orders: optimize queries by company_id, status, and created_at
CREATE INDEX IF NOT EXISTS idx_purchase_orders_company_status_created ON purchase_orders(company_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_company_deleted ON purchase_orders(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Transactions: optimize queries by company_id, type, and date
CREATE INDEX IF NOT EXISTS idx_transactions_company_type_date ON transactions(company_id, type, date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_company_deleted ON transactions(company_id, deleted_at) WHERE deleted_at IS NULL;

-- Invitations: optimize queries by company_id and status
CREATE INDEX IF NOT EXISTS idx_invitations_company_status ON invitations(company_id, status);

-- Media: optimize queries by company_id and deleted_at
CREATE INDEX IF NOT EXISTS idx_media_company_deleted ON media(company_id, deleted_at) WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop composite indexes
DROP INDEX IF EXISTS idx_orders_company_status_created;
DROP INDEX IF EXISTS idx_orders_company_deleted;
DROP INDEX IF EXISTS idx_products_company_active_deleted;
DROP INDEX IF EXISTS idx_products_company_deleted;
DROP INDEX IF EXISTS idx_ingredients_company_active_deleted;
DROP INDEX IF EXISTS idx_ingredients_company_deleted;
DROP INDEX IF EXISTS idx_users_company_deleted;
DROP INDEX IF EXISTS idx_companies_deleted;
DROP INDEX IF EXISTS idx_stock_movements_company_type_created;
DROP INDEX IF EXISTS idx_stock_movements_company_deleted;
DROP INDEX IF EXISTS idx_purchase_orders_company_status_created;
DROP INDEX IF EXISTS idx_purchase_orders_company_deleted;
DROP INDEX IF EXISTS idx_transactions_company_type_date;
DROP INDEX IF EXISTS idx_transactions_company_deleted;
DROP INDEX IF EXISTS idx_invitations_company_status;
DROP INDEX IF EXISTS idx_media_company_deleted;

-- +goose StatementEnd
