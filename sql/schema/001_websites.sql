-- +goose Up
CREATE TABLE IF NOT EXISTS websites (
    uuid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    created_at TEXT
);

-- +goose Down
DROP TABLE websites;