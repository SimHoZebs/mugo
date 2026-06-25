package routes_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/simhozebs/mugo/internal/db/mocks"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	repomocks "github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/simhozebs/mugo/internal/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateMeal(t *testing.T) {
	dbMock := new(mocks.DBMock)
	mealRepoMock := new(repomocks.MealLogRepositoryMock)
	dbMock.On("Meals").Return(mealRepoMock)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := pgutil.ParseUUID(userID)
	convUUID, _ := pgutil.ParseUUID("550e8400-e29b-41d4-a716-446655440001")

	payload := models.NutritionPayload{
		Name:     "Chicken Sandwich",
		MealType: models.MealTypeLunch,
		Date:     time.Now().Format("2006-01-02"),
		Macros: models.Macros{
			Calories: 450,
			Protein:  35,
			Carbs:    40,
			Fat:      15,
		},
		Assumptions: []models.Assumption{
			{Category: "portion", Field: "weight", AssumedValue: 150, Unit: "g"},
		},
	}
	batch := models.MealsBatchPayload{Meals: []models.NutritionPayload{payload}}
	payloadJSON, _ := json.Marshal(batch)

	expectedMeal := &models.MealLog{
		ID:       "meal-123",
		UserID:   userID,
		FoodName: payload.Name,
		MealType: string(payload.MealType),
		Macros:   payload.Macros,
	}
	mealRepoMock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedMeal, nil)

	runResult := &runner.RunResult{FinalText: string(payloadJSON)}

	meals, err := routes.CreateMeal(context.Background(), userUUID, convUUID, runResult, dbMock)

	require.NoError(t, err)
	require.Len(t, meals, 1)
	assert.Equal(t, "meal-123", meals[0].ID)
	assert.Equal(t, "Chicken Sandwich", meals[0].FoodName)

	mealRepoMock.AssertExpectations(t)
}

func TestCreateMeal_InvalidJSON(t *testing.T) {
	dbMock := new(mocks.DBMock)

	userUUID, _ := pgutil.ParseUUID("550e8400-e29b-41d4-a716-446655440000")
	convUUID, _ := pgutil.ParseUUID("550e8400-e29b-41d4-a716-446655440001")

	runResult := &runner.RunResult{FinalText: "not json"}

	_, err := routes.CreateMeal(context.Background(), userUUID, convUUID, runResult, dbMock)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse nutrition response")
}

func TestCreateMeal_EmptyMeals(t *testing.T) {
	dbMock := new(mocks.DBMock)

	userUUID, _ := pgutil.ParseUUID("550e8400-e29b-41d4-a716-446655440000")
	convUUID, _ := pgutil.ParseUUID("550e8400-e29b-41d4-a716-446655440001")

	runResult := &runner.RunResult{FinalText: `{"meals":[]}`}

	meals, err := routes.CreateMeal(context.Background(), userUUID, convUUID, runResult, dbMock)

	require.NoError(t, err)
	assert.Empty(t, meals)
}
