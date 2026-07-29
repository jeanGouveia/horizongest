-- Migration: Migrar campos monetários de Order e OrderItem de REAL para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: orders
-- Adicionar coluna temporária em BIGINT
ALTER TABLE orders ADD COLUMN total_price_cents BIGINT DEFAULT 0;

-- Migrar dados: multiplicar por 100 para converter de reais para centavos
UPDATE orders SET total_price_cents = ROUND(total_price * 100);

-- Remover coluna antiga (REAL)
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

-- Remover colunas antigas (REAL)
ALTER TABLE order_items DROP COLUMN unit_price;
ALTER TABLE order_items DROP COLUMN product_promotion_price;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE order_items RENAME COLUMN unit_price_cents TO unit_price;
ALTER TABLE order_items RENAME COLUMN product_promotion_price_cents TO product_promotion_price;
