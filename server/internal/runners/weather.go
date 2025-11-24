package runners

import (
	"log"
	"server/internal/agents"
	"server/internal/config"
	"server/internal/shared"

	"google.golang.org/adk/runner"
)

func NewWeather() (*shared.AgentService, error) {
	agent, err := agents.Weather()
	if err != nil {
		log.Fatalf("failed to create weather agent: %v", err)
	}

	runnerConfig := runner.Config{
		Agent:          agent,
		AppName:        config.AppName,
		SessionService: shared.GetGlobalInMemorySessionService(),
	}

	agentRunner, err := runner.New(
		runnerConfig,
	)
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	return &shared.AgentService{
		Runner: agentRunner,
		Config: runnerConfig,
	}, nil
}
