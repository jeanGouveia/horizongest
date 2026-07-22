-- +goose Up
-- +goose StatementBegin

-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    price REAL DEFAULT 0,
    currency TEXT DEFAULT 'BRL',
    interval TEXT DEFAULT 'monthly',
    max_users INTEGER DEFAULT 1,
    max_products INTEGER DEFAULT 100,
    features TEXT,
    active INTEGER DEFAULT 1,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
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
