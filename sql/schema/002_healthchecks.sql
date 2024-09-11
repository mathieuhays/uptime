-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS healthchecks (
    id TEXT PRIMARY KEY,
    website_id TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (website_id) REFERENCES websites(uuid) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE healthchecks;