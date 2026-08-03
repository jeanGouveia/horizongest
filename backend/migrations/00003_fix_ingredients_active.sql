-- +goose Up

ALTER TABLE ingredients
ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down

ALTER TABLE ingredients DROP COLUMN active;