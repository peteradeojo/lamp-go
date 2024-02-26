-- name: SaveLogs :execresult
INSERT INTO logs (appId, text, createdAt, updatedAt, level, saved) VALUES (?, ?, NOW(), NOW(), ?, 0);

-- name: GetAppWithToken :one
SELECT * FROM apps WHERE token = ?;