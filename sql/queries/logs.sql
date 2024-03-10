-- name: SaveLogs :one
INSERT INTO logs (appToken, text, createdAt, updatedAt, level, saved, context,ip, tags) VALUES ($1, $2, now(), now(), $3, 0, $4, $5, $6) RETURNING 1;

-- name: GetAppWithToken :one
SELECT * FROM apps WHERE token = $1;

-- name: GetLogs :many
SELECT * FROM logs WHERE appToken = $1;