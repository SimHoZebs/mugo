package adk

import (
	"context"
	"fmt"
	"log"

	"github.com/simhozebs/mugo/internal/config"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/session/database"
	"google.golang.org/genai"
	"gorm.io/driver/postgres"
)

type AgentRunner interface {
	Run(ctx context.Context, userID, sessionID, text string) (*RunResult, error)
	CreateSession(ctx context.Context, userID, sessionID string) error
	Name() string
}

type agentRunner struct {
	runner         *runner.Runner
	sessionService session.Service
	appName        string
}

func (r *agentRunner) Name() string {
	return r.appName
}

type RunResult struct {
	Events    []*session.Event
	FinalText string
}

func NewAgentRunner(appName string, ag agent.Agent, sessionService session.Service) (AgentRunner, error) {
	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          ag,
		SessionService: sessionService,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create runner for %s: %w", appName, err)
	}

	return &agentRunner{
		runner:         r,
		sessionService: sessionService,
		appName:        appName,
	}, nil
}

func (r *agentRunner) Run(ctx context.Context, userID, sessionID, text string) (*RunResult, error) {
	msg := &genai.Content{
		Role:  string(genai.RoleUser),
		Parts: []*genai.Part{{Text: text}},
	}

	events := []*session.Event{}
	var lastText string

	for event, err := range r.runner.Run(ctx, userID, sessionID, msg, agent.RunConfig{}) {
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

// CreateSession creates a new ADK session. Call this when starting a new
// conversation, before the first Run() call. If the session already exists
// or the ID is invalid, it returns an error.
func (r *agentRunner) CreateSession(ctx context.Context, userID, sessionID string) error {
	_, err := r.sessionService.Create(ctx, &session.CreateRequest{
		AppName:   r.appName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to create ADK session: %w", err)
	}
	return nil
}

func (r *agentRunner) GetSession(ctx context.Context, userID, sessionID string) (session.Session, error) {
	resp, err := r.sessionService.Get(ctx, &session.GetRequest{
		AppName:   r.appName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Session, nil
}

func CreateSessionService() (session.Service, error) {
	dbURL := config.GetDatabaseURL()
	if dbURL == "" {
		return nil, fmt.Errorf("no database URL configured for session storage")
	}

	log.Println("Initializing database-backed session service for persistent sessions")

	// Pass the dialector directly — postgres.Open() returns a gorm.Dialector
	// without opening a connection. NewSessionService opens the one it needs internally.
	svc, err := database.NewSessionService(postgres.Open(dbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create session service: %w", err)
	}

	if err := database.AutoMigrate(svc); err != nil {
		return nil, fmt.Errorf("failed to migrate session tables: %w", err)
	}

	log.Println("ADK session tables migrated successfully")
	return svc, nil
}
