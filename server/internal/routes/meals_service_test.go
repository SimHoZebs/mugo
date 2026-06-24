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
	convRepoMock := new(repomocks.ConversationRepositoryMock)

	dbMock.On("Meals").Return(mealRepoMock)
	dbMock.On("Conversations").Return(convRepoMock)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := pgutil.ParseUUID(userID)
	sessionID := "session-456"

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

	var createSessionCalled bool
	mockRun := func(ctx context.Context, uid, sid, text string) (*runner.RunResult, error) {
		return &runner.RunResult{FinalText: string(payloadJSON)}, nil
	}
	mockCreateSession := func(ctx context.Context, uid, sid string) error {
		assert.Equal(t, userID, uid)
		assert.Equal(t, sessionID, sid)
		createSessionCalled = true
		return nil
	}

	result, err := routes.CreateMeal(context.Background(), routes.CreateMealInput{
		UserUUID:    userUUID,
		UserID:      userID,
		SessionID:   sessionID,
		Description: "I ate a chicken sandwich",
	}, mockRun, mockCreateSession, dbMock)

	require.NoError(t, err)
	assert.Equal(t, sessionID, result.SessionID)
	require.Len(t, result.Meals, 1)
	assert.Equal(t, "meal-123", result.Meals[0].ID)
	assert.Equal(t, "Chicken Sandwich", result.Meals[0].FoodName)
	assert.True(t, createSessionCalled, "createSession should be called")

	convRepoMock.AssertExpectations(t)
	mealRepoMock.AssertExpectations(t)
}

func TestCreateMeal_AgentReturnsInvalidJSON(t *testing.T) {
	dbMock := new(mocks.DBMock)
	convRepoMock := new(repomocks.ConversationRepositoryMock)
	dbMock.On("Conversations").Return(convRepoMock)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := pgutil.ParseUUID(userID)

	convRepoMock.On("Create", mock.Anything, userUUID, "", "New Meal Log").
		Return(&models.Conversation{ID: "550e8400-e29b-41d4-a716-446655440002", SessionID: "sess-1"}, nil)

	mockRun := func(ctx context.Context, uid, sid, text string) (*runner.RunResult, error) {
		return &runner.RunResult{FinalText: "not json"}, nil
	}
	mockCreateSession := func(ctx context.Context, uid, sid string) error { return nil }

	_, err := routes.CreateMeal(context.Background(), routes.CreateMealInput{
		UserUUID: userUUID,
		UserID:   userID,
	}, mockRun, mockCreateSession, dbMock)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse nutrition response")
}

func TestCreateMeal_CreateSessionFails(t *testing.T) {
	dbMock := new(mocks.DBMock)
	convRepoMock := new(repomocks.ConversationRepositoryMock)
	dbMock.On("Conversations").Return(convRepoMock)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID, _ := pgutil.ParseUUID(userID)

	convRepoMock.On("Create", mock.Anything, userUUID, "", "New Meal Log").
		Return(&models.Conversation{ID: "550e8400-e29b-41d4-a716-446655440002", SessionID: "sess-1"}, nil)

	mockRun := func(ctx context.Context, uid, sid, text string) (*runner.RunResult, error) {
		t.Fatal("run should not be called when createSession fails")
		return nil, nil
	}
	mockCreateSession := func(ctx context.Context, uid, sid string) error {
		return assert.AnError
	}

	_, err := routes.CreateMeal(context.Background(), routes.CreateMealInput{
		UserUUID: userUUID,
		UserID:   userID,
	}, mockRun, mockCreateSession, dbMock)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create ADK session")
}
