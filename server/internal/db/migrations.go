package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// runMigrations executes all SQL files in the migrations directory.
// It uses a simple approach: read all files, sort them by name, and execute them.
// Note: For a production system, we should use a proper migration tool that tracks
// which migrations have already been run in a `schema_migrations` table.
func (ld *LazyDatabase) runMigrations(ctx context.Context, db *Database) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)

	return db.WithTx(ctx, func(ctx context.Context, txDB *TxDatabase) error {
		// Ensure the extensions and tables are created
		for _, filename := range filenames {
			log.Printf("Applying migration: %s", filename)
			content, err := migrationsFS.ReadFile("migrations/" + filename)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filename, err)
			}

			// Split the file into "Up" and "Down" sections if present, and only run "Up"
			// This is a simple parser for the +migrate format
			sqlContent := string(content)
			upIndex := strings.Index(sqlContent, "-- +migrate Up")
			downIndex := strings.Index(sqlContent, "-- +migrate Down")

			var toExecute string
			if upIndex != -1 {
				if downIndex != -1 && downIndex > upIndex {
					toExecute = sqlContent[upIndex:downIndex]
				} else {
					toExecute = sqlContent[upIndex:]
				}
			} else {
				// If no +migrate markers, execute the whole file
				toExecute = sqlContent
			}

			// Remove the markers themselves
			toExecute = strings.ReplaceAll(toExecute, "-- +migrate Up", "")
			toExecute = strings.ReplaceAll(toExecute, "-- +migrate StatementBegin", "")
			toExecute = strings.ReplaceAll(toExecute, "-- +migrate StatementEnd", "")

			if strings.TrimSpace(toExecute) == "" {
				continue
			}

			_, err = txDB.tx.Exec(ctx, toExecute)
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", filename, err)
			}
		}
		return nil
	})
}

// Helper to check if a table exists (could be used for more robust migration tracking)
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
