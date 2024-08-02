-- name: AddComment :exec
INSERT INTO task_comments (message, sender_id, task_id, createdat, updatedat) VALUES ($1, $2, $3, now(), now());