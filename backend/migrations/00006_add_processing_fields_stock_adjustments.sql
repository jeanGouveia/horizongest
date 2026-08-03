-- +goose Up
-- +goose StatementBegin

-- Adicionar campos para rastreabilidade operacional
ALTER TABLE stock_adjustments_pending 
ADD COLUMN processed_by BIGINT;

ALTER TABLE stock_adjustments_pending 
ADD COLUMN processing_notes TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE stock_adjustments_pending 
DROP COLUMN processing_notes;

ALTER TABLE stock_adjustments_pending 
DROP COLUMN processed_by;

-- +goose StatementEnd
