package main

import (
	"context"
	"github.com/simhozebs/mugo/internal/agents"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/server/restapi/services"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	weatherAgent, err := agents.Weather()
	echoAgent, err := agents.NewEchoAgent()
	nutritionAgent, err := agents.Nutrition()

	agentLoader, err := services.NewMultiAgentLoader(weatherAgent, echoAgent, nutritionAgent)
	if err != nil {
		log.Fatalf("Failed to create agent loader: %v", err)
	}
	config := &adk.Config{
		AgentLoader: agentLoader,
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
