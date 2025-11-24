package runners

import (
	"log"
	"github.com/simhozebs/mugo/internal/agents"
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/shared"

	"google.golang.org/adk/runner"
)

func NewNutrition() (*shared.AgentService, error) {
	agent, err := agents.Nutrition()
	if err != nil {
		log.Fatalf("failed to create nutrition agent: %v", err)
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
