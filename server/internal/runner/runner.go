package runner

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type RunFunc func(ctx context.Context, userID, sessionID, text string) (*RunResult, error)

type CreateSessionFunc func(ctx context.Context, userID, sessionID string) error

type RunResult struct {
	Events    []*session.Event
	FinalText string
}

func NewRunner(appName string, ag agent.Agent, ss session.Service) (*runner.Runner, error) {
	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          ag,
		SessionService: ss,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create runner for %s: %w", appName, err)
	}
	return r, nil
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
