package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/agents"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
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

	sessionService := adk.CreateSessionService()

	// 1. Initialize Individual Agents
	// We do this manually for now to handle dependencies (like Orchestrator needing Macro)

	macroAgent, err := agents.MacroEstimator()
	if err != nil {
		log.Printf("Macro Estimator could not be initialized: %v", err)
	}

	orchestratorAgent, err := agents.MealOrchestrator(macroAgent)
	if err != nil {
		log.Printf("Meal Orchestrator could not be initialized: %v", err)
	}

	echoAgent, err := agents.NewEchoAgent()
	if err != nil {
		log.Printf("Warning: Echo Agent could not be initialized: %v", err)
	}

	weatherAgent, err := agents.Weather()
	if err != nil {
		log.Printf("Weather Agent could not be initialized: %v", err)
	}

	// 2. Create Runners and Populate Registry
	var runners []adk.AgentRunner
	addRunner := func(id string, a agent.Agent) {
		if a == nil {
			return
		}
		r, err := adk.NewAgentRunner(id, a, sessionService)
		if err != nil {
			log.Printf("Warning: Failed to create runner for %s: %v", id, err)
			return
		}
		runners = append(runners, r)
	}

	addRunner("macro_estimator", macroAgent)
	addRunner("meal_orchestrator", orchestratorAgent)
	addRunner("echo_agent", echoAgent)
	addRunner("hello_time_agent", weatherAgent)

	runnersRegistry := adk.NewRunnerRegistry(runners...)

	lazyDB := db.NewLazyDatabase(ctx)
	defer lazyDB.Close()
	log.Println("Lazy database initialized - will connect on first use")

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

	mealRunner, _ := runnersRegistry.Get("meal_orchestrator")
	echoRunner, _ := runnersRegistry.Get("echo_agent")
	weatherRunner, _ := runnersRegistry.Get("hello_time_agent")

	routes.RegisterAgentEndpoints(api, "/agents", echoRunner, weatherRunner)
	routes.RegisterDebugEndpoints(api, "/debug", sessionService, lazyDB)
	routes.RegisterUserEndpoints(api, "/users", lazyDB)
	routes.RegisterMealEndpoints(api, "/meals", mealRunner, lazyDB)
	routes.RegisterAnalyticsEndpoints(api, "/analytics", lazyDB)
	routes.RegisterConversationEndpoints(api, "/conversations", lazyDB)
	routes.RegisterTranscriptionEndpoints(api, "/transcription")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	srv := &http.Server{
		Addr:    ":" + port,
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
		log.Printf("Received signal %v, starting graceful shutdown...", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
			srv.Close()
		}

		log.Println("Server stopped gracefully")
	}
}
