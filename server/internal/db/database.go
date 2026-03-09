package db

import (
	"context"
	"fmt"

	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/db/repository"
)

type Database struct {
	UserRepository         repository.UserRepository
	ConversationRepository repository.ConversationRepository
	MealLogRepository      repository.MealLogRepository
	NutritionRepository    repository.NutritionSummaryRepository
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

func (d *Database) Users() repository.UserRepository {
	return d.UserRepository
}

func (d *Database) Conversations() repository.ConversationRepository {
	return d.ConversationRepository
}

func (d *Database) Meals() repository.MealLogRepository {
	return d.MealLogRepository
}

func (d *Database) Nutrition() repository.NutritionSummaryRepository {
	return d.NutritionRepository
}

func (d *Database) WithTx(ctx context.Context, fn func(ctx context.Context, db DB) error) error {
	// If already in a transaction, just reuse it
	if _, ok := pgutil.TxFromContext(ctx); ok {
		return fn(ctx, d)
	}

	tx, err := d.pool.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	// Inject transaction into context
	ctxWithTx := pgutil.WithTx(ctx, tx)

	err = fn(ctxWithTx, d)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
