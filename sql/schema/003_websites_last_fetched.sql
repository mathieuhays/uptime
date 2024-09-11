-- +goose Up
ALTER TABLE websites ADD COLUMN last_fetched_at INTEGER AFTER 'url';

-- +goose Down
ALTER TABLE websites DROP COLUMN last_fetched_at;