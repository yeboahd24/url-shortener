-- queries.sql
-- name: CreateUser :one
INSERT INTO users (user_id, username, email, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE key = $1;

-- name: CreateAPIKey :one
INSERT INTO api_keys (key, user_id, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateURL :one
INSERT INTO urls (short_id, long_url, user_id, created_at, expires_at, click_limit)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetURL :one
SELECT * FROM urls WHERE short_id = $1;

-- name: LogClick :exec
INSERT INTO clicks (short_id, ip_address, user_agent, clicked_at)
VALUES ($1, $2, $3, $4);

-- name: ListClicks :many
SELECT * FROM clicks WHERE short_id = $1;

-- name: ListUserURLs :many
SELECT * FROM urls WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteURL :exec
DELETE FROM urls WHERE short_id = $1 AND user_id = $2;

-- name: UpdateURL :one
UPDATE urls
SET long_url = COALESCE($2, long_url),
    expires_at = COALESCE($3, expires_at),
    click_limit = COALESCE($4, click_limit)
WHERE short_id = $1 AND user_id = $5
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUserAPIKeys :many
SELECT * FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE key = $1 AND user_id = $2;

-- name: GetTotalURLs :one
SELECT COUNT(*) as total FROM urls;

-- name: GetTotalClicks :one
SELECT COUNT(*) as total FROM clicks;

-- name: GetTotalUsers :one
SELECT COUNT(*) as total FROM users;
