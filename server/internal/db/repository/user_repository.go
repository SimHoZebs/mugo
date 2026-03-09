package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type userRepository struct {
	queries *dbgenerated.Queries
}

func (repo *userRepository) q(ctx context.Context) *dbgenerated.Queries {
	if tx, ok := pgutil.TxFromContext(ctx); ok {
		return repo.queries.WithTx(tx)
	}
	return repo.queries
}

func NewUserRepository(queries *dbgenerated.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (repo *userRepository) Create(ctx context.Context, username string) (*models.User, error) {
	result, err := repo.q(ctx).CreateUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return mapToUser(result), nil
}

func (repo *userRepository) GetByID(ctx context.Context, id pgtype.UUID) (*models.User, error) {
	result, err := repo.q(ctx).GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return mapToUser(result), nil
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	result, err := repo.q(ctx).GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return mapToUser(result), nil
}

func (repo *userRepository) Exists(ctx context.Context, username string) (bool, error) {
	return repo.q(ctx).UserExists(ctx, username)
}

func (repo *userRepository) List(ctx context.Context) ([]*models.User, error) {
	results, err := repo.q(ctx).ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	users := make([]*models.User, len(results))
	for i, u := range results {
		users[i] = mapToUser(u)
	}
	return users, nil
}

func (repo *userRepository) Update(ctx context.Context, id pgtype.UUID, username string) (*models.User, error) {
	arg := dbgenerated.UpdateUserParams{
		ID:       id,
		Username: username,
	}

	result, err := repo.q(ctx).UpdateUser(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return mapToUser(result), nil
}

func (repo *userRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	err := repo.q(ctx).DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func mapToUser(u dbgenerated.User) *models.User {
	return &models.User{
		ID:        u.ID.String(),
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}
