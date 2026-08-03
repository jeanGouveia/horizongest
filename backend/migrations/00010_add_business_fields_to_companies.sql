-- +goose Up
-- +goose StatementBegin

-- Migration: Add Business Engine fields to companies table
-- Version: 2.0.3
-- Description: Adds business_type, locale, currency, and timezone fields to support Business Engine

-- Add business_type field (nullable for backward compatibility)
ALTER TABLE companies ADD COLUMN business_type VARCHAR(50);

-- Add locale field (nullable for backward compatibility, defaults to 'pt-BR')
ALTER TABLE companies ADD COLUMN locale VARCHAR(10) DEFAULT 'pt-BR';

-- Add currency field (nullable for backward compatibility, defaults to 'BRL')
ALTER TABLE companies ADD COLUMN currency VARCHAR(10) DEFAULT 'BRL';

-- Add timezone field (nullable for backward compatibility, defaults to 'America/Sao_Paulo')
ALTER TABLE companies ADD COLUMN timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo';

-- Update existing companies with default values
UPDATE companies
SET
    business_type = 'generic',
    locale = 'pt-BR',
    currency = 'BRL',
    timezone = 'America/Sao_Paulo'
WHERE business_type IS NULL OR locale IS NULL OR currency IS NULL OR timezone IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE companies DROP COLUMN timezone;
ALTER TABLE companies DROP COLUMN currency;
ALTER TABLE companies DROP COLUMN locale;
ALTER TABLE companies DROP COLUMN business_type;

-- +goose StatementEnd
