package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/adk/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T, s *TestSuite, username string) string {
	s.API.Delete("/users/by-username/" + username)
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

func setupMockMealRunner() adk.AgentRunner {
	return &mocks.MockAgentRunner{
		RunFunc: func(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error) {
			return &adk.RunResult{
				FinalText: `{
					"meals": [
						{
							"name": "Chicken Salad",
							"meal_type": "lunch",
							"date": "2025-05-20",
							"macros": {"calories": 450, "protein": 35, "carbs": 10, "fat": 25},
							"assumptions": ["standard portions"],
							"confidence": 0.9
						}
					]
				}`,
			}, nil
		},
	}
}

func TestMeals_Create(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	userID := createTestUser(t, s, "meal_create_user")
	s.RegisterMeals(setupMockMealRunner())

	resp := s.API.Post("/meals", struct {
		UserID      string `json:"user_id"`
		Description string `json:"description"`
	}{
		UserID:      userID,
		Description: "I had a chicken salad",
	})

	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Meals []*models.MealLog `json:"meals"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Len(t, body.Meals, 1)
	assert.Equal(t, "Chicken Salad", body.Meals[0].FoodName)
}

func TestMeals_List(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	userID := createTestUser(t, s, "meal_list_user")
	s.RegisterMeals(setupMockMealRunner())

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
	assert.NotEmpty(t, body.Meals)
	assert.Equal(t, userID, body.Meals[0].UserID)
}
