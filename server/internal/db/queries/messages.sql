-- name: CreateMessage :one
INSERT INTO conversation_messages (conversation_id, role, content, metadata)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM conversation_messages WHERE id = $1;

-- name: ListMessagesByConversation :many
SELECT * FROM conversation_messages WHERE conversation_id = $1 ORDER BY created_at ASC;

-- name: ListRecentMessagesByConversation :many
SELECT * FROM conversation_messages 
WHERE conversation_id = $1 
ORDER BY created_at DESC
LIMIT $2;

-- name: DeleteMessagesByConversation :exec
DELETE FROM conversation_messages WHERE conversation_id = $1;
