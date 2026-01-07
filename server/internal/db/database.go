package db

import (
	"context"

	"github.com/simhozebs/mugo/internal/db/repository"
)

type Database struct {
	UserRepository         *repository.UserRepository
	ConversationRepository *repository.ConversationRepository
	MealLogRepository      *repository.MealLogRepository
	NutritionRepository    *repository.NutritionSummaryRepository
	pool                   *Pool
}

func NewDatabase(ctx context.Context) (*Database, error) {
	pool, err := NewPool(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{
		UserRepository:         repository.NewUserRepository(pool.Queries),
		ConversationRepository: repository.NewConversationRepository(pool.Queries),
		MealLogRepository:      repository.NewMealLogRepository(pool.Queries),
		NutritionRepository:    repository.NewNutritionSummaryRepository(pool.Queries),
		pool:                   pool,
	}, nil
}

func (d *Database) Close() {
	d.pool.Close()
}
