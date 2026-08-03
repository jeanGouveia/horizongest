-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS media (
    id            BIGSERIAL PRIMARY KEY,
    file_name     VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path     TEXT NOT NULL,
    thumbnail_path TEXT,
    file_size     BIGINT NOT NULL,
    mime_type     VARCHAR(100) NOT NULL,
    width         BIGINT,
    height        BIGINT,
    alt_text      TEXT,
    entity_type   VARCHAR(50) NOT NULL,
    entity_id     BIGINT,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_media_entity ON media(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_media_deleted_at ON media(deleted_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_media_deleted_at;
DROP INDEX IF EXISTS idx_media_entity;
DROP TABLE IF EXISTS media;

-- +goose StatementEnd
