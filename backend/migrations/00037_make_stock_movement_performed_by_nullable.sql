-- +goose Up
-- +goose StatementBegin

-- Sprint 2.1: Make StockMovement.performed_by nullable
-- This allows system operations (e.g., iFood webhooks) to create stock movements
-- without requiring a user ID, since webhooks are automated system events

-- Drop the NOT NULL constraint and FOREIGN KEY constraint
ALTER TABLE stock_movements ALTER COLUMN performed_by DROP NOT NULL;
ALTER TABLE stock_movements DROP CONSTRAINT stock_movements_performed_by_fkey;

-- Add index for performed_by (nullable queries)
CREATE INDEX IF NOT EXISTS idx_stock_movements_performed_by ON stock_movements(performed_by);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Restore the NOT NULL constraint and FOREIGN KEY constraint
DROP INDEX IF EXISTS idx_stock_movements_performed_by;
ALTER TABLE stock_movements ADD CONSTRAINT stock_movements_performed_by_fkey FOREIGN KEY (performed_by) REFERENCES users(id);
ALTER TABLE stock_movements ALTER COLUMN performed_by SET NOT NULL;

-- +goose StatementEnd
