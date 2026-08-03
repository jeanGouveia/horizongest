-- +goose Up
-- +goose StatementBegin

-- Migration: Migrar campos monetários de Product de NUMERIC(10,2) para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: products
-- Adicionar colunas temporárias em BIGINT
ALTER TABLE products ADD COLUMN price_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN cost_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN cmv_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN profit_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN suggested_price_cents BIGINT DEFAULT 0;

-- Migrar dados: multiplicar por 100 para converter de reais para centavos
UPDATE products SET 
  price_cents = ROUND(price * 100),
  cost_cents = ROUND(cost * 100),
  cmv_cents = ROUND(cmv * 100),
  profit_cents = ROUND(profit * 100),
  suggested_price_cents = ROUND(suggested_price * 100);

-- Remover colunas antigas (NUMERIC)
ALTER TABLE products DROP COLUMN price;
ALTER TABLE products DROP COLUMN cost;
ALTER TABLE products DROP COLUMN cmv;
ALTER TABLE products DROP COLUMN profit;
ALTER TABLE products DROP COLUMN suggested_price;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE products RENAME COLUMN price_cents TO price;
ALTER TABLE products RENAME COLUMN cost_cents TO cost;
ALTER TABLE products RENAME COLUMN cmv_cents TO cmv;
ALTER TABLE products RENAME COLUMN profit_cents TO profit;
ALTER TABLE products RENAME COLUMN suggested_price_cents TO suggested_price;

-- Tabela: product_ingredients
-- Adicionar colunas temporárias em BIGINT
ALTER TABLE product_ingredients ADD COLUMN unit_cost_cents BIGINT DEFAULT 0;
ALTER TABLE product_ingredients ADD COLUMN total_cost_cents BIGINT DEFAULT 0;

-- Migrar dados
UPDATE product_ingredients SET 
  unit_cost_cents = ROUND(unit_cost * 100),
  total_cost_cents = ROUND(total_cost * 100);

-- Remover colunas antigas (NUMERIC)
ALTER TABLE product_ingredients DROP COLUMN unit_cost;
ALTER TABLE product_ingredients DROP COLUMN total_cost;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE product_ingredients RENAME COLUMN unit_cost_cents TO unit_cost;
ALTER TABLE product_ingredients RENAME COLUMN total_cost_cents TO total_cost;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert migration: migrar campos monetários de Product de BIGINT para NUMERIC(10,2)
-- Tabela: product_ingredients
ALTER TABLE product_ingredients ADD COLUMN unit_cost NUMERIC(10,2) DEFAULT 0;
ALTER TABLE product_ingredients ADD COLUMN total_cost NUMERIC(10,2) DEFAULT 0;

UPDATE product_ingredients SET 
  unit_cost = unit_cost_cents / 100.0,
  total_cost = total_cost_cents / 100.0;

ALTER TABLE product_ingredients DROP COLUMN unit_cost_cents;
ALTER TABLE product_ingredients DROP COLUMN total_cost_cents;

ALTER TABLE product_ingredients RENAME COLUMN unit_cost TO unit_cost_cents;
ALTER TABLE product_ingredients RENAME COLUMN total_cost TO total_cost_cents;

-- Tabela: products
ALTER TABLE products ADD COLUMN price NUMERIC(10,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN cost NUMERIC(10,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN cmv NUMERIC(10,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN profit NUMERIC(10,2) DEFAULT 0;
ALTER TABLE products ADD COLUMN suggested_price NUMERIC(10,2) DEFAULT 0;

UPDATE products SET 
  price = price_cents / 100.0,
  cost = cost_cents / 100.0,
  cmv = cmv_cents / 100.0,
  profit = profit_cents / 100.0,
  suggested_price = suggested_price_cents / 100.0;

ALTER TABLE products DROP COLUMN price_cents;
ALTER TABLE products DROP COLUMN cost_cents;
ALTER TABLE products DROP COLUMN cmv_cents;
ALTER TABLE products DROP COLUMN profit_cents;
ALTER TABLE products DROP COLUMN suggested_price_cents;

ALTER TABLE products RENAME COLUMN price TO price_cents;
ALTER TABLE products RENAME COLUMN cost TO cost_cents;
ALTER TABLE products RENAME COLUMN cmv TO cmv_cents;
ALTER TABLE products RENAME COLUMN profit TO profit_cents;
ALTER TABLE products RENAME COLUMN suggested_price TO suggested_price_cents;

-- +goose StatementEnd
