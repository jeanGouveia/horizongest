-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_sessions (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    platform_user_id INTEGER NOT NULL REFERENCES platform_users(id) ON DELETE CASCADE,
    token         TEXT     NOT NULL UNIQUE,
    expires_at    INTEGER  NOT NULL,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
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
