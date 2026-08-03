-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          VARCHAR(50) NOT NULL DEFAULT 'PlatformSupport',
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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
