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
	"gorm.io/gorm"
)

type AgentRunner interface {
	Run(ctx context.Context, userID, sessionID, text string) (*RunResult, error)
}

type agentRunner struct {
	runner         *runner.Runner
	sessionService session.Service
	appName        string
}

type RunResult struct {
	Events    []*session.Event
	FinalText string
}

type RunnerRegistry struct {
	runners map[string]AgentRunner
}

func NewRunnerRegistry(runners ...AgentRunner) *RunnerRegistry {
	r := &RunnerRegistry{
		runners: make(map[string]AgentRunner),
	}
	for _, runner := range runners {
		if ar, ok := runner.(*agentRunner); ok {
			r.runners[ar.appName] = runner
		}
	}
	return r
}

func (r *RunnerRegistry) Get(appName string) (AgentRunner, bool) {
	runner, ok := r.runners[appName]
	return runner, ok
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

func (r *RunnerRegistry) GetSessionService() session.Service {
	for _, runner := range r.runners {
		if ar, ok := runner.(*agentRunner); ok {
			return ar.sessionService
		}
	}
	return nil
}

func CreateSessionService() (session.Service, error) {
	dbURL := config.GetADKSessionDatabaseURL()
	if dbURL == "" {
		log.Println("No database URL configured, using in-memory session service (sessions will be lost on restart)")
		return session.InMemoryService(), nil
	}

	log.Println("Initializing database-backed session service for persistent sessions")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	svc, err := database.NewSessionService(db.Dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to create session service: %w", err)
	}

	if err := database.AutoMigrate(svc); err != nil {
		return nil, fmt.Errorf("failed to migrate session tables: %w", err)
	}

	log.Println("ADK session tables migrated successfully")
	return svc, nil
}
