package mocks

import (
	"context"

	"github.com/simhozebs/mugo/internal/adk"
	"github.com/stretchr/testify/mock"
	"google.golang.org/adk/server/restapi/models"
)

// AgentClientMock is a mock implementation of the AgentClient interface.
type AgentClientMock struct {
	mock.Mock
}

func (m *AgentClientMock) ListApps(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *AgentClientMock) CreateSession(ctx context.Context, appName, userID, sessionID string, state map[string]any) (*models.Session, error) {
	args := m.Called(ctx, appName, userID, sessionID, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *AgentClientMock) GetSession(ctx context.Context, appName, userID, sessionID string) (*models.Session, error) {
	args := m.Called(ctx, appName, userID, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *AgentClientMock) DeleteSession(ctx context.Context, appName, userID, sessionID string) error {
	args := m.Called(ctx, appName, userID, sessionID)
	return args.Error(0)
}

func (m *AgentClientMock) Run(ctx context.Context, runReq models.RunAgentRequest) (*adk.RunResult, error) {
	args := m.Called(ctx, runReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adk.RunResult), args.Error(1)
}

func (m *AgentClientMock) RunWithAutoSession(ctx context.Context, runReq models.RunAgentRequest) (*adk.RunResult, error) {
	args := m.Called(ctx, runReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*adk.RunResult), args.Error(1)
}
