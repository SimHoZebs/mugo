package mocks

import (
	"context"

	"github.com/simhozebs/mugo/internal/adk"
)

type MockAgentRunner struct {
	adk.AgentRunner
	RunFunc func(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error)
}

func (m *MockAgentRunner) Run(ctx context.Context, userID, sessionID, text string) (*adk.RunResult, error) {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, userID, sessionID, text)
	}
	return &adk.RunResult{}, nil
}
