-- +goose Up
-- +goose StatementBegin

-- Sprint 3.4 - Security Hardening: Add ON DELETE clauses to foreign keys
-- This prevents orphaned records when parent records are deleted
-- PostgreSQL supports ALTER TABLE to add ON DELETE to existing foreign keys

-- For platform_audit: platform_user_id should be SET NULL on user deletion
ALTER TABLE platform_audit DROP CONSTRAINT IF EXISTS platform_audit_platform_user_id_fkey;
ALTER TABLE platform_audit ADD CONSTRAINT platform_audit_platform_user_id_fkey 
    FOREIGN KEY (platform_user_id) REFERENCES platform_users(id) ON DELETE SET NULL;

-- For invitations: created_by should be SET NULL on user deletion
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_created_by_fkey;
ALTER TABLE invitations ADD CONSTRAINT invitations_created_by_fkey 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert ON DELETE changes
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS invitations_created_by_fkey;
ALTER TABLE invitations ADD CONSTRAINT invitations_created_by_fkey 
    FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE platform_audit DROP CONSTRAINT IF EXISTS platform_audit_platform_user_id_fkey;
ALTER TABLE platform_audit ADD CONSTRAINT platform_audit_platform_user_id_fkey 
    FOREIGN KEY (platform_user_id) REFERENCES platform_users(id);

-- +goose StatementEnd
