-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS platform_brand_config (
    id                   BIGSERIAL PRIMARY KEY,
    platform_name        VARCHAR(255) NOT NULL,
    platform_short_name  VARCHAR(100) NOT NULL,
    owner_company_name   VARCHAR(255) NOT NULL,
    owner_document       VARCHAR(50),
    website              VARCHAR(500) NOT NULL,
    support_email        VARCHAR(255) NOT NULL,
    support_url          VARCHAR(500) NOT NULL,
    logo_path            TEXT,
    favicon_path         TEXT,
    logo_light           TEXT,
    logo_dark            TEXT,
    icon                 TEXT,
    login_background     TEXT,
    login_illustration   TEXT,
    copyright            VARCHAR(500) NOT NULL,
    privacy_policy_url   TEXT,
    terms_url            TEXT,
    instagram_url        TEXT,
    facebook_url         TEXT,
    linkedin_url         TEXT,
    youtube_url          TEXT,
    default_language     VARCHAR(10),
    default_timezone     VARCHAR(50),
    maintenance_mode     BOOLEAN DEFAULT FALSE,
    maintenance_message  TEXT,
    primary_color        VARCHAR(7) NOT NULL,
    secondary_color      VARCHAR(7) NOT NULL,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by           BIGINT
);

-- Insert default HorizonGest branding (idempotent)
INSERT INTO platform_brand_config (
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
    'America/Sao_Paulo', FALSE, '', '#0f172a', '#6366f1'
) ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS platform_brand_config;

-- +goose StatementEnd
