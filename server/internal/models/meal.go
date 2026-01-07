package models

type MealLog struct {
	ID             string       `json:"id"`
	UserID         string       `json:"user_id"`
	ConversationID *string      `json:"conversation_id,omitempty"`
	FoodName       string       `json:"food_name"`
	MealType       string       `json:"meal_type"`
	RecordedAt     string       `json:"recorded_at"`
	Macros         Macros       `json:"macros"`
	Assumptions    []Assumption `json:"assumptions"`
	FoodSource     string       `json:"food_source"`
	RawResponse    interface{}  `json:"raw_response,omitempty"`
	CreatedAt      string       `json:"created_at"`
}
