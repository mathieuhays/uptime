-- name: CreateSession :one
INSERT INTO sessions (id, user_id, refresh_token, expire_at, created_at, updated_at)
VALUES ($1, $2, encode(sha256(random()::text::bytea), 'hex'), $3, $4, $5)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions WHERE user_id = $1 ORDER BY id;