-- +goose Up
-- +goose StatementBegin

-- Sprint 3.4 - Security Hardening: Add composite indexes for performance
-- Note: orders table will be created in migration 00027 with composite indexes already included
-- Other indexes will be added in future migrations when columns are available

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- No indexes to drop in this migration

-- +goose StatementEnd
