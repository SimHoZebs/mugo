package agents

import (
	"testing"

	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeMealsBatchResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantNil  bool
		validate func(t *testing.T, batch *models.MealsBatchPayload)
	}{
		{
			name:    "valid single meal batch",
			input:   `{"meals":[{"name":"Chicken Sandwich","date":"2025-01-07","macros":{"calories":450,"protein":35,"carbs":40,"fat":15},"assumptions":[]}]}`,
			wantErr: false,
			wantNil: false,
			validate: func(t *testing.T, batch *models.MealsBatchPayload) {
				require.Len(t, batch.Meals, 1)
				assert.Equal(t, "Chicken Sandwich", batch.Meals[0].Name)
			},
		},
		{
			name:    "valid multiple meals batch",
			input:   `{"meals":[{"name":"Breakfast","date":"2025-01-07","macros":{"calories":300,"protein":20,"carbs":30,"fat":10},"assumptions":[]},{"name":"Lunch","date":"2025-01-07","macros":{"calories":500,"protein":30,"carbs":50,"fat":20},"assumptions":[]}]}`,
			wantErr: false,
			wantNil: false,
			validate: func(t *testing.T, batch *models.MealsBatchPayload) {
				require.Len(t, batch.Meals, 2)
				assert.Equal(t, "Breakfast", batch.Meals[0].Name)
				assert.Equal(t, "Lunch", batch.Meals[1].Name)
			},
		},
		{
			name:    "strips markdown code blocks",
			input:   "```json\n{\"meals\":[{\"name\":\"Test\",\"date\":\"2025-01-07\",\"macros\":{\"calories\":100,\"protein\":10,\"carbs\":10,\"fat\":5},\"assumptions\":[]}]}\n```",
			wantErr: false,
			wantNil: false,
			validate: func(t *testing.T, batch *models.MealsBatchPayload) {
				require.Len(t, batch.Meals, 1)
				assert.Equal(t, "Test", batch.Meals[0].Name)
			},
		},
		{
			name:    "returns nil for non-JSON text",
			input:   "This is just plain text, not JSON",
			wantErr: false,
			wantNil: true,
		},
		{
			name:    "returns nil for text starting with bracket",
			input:   "[1, 2, 3]",
			wantErr: false,
			wantNil: true,
		},
		{
			name:    "returns nil for empty string",
			input:   "",
			wantErr: false,
			wantNil: true,
		},
		{
			name:    "returns nil for whitespace only",
			input:   "   \n\t  ",
			wantErr: false,
			wantNil: true,
		},
		{
			name:    "invalid JSON returns error",
			input:   `{"meals":[{"name":}]}`,
			wantErr: true,
		},
		{
			name:    "valid JSON but wrong schema returns error",
			input:   `{"wrong_field":"value"}`,
			wantErr: false,
			wantNil: false,
			validate: func(t *testing.T, batch *models.MealsBatchPayload) {
				assert.Empty(t, batch.Meals)
			},
		},
		{
			name:    "empty meals array",
			input:   `{"meals":[]}`,
			wantErr: false,
			wantNil: false,
			validate: func(t *testing.T, batch *models.MealsBatchPayload) {
				assert.Empty(t, batch.Meals)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeMealsBatchResponse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, result)
			} else if !tt.wantErr {
				require.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}
