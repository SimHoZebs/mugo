package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/db/pgutil"
	"github.com/simhozebs/mugo/internal/models"
)

type ListConversationsResponse struct {
	Body struct {
		Conversations []*models.Conversation `json:"conversations"`
	}
}

type GetConversationResponse struct {
	Body struct {
		Conversation *models.Conversation `json:"conversation"`
	}
}

func RegisterConversationEndpoints(humaAPI huma.API, prefix string, provider db.DBProvider) {
	conversationsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(conversationsGroup, huma.Operation{
		OperationID: "list-conversations",
		Method:      "GET",
		Path:        "/{user_id}",
		Summary:     "List conversations for a user",
		Tags:        []string{"Conversations"},
	}, func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	}) (*ListConversationsResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		conversations, err := database.Conversations().ListByUser(ctx, userUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to list conversations: %w", err)
		}

		resp := &ListConversationsResponse{}
		resp.Body.Conversations = conversations
		return resp, nil
	})

	huma.Register(conversationsGroup, huma.Operation{
		OperationID: "get-conversation-by-session",
		Method:      "GET",
		Path:        "/{user_id}/session/{session_id}",
		Summary:     "Get a conversation by session ID",
		Tags:        []string{"Conversations"},
	}, func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		SessionID string `path:"session_id" example:"session_12345" doc:"Session ID"`
	}) (*GetConversationResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		userUUID, err := pgutil.ParseUUID(input.UserID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user ID", err)
		}

		conversation, err := database.Conversations().GetBySessionID(ctx, userUUID, input.SessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}

		resp := &GetConversationResponse{}
		resp.Body.Conversation = conversation
		return resp, nil
	})

	huma.Register(conversationsGroup, huma.Operation{
		OperationID: "get-conversation-by-id",
		Method:      "GET",
		Path:        "/{conversation_id}",
		Summary:     "Get a conversation by ID",
		Tags:        []string{"Conversations"},
	}, func(ctx context.Context, input *struct {
		ConversationID string `path:"conversation_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Conversation ID"`
	}) (*GetConversationResponse, error) {
		database, err := GetDB(provider)
		if err != nil {
			return nil, err
		}

		convUUID, err := pgutil.ParseUUID(input.ConversationID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid conversation ID", err)
		}

		conversation, err := database.Conversations().GetByID(ctx, convUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}

		resp := &GetConversationResponse{}
		resp.Body.Conversation = conversation
		return resp, nil
	})
}
