-- +goose Up
-- +goose StatementBegin

-- Migration: Migrar campos monetários de Purchase de NUMERIC(10,2) para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: purchase_orders
ALTER TABLE purchase_orders ADD COLUMN subtotal_cents BIGINT DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN tax_cents BIGINT DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN discount_cents BIGINT DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN total_cents BIGINT DEFAULT 0;

UPDATE purchase_orders SET 
  subtotal_cents = ROUND(subtotal * 100),
  tax_cents = ROUND(tax * 100),
  discount_cents = ROUND(discount * 100),
  total_cents = ROUND(total * 100);

ALTER TABLE purchase_orders DROP COLUMN subtotal;
ALTER TABLE purchase_orders DROP COLUMN tax;
ALTER TABLE purchase_orders DROP COLUMN discount;
ALTER TABLE purchase_orders DROP COLUMN total;

ALTER TABLE purchase_orders RENAME COLUMN subtotal_cents TO subtotal;
ALTER TABLE purchase_orders RENAME COLUMN tax_cents TO tax;
ALTER TABLE purchase_orders RENAME COLUMN discount_cents TO discount;
ALTER TABLE purchase_orders RENAME COLUMN total_cents TO total;

-- Tabela: purchase_order_items
ALTER TABLE purchase_order_items ADD COLUMN unit_price_cents BIGINT DEFAULT 0;
ALTER TABLE purchase_order_items ADD COLUMN subtotal_cents BIGINT DEFAULT 0;

UPDATE purchase_order_items SET 
  unit_price_cents = ROUND(unit_price * 100),
  subtotal_cents = ROUND(subtotal * 100);

ALTER TABLE purchase_order_items DROP COLUMN unit_price;
ALTER TABLE purchase_order_items DROP COLUMN subtotal;

ALTER TABLE purchase_order_items RENAME COLUMN unit_price_cents TO unit_price;
ALTER TABLE purchase_order_items RENAME COLUMN subtotal_cents TO subtotal;

-- Tabela: purchase_receiving_items
ALTER TABLE purchase_receiving_items ADD COLUMN unit_price_cents BIGINT DEFAULT 0;
ALTER TABLE purchase_receiving_items ADD COLUMN subtotal_cents BIGINT DEFAULT 0;

UPDATE purchase_receiving_items SET 
  unit_price_cents = ROUND(unit_price * 100),
  subtotal_cents = ROUND(subtotal * 100);

ALTER TABLE purchase_receiving_items DROP COLUMN unit_price;
ALTER TABLE purchase_receiving_items DROP COLUMN subtotal;

ALTER TABLE purchase_receiving_items RENAME COLUMN unit_price_cents TO unit_price;
ALTER TABLE purchase_receiving_items RENAME COLUMN subtotal_cents TO subtotal;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert migration: migrar campos monetários de Purchase de BIGINT para NUMERIC(10,2)
-- Tabela: purchase_receiving_items
ALTER TABLE purchase_receiving_items ADD COLUMN unit_price NUMERIC(10,2) DEFAULT 0;
ALTER TABLE purchase_receiving_items ADD COLUMN subtotal NUMERIC(10,2) DEFAULT 0;

UPDATE purchase_receiving_items SET 
  unit_price = unit_price_cents / 100.0,
  subtotal = subtotal_cents / 100.0;

ALTER TABLE purchase_receiving_items DROP COLUMN unit_price_cents;
ALTER TABLE purchase_receiving_items DROP COLUMN subtotal_cents;

ALTER TABLE purchase_receiving_items RENAME COLUMN unit_price TO unit_price_cents;
ALTER TABLE purchase_receiving_items RENAME COLUMN subtotal TO subtotal_cents;

-- Tabela: purchase_order_items
ALTER TABLE purchase_order_items ADD COLUMN unit_price NUMERIC(10,2) DEFAULT 0;
ALTER TABLE purchase_order_items ADD COLUMN subtotal NUMERIC(10,2) DEFAULT 0;

UPDATE purchase_order_items SET 
  unit_price = unit_price_cents / 100.0,
  subtotal = subtotal_cents / 100.0;

ALTER TABLE purchase_order_items DROP COLUMN unit_price_cents;
ALTER TABLE purchase_order_items DROP COLUMN subtotal_cents;

ALTER TABLE purchase_order_items RENAME COLUMN unit_price TO unit_price_cents;
ALTER TABLE purchase_order_items RENAME COLUMN subtotal TO subtotal_cents;

-- Tabela: purchase_orders
ALTER TABLE purchase_orders ADD COLUMN subtotal NUMERIC(10,2) DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN tax NUMERIC(10,2) DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN discount NUMERIC(10,2) DEFAULT 0;
ALTER TABLE purchase_orders ADD COLUMN total NUMERIC(10,2) DEFAULT 0;

UPDATE purchase_orders SET 
  subtotal = subtotal_cents / 100.0,
  tax = tax_cents / 100.0,
  discount = discount_cents / 100.0,
  total = total_cents / 100.0;

ALTER TABLE purchase_orders DROP COLUMN subtotal_cents;
ALTER TABLE purchase_orders DROP COLUMN tax_cents;
ALTER TABLE purchase_orders DROP COLUMN discount_cents;
ALTER TABLE purchase_orders DROP COLUMN total_cents;

ALTER TABLE purchase_orders RENAME COLUMN subtotal TO subtotal_cents;
ALTER TABLE purchase_orders RENAME COLUMN tax TO tax_cents;
ALTER TABLE purchase_orders RENAME COLUMN discount TO discount_cents;
ALTER TABLE purchase_orders RENAME COLUMN total TO total_cents;

-- +goose StatementEnd
