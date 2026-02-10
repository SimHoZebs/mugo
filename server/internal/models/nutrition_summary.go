package models

type DailyNutritionSummary struct {
	ID            string  `json:"id" doc:"Unique summary ID" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        string  `json:"user_id" doc:"User ID" example:"user-123"`
	Date          string  `json:"date" doc:"Date of the summary" example:"2025-01-07"`
	TotalCalories float64 `json:"total_calories" doc:"Total calories consumed" example:"2100.5"`
	TotalProtein  float64 `json:"total_protein" doc:"Total protein in grams" example:"145.2"`
	TotalCarbs    float64 `json:"total_carbs" doc:"Total carbohydrates in grams" example:"210.8"`
	TotalFat      float64 `json:"total_fat" doc:"Total fat in grams" example:"65.4"`
	MealCount     int     `json:"meal_count" doc:"Number of meals recorded" example:"4"`
	CreatedAt     string  `json:"created_at" doc:"When the summary was created" example:"2025-01-07T23:59:59Z"`
	UpdatedAt     string  `json:"updated_at" doc:"When the summary was last updated" example:"2025-01-07T23:59:59Z"`
}

type WeeklyNutritionSummary struct {
	ID               string  `json:"id" doc:"Unique summary ID" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID           string  `json:"user_id" doc:"User ID" example:"user-123"`
	WeekStartDate    string  `json:"week_start_date" doc:"Start date of the week" example:"2025-01-06"`
	TotalCalories    float64 `json:"total_calories" doc:"Total calories for the week" example:"14700.5"`
	TotalProtein     float64 `json:"total_protein" doc:"Total protein for the week" example:"1015.2"`
	TotalCarbs       float64 `json:"total_carbs" doc:"Total carbohydrates for the week" example:"1470.8"`
	TotalFat         float64 `json:"total_fat" doc:"Total fat for the week" example:"455.4"`
	AvgDailyCalories float64 `json:"avg_daily_calories" doc:"Average daily calories" example:"2100.1"`
	AvgDailyProtein  float64 `json:"avg_daily_protein" doc:"Average daily protein" example:"145.0"`
	AvgDailyCarbs    float64 `json:"avg_daily_carbs" doc:"Average daily carbohydrates" example:"210.1"`
	AvgDailyFat      float64 `json:"avg_daily_fat" doc:"Average daily fat" example:"65.1"`
	MealCount        int     `json:"meal_count" doc:"Total number of meals in the week" example:"28"`
	CreatedAt        string  `json:"created_at" doc:"When the summary was created" example:"2025-01-12T23:59:59Z"`
	UpdatedAt        string  `json:"updated_at" doc:"When the summary was last updated" example:"2025-01-12T23:59:59Z"`
}
