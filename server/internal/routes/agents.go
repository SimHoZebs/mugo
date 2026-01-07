package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/api"
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/models"
	adkmodels "google.golang.org/adk/server/restapi/models"
	"google.golang.org/genai"
)

// RegisterAgentEndpoints registers all agent-related endpoints.
func RegisterAgentEndpoints(humaAPI huma.API, prefix string, adkClient *adk.Client, database *db.Database) {
	agentsGroup := huma.NewGroup(humaAPI, prefix)

	// Weather endpoint
	huma.Post(agentsGroup, "/weather", func(ctx context.Context, input *api.WeatherRequest) (*api.WeatherResponse, error) {
		appName, ok := config.AgentMapping["weather"]
		if !ok {
			return nil, fmt.Errorf("weather agent not configured")
		}

		fmt.Printf("Received weather request for city: %s (user: %s, session: %s)\n",
			input.Body.City, input.Body.UserID, input.Body.SessionID)

		result, err := adkClient.RunWithAutoSession(ctx, adkmodels.RunAgentRequest{
			AppName:   appName,
			UserId:    input.Body.UserID,
			SessionId: input.Body.SessionID,
			NewMessage: genai.Content{
				Role:  string(genai.RoleUser),
				Parts: []*genai.Part{{Text: input.Body.City}},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("weather agent processing failed: %w", err)
		}

		resp := &api.WeatherResponse{}
		resp.Body.Forecast = result.FinalText
		return resp, nil
	})

	// Nutrition endpoint
	huma.Post(agentsGroup, "/nutrition", func(ctx context.Context, input *api.NutritionRequest) (*api.NutritionResponse, error) {
		appName, ok := config.AgentMapping["nutrition"]
		if !ok {
			return nil, fmt.Errorf("nutrition agent not configured")
		}

		fmt.Printf("Received nutrition request: %s (user: %s, session: %s)\n",
			input.Body.Text, input.Body.UserID, input.Body.SessionID)

		result, err := adkClient.RunWithAutoSession(ctx, adkmodels.RunAgentRequest{
			AppName:   appName,
			UserId:    input.Body.UserID,
			SessionId: input.Body.SessionID,
			NewMessage: genai.Content{
				Role:  string(genai.RoleUser),
				Parts: []*genai.Part{{Text: input.Body.Text}},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("nutrition agent processing failed: %w", err)
		}

		resp := &api.NutritionResponse{}
		var payload models.NutritionPayload
		if err := json.Unmarshal([]byte(result.FinalText), &payload); err != nil {
			return nil, fmt.Errorf("failed to parse nutrition response: %w", err)
		}
		resp.Body.Analysis = payload
		resp.Body.SessionID = input.Body.SessionID

		// Persist to database if available
		if database != nil {
			_, _ = database.MealLogRepository.Create(ctx,
				input.Body.UserID,
				input.Body.SessionID,
				payload.Name,
				string(payload.MealType),
				time.Now(),
				payload.Macros,
				payload.Assumptions,
				"ai_estimated",
				payload,
			)
		}

		return resp, nil
	})
}
