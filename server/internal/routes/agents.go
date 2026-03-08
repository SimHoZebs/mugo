package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/api"
)

func RegisterAgentEndpoints(humaAPI huma.API, prefix string, echoRunner adk.AgentRunner, weatherRunner adk.AgentRunner) {
	agentsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(agentsGroup, huma.Operation{
		OperationID: "agent-echo",
		Method:      "POST",
		Path:        "/echo",
		Summary:     "Echo agent",
		Description: "Tests ADK server response without LLM",
		Tags:        []string{"Agents"},
	}, func(ctx context.Context, input *api.EchoRequest) (*api.EchoResponse, error) {
		fmt.Printf("Echo request: %s (user: %s, session: %s)\n",
			input.Body.Message, input.Body.UserID, input.Body.SessionID)

		if echoRunner == nil {
			return nil, huma.Error503ServiceUnavailable("echo_agent not found")
		}

		// Ensure ADK session exists
		if err := echoRunner.CreateSession(ctx, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create echo session: %v", err))
		}

		result, err := echoRunner.Run(ctx, input.Body.UserID, input.Body.SessionID, input.Body.Message)
		if err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("echo agent processing failed: %v", err))
		}

		resp := &api.EchoResponse{}
		resp.Body.Echo = result.FinalText
		return resp, nil
	})

	huma.Register(agentsGroup, huma.Operation{
		OperationID: "agent-weather",
		Method:      "POST",
		Path:        "/weather",
		Summary:     "Weather agent",
		Description: "Tests ADK + LLM integration",
		Tags:        []string{"Agents"},
	}, func(ctx context.Context, input *api.WeatherRequest) (*api.WeatherResponse, error) {
		fmt.Printf("Weather request for city: %s (user: %s, session: %s)\n",
			input.Body.City, input.Body.UserID, input.Body.SessionID)

		if weatherRunner == nil {
			return nil, huma.Error503ServiceUnavailable("weather_agent not found")
		}

		// Ensure ADK session exists
		if err := weatherRunner.CreateSession(ctx, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create weather session: %v", err))
		}

		result, err := weatherRunner.Run(ctx, input.Body.UserID, input.Body.SessionID, input.Body.City)
		if err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("weather agent processing failed: %v", err))
		}

		resp := &api.WeatherResponse{}
		resp.Body.Forecast = result.FinalText
		return resp, nil
	})
}
