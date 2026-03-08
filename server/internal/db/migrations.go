package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// runMigrations executes all SQL files in the migrations directory using golang-migrate.
// It uses a schema_migrations table to track applied migrations.
func (ld *LazyDatabase) runMigrations(ctx context.Context, db *Database) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	// Correctly format the database URL for the pgx5 driver.
	// The driver expects pgx5://user:pass@host:port/db
	databaseURL := ToMigrateURL(db.pool.Config().ConnConfig.ConnString())

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Check if we need to force version 1 (migration from legacy system)
	// If the users table exists but schema_migrations doesn't, we force 1.
	var exists bool
	err = db.pool.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')").Scan(&exists)
	if err == nil && exists {
		var schemaExists bool
		err = db.pool.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'schema_migrations')").Scan(&schemaExists)
		if err == nil && !schemaExists {
			log.Println("Legacy migration detected: 'users' table exists but 'schema_migrations' doesn't. Forcing version 1.")
			if err := m.Force(1); err != nil {
				return fmt.Errorf("failed to force migration version 1: %w", err)
			}
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}

// Helper to check if a table exists (preserved for potential future use)
func tableExists(ctx context.Context, tx pgx.Tx, tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = $1
	)`
	err := tx.QueryRow(ctx, query, tableName).Scan(&exists)
	return exists, err
}
