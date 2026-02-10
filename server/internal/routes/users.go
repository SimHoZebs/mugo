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

type UpdateUserRequest struct {
	UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	Body   struct {
		Username string                 `json:"username" example:"johndoe" doc:"Unique username"`
		Metadata map[string]interface{} `json:"metadata,omitempty" doc:"Optional user metadata"`
	}
}

type UpdateUserResponse struct {
	Body struct {
		User *models.User `json:"user"`
	}
}

type DeleteUserRequest struct {
	UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
}

// RegisterUserEndpoints registers user management endpoints.
func RegisterUserEndpoints(humaAPI huma.API, prefix string, lazyDB *db.LazyDatabase) {
	usersGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(usersGroup, huma.Operation{
		OperationID: "create-user",
		Method:      "POST",
		Path:        "",
		Summary:     "Create a new user",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *CreateUserRequest) (*CreateUserResponse, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		exists, err := database.UserRepository.Exists(ctx, input.Body.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check user existence: %w", err)
		}
		if exists {
			return nil, huma.Error409Conflict(fmt.Sprintf("Username '%s' already exists", input.Body.Username))
		}

		user, err := database.UserRepository.Create(ctx, input.Body.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		resp := &CreateUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Register(usersGroup, huma.Operation{
		OperationID: "list-users",
		Method:      "GET",
		Path:        "",
		Summary:     "List all users",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct{}) (*ListUsersResponse, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		users, err := database.UserRepository.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %w", err)
		}

		resp := &ListUsersResponse{}
		resp.Body.Users = users
		return resp, nil
	})

	huma.Register(usersGroup, huma.Operation{
		OperationID: "get-user-by-id",
		Method:      "GET",
		Path:        "/{user_id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	}) (*GetUserResponse, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		user, err := database.UserRepository.GetByID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		resp := &GetUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Register(usersGroup, huma.Operation{
		OperationID: "get-user-by-username",
		Method:      "GET",
		Path:        "/by-username/{username}",
		Summary:     "Get a user by username",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *struct {
		Username string `path:"username" example:"johndoe" doc:"Username"`
	}) (*GetUserResponse, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		user, err := database.UserRepository.GetByUsername(ctx, input.Username)
		if err != nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("User '%s' not found", input.Username))
		}

		resp := &GetUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Register(usersGroup, huma.Operation{
		OperationID: "update-user",
		Method:      "PUT",
		Path:        "/{user_id}",
		Summary:     "Update a user",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *UpdateUserRequest) (*UpdateUserResponse, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		user, err := database.UserRepository.Update(ctx, input.UserID, input.Body.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		resp := &UpdateUserResponse{}
		resp.Body.User = user
		return resp, nil
	})

	huma.Register(usersGroup, huma.Operation{
		OperationID: "delete-user",
		Method:      "DELETE",
		Path:        "/{user_id}",
		Summary:     "Delete a user",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *DeleteUserRequest) (*struct{}, error) {
		database, err := lazyDB.GetDatabase()
		if err != nil {
			return nil, huma.Error503ServiceUnavailable("Database temporarily unavailable", err)
		}

		err = database.UserRepository.Delete(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete user: %w", err)
		}

		return nil, nil
	})
}
