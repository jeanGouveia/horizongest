-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_brand_config (
    id                   INTEGER  PRIMARY KEY AUTOINCREMENT,
    platform_name        TEXT     NOT NULL,
    platform_short_name  TEXT     NOT NULL,
    owner_company_name   TEXT     NOT NULL,
    owner_document       TEXT,
    website              TEXT     NOT NULL,
    support_email        TEXT     NOT NULL,
    support_url          TEXT     NOT NULL,
    logo_path            TEXT,
    favicon_path         TEXT,
    logo_light           TEXT,
    logo_dark            TEXT,
    icon                 TEXT,
    login_background     TEXT,
    login_illustration   TEXT,
    copyright            TEXT     NOT NULL,
    privacy_policy_url   TEXT,
    terms_url            TEXT,
    instagram_url        TEXT,
    facebook_url         TEXT,
    linkedin_url         TEXT,
    youtube_url          TEXT,
    default_language     TEXT,
    default_timezone     TEXT,
    maintenance_mode     INTEGER  DEFAULT 0,
    maintenance_message  TEXT,
    primary_color        TEXT     NOT NULL,
    secondary_color      TEXT     NOT NULL,
    updated_at           INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_by           INTEGER
);

-- Insert default HorizonGest branding (idempotent)
INSERT OR IGNORE INTO platform_brand_config (
    id, platform_name, platform_short_name, owner_company_name, owner_document, website,
    support_email, support_url, logo_path, favicon_path, logo_light, logo_dark, icon,
    login_background, login_illustration, copyright, privacy_policy_url, terms_url,
    instagram_url, facebook_url, linkedin_url, youtube_url, default_language,
    default_timezone, maintenance_mode, maintenance_message, primary_color, secondary_color
) VALUES (
    1, 'HorizonGest', 'Horizon', 'HorizonGest Inc.', '', 'https://horizongest.com',
    'support@horizongest.com', 'https://help.horizongest.com',
    '/assets/platform/logo.svg', '/assets/platform/favicon.ico', '', '', '',
    '', '', '© 2024 HorizonGest Inc. All rights reserved.', '', '',
    '', '', '', '', 'pt-BR',
    'America/Sao_Paulo', 0, '', '#0f172a', '#6366f1'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS platform_brand_config;

-- +goose StatementEnd
