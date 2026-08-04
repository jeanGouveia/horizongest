-- HorizonGest Database Reset Script
-- Purpose: Truncate all domain tables while preserving schema, indexes, constraints, and system data
-- Usage: psql -U horizongest_user -h localhost -d horizongest -f scripts/reset_database.sql

BEGIN;

-- Truncate domain tables (business data)
-- Order matters due to foreign key constraints - CASCADE handles this automatically

TRUNCATE TABLE order_items RESTART IDENTITY CASCADE;
TRUNCATE TABLE orders RESTART IDENTITY CASCADE;
TRUNCATE TABLE product_ingredients RESTART IDENTITY CASCADE;
TRUNCATE TABLE products RESTART IDENTITY CASCADE;
TRUNCATE TABLE ingredients RESTART IDENTITY CASCADE;
TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
TRUNCATE TABLE stock_adjustments_pending RESTART IDENTITY CASCADE;
TRUNCATE TABLE media RESTART IDENTITY CASCADE;
TRUNCATE TABLE invitations RESTART IDENTITY CASCADE;
TRUNCATE TABLE companies RESTART IDENTITY CASCADE;
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- Truncate auth/session tables (temporary data)
TRUNCATE TABLE gorm_token_blacklists RESTART IDENTITY CASCADE;
TRUNCATE TABLE password_reset_tokens RESTART IDENTITY CASCADE;
TRUNCATE TABLE platform_sessions RESTART IDENTITY CASCADE;

-- Truncate audit tables (optional - comment out if you want to preserve audit history)
TRUNCATE TABLE impersonation_audit RESTART IDENTITY CASCADE;
TRUNCATE TABLE platform_audit RESTART IDENTITY CASCADE;

-- Truncate outbox events (event sourcing - comment out if you want to preserve event history)
TRUNCATE TABLE outbox_events RESTART IDENTITY CASCADE;

-- Preserve platform configuration and users (bootstrap data)
-- These tables are NOT truncated:
-- - platform_brand_config (platform branding configuration)
-- - platform_users (platform admin users)

COMMIT;

-- Verification query (optional)
SELECT 'Database reset completed successfully' AS status;
