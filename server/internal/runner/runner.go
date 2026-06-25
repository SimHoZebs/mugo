package runner

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type RunResult struct {
	Events    []*session.Event
	FinalText string
}

type AgentRunner struct {
	Runner         *runner.Runner
	SessionService session.Service
	AppName        string
}

func NewAgentRunner(appName string, ag agent.Agent, ss session.Service) (*AgentRunner, error) {
	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          ag,
		SessionService: ss,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create runner for %s: %w", appName, err)
	}
	return &AgentRunner{
		Runner:         r,
		SessionService: ss,
		AppName:        appName,
	}, nil
}

func (a *AgentRunner) Run(ctx context.Context, userID, sessionID, text string) (*RunResult, error) {
	return Run(ctx, a.Runner, userID, sessionID, text)
}

func (a *AgentRunner) CreateSession(ctx context.Context, userID, sessionID string) error {
	return CreateSession(ctx, a.SessionService, a.AppName, userID, sessionID)
}

func Run(ctx context.Context, r *runner.Runner, userID, sessionID, text string) (*RunResult, error) {
	msg := &genai.Content{
		Role:  string(genai.RoleUser),
		Parts: []*genai.Part{{Text: text}},
	}

	events := []*session.Event{}
	var lastText string

	for event, err := range r.Run(ctx, userID, sessionID, msg, agent.RunConfig{}) {
		if err != nil {
			return nil, fmt.Errorf("runner error: %w", err)
		}
		events = append(events, event)

		if event.Content != nil && event.Content.Role == "model" {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					lastText = part.Text
				}
			}
		}
	}

	return &RunResult{
		Events:    events,
		FinalText: lastText,
	}, nil
}

func CreateSession(ctx context.Context, ss session.Service, appName, userID, sessionID string) error {
	_, err := ss.Create(ctx, &session.CreateRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to create ADK session: %w", err)
	}
	return nil
}
