-- name: CreateUser :one
INSERT INTO users (username, metadata)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS exists;
