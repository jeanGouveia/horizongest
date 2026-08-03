-- +goose Up
-- +goose StatementBegin

-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    price NUMERIC(10,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'BRL',
    interval VARCHAR(20) DEFAULT 'monthly',
    max_users INTEGER DEFAULT 1,
    max_products INTEGER DEFAULT 100,
    features TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_plans_slug ON plans(slug);
CREATE INDEX IF NOT EXISTS idx_plans_active ON plans(active);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_plans_active;
DROP INDEX IF EXISTS idx_plans_slug;
DROP TABLE IF EXISTS plans;

-- +goose StatementEnd
