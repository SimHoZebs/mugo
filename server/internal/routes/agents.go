package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/runner"
	"github.com/simhozebs/mugo/internal/api"
)

func RegisterAgentEndpoints(humaAPI huma.API, prefix string, echoRun runner.RunFunc, echoCreateSession runner.CreateSessionFunc, weatherRun runner.RunFunc, weatherCreateSession runner.CreateSessionFunc) {
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

		if echoRun == nil {
			return nil, huma.Error503ServiceUnavailable("echo_agent not found")
		}

		// Ensure ADK session exists
		if err := echoCreateSession(ctx, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create echo session: %v", err))
		}

		result, err := echoRun(ctx, input.Body.UserID, input.Body.SessionID, input.Body.Message)
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

		if weatherRun == nil {
			return nil, huma.Error503ServiceUnavailable("weather_agent not found")
		}

		// Ensure ADK session exists
		if err := weatherCreateSession(ctx, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create weather session: %v", err))
		}

		result, err := weatherRun(ctx, input.Body.UserID, input.Body.SessionID, input.Body.City)
		if err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("weather agent processing failed: %v", err))
		}

		resp := &api.WeatherResponse{}
		resp.Body.Forecast = result.FinalText
		return resp, nil
	})
}
