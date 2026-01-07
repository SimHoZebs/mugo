package models

type DailyNutritionSummary struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Date          string  `json:"date"`
	TotalCalories float64 `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
	MealCount     int     `json:"meal_count"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type WeeklyNutritionSummary struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	WeekStartDate    string  `json:"week_start_date"`
	TotalCalories    float64 `json:"total_calories"`
	TotalProtein     float64 `json:"total_protein"`
	TotalCarbs       float64 `json:"total_carbs"`
	TotalFat         float64 `json:"total_fat"`
	AvgDailyCalories float64 `json:"avg_daily_calories"`
	AvgDailyProtein  float64 `json:"avg_daily_protein"`
	AvgDailyCarbs    float64 `json:"avg_daily_carbs"`
	AvgDailyFat      float64 `json:"avg_daily_fat"`
	MealCount        int     `json:"meal_count"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
