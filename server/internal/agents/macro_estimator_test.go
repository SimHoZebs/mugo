package agents

import (
	"testing"

	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeNutritionResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, payload *models.NutritionPayload)
	}{
		{
			name:    "valid response with all fields",
			input:   `{"name":"Chicken Sandwich","date":"2025-01-07","macros":{"calories":450,"protein":35,"carbs":40,"fat":15},"assumptions":[{"assumed_value":150}],"meal_type":"lunch"}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				assert.Equal(t, "Chicken Sandwich", payload.Name)
				assert.Equal(t, "2025-01-07", payload.Date)
				assert.Equal(t, models.MealTypeLunch, payload.MealType)
				assert.Equal(t, 450.0, payload.Macros.Calories)
				assert.Equal(t, 35.0, payload.Macros.Protein)
				assert.Equal(t, 40.0, payload.Macros.Carbs)
				assert.Equal(t, 15.0, payload.Macros.Fat)
				require.Len(t, payload.Assumptions, 1)
				assert.Equal(t, 150.0, payload.Assumptions[0].AssumedValue)
			},
		},
		{
			name:    "strips markdown code blocks",
			input:   "```json\n{\"name\":\"Test\",\"date\":\"2025-01-07\",\"macros\":{\"calories\":100,\"protein\":10,\"carbs\":10,\"fat\":5},\"assumptions\":[],\"meal_type\":\"brunch\"}\n```",
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				assert.Equal(t, "Test", payload.Name)
				assert.Equal(t, models.MealTypeBrunch, payload.MealType)
			},
		},
		{
			name:    "invalid JSON returns error",
			input:   `{"name":}`,
			wantErr: true,
		},
		{
			name:    "empty string returns error",
			input:   "",
			wantErr: true,
		},
		{
			name:    "JSON without required fields - macros missing",
			input:   `{"name":"Test","date":"2025-01-07","assumptions":[]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				assert.Equal(t, "Test", payload.Name)
				assert.Equal(t, 0.0, payload.Macros.Calories)
			},
		},
		{
			name:    "empty assumptions array",
			input:   `{"name":"Test","date":"2025-01-07","macros":{"calories":100,"protein":10,"carbs":10,"fat":5},"assumptions":[]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				assert.Empty(t, payload.Assumptions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanText, result, err := NormalizeNutritionResponse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, cleanText)
				assert.NotContains(t, cleanText, "```")
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}
