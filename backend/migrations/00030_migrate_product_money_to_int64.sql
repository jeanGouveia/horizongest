-- Migration: Migrar campos monetários de Product de REAL para BIGINT (centavos)
-- Sprint 5A: Eliminação completa de float64 monetário

-- Tabela: products
-- Adicionar colunas temporárias em BIGINT
ALTER TABLE products ADD COLUMN price_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN promotion_price_cents BIGINT;
ALTER TABLE products ADD COLUMN cost_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN cmv_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN profit_cents BIGINT DEFAULT 0;
ALTER TABLE products ADD COLUMN suggested_price_cents BIGINT DEFAULT 0;

-- Migrar dados: multiplicar por 100 para converter de reais para centavos
UPDATE products SET 
  price_cents = ROUND(price * 100),
  promotion_price_cents = CASE WHEN promotion_price IS NOT NULL THEN ROUND(promotion_price * 100) ELSE NULL END,
  cost_cents = ROUND(cost * 100),
  cmv_cents = ROUND(cmv * 100),
  profit_cents = ROUND(profit * 100),
  suggested_price_cents = ROUND(suggested_price * 100);

-- Remover colunas antigas (REAL)
ALTER TABLE products DROP COLUMN price;
ALTER TABLE products DROP COLUMN promotion_price;
ALTER TABLE products DROP COLUMN cost;
ALTER TABLE products DROP COLUMN cmv;
ALTER TABLE products DROP COLUMN profit;
ALTER TABLE products DROP COLUMN suggested_price;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE products RENAME COLUMN price_cents TO price;
ALTER TABLE products RENAME COLUMN promotion_price_cents TO promotion_price;
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

-- Remover colunas antigas (REAL)
ALTER TABLE product_ingredients DROP COLUMN unit_cost;
ALTER TABLE product_ingredients DROP COLUMN total_cost;

-- Renomear colunas temporárias para nomes originais
ALTER TABLE product_ingredients RENAME COLUMN unit_cost_cents TO unit_cost;
ALTER TABLE product_ingredients RENAME COLUMN total_cost_cents TO total_cost;
