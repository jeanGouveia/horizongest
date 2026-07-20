-- Sprint 3.4 - Security Hardening: Add composite indexes for performance
-- These indexes improve query performance for common tenant-scoped queries

-- Index for products: company_id + active (for listing active products by company)
CREATE INDEX IF NOT EXISTS idx_products_company_active ON products(company_id, active) WHERE deleted_at IS NULL;

-- Index for products: company_id + slug (for finding product by slug within company)
CREATE INDEX IF NOT EXISTS idx_products_company_slug ON products(company_id, slug) WHERE deleted_at IS NULL;

-- Index for orders: company_id + status (for filtering orders by company and status)
CREATE INDEX IF NOT EXISTS idx_orders_company_status ON orders(company_id, status) WHERE deleted_at IS NULL;

-- Index for orders: company_id + created_at (for chronological order listing by company)
CREATE INDEX IF NOT EXISTS idx_orders_company_created_at ON orders(company_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index for users: company_id + active (for listing active users by company)
CREATE INDEX IF NOT EXISTS idx_users_company_active ON users(company_id, active) WHERE deleted_at IS NULL;

-- Index for ingredients: company_id + active (for listing active ingredients by company)
CREATE INDEX IF NOT EXISTS idx_ingredients_company_active ON ingredients(company_id, active) WHERE deleted_at IS NULL;

-- Index for companies: active + created_at (for listing active companies chronologically)
CREATE INDEX IF NOT EXISTS idx_companies_active_created ON companies(active, created_at DESC) WHERE deleted_at IS NULL;
