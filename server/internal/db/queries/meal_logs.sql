-- name: CreateMealLog :one
INSERT INTO meal_logs (
    user_id, conversation_id, food_name, meal_type, recorded_at,
    macros, assumptions, food_source, raw_response
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetMealLog :one
SELECT * FROM meal_logs WHERE id = $1;

-- name: GetMealLogWithSession :one
SELECT 
    m.*,
    c.session_id as adk_session_id
FROM meal_logs m
LEFT JOIN conversations c ON m.conversation_id = c.id
WHERE m.id = $1;

-- name: ListMealLogsByUser :many
SELECT * FROM meal_logs 
WHERE user_id = $1 
ORDER BY recorded_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMealLogsByUserAndDate :many
SELECT * FROM meal_logs 
WHERE user_id = $1 AND recorded_at::date = $2
ORDER BY recorded_at ASC;

-- name: ListMealLogsByUserAndDateRange :many
SELECT * FROM meal_logs 
WHERE user_id = $1 
AND recorded_at >= sqlc.arg('start_date')
AND recorded_at < sqlc.arg('end_date')
ORDER BY recorded_at ASC;

-- name: ListMealLogsByConversation :many
SELECT * FROM meal_logs WHERE conversation_id = $1 ORDER BY recorded_at ASC;

-- name: CountMealLogsByUserAndDate :one
SELECT COUNT(*) FROM meal_logs 
WHERE user_id = $1 AND recorded_at::date = $2;

-- name: DeleteMealLog :exec
DELETE FROM meal_logs WHERE id = $1;

-- name: UpdateMealLog :one
UPDATE meal_logs
SET food_name = $2,
    meal_type = $3,
    macros = $4,
    assumptions = $5,
    raw_response = $6
WHERE id = $1
RETURNING *;

-- name: CreateMealLogs :copyfrom
INSERT INTO meal_logs (
    user_id, conversation_id, food_name, meal_type, recorded_at,
    macros, assumptions, food_source, raw_response
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
