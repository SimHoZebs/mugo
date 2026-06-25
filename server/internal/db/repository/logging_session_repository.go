package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	dbgenerated "github.com/simhozebs/mugo/internal/db/dbgenerated"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type loggingSessionRepository struct {
	queries *dbgenerated.Queries
}

func (repo *loggingSessionRepository) q(ctx context.Context) *dbgenerated.Queries {
	if tx, ok := pgutil.TxFromContext(ctx); ok {
		return repo.queries.WithTx(tx)
	}
	return repo.queries
}

func NewLoggingSessionRepository(queries *dbgenerated.Queries) LoggingSessionRepository {
	return &loggingSessionRepository{queries: queries}
}

func (repo *loggingSessionRepository) Create(ctx context.Context, userID pgtype.UUID, sessionID, title string) (*models.LoggingSession, error) {
	if sessionID == "" {
		sessionID = pgutil.GenerateUUID()
	}

	arg := dbgenerated.CreateLoggingSessionParams{
		UserID:    userID,
		SessionID: sessionID,
		Title:     pgutil.Text(title),
	}
	result, err := repo.q(ctx).CreateLoggingSession(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create logging session: %w", err)
	}
	return mapToLoggingSession(result), nil
}

func (repo *loggingSessionRepository) GetByID(ctx context.Context, id pgtype.UUID) (*models.LoggingSession, error) {
	result, err := repo.q(ctx).GetLoggingSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get logging session: %w", err)
	}
	return mapToLoggingSession(result), nil
}

func (repo *loggingSessionRepository) GetBySessionID(ctx context.Context, userID pgtype.UUID, sessionID string) (*models.LoggingSession, error) {
	arg := dbgenerated.GetLoggingSessionBySessionIDParams{
		UserID:    userID,
		SessionID: sessionID,
	}
	result, err := repo.q(ctx).GetLoggingSessionBySessionID(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get logging session by session ID: %w", err)
	}
	return mapToLoggingSession(result), nil
}

func (repo *loggingSessionRepository) ListByUser(ctx context.Context, userID pgtype.UUID) ([]*models.LoggingSession, error) {
	results, err := repo.q(ctx).ListLoggingSessionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list logging sessions: %w", err)
	}
	sessions := make([]*models.LoggingSession, len(results))
	for i, c := range results {
		sessions[i] = mapToLoggingSession(c)
	}
	return sessions, nil
}

func (repo *loggingSessionRepository) UpdateTitle(ctx context.Context, id pgtype.UUID, title string) (*models.LoggingSession, error) {
	arg := dbgenerated.UpdateLoggingSessionTitleParams{
		ID:    id,
		Title: pgutil.Text(title),
	}
	result, err := repo.q(ctx).UpdateLoggingSessionTitle(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to update logging session title: %w", err)
	}
	return mapToLoggingSession(result), nil
}

func (repo *loggingSessionRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return repo.q(ctx).DeleteLoggingSession(ctx, id)
}

func mapToLoggingSession(c dbgenerated.LoggingSession) *models.LoggingSession {
	return &models.LoggingSession{
		ID:        c.ID.String(),
		UserID:    c.UserID.String(),
		SessionID: c.SessionID,
		Title:     pgutil.FromText(c.Title),
		CreatedAt: c.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Time.Format(time.RFC3339),
	}
}
