-- Add role column to users table for RBAC (Sprint 6)
-- Nullable for Core V1 compatibility (users without CompanyID have Role == null)

ALTER TABLE users ADD COLUMN role TEXT NULL;

-- Add comment to document the purpose
-- Note: SQLite doesn't support COMMENT, but the column name is self-explanatory
