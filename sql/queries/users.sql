-- name: GetUser :one
SELECT * FROM users u LEFT JOIN accounts acc ON acc.userid = u.id WHERE u.id = $1;