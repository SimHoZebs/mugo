package db

import (
	"context"

	"github.com/simhozebs/mugo/internal/db/repository"
)

// DB interface defines the accessors for all repositories.
type DB interface {
	Users() repository.UserRepository
	Conversations() repository.ConversationRepository
	Meals() repository.MealLogRepository
	Nutrition() repository.NutritionSummaryRepository
	WithTx(ctx context.Context, fn func(ctx context.Context, txDB DB) error) error
}

// DBProvider interface defines how to obtain a DB connection.
type DBProvider interface {
	GetDatabase() (DB, error)
}
