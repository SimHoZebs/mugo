package routes_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/adk"
	adkmocks "github.com/simhozebs/mugo/internal/adk/mocks"
	"github.com/simhozebs/mugo/internal/db/mocks"
	repomocks "github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMealLog(t *testing.T) {
	_, api := humatest.New(t)

	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	mealRepoMock := new(repomocks.MealLogRepositoryMock)
	adkClientMock := new(adkmocks.AgentClientMock)

	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Meals").Return(mealRepoMock)

	userID := "user-123"
	sessionID := "session-456"
	foodDescription := "I ate a chicken sandwich"

	// Mock ADK response
	payload := models.NutritionPayload{
		Name:     "Chicken Sandwich",
		MealType: models.MealTypeLunch,
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
	payloadJSON, _ := json.Marshal(payload)

	adkResult := &adk.RunResult{
		FinalText: string(payloadJSON),
	}

	adkClientMock.On("RunWithAutoSession", mock.Anything, mock.MatchedBy(func(req interface{}) bool {
		// You can add more specific matching here if needed
		return true
	})).Return(adkResult, nil)

	// Mock DB Create
	expectedMeal := &models.MealLog{
		ID:       "meal-123",
		UserID:   userID,
		FoodName: payload.Name,
		MealType: string(payload.MealType),
		Macros:   payload.Macros,
	}

	mealRepoMock.On("Create", mock.Anything, userID, sessionID, payload.Name, string(payload.MealType), mock.Anything, payload.Macros, payload.Assumptions, "ai_estimated", payload).
		Return(expectedMeal, nil)

	routes.RegisterMealEndpoints(api, "/meals", adkClientMock, dbProviderMock)

	resp := api.Post("/meals", struct {
		UserID      string `json:"user_id"`
		SessionID   string `json:"session_id"`
		Description string `json:"description"`
	}{
		UserID:      userID,
		SessionID:   sessionID,
		Description: foodDescription,
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "meal-123")
	assert.Contains(t, resp.Body.String(), "Chicken Sandwich")

	dbProviderMock.AssertExpectations(t)
	adkClientMock.AssertExpectations(t)
	mealRepoMock.AssertExpectations(t)
}
