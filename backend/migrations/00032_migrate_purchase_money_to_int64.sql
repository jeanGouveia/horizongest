-- Migration: Migrar campos monetários de Purchase de REAL para BIGINT (centavos)
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
