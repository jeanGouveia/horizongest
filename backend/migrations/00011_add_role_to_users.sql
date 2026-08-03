-- +goose Up
-- +goose StatementBegin

-- Add role column to users table for RBAC (Sprint 6)
-- Nullable for Core V1 compatibility (users without CompanyID have Role == null)

ALTER TABLE users ADD COLUMN role VARCHAR(50) NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN role;

-- +goose StatementEnd
