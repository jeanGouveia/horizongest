-- +goose Up
-- +goose StatementBegin

-- Sprint 4 - Ficha Técnica Avançada: Add fields for advanced recipe calculations

-- Add fields to product_ingredients table
ALTER TABLE product_ingredients ADD COLUMN loss REAL DEFAULT 0.0;
ALTER TABLE product_ingredients ADD COLUMN yield REAL DEFAULT 1.0;
ALTER TABLE product_ingredients ADD COLUMN unit_cost REAL DEFAULT 0.0;
ALTER TABLE product_ingredients ADD COLUMN total_cost REAL DEFAULT 0.0;

-- Add fields to products table
ALTER TABLE products ADD COLUMN cost REAL DEFAULT 0.0;
ALTER TABLE products ADD COLUMN cmv REAL DEFAULT 0.0;
ALTER TABLE products ADD COLUMN margin REAL DEFAULT 0.0;
ALTER TABLE products ADD COLUMN profit REAL DEFAULT 0.0;
ALTER TABLE products ADD COLUMN suggested_price REAL DEFAULT 0.0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE products DROP COLUMN suggested_price;
ALTER TABLE products DROP COLUMN profit;
ALTER TABLE products DROP COLUMN margin;
ALTER TABLE products DROP COLUMN cmv;
ALTER TABLE products DROP COLUMN cost;
ALTER TABLE product_ingredients DROP COLUMN total_cost;
ALTER TABLE product_ingredients DROP COLUMN unit_cost;
ALTER TABLE product_ingredients DROP COLUMN yield;
ALTER TABLE product_ingredients DROP COLUMN loss;

-- +goose StatementEnd
