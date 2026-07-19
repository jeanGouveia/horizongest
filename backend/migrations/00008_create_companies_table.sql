-- +goose Up
-- Migration: Create companies table for multi-tenant support
-- This migration adds the Company (Tenant) entity to support Platform 2.0 multi-tenant architecture
-- Core V1 compatibility: All company_id fields are nullable to maintain backward compatibility

CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    active BOOLEAN DEFAULT 1,
    logo_url VARCHAR(500),
    primary_color VARCHAR(7) DEFAULT '#3b82f6',
    secondary_color VARCHAR(7) DEFAULT '#1e40af',
    deleted_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on slug for faster lookups
CREATE INDEX IF NOT EXISTS idx_companies_slug ON companies(slug);

-- Create index on active for filtering active companies
CREATE INDEX IF NOT EXISTS idx_companies_active ON companies(active);

-- +goose Down
DROP INDEX IF EXISTS idx_companies_active;
DROP INDEX IF EXISTS idx_companies_slug;
DROP TABLE IF EXISTS companies;
