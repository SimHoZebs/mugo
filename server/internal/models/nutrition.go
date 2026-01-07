package models

type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeDinner    MealType = "dinner"
	MealTypeSnack     MealType = "snack"
	MealTypeUnknown   MealType = "unknown"
)

// NutritionPayload is the structured response from the nutrition agent.
type NutritionPayload struct {
	Name        string       `json:"name"`
	MealType    MealType     `json:"meal_type"`
	Macros      Macros       `json:"macros"`
	Assumptions []Assumption `json:"assumptions"`
}

// Assumption represents an assumption made during nutritional analysis.
type Assumption struct {
	ID           string  `json:"id,omitempty"`
	Category     string  `json:"category,omitempty"`
	Field        string  `json:"field,omitempty"`
	AssumedValue float64 `json:"assumed_value"`
	Unit         string  `json:"unit,omitempty"`
	Confidence   string  `json:"confidence,omitempty"`
	Rationale    string  `json:"rationale,omitempty"`
}

// Macros represents the macronutrient values.
type Macros struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}
