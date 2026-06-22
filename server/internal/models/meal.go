package models

import "time"

type MealLog struct {
	ID             string       `json:"id" doc:"Unique meal log ID" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID         string       `json:"user_id" doc:"User ID" example:"user-123"`
	ConversationID *string      `json:"conversation_id,omitempty" doc:"Associated conversation ID" example:"session-456"`
	FoodName       string       `json:"food_name" doc:"Name of the food" example:"Chicken Sandwich"`
	MealType       string       `json:"meal_type" doc:"Type of meal" example:"lunch"`
	RecordedAt     string       `json:"recorded_at" doc:"When the meal was recorded" example:"2025-01-07T12:00:00Z"`
	Macros         Macros       `json:"macros" doc:"Macronutrient breakdown"`
	Assumptions    []Assumption `json:"assumptions" doc:"List of assumptions made"`
	FoodSource     string       `json:"food_source" doc:"Source of the data" example:"ai_estimated"`
	RawResponse    interface{}  `json:"raw_response,omitempty" doc:"Raw data from the source"`
	CreatedAt      string       `json:"created_at" doc:"When the record was created" example:"2025-01-07T12:05:00Z"`
}

type MealLogParams struct {
	FoodName    string
	MealType    string
	RecordedAt  time.Time
	Macros      Macros
	Assumptions []Assumption
	FoodSource  string
	RawResponse interface{}
}
