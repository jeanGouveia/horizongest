-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS global_config (
    id                   INTEGER  PRIMARY KEY AUTOINCREMENT,
    default_timezone     TEXT     DEFAULT 'America/Sao_Paulo',
    default_locale       TEXT     DEFAULT 'pt-BR',
    monetary_format      TEXT     DEFAULT 'BRL R$ 1.000,00',
    date_format          TEXT     DEFAULT 'DD/MM/YYYY',
    time_format          TEXT     DEFAULT 'HH:mm',
    max_upload_size_mb   INTEGER  DEFAULT 10,
    max_image_size_mb    INTEGER  DEFAULT 5,
    allowed_image_types  TEXT     DEFAULT 'jpg,png,webp,gif',
    allowed_file_types   TEXT     DEFAULT 'pdf,doc,docx,xlsx,xls,txt',
    maintenance_mode     INTEGER  DEFAULT 0,
    maintenance_message  TEXT,
    enable_finance       INTEGER  DEFAULT 1,
    enable_purchasing    INTEGER  DEFAULT 1,
    enable_inventory     INTEGER  DEFAULT 1,
    enable_crm           INTEGER  DEFAULT 0,
    enable_calendar      INTEGER  DEFAULT 0,
    enable_pos           INTEGER  DEFAULT 0,
    enable_ai            INTEGER  DEFAULT 0,
    enable_delivery      INTEGER  DEFAULT 0,
    enable_marketplace   INTEGER  DEFAULT 0,
    updated_at           INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_by           INTEGER
);

-- Insert default global configuration (idempotent)
INSERT OR IGNORE INTO global_config (
    id, default_timezone, default_locale, monetary_format, date_format, time_format,
    max_upload_size_mb, max_image_size_mb, allowed_image_types, allowed_file_types,
    maintenance_mode, maintenance_message, enable_finance, enable_purchasing,
    enable_inventory, enable_crm, enable_calendar, enable_pos, enable_ai,
    enable_delivery, enable_marketplace
) VALUES (
    1, 'America/Sao_Paulo', 'pt-BR', 'BRL R$ 1.000,00', 'DD/MM/YYYY', 'HH:mm',
    10, 5, 'jpg,png,webp,gif', 'pdf,doc,docx,xlsx,xls,txt',
    0, '', 1, 1,
    1, 0, 0, 0, 0,
    0, 0
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS global_config;

-- +goose StatementEnd
