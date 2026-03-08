package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/simhozebs/mugo/internal/db"
)

func main() {
	forceCmd := flag.NewFlagSet("force", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'up', 'down' or 'force' subcommands")
		os.Exit(1)
	}

	ctx := context.Background()
	lazyDB := db.NewLazyDatabase(ctx)
	defer lazyDB.Close()

	// We need to trigger the connection to get the URL
	_, err := lazyDB.GetDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// This is a bit hacky but we need the raw URL in pgx5 format
	// The internal/db package doesn't expose the formatting logic, so we replicate it or
	// we could export a helper. For now, let's replicate to keep internal/db clean.

	// We'll use the filesystem for migrations in the CLI for more flexibility,
	// though we could also use the embedded ones.
	migrationsPath := "file://internal/db/migrations"

	// Get database URL from environment or fallback to formatted one from pool
	// Note: In production/actual use, we'd probably want to pass the URL directly.
	// But since this runs in the same environment as the app, we can use the same logic.

	// Extract raw URL from the database instance if possible, or just look for DATABASE_URL
	rawURL := os.Getenv("DATABASE_URL")
	if rawURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	databaseURL := db.ToMigrateURL(rawURL)

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch os.Args[1] {
	case "up":
		upCmd.Parse(os.Args[2:])
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to run up migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		downCmd.Parse(os.Args[2:])
		args := downCmd.Args()
		if len(args) == 0 {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("failed to run down migrations: %v", err)
			}
		} else {
			steps, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("invalid steps: %v", err)
			}
			if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("failed to run down migrations: %v", err)
			}
		}
		fmt.Println("Down migrations applied successfully")

	case "force":
		forceCmd.Parse(os.Args[2:])
		args := forceCmd.Args()
		if len(args) != 1 {
			log.Fatal("expected exactly one version argument for force")
		}
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("failed to force migration version: %v", err)
		}
		fmt.Println("Migration version forced successfully")

	default:
		fmt.Println("expected 'up', 'down' or 'force' subcommands")
		os.Exit(1)
	}
}
