-- Sprint 3.4 - Security Hardening: Add ON DELETE clauses to foreign keys
-- This prevents orphaned records when parent records are deleted

-- Note: SQLite does not support ALTER TABLE to add ON DELETE to existing foreign keys
-- Tables need to be recreated. For now, this migration documents the required changes.
-- In production with PostgreSQL/MySQL, these would be ALTER TABLE statements.

-- For platform_audit: platform_user_id should be SET NULL on user deletion
-- ALTER TABLE platform_audit DROP CONSTRAINT platform_audit_platform_user_id_fkey;
-- ALTER TABLE platform_audit ADD CONSTRAINT platform_audit_platform_user_id_fkey 
--     FOREIGN KEY (platform_user_id) REFERENCES platform_users(id) ON DELETE SET NULL;

-- For products: category_id should be SET NULL on category deletion
-- ALTER TABLE products DROP CONSTRAINT products_category_id_fkey;
-- ALTER TABLE products ADD CONSTRAINT products_category_id_fkey 
--     FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;

-- For order_items: product_id should be SET NULL on product deletion (keep order history)
-- ALTER TABLE order_items DROP CONSTRAINT order_items_product_id_fkey;
-- ALTER TABLE order_items ADD CONSTRAINT order_items_product_id_fkey 
--     FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;

-- For invitations: created_by should be SET NULL on user deletion
-- ALTER TABLE invitations DROP CONSTRAINT invitations_created_by_fkey;
-- ALTER TABLE invitations ADD CONSTRAINT invitations_created_by_fkey 
--     FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

-- For stock_adjustments_pending: order_id should be CASCADE on order deletion
-- ALTER TABLE stock_adjustments_pending DROP CONSTRAINT stock_adjustments_pending_order_id_fkey;
-- ALTER TABLE stock_adjustments_pending ADD CONSTRAINT stock_adjustments_pending_order_id_fkey 
--     FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

-- For stock_adjustments_pending: ingredient_id should be CASCADE on ingredient deletion
-- ALTER TABLE stock_adjustments_pending DROP CONSTRAINT stock_adjustments_pending_ingredient_id_fkey;
-- ALTER TABLE stock_adjustments_pending ADD CONSTRAINT stock_adjustments_pending_ingredient_id_fkey 
--     FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE;

-- For SQLite, these changes require recreating tables with the proper FK constraints
-- This is documented for future migration to PostgreSQL/MySQL
