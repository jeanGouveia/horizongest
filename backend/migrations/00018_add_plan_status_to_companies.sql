-- Add plan and status columns to companies table
ALTER TABLE companies 
ADD COLUMN plan_id BIGINT UNSIGNED NULL,
ADD COLUMN status VARCHAR(20) DEFAULT 'active',
ADD COLUMN trial_ends_at TIMESTAMP NULL,
ADD INDEX idx_plan_id (plan_id),
ADD CONSTRAINT fk_companies_plan_id FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE SET NULL;
