package routes

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/simhozebs/mugo/internal/api"
	"github.com/simhozebs/mugo/internal/runner"
	adkrunner "google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

func RegisterAgentEndpoints(humaAPI huma.API, prefix string, echoRunner *adkrunner.Runner, echoSessionService session.Service, echoAppName string, weatherRunner *adkrunner.Runner, weatherSessionService session.Service, weatherAppName string) {
	agentsGroup := huma.NewGroup(humaAPI, prefix)

	huma.Register(agentsGroup, huma.Operation{
		OperationID: "agent-echo",
		Method:      "POST",
		Path:        "/echo",
		Summary:     "Echo agent",
		Description: "Tests ADK server response without LLM",
		Tags:        []string{"Agents"},
	}, func(ctx context.Context, input *api.EchoRequest) (*api.EchoResponse, error) {
		if echoRunner == nil {
			return nil, huma.Error503ServiceUnavailable("echo_agent not found")
		}

		if err := runner.CreateSession(ctx, echoSessionService, echoAppName, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create echo session: %v", err))
		}

		result, err := runner.Run(ctx, echoRunner, input.Body.UserID, input.Body.SessionID, input.Body.Message)
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
		if weatherRunner == nil {
			return nil, huma.Error503ServiceUnavailable("weather_agent not found")
		}

		if err := runner.CreateSession(ctx, weatherSessionService, weatherAppName, input.Body.UserID, input.Body.SessionID); err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create weather session: %v", err))
		}

		result, err := runner.Run(ctx, weatherRunner, input.Body.UserID, input.Body.SessionID, input.Body.City)
		if err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("weather agent processing failed: %v", err))
		}

		resp := &api.WeatherResponse{}
		resp.Body.Forecast = result.FinalText
		return resp, nil
	})
}
