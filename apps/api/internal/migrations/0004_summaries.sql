-- +goose Up
ALTER TABLE usage ADD COLUMN summaries bigint NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE usage DROP COLUMN summaries;
