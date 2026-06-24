package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/simhozebs/mugo/internal/agents"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
	"github.com/simhozebs/mugo/internal/runner"
	"google.golang.org/adk/agent"
)

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	lazyDB := db.NewLazyDatabase(ctx)
	defer lazyDB.Close()
	log.Println("Lazy database initialized - will connect on first use")

	sessionService, err := runner.CreateSessionService()
	if err != nil {
		log.Printf("CRITICAL: ADK Session Service could not be initialized: %v", err)
		log.Println("Warning: AI agent features will be disabled")
	}

	model, err := agents.NewGeminiModel()
	if err != nil {
		log.Printf("CRITICAL: Shared Gemini Model could not be initialized: %v", err)
	}

	macroAgent, err := agents.CreateMacroEstimator(model)
	if err != nil {
		log.Printf("Macro Estimator could not be initialized: %v", err)
	}

	orchestratorAgent, err := agents.CreateMealOrchestrator(model, macroAgent)
	if err != nil {
		log.Printf("Meal Orchestrator could not be initialized: %v", err)
	}

	echoAgent, err := agents.NewEchoAgent()
	if err != nil {
		log.Printf("Warning: Echo Agent could not be initialized: %v", err)
	}

	weatherAgent, err := agents.Weather(model)
	if err != nil {
		log.Printf("Weather Agent could not be initialized: %v", err)
	}

	createRunner := func(id string, a agent.Agent) (runner.RunFunc, runner.CreateSessionFunc) {
		if a == nil || sessionService == nil {
			return nil, nil
		}
		r, err := runner.NewRunner(id, a, sessionService)
		if err != nil {
			log.Printf("Warning: Failed to create runner for %s: %v", id, err)
			return nil, nil
		}
		run := func(ctx context.Context, userID, sessionID, text string) (*runner.RunResult, error) {
			return runner.Run(ctx, r, userID, sessionID, text)
		}
		createSession := func(ctx context.Context, userID, sessionID string) error {
			return runner.CreateSession(ctx, sessionService, id, userID, sessionID)
		}
		return run, createSession
	}

	mealRun, mealCreateSession := createRunner("meal_orchestrator", orchestratorAgent)
	echoRun, echoCreateSession := createRunner("echo_agent", echoAgent)
	weatherRun, weatherCreateSession := createRunner("weather_agent", weatherAgent)

	r := chi.NewMux()
	api := humachi.New(r, huma.DefaultConfig("Mugo API", "0.1.0"))

	huma.Register(api, huma.Operation{
		OperationID: "greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Greeting",
		Description: "Returns a greeting message",
		Tags:        []string{"General"},
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	routes.RegisterAgentEndpoints(api, "/agents", echoRun, echoCreateSession, weatherRun, weatherCreateSession)
	routes.RegisterDebugEndpoints(api, "/debug", sessionService, lazyDB)
	routes.RegisterUserEndpoints(api, "/users", lazyDB)
	routes.RegisterMealEndpoints(api, "/meals", mealRun, mealCreateSession, lazyDB)
	routes.RegisterAnalyticsEndpoints(api, "/analytics", lazyDB)
	routes.RegisterConversationEndpoints(api, "/conversations", lazyDB)
	routes.RegisterTranscriptionEndpoints(api, "/transcription")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down", sig)
	}
}
