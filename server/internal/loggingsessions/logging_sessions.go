package loggingsessions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

func Create(ctx context.Context, database db.DB, userUUID pgtype.UUID, sessionID, title string) (*models.LoggingSession, pgtype.UUID, error) {
	session, err := database.LoggingSessions().Create(ctx, userUUID, sessionID, title)
	if err != nil {
		return nil, pgtype.UUID{}, fmt.Errorf("failed to create logging session: %w", err)
	}
	sessionUUID, err := pgutil.ParseUUID(session.ID)
	if err != nil {
		return nil, pgtype.UUID{}, fmt.Errorf("failed to parse logging session ID: %w", err)
	}
	return session, sessionUUID, nil
}
