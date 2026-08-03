-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS global_config (
    id                   BIGSERIAL PRIMARY KEY,
    default_timezone     VARCHAR(50) DEFAULT 'America/Sao_Paulo',
    default_locale       VARCHAR(10) DEFAULT 'pt-BR',
    monetary_format      VARCHAR(50) DEFAULT 'BRL R$ 1.000,00',
    date_format          VARCHAR(20) DEFAULT 'DD/MM/YYYY',
    time_format          VARCHAR(20) DEFAULT 'HH:mm',
    max_upload_size_mb   INTEGER DEFAULT 10,
    max_image_size_mb    INTEGER DEFAULT 5,
    allowed_image_types  TEXT DEFAULT 'jpg,png,webp,gif',
    allowed_file_types   TEXT DEFAULT 'pdf,doc,docx,xlsx,xls,txt',
    maintenance_mode     BOOLEAN DEFAULT FALSE,
    maintenance_message  TEXT,
    enable_finance       BOOLEAN DEFAULT TRUE,
    enable_purchasing    BOOLEAN DEFAULT TRUE,
    enable_inventory     BOOLEAN DEFAULT TRUE,
    enable_crm           BOOLEAN DEFAULT FALSE,
    enable_calendar      BOOLEAN DEFAULT FALSE,
    enable_pos           BOOLEAN DEFAULT FALSE,
    enable_ai            BOOLEAN DEFAULT FALSE,
    enable_delivery      BOOLEAN DEFAULT FALSE,
    enable_marketplace   BOOLEAN DEFAULT FALSE,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by           BIGINT
);

-- Insert default global configuration (idempotent)
INSERT INTO global_config (
    id, default_timezone, default_locale, monetary_format, date_format, time_format,
    max_upload_size_mb, max_image_size_mb, allowed_image_types, allowed_file_types,
    maintenance_mode, maintenance_message, enable_finance, enable_purchasing,
    enable_inventory, enable_crm, enable_calendar, enable_pos, enable_ai,
    enable_delivery, enable_marketplace
) VALUES (
    1, 'America/Sao_Paulo', 'pt-BR', 'BRL R$ 1.000,00', 'DD/MM/YYYY', 'HH:mm',
    10, 5, 'jpg,png,webp,gif', 'pdf,doc,docx,xlsx,xls,txt',
    FALSE, '', TRUE, TRUE,
    TRUE, FALSE, FALSE, FALSE, FALSE,
    FALSE, FALSE
) ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS global_config;

-- +goose StatementEnd
