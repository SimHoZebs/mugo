package routes_test

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/simhozebs/mugo/internal/db/mocks"
	repomocks "github.com/simhozebs/mugo/internal/db/repository/mocks"
	"github.com/simhozebs/mugo/internal/models"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Create the test API
	_, api := humatest.New(t)

	// Create mocks
	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	userRepoMock := new(repomocks.UserRepositoryMock)

	// Set up expectations
	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Users").Return(userRepoMock)

	username := "testuser"
	expectedUser := &models.User{
		ID:       "user-123",
		Username: username,
	}

	userRepoMock.On("Exists", mock.Anything, username).Return(false, nil)
	userRepoMock.On("Create", mock.Anything, username).Return(expectedUser, nil)

	// Register the endpoints
	routes.RegisterUserEndpoints(api, "/users", dbProviderMock)

	// Perform the request
	resp := api.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), username)
	assert.Contains(t, resp.Body.String(), "user-123")

	dbProviderMock.AssertExpectations(t)
	dbMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
}

func TestCreateUser_Conflict(t *testing.T) {
	_, api := humatest.New(t)

	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	userRepoMock := new(repomocks.UserRepositoryMock)

	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Users").Return(userRepoMock)

	username := "existinguser"
	userRepoMock.On("Exists", mock.Anything, username).Return(true, nil)

	routes.RegisterUserEndpoints(api, "/users", dbProviderMock)

	resp := api.Post("/users", struct {
		Username string `json:"username"`
	}{
		Username: username,
	})

	assert.Equal(t, http.StatusConflict, resp.Code)
	dbProviderMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	_, api := humatest.New(t)

	dbProviderMock := new(mocks.DBProviderMock)
	dbMock := new(mocks.DBMock)
	userRepoMock := new(repomocks.UserRepositoryMock)

	dbProviderMock.On("GetDatabase").Return(dbMock, nil)
	dbMock.On("Users").Return(userRepoMock)

	userID := "user-123"
	expectedUser := &models.User{
		ID:       userID,
		Username: "testuser",
	}

	userRepoMock.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	routes.RegisterUserEndpoints(api, "/users", dbProviderMock)

	resp := api.Get("/users/" + userID)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), userID)
	dbProviderMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
}
