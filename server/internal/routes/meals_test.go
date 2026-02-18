package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/db/mocks"
	repomocks "github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAgentRunner struct {
	adk.AgentRunner
	runFunc func(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error)
}

func (m *mockAgentRunner) Run(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error) {
	if m.runFunc != nil {
		return m.runFunc(ctx, userID, sessionID, text)
	}
	return &adk.RunResult{}, nil
}

func TestCreateMealLog(t *testing.T) {
	_, api := humatest.New(t)

	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	mealRepoMock := new(repomocks.MealLogRepositoryMock)

	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Meals").Return(mealRepoMock)

	userID := "user-123"
	sessionID := "session-456"
	foodDescription := "I ate a chicken sandwich"

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

	mockRunner := &mockAgentRunner{
		runFunc: func(ctx context.Context, uid, sid, text string) (*adk.RunResult, error) {
			return &adk.RunResult{
				FinalText: string(payloadJSON),
			}, nil
		},
	}

	expectedMeal := &models.MealLog{
		ID:       "meal-123",
		UserID:   userID,
		FoodName: payload.Name,
		MealType: string(payload.MealType),
		Macros:   payload.Macros,
	}

	mealRepoMock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedMeal, nil)

	routes.RegisterMealEndpoints(api, "/meals", mockRunner, dbProviderMock)

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
	mealRepoMock.AssertExpectations(t)
}
