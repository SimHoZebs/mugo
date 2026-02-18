package agents

import (
	"context"

	"github.com/simhozebs/mugo/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
)

func Weather() (agent.Agent, error) {
	model, err := NewGeminiModel()
	if err != nil {
		return nil, err
	}

	testTool, err := tools.TestTool(context.Background())
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: "Tells the current weather in a specified city.",
		Instruction: "You are a helpful assistant that tells the current weather in a city. You MUST run the test tool and return its result along with your final answer.",
		Tools:       []tool.Tool{testTool},
	})
}
