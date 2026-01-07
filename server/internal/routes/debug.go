package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/db"
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

// RegisterDebugEndpoints registers debug endpoints.
// Note: These endpoints now proxy to the ADK server for session information.
func RegisterDebugEndpoints(humaAPI huma.API, prefix string, adkClient *adk.Client, database *db.Database) {
	debugGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(
		debugGroup,
		huma.Operation{
			OperationID: "debug_list_sessions",
			Summary:     "List all session IDs for a user",
			Method:      http.MethodGet,
			Path:        "/sessions/{user_id}",
		},
		func(ctx context.Context, input *struct {
			UserId string `path:"user_id" example:"user_12345" doc:"User ID to list sessions for"`
		}) (response *debugListSessionsResponse, err error) {
			resp := &debugListSessionsResponse{}
			if database != nil {
				conversations, err := database.ConversationRepository.ListByUser(ctx, input.UserId)
				if err == nil {
					for _, c := range conversations {
						resp.Body.SessionIds = append(resp.Body.SessionIds, c.SessionID)
					}
					return resp, nil
				}
			}
			resp.Body.SessionIds = []string{fmt.Sprintf("Could not retrieve sessions for user: %s", input.UserId)}
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
			Responses: map[string]*huma.Response{
				"400": {
					Description: "Bad Request - Error retrieving session",
				},
			},
		},
		func(ctx context.Context, input *DebugGetMessagesRequest) (response *debugGetMessagesResponse, err error) {
			appName := config.AgentMapping["nutrition"]

			session, err := adkClient.GetSession(ctx, appName, input.UserId, input.SessionId)
			if err != nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("Error retrieving session: %v", err))
			}

			if session == nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("Session not found: %s", input.SessionId))
			}

			var messages []string
			for _, ev := range session.Events {
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
