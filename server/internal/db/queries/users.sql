-- name: CreateUser :one
INSERT INTO users (username)
VALUES ($1)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) AS exists;

-- name: UpdateUser :one
UPDATE users
SET username = $2, metadata = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
