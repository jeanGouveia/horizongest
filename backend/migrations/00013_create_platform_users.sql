-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_users (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    name          TEXT     NOT NULL,
    email         TEXT     NOT NULL UNIQUE,
    password_hash TEXT     NOT NULL,
    role          TEXT     NOT NULL DEFAULT 'PlatformSupport',
    active        INTEGER  NOT NULL DEFAULT 1 CHECK(active IN (0,1)),
    deleted_at    INTEGER,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_platform_users_email ON platform_users(email);
CREATE INDEX IF NOT EXISTS idx_platform_users_role ON platform_users(role);
CREATE INDEX IF NOT EXISTS idx_platform_users_deleted_at ON platform_users(deleted_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_platform_users_deleted_at;
DROP INDEX IF EXISTS idx_platform_users_role;
DROP INDEX IF EXISTS idx_platform_users_email;
DROP TABLE IF EXISTS platform_users;

-- +goose StatementEnd
