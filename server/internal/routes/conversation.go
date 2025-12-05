package routes

import (
	"context"
	"fmt"

	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/config"
	"google.golang.org/adk/server/restapi/models"
	"google.golang.org/genai"
)

// ConversationRequest is the request body for conversation endpoint.
type ConversationRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// ConversationResponse is the response body for conversation endpoint.
type ConversationResponse struct {
	Body struct {
		Text string `json:"text"`
	}
}

// ConversationHandler handles conversation requests using the echo agent.
func ConversationHandler(ctx context.Context, adkClient *adk.Client, input *struct {
	Body ConversationRequest `body:""`
}) (*ConversationResponse, error) {
	appName, ok := config.AgentMapping["echo"]
	if !ok {
		return nil, fmt.Errorf("echo agent not configured")
	}

	fmt.Printf("Conversation request: %s (user: %s, session: %s)\n",
		input.Body.Message, input.Body.UserID, input.Body.SessionID)

	result, err := adkClient.RunWithAutoSession(ctx, models.RunAgentRequest{
		AppName:   appName,
		UserId:    input.Body.UserID,
		SessionId: input.Body.SessionID,
		NewMessage: genai.Content{
			Role:  string(genai.RoleUser),
			Parts: []*genai.Part{{Text: input.Body.Message}},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("agent processing failed: %w", err)
	}

	resp := &ConversationResponse{}
	resp.Body.Text = result.FinalText
	return resp, nil
}
