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

func TestUser_Create(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	username := "test_create_user_" + fmt.Sprintf("%d", time.Now().UnixNano())
	s.API.Delete("/users/by-username/" + username)

	resp := s.API.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})

	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		User *models.User `json:"user"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, username, body.User.Username)
}

func TestUser_Get(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	// Setup: Create a user to get
	username := "test_get_user_" + fmt.Sprintf("%d", time.Now().UnixNano())
	s.API.Delete("/users/by-username/" + username)
	createResp := s.API.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})
	var createBody struct {
		User *models.User `json:"user"`
	}
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	userID := createBody.User.ID

	// Test: Get the user
	resp := s.API.Get("/users/" + userID)
	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		User *models.User `json:"user"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, userID, body.User.ID)
}

func TestUser_Update(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	// Setup: Create a user to update
	username := "test_update_user_" + fmt.Sprintf("%d", time.Now().UnixNano())
	s.API.Delete("/users/by-username/" + username)
	createResp := s.API.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})
	var createBody struct {
		User *models.User `json:"user"`
	}
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	userID := createBody.User.ID

	// Test: Update the user
	newUsername := "updated_user_" + t.Name()
	resp := s.API.Put("/users/"+userID, struct {
		Username string `json:"username"`
	}{
		Username: newUsername,
	})
	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		User *models.User `json:"user"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, newUsername, body.User.Username)
}

func TestUser_Delete(t *testing.T) {
	s := SetupTestSuite(t)
	defer s.Teardown()

	// Setup: Create a user to delete
	username := "test_delete_user_" + fmt.Sprintf("%d", time.Now().UnixNano())
	s.API.Delete("/users/by-username/" + username)
	createResp := s.API.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})
	var createBody struct {
		User *models.User `json:"user"`
	}
	json.Unmarshal(createResp.Body.Bytes(), &createBody)
	userID := createBody.User.ID

	// Test: Delete the user
	resp := s.API.Delete("/users/" + userID)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Verify deletion
	verifyResp := s.API.Get("/users/" + userID)
	assert.NotEqual(t, http.StatusOK, verifyResp.Code)
}
