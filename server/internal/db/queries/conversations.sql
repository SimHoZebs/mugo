-- name: CreateConversation :one
INSERT INTO conversations (user_id, session_id, title)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetConversation :one
SELECT * FROM conversations WHERE id = $1;

-- name: GetConversationBySessionID :one
SELECT * FROM conversations WHERE user_id = $1 AND session_id = $2;

-- name: ListConversationsByUser :many
SELECT * FROM conversations WHERE user_id = $1 ORDER BY updated_at DESC;

-- name: UpdateConversationTitle :one
UPDATE conversations
SET title = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteConversation :exec
DELETE FROM conversations WHERE id = $1;
