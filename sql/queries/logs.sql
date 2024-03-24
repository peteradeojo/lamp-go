-- name: SaveLogs :one
INSERT INTO logs (apptoken, text, createdat, updatedat, level, context,ip, tags) VALUES ($1, $2, now(), now(), $3, $4, $5, $6) RETURNING 1;

-- name: GetAppWithToken :one
SELECT * FROM apps WHERE token = $1;

-- name: GetLogs :many
SELECT * FROM logs WHERE appToken = $1;

-- name: CreateSystemLog :exec
INSERT INTO system_logs (id, text, stack, context, level, from_system, createdat, updatedat, origin) VALUES (
  uuid_generate_v4(), 
  $1, 
  $2, 
  $3, 
  $4, 
  B'1', 
  now(), 
  now(), 
  'go-api'
);