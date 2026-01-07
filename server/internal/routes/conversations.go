package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/db"
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

// RegisterConversationEndpoints registers conversation endpoints.
func RegisterConversationEndpoints(humaAPI huma.API, prefix string, database *db.Database) {
	conversationsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Get(conversationsGroup, "/{user_id}", func(ctx context.Context, input *struct {
		UserID string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	}) (*ListConversationsResponse, error) {
		conversations, err := database.ConversationRepository.ListByUser(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to list conversations: %w", err)
		}

		resp := &ListConversationsResponse{}
		resp.Body.Conversations = conversations
		return resp, nil
	})

	huma.Get(conversationsGroup, "/{user_id}/session/{session_id}", func(ctx context.Context, input *struct {
		UserID    string `path:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
		SessionID string `path:"session_id" example:"session_12345" doc:"Session ID"`
	}) (*GetConversationResponse, error) {
		conversation, err := database.ConversationRepository.GetBySessionID(ctx, input.UserID, input.SessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}

		resp := &GetConversationResponse{}
		resp.Body.Conversation = conversation
		return resp, nil
	})

	huma.Get(conversationsGroup, "/{conversation_id}", func(ctx context.Context, input *struct {
		ConversationID string `path:"conversation_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Conversation ID"`
	}) (*GetConversationResponse, error) {
		conversation, err := database.ConversationRepository.GetByID(ctx, input.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}

		resp := &GetConversationResponse{}
		resp.Body.Conversation = conversation
		return resp, nil
	})
}
