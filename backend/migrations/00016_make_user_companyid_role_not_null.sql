-- +goose Up
-- +goose StatementBegin

-- Sprint 3: Make CompanyID and Role NOT NULL in users table
-- This migration enforces that all users must belong to a company

-- Step 1: Add active column if it doesn't exist (ignore error if already exists)
-- SQLite doesn't support IF NOT EXISTS for ALTER TABLE, so we try and ignore
ALTER TABLE users ADD COLUMN active INTEGER NOT NULL DEFAULT 1 CHECK(active IN (0,1));

-- Step 2: Add deleted_at column if it doesn't exist
ALTER TABLE users ADD COLUMN deleted_at INTEGER;

-- Step 3: Add default values for existing NULL records (if any)
UPDATE users SET company_id = 0 WHERE company_id IS NULL;
UPDATE users SET role = 'owner' WHERE role IS NULL;

-- Step 4: Create new table with NOT NULL constraints
-- SQLite doesn't support ALTER COLUMN directly, so we need to recreate the table

-- Create new users table with NOT NULL constraints
CREATE TABLE IF NOT EXISTS users_new (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    name          TEXT     NOT NULL,
    email         TEXT     NOT NULL UNIQUE,
    password_hash TEXT     NOT NULL,
    active        INTEGER  NOT NULL DEFAULT 1 CHECK(active IN (0,1)),
    company_id    INTEGER  NOT NULL DEFAULT 0,
    role          TEXT     NOT NULL DEFAULT 'owner',
    deleted_at    INTEGER,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- Copy data from old table to new table
INSERT INTO users_new (id, name, email, password_hash, active, company_id, role, deleted_at, created_at, updated_at)
SELECT id, name, email, password_hash,
       COALESCE(active, 1) as active,
       COALESCE(company_id, 0) as company_id,
       COALESCE(role, 'owner') as role,
       deleted_at,
       created_at,
       updated_at
FROM users;

-- Drop old table
DROP TABLE users;

-- Rename new table to users
ALTER TABLE users_new RENAME TO users;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Revert to nullable CompanyID and Role
CREATE TABLE IF NOT EXISTS users_new (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    name          TEXT     NOT NULL,
    email         TEXT     NOT NULL UNIQUE,
    password_hash TEXT     NOT NULL,
    active        INTEGER  NOT NULL DEFAULT 1 CHECK(active IN (0,1)),
    company_id    INTEGER,
    role          TEXT,
    deleted_at    INTEGER,
    created_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at    INTEGER  NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO users_new (id, name, email, password_hash, active, company_id, role, deleted_at, created_at, updated_at)
SELECT id, name, email, password_hash, active, company_id, role, deleted_at, created_at, updated_at
FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- +goose StatementEnd
