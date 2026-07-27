package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/simhozebs/mugo/internal/config"
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

	// We'll use the filesystem for migrations in the CLI for more flexibility,
	// though we could also use the embedded ones.
	migrationsPath := "file://internal/db/migrations"

	// Get database URL from environment
	rawURL := config.GetDatabaseURL()
	if rawURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Format it specifically for golang-migrate and pgx/v5
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
