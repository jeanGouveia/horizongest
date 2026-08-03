-- +goose Up
-- +goose StatementBegin

-- Add plan and status columns to companies table
ALTER TABLE companies ADD COLUMN plan_id BIGINT;
ALTER TABLE companies ADD COLUMN status VARCHAR(20) DEFAULT 'active';
ALTER TABLE companies ADD COLUMN trial_ends_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_companies_plan_id ON companies(plan_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_companies_plan_id;
ALTER TABLE companies DROP COLUMN trial_ends_at;
ALTER TABLE companies DROP COLUMN status;
ALTER TABLE companies DROP COLUMN plan_id;

-- +goose StatementEnd
