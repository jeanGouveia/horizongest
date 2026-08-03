-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_sessions (
    id            BIGSERIAL PRIMARY KEY,
    platform_user_id BIGINT NOT NULL REFERENCES platform_users(id) ON DELETE CASCADE,
    token         VARCHAR(255) NOT NULL UNIQUE,
    expires_at    TIMESTAMP NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_platform_sessions_token ON platform_sessions(token);
CREATE INDEX IF NOT EXISTS idx_platform_sessions_platform_user_id ON platform_sessions(platform_user_id);
CREATE INDEX IF NOT EXISTS idx_platform_sessions_expires_at ON platform_sessions(expires_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_platform_sessions_expires_at;
DROP INDEX IF EXISTS idx_platform_sessions_platform_user_id;
DROP INDEX IF EXISTS idx_platform_sessions_token;
DROP TABLE IF EXISTS platform_sessions;

-- +goose StatementEnd
