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
			input:   `{"name":"Chicken Sandwich","date":"2025-01-07","macros":{"calories":450,"protein":35,"carbs":40,"fat":15},"assumptions":[{"id":"A1","assumed_value":150,"unit":"g"}],"meal_type":"lunch"}`,
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
				assert.Equal(t, "A1", payload.Assumptions[0].ID)
				assert.Equal(t, "g", payload.Assumptions[0].Unit)
			},
		},
		{
			name:    "assigns ID when missing",
			input:   `{"name":"Test","date":"2025-01-07","macros":{"calories":100,"protein":10,"carbs":10,"fat":5},"assumptions":[{"assumed_value":150}]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				require.Len(t, payload.Assumptions, 1)
				assert.Equal(t, "A1", payload.Assumptions[0].ID)
			},
		},
		{
			name:    "assigns default unit when missing",
			input:   `{"name":"Test","date":"2025-01-07","macros":{"calories":100,"protein":10,"carbs":10,"fat":5},"assumptions":[{"assumed_value":150}]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				require.Len(t, payload.Assumptions, 1)
				assert.Equal(t, "g", payload.Assumptions[0].Unit)
			},
		},
		{
			name:    "assigns sequential IDs for multiple assumptions",
			input:   `{"name":"Test","date":"2025-01-07","macros":{"calories":100,"protein":10,"carbs":10,"fat":5},"assumptions":[{"assumed_value":150},{"assumed_value":200}]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				require.Len(t, payload.Assumptions, 2)
				assert.Equal(t, "A1", payload.Assumptions[0].ID)
				assert.Equal(t, "A2", payload.Assumptions[1].ID)
			},
		},
		{
			name:    "preserves existing ID and unit",
			input:   `{"name":"Test","date":"2025-01-07","macros":{"calories":100,"protein":10,"carbs":10,"fat":5},"assumptions":[{"id":"CUSTOM","assumed_value":150,"unit":"ml"}]}`,
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				require.Len(t, payload.Assumptions, 1)
				assert.Equal(t, "CUSTOM", payload.Assumptions[0].ID)
				assert.Equal(t, "ml", payload.Assumptions[0].Unit)
			},
		},
		{
			name:    "strips markdown code blocks",
			input:   "```json\n{\"name\":\"Test\",\"date\":\"2025-01-07\",\"macros\":{\"calories\":100,\"protein\":10,\"carbs\":10,\"fat\":5},\"assumptions\":[]}\n```",
			wantErr: false,
			validate: func(t *testing.T, payload *models.NutritionPayload) {
				assert.Equal(t, "Test", payload.Name)
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
			result, err := NormalizeNutritionResponse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}
