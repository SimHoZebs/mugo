package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/simhozebs/mugo/internal/config"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
)

type Pool struct {
	*pgxpool.Pool
	*dbgenerated.Queries
}

func NewPool(ctx context.Context) (*Pool, error) {
	databaseURL := config.GetDatabaseURL()
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	pgxConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pgxConfig.MinConns = int32(config.GetDatabaseMinConns())
	pgxConfig.MaxConns = int32(config.GetDatabaseMaxConns())
	pgxConfig.MaxConnLifetime = config.GetDatabaseMaxConnLifetime()
	pgxConfig.MaxConnIdleTime = config.GetDatabaseMaxConnIdleTime()
	pgxConfig.HealthCheckPeriod = config.GetDatabaseHealthCheckPeriod()
	pgxConfig.ConnConfig.ConnectTimeout = config.GetDatabaseConnectTimeout()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
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
