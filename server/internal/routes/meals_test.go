package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/adk"
	adkmocks "github.com/simhozebs/mugo/internal/adk/mocks"
	"github.com/simhozebs/mugo/internal/db/mocks"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	repomocks "github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateMealLog(t *testing.T) {
	_, api := humatest.New(t)

	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	mealRepoMock := new(repomocks.MealLogRepositoryMock)
	convRepoMock := new(repomocks.ConversationRepositoryMock)

	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Meals").Return(mealRepoMock)
	dbMock.On("Conversations").Return(convRepoMock)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := pgutil.ParseUUID(userID)
	sessionID := "session-456"
	foodDescription := "I ate a chicken sandwich"

	expectedConv := &models.Conversation{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		UserID:    userID,
		SessionID: sessionID,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	convRepoMock.On("Create", mock.Anything, userUUID, sessionID, "New Meal Log").
		Return(expectedConv, nil)

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

	var createSessionCalled bool
	mockRunner := &adkmocks.MockAgentRunner{
		CreateSessionFunc: func(ctx context.Context, uid, sid string) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, sessionID, sid)
			createSessionCalled = true
			return nil
		},
		RunFunc: func(ctx context.Context, uid, sid, text string) (*adk.RunResult, error) {
			return &adk.RunResult{
				FinalText: string(payloadJSON),
			}, nil
		},
		NameFunc: func() string { return "meal_orchestrator" },
	}

	expectedMeal := &models.MealLog{
		ID:       "meal-123",
		UserID:   userID,
		FoodName: payload.Name,
		MealType: string(payload.MealType),
		Macros:   payload.Macros,
	}

	mealRepoMock.On("CreateBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(([]*models.MealLog)(nil), nil)
	mealRepoMock.On("ListByConversation", mock.Anything, mock.Anything).Return([]*models.MealLog{expectedMeal}, nil)
	mealRepoMock.On("ListByConversation", mock.Anything, mock.Anything).
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
	assert.True(t, createSessionCalled, "CreateSession should be called for new conversations")

	var body struct {
		SessionID string            `json:"session_id"`
		Meals     []*models.MealLog `json:"meals"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, sessionID, body.SessionID)
	require.Len(t, body.Meals, 1)
	assert.Equal(t, "meal-123", body.Meals[0].ID)
	assert.Equal(t, "Chicken Sandwich", body.Meals[0].FoodName)

	dbProviderMock.AssertExpectations(t)
	convRepoMock.AssertExpectations(t)
	mealRepoMock.AssertExpectations(t)
}
