package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/models"
)

type UserRepository struct {
	queries *dbgenerated.Queries
}

func NewUserRepository(queries *dbgenerated.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, username string, metadata map[string]interface{}) (*models.User, error) {
	metadataJSON, _ := json.Marshal(metadata)
	arg := dbgenerated.CreateUserParams{
		Username: username,
		Metadata: metadataJSON,
	}
	result, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return mapToUser(result), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	pgUUID := pgtype.UUID{
		Bytes: [16]byte(parsedUUID),
		Valid: true,
	}
	result, err := r.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return mapToUser(result), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	result, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return mapToUser(result), nil
}

func (r *UserRepository) Exists(ctx context.Context, username string) (bool, error) {
	return r.queries.UserExists(ctx, username)
}

func (r *UserRepository) List(ctx context.Context) ([]*models.User, error) {
	results, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	users := make([]*models.User, len(results))
	for i, u := range results {
		users[i] = mapToUser(u)
	}
	return users, nil
}

func mapToUser(u dbgenerated.User) *models.User {
	var metadata map[string]interface{}
	if u.Metadata != nil {
		json.Unmarshal(u.Metadata, &metadata)
	}
	return &models.User{
		ID:        u.ID.String(),
		Username:  u.Username,
		Metadata:  metadata,
		CreatedAt: u.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}
