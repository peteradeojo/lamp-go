-- name: GetApps :many
SELECT * FROM apps LIMIT $1;