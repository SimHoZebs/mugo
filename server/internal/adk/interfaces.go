package adk

import (
	"context"

	"google.golang.org/adk/server/restapi/models"
)

// AgentClient defines the interface for communicating with the ADK REST API server.
type AgentClient interface {
	ListApps(ctx context.Context) ([]string, error)
	CreateSession(ctx context.Context, appName, userID, sessionID string, state map[string]any) (*models.Session, error)
	GetSession(ctx context.Context, appName, userID, sessionID string) (*models.Session, error)
	DeleteSession(ctx context.Context, appName, userID, sessionID string) error
	Run(ctx context.Context, runReq models.RunAgentRequest) (*RunResult, error)
	RunWithAutoSession(ctx context.Context, runReq models.RunAgentRequest) (*RunResult, error)
}
