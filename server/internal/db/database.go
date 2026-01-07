package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

type TxDatabase struct {
	UserRepository         *repository.UserRepository
	ConversationRepository *repository.ConversationRepository
	MealLogRepository      *repository.MealLogRepository
	NutritionRepository    *repository.NutritionSummaryRepository
	tx                     pgx.Tx
}

func (d *Database) WithTx(ctx context.Context, fn func(ctx context.Context, txDB *TxDatabase) error) error {
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

	txDB := &TxDatabase{
		UserRepository:         repository.NewUserRepository(d.pool.Queries.WithTx(tx)),
		ConversationRepository: repository.NewConversationRepository(d.pool.Queries.WithTx(tx)),
		MealLogRepository:      repository.NewMealLogRepository(d.pool.Queries.WithTx(tx)),
		NutritionRepository:    repository.NewNutritionSummaryRepository(d.pool.Queries.WithTx(tx)),
		tx:                     tx,
	}

	err = fn(ctx, txDB)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
