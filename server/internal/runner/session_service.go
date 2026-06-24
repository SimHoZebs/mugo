package runner

import (
	"fmt"
	"log"

	"github.com/simhozebs/mugo/internal/config"
	"google.golang.org/adk/session"
	"google.golang.org/adk/session/database"
	"gorm.io/driver/postgres"
)

func CreateSessionService() (session.Service, error) {
	dbURL := config.GetDatabaseURL()
	if dbURL == "" {
		return nil, fmt.Errorf("no database URL configured for session storage")
	}

	log.Println("Initializing database-backed session service for persistent sessions")

	// Pass the dialector directly — postgres.Open() returns a gorm.Dialector
	// without opening a connection. NewSessionService opens the one it needs internally.
	svc, err := database.NewSessionService(postgres.Open(dbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create session service: %w", err)
	}

	if err := database.AutoMigrate(svc); err != nil {
		return nil, fmt.Errorf("failed to migrate session tables: %w", err)
	}

	log.Println("ADK session tables migrated successfully")
	return svc, nil
}
