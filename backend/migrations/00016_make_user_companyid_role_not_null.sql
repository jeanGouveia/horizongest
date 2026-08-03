-- +goose Up
-- +goose StatementBegin

-- Sprint 3: Make CompanyID and Role NOT NULL in users table
-- This migration enforces that all users must belong to a company

-- Step 1: Add active column if it doesn't exist (PostgreSQL supports IF NOT EXISTS)
ALTER TABLE users ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT TRUE;

-- Step 2: Add deleted_at column if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Step 3: Add default values for existing NULL records (if any)
UPDATE users SET company_id = 0 WHERE company_id IS NULL;
UPDATE users SET role = 'owner' WHERE role IS NULL;

-- Step 4: Make columns NOT NULL (PostgreSQL supports ALTER COLUMN directly)
ALTER TABLE users ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN company_id SET DEFAULT 0;
ALTER TABLE users ALTER COLUMN role SET NOT NULL;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'owner';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Revert to nullable CompanyID and Role
ALTER TABLE users ALTER COLUMN company_id DROP NOT NULL;
ALTER TABLE users ALTER COLUMN company_id DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role DROP NOT NULL;
ALTER TABLE users ALTER COLUMN role DROP DEFAULT;

-- +goose StatementEnd
