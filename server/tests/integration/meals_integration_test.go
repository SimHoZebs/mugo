package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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

func TestMeals_List(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	userID := createTestUser(t, s, "meal_list_user")

	// Test: List meals for a user with no meals
	resp := s.API.Get("/meals/" + userID)
	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Meals []*models.MealLog `json:"meals"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Empty(t, body.Meals)
}
