-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS media (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    file_name     TEXT     NOT NULL,
    original_name TEXT     NOT NULL,
    file_path     TEXT     NOT NULL,
    thumbnail_path TEXT,
    file_size     INTEGER  NOT NULL,
    mime_type     TEXT     NOT NULL,
    width         INTEGER,
    height        INTEGER,
    alt_text      TEXT,
    entity_type   TEXT     NOT NULL,
    entity_id     INTEGER,
    deleted_at    INTEGER,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
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
