-- name: UpsertWeeklyNutritionSummary :one
INSERT INTO weekly_nutrition_summaries (
    user_id, week_start_date, total_calories, total_protein, total_carbs, total_fat,
    avg_daily_calories, avg_daily_protein, avg_daily_carbs, avg_daily_fat, meal_count
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (user_id, week_start_date) 
DO UPDATE SET 
    total_calories = EXCLUDED.total_calories,
    total_protein = EXCLUDED.total_protein,
    total_carbs = EXCLUDED.total_carbs,
    total_fat = EXCLUDED.total_fat,
    avg_daily_calories = EXCLUDED.avg_daily_calories,
    avg_daily_protein = EXCLUDED.avg_daily_protein,
    avg_daily_carbs = EXCLUDED.avg_daily_carbs,
    avg_daily_fat = EXCLUDED.avg_daily_fat,
    meal_count = EXCLUDED.meal_count,
    updated_at = NOW()
RETURNING *;

-- name: GetWeeklyNutritionSummary :one
SELECT * FROM weekly_nutrition_summaries WHERE user_id = $1 AND week_start_date = $2;

-- name: ListWeeklyNutritionSummariesByUser :many
SELECT * FROM weekly_nutrition_summaries 
WHERE user_id = $1 
ORDER BY week_start_date DESC
LIMIT $2 OFFSET $3;

-- name: ListWeeklyNutritionSummariesByUserAndDateRange :many
SELECT * FROM weekly_nutrition_summaries 
WHERE user_id = $1 
AND week_start_date >= $2 
AND week_start_date <= $3
ORDER BY week_start_date ASC;
