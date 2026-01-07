package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
)

type CreateUserRequest struct {
	Body struct {
		Username string                 `json:"username" example:"johndoe" doc:"Unique username"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Optional user metadata"`
	}
}

type CreateUserResponse struct {
	Body struct {
		User *models.User `json:"user"`
	}
}

type GetUserResponse struct {
	Body struct {
		User *models.User `json:"user"`
	}
}

type ListUsersResponse struct {
	Body struct {
		Users []*models.User `json:"users"`
	}
}

// RegisterUserEndpoints registers user management endpoints.
func RegisterUserEndpoints(humaAPI huma.API, prefix string, database *db.Database) {
	usersGroup := huma.NewGroup(humaAPI, prefix)

	huma.Post(usersGroup, "", func(ctx context.Context, input *CreateUserRequest) (*CreateUserResponse, error) {
		exists, err := database.UserRepository.Exists(ctx, input.Body.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check user existence: %w", err)
		}
		if exists {
			return nil, huma.Error409Conflict(fmt.Sprintf("Username '%s' already exists", input.Body.Username))
		}

		user, err := database.UserRepository.Create(ctx, input.Body.Username, input.Body.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		resp := &CreateUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Get(usersGroup, "", func(ctx context.Context, input *struct{}) (*ListUsersResponse, error) {
		users, err := database.UserRepository.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %w", err)
		}

		resp := &ListUsersResponse{}
		resp.Body.Users = users
		return resp, nil
	})

	huma.Get(usersGroup, "/{user_id}", func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	}) (*GetUserResponse, error) {
		user, err := database.UserRepository.GetByID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		resp := &GetUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Get(usersGroup, "/by-username/{username}", func(ctx context.Context, input *struct {
		Username string `path:"username" example:"johndoe" doc:"Username"`
	}) (*GetUserResponse, error) {
		user, err := database.UserRepository.GetByUsername(ctx, input.Username)
		if err != nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("User '%s' not found", input.Username))
		}

		resp := &GetUserResponse{}
		resp.Body.User = user
		return resp, nil
	})
}
