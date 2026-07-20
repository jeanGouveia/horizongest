-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_audit (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    platform_user_id INTEGER,
    action        TEXT     NOT NULL,
    entity_type   TEXT     NOT NULL,
    entity_id     INTEGER,
    changes       TEXT,    -- JSON string of changes
    ip_address    TEXT,
    user_agent    TEXT,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_platform_audit_platform_user_id ON platform_audit(platform_user_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_entity ON platform_audit(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_platform_audit_created_at ON platform_audit(created_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_platform_audit_created_at;
DROP INDEX IF EXISTS idx_platform_audit_entity;
DROP INDEX IF EXISTS idx_platform_audit_platform_user_id;
DROP TABLE IF EXISTS platform_audit;

-- +goose StatementEnd
