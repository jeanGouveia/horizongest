-- +goose Up
-- +goose StatementBegin

-- Migration: Migrar campos monetários de Order e OrderItem de NUMERIC(10,2) para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: orders
-- Adicionar coluna temporária em BIGINT
ALTER TABLE orders ADD COLUMN total_price_cents BIGINT DEFAULT 0;

-- Migrar dados: multiplicar por 100 para converter de reais para centavos
UPDATE orders SET total_price_cents = ROUND(total_price * 100);

-- Remover coluna antiga (NUMERIC)
ALTER TABLE orders DROP COLUMN total_price;

-- Renomear coluna temporária para nome original
ALTER TABLE orders RENAME COLUMN total_price_cents TO total_price;

-- Tabela: order_items
-- Adicionar colunas temporárias em BIGINT
ALTER TABLE order_items ADD COLUMN unit_price_cents BIGINT DEFAULT 0;
ALTER TABLE order_items ADD COLUMN product_promotion_price_cents BIGINT;

-- Migrar dados
UPDATE order_items SET 
  unit_price_cents = ROUND(unit_price * 100),
  product_promotion_price_cents = CASE WHEN product_promotion_price IS NOT NULL THEN ROUND(product_promotion_price * 100) ELSE NULL END;

-- Remover colunas antigas (NUMERIC)
ALTER TABLE order_items DROP COLUMN unit_price;
ALTER TABLE order_items DROP COLUMN product_promotion_price;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE order_items RENAME COLUMN unit_price_cents TO unit_price;
ALTER TABLE order_items RENAME COLUMN product_promotion_price_cents TO product_promotion_price;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert migration: migrar campos monetários de Order e OrderItem de BIGINT para NUMERIC(10,2)
-- Tabela: order_items
ALTER TABLE order_items ADD COLUMN unit_price NUMERIC(10,2) DEFAULT 0;
ALTER TABLE order_items ADD COLUMN product_promotion_price NUMERIC(10,2);

UPDATE order_items SET 
  unit_price = unit_price_cents / 100.0,
  product_promotion_price = CASE WHEN product_promotion_price_cents IS NOT NULL THEN product_promotion_price_cents / 100.0 ELSE NULL END;

ALTER TABLE order_items DROP COLUMN unit_price_cents;
ALTER TABLE order_items DROP COLUMN product_promotion_price_cents;

ALTER TABLE order_items RENAME COLUMN unit_price TO unit_price_cents;
ALTER TABLE order_items RENAME COLUMN product_promotion_price TO product_promotion_price_cents;

-- Tabela: orders
ALTER TABLE orders ADD COLUMN total_price NUMERIC(10,2) DEFAULT 0;

UPDATE orders SET total_price = total_price_cents / 100.0;

ALTER TABLE orders DROP COLUMN total_price_cents;

ALTER TABLE orders RENAME COLUMN total_price TO total_price_cents;

-- +goose StatementEnd
