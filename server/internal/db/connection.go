package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
)

type Pool struct {
	*pgxpool.Pool
	*dbgenerated.Queries
}

func NewPool(ctx context.Context) (*Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MinConns = 5
	config.MaxConns = 25

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{
		Pool:    pool,
		Queries: dbgenerated.New(pool),
	}, nil
}

func (p *Pool) Close() {
	p.Pool.Close()
}
