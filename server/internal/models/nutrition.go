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
	Name        string       `json:"name" doc:"Name of the food item" example:"Chicken Sandwich"`
	MealType    MealType     `json:"meal_type" doc:"Type of meal" example:"lunch"`
	Macros      Macros       `json:"macros" doc:"Macronutrient breakdown"`
	Assumptions []Assumption `json:"assumptions" doc:"List of assumptions made during analysis"`
}

// Assumption represents an assumption made during nutritional analysis.
type Assumption struct {
	ID           string  `json:"id,omitempty" doc:"Unique identifier for the assumption"`
	Category     string  `json:"category,omitempty" doc:"Category of assumption (e.g., portion, ingredient)" example:"portion"`
	Field        string  `json:"field,omitempty" doc:"The field being assumed (e.g., weight, quantity)" example:"weight"`
	AssumedValue float64 `json:"assumed_value" doc:"The value used for the assumption" example:"150"`
	Unit         string  `json:"unit,omitempty" doc:"Unit of measurement" example:"g"`
	Confidence   string  `json:"confidence,omitempty" doc:"Confidence level of the assumption" example:"high"`
	Rationale    string  `json:"rationale,omitempty" doc:"Reasoning behind the assumption" example:"Standard chicken breast weight"`
}

// Macros represents the macronutrient values.
type Macros struct {
	Calories float64 `json:"calories" doc:"Total calories in kcal" example:"450"`
	Protein  float64 `json:"protein" doc:"Protein in grams" example:"35"`
	Carbs    float64 `json:"carbs" doc:"Carbohydrates in grams" example:"40"`
	Fat      float64 `json:"fat" doc:"Fat in grams" example:"15"`
}
