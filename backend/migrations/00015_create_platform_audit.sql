-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_audit (
    id            BIGSERIAL PRIMARY KEY,
    platform_user_id BIGINT,
    action        VARCHAR(50) NOT NULL,
    entity_type   VARCHAR(50) NOT NULL,
    entity_id     BIGINT,
    changes       TEXT,    -- JSON string of changes
    ip_address    VARCHAR(50),
    user_agent    TEXT,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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
