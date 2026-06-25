-- name: CreateLoggingSession :one
INSERT INTO logging_sessions (user_id, session_id, title)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLoggingSession :one
SELECT * FROM logging_sessions WHERE id = $1;

-- name: GetLoggingSessionBySessionID :one
SELECT * FROM logging_sessions WHERE user_id = $1 AND session_id = $2;

-- name: ListLoggingSessionsByUser :many
SELECT * FROM logging_sessions WHERE user_id = $1 ORDER BY updated_at DESC;

-- name: UpdateLoggingSessionTitle :one
UPDATE logging_sessions
SET title = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLoggingSession :exec
DELETE FROM logging_sessions WHERE id = $1;
