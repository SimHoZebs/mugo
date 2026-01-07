-- name: UpsertDailyNutritionSummary :one
INSERT INTO daily_nutrition_summaries (
    user_id, date, total_calories, total_protein, total_carbs, total_fat, meal_count
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id, date) 
DO UPDATE SET 
    total_calories = EXCLUDED.total_calories,
    total_protein = EXCLUDED.total_protein,
    total_carbs = EXCLUDED.total_carbs,
    total_fat = EXCLUDED.total_fat,
    meal_count = EXCLUDED.meal_count,
    updated_at = NOW()
RETURNING *;

-- name: GetDailyNutritionSummary :one
SELECT * FROM daily_nutrition_summaries WHERE user_id = $1 AND date = $2;

-- name: ListDailyNutritionSummariesByUser :many
SELECT * FROM daily_nutrition_summaries 
WHERE user_id = $1 
ORDER BY date DESC
LIMIT $2 OFFSET $3;

-- name: ListDailyNutritionSummariesByUserAndDateRange :many
SELECT * FROM daily_nutrition_summaries 
WHERE user_id = $1 
AND date >= $2 
AND date <= $3
ORDER BY date ASC;
