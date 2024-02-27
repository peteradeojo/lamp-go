-- name: SaveLogs :execresult
INSERT INTO logs (appToken, text, createdAt, updatedAt, level, saved, context) VALUES (?, ?, NOW(), NOW(), ?, 0, ?);

-- name: GetAppWithToken :one
SELECT * FROM apps WHERE token = ?;