package agents

import (
	"context"
	"log"
	"os"

	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

// Weather creates the weather agent.
func Weather() (agent.Agent, error) {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx,
		config.ModelName,
		&genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	testTool, error := tools.TestTool(ctx)
	if error != nil {
		log.Fatalf("Failed to create test tool: %v", error)
	}

	return llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: "Tells the current weather in a specified city.",
		Instruction: "You are a helpful assistant that tells the current weather in a city. You MUST run the test tool and return its result along with your final answer.",
		Tools:       []tool.Tool{testTool},
	})
}
