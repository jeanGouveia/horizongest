-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS products (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    price          NUMERIC(10,2) NOT NULL DEFAULT 0,
    is_composto    BOOLEAN NOT NULL DEFAULT FALSE,
    stock_quantity NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ingredients (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    unit           VARCHAR(50) NOT NULL DEFAULT 'un',
    stock_quantity NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_ingredients (
    id            BIGSERIAL PRIMARY KEY,
    product_id    BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      NUMERIC(10,2) NOT NULL CHECK(quantity > 0),
    UNIQUE(product_id, ingredient_id)
);

CREATE INDEX IF NOT EXISTS idx_product_ingredients_product ON product_ingredients(product_id);

CREATE TABLE IF NOT EXISTS product_compositions (
    id                   BIGSERIAL PRIMARY KEY,
    parent_product_id    BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    component_product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity             NUMERIC(10,2) NOT NULL CHECK(quantity > 0),
    UNIQUE(parent_product_id, component_product_id),
    CHECK(parent_product_id != component_product_id)
);

CREATE INDEX IF NOT EXISTS idx_compositions_parent ON product_compositions(parent_product_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_compositions_parent;
DROP TABLE IF EXISTS product_compositions;
DROP INDEX IF EXISTS idx_product_ingredients_product;
DROP TABLE IF EXISTS product_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
