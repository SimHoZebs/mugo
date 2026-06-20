package routes

import (
	"context"
	"fmt"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"google.golang.org/adk/session"
)

type DebugGetMessagesRequest struct {
	UserId    string `path:"user_id" example:"user_12345" doc:"User ID associated with the session"`
	SessionId string `path:"session_id" example:"session_12345" doc:"Session ID to retrieve messages from"`
}

type debugGetMessagesResponse struct {
	Body struct {
		Messages []string `json:"messages"`
	}
}

type debugListSessionsResponse struct {
	Body struct {
		SessionIds []string `json:"session_ids"`
	}
}

func RegisterDebugEndpoints(humaAPI huma.API, prefix string, sessionService session.Service, provider db.DBProvider) {
	debugGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(
		debugGroup,
		huma.Operation{
			OperationID: "debug_list_sessions",
			Summary:     "List all session IDs for a user",
			Method:      http.MethodGet,
			Path:        "/sessions/{user_id}",
			Tags:        []string{"Debug"},
		},
		func(ctx context.Context, input *struct {
			UserId string `path:"user_id" example:"user_12345" doc:"User ID to list sessions for"`
		}) (response *debugListSessionsResponse, err error) {
			database, err := GetDB(provider)
			if err != nil {
				return nil, huma.Error500InternalServerError(fmt.Sprintf("Database unavailable: %v", err))
			}

			userUUID, err := pgutil.ParseUUID(input.UserId)
			if err != nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("Invalid user ID: %v", err))
			}

			conversations, err := database.Conversations().ListByUser(ctx, userUUID)
			if err != nil {
				return nil, huma.Error500InternalServerError(fmt.Sprintf("Could not retrieve sessions for user: %s", input.UserId))
			}

			resp := &debugListSessionsResponse{}
			for _, c := range conversations {
				resp.Body.SessionIds = append(resp.Body.SessionIds, c.SessionID)
			}
			return resp, nil
		},
	)

	huma.Register(
		debugGroup,
		huma.Operation{
			OperationID: "debug_get_messages",
			Method:      http.MethodGet,
			Path:        "/messages/{user_id}/{session_id}",
			Summary:     "Retrieve messages from a user session",
			Tags:        []string{"Debug"},
			Responses: map[string]*huma.Response{
				"400": {
					Description: "Bad Request - Error retrieving session",
				},
			},
		},
		func(ctx context.Context, input *DebugGetMessagesRequest) (response *debugGetMessagesResponse, err error) {
			if sessionService == nil {
				return nil, huma.Error500InternalServerError("Session service not available")
			}

			sess, err := sessionService.Get(ctx, &session.GetRequest{
				AppName:   "macro_estimator",
				UserID:    input.UserId,
				SessionID: input.SessionId,
			})
			if err != nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("Error retrieving session: %v", err))
			}

			if sess == nil || sess.Session == nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("Session not found: %s", input.SessionId))
			}

			var messages []string
			for ev := range sess.Session.Events().All() {
				if ev.Content == nil {
					continue
				}
				for _, p := range ev.Content.Parts {
					if p != nil && p.Text != "" {
						messages = append(messages, p.Text)
					}
				}
			}

			resp := &debugGetMessagesResponse{}
			resp.Body.Messages = messages
			return resp, nil
		},
	)

}
