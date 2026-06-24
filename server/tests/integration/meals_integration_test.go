package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/simhozebs/mugo/internal/runner"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, s *TestSuite, username string) string {
	username = fmt.Sprintf("%s_%d", username, time.Now().UnixNano())
	resp := s.API.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})
	require.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		User *models.User `json:"user"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	return body.User.ID
}

func setupMockMealRun() runner.RunFunc {
	return func(ctx context.Context, userID, sessionID, text string) (*runner.RunResult, error) {
		return &runner.RunResult{
			FinalText: `{
				"meals": [
					{
						"name": "Chicken Salad",
						"meal_type": "lunch",
						"date": "2025-05-20",
						"macros": {"calories": 450, "protein": 35, "carbs": 10, "fat": 25},
						"assumptions": [{"category":"portion","field":"serving_size","assumed_value":1,"unit":"serving","confidence":"medium","rationale":"standard portions"}],
						"confidence": 0.9
					}
				]
			}`,
		}, nil
	}
}

func setupMockMealCreateSession() runner.CreateSessionFunc {
	return func(ctx context.Context, userID, sessionID string) error { return nil }
}

func TestMeals_Create(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	userID := createTestUser(t, s, "meal_create_user")
	s.RegisterMeals(setupMockMealRun(), setupMockMealCreateSession())

	resp := s.API.Post("/meals", struct {
		UserID      string `json:"user_id"`
		Description string `json:"description"`
	}{
		UserID:      userID,
		Description: "I had a chicken salad",
	})

	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		SessionID string            `json:"session_id"`
		Meals     []*models.MealLog `json:"meals"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.NotEmpty(t, body.SessionID)
	require.Len(t, body.Meals, 1)
	assert.Equal(t, "Chicken Salad", body.Meals[0].FoodName)
}

func TestMeals_List(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	userID := createTestUser(t, s, "meal_list_user")
	s.RegisterMeals(setupMockMealRun(), setupMockMealCreateSession())

	// Create a meal first
	s.API.Post("/meals", struct {
		UserID      string `json:"user_id"`
		Description string `json:"description"`
	}{
		UserID:      userID,
		Description: "I had a chicken salad",
	})

	// Test: List meals
	resp := s.API.Get("/meals/" + userID)
	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Meals []*models.MealLog `json:"meals"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	require.NotEmpty(t, body.Meals)
	assert.Equal(t, userID, body.Meals[0].UserID)
}
