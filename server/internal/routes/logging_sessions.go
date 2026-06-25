package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type ListLoggingSessionsResponse struct {
	Body struct {
		LoggingSessions []*models.LoggingSession `json:"logging_sessions"`
	}
}

type GetLoggingSessionResponse struct {
	Body struct {
		LoggingSession *models.LoggingSession `json:"logging_session"`
	}
}

func RegisterLoggingSessionEndpoints(humaAPI huma.API, prefix string, provider db.DBProvider) {
	sessionsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(sessionsGroup, huma.Operation{
		OperationID: "list-logging-sessions",
		Method:      "GET",
		Path:        "/{user_id}",
		Summary:     "List logging sessions for a user",
		Tags:        []string{"Logging Sessions"},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	}) (*ListLoggingSessionsResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		sessions, err := database.LoggingSessions().ListByUser(ctx, userUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to list logging sessions: %w", err)
		}

		resp := &ListLoggingSessionsResponse{}
		resp.Body.LoggingSessions = sessions
		return resp, nil
	})

	huma.Register(sessionsGroup, huma.Operation{
		OperationID: "get-logging-session-by-session-id",
		Method:      "GET",
		Path:        "/{user_id}/session/{session_id}",
		Summary:     "Get a logging session by session ID",
		Tags:        []string{"Logging Sessions"},
	}, func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		SessionID string `path:"session_id" example:"session_12345" doc:"Session ID"`
	}) (*GetLoggingSessionResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		session, err := database.LoggingSessions().GetBySessionID(ctx, userUUID, input.SessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get logging session: %w", err)
		}

		resp := &GetLoggingSessionResponse{}
		resp.Body.LoggingSession = session
		return resp, nil
	})

	huma.Register(sessionsGroup, huma.Operation{
		OperationID: "get-logging-session-by-id",
		Method:      "GET",
		Path:        "/{logging_session_id}",
		Summary:     "Get a logging session by ID",
		Tags:        []string{"Logging Sessions"},
	}, func(ctx context.Context, input *struct {
		LoggingSessionID string `path:"logging_session_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Logging session ID"`
	}) (*GetLoggingSessionResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		sessionUUID, err := pgutil.ParseUUID(input.LoggingSessionID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid logging session ID", err)
		}

		session, err := database.LoggingSessions().GetByID(ctx, sessionUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get logging session: %w", err)
		}

		resp := &GetLoggingSessionResponse{}
		resp.Body.LoggingSession = session
		return resp, nil
	})
}
