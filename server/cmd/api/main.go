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

	macroAgent, err := agents.MacroEstimator()
	if err != nil {
		log.Fatalf("Failed to create macro estimator agent: %v", err)
	}

	orchestratorAgent, err := agents.MealOrchestrator(macroAgent)
	if err != nil {
		log.Fatalf("Failed to create meal orchestrator agent: %v", err)
	}

	echoAgent, err := agents.NewEchoAgent()
	if err != nil {
		log.Fatalf("Failed to create echo agent: %v", err)
	}

	weatherAgent, err := agents.Weather()
	if err != nil {
		log.Fatalf("Failed to create weather agent: %v", err)
	}

	macroRunner, err := adk.NewAgentRunner("macro_estimator", macroAgent, sessionService)
	if err != nil {
		log.Fatalf("Failed to create macro runner: %v", err)
	}

	orchestratorRunner, err := adk.NewAgentRunner("meal_orchestrator", orchestratorAgent, sessionService)
	if err != nil {
		log.Fatalf("Failed to create orchestrator runner: %v", err)
	}

	echoRunner, err := adk.NewAgentRunner("echo_agent", echoAgent, sessionService)
	if err != nil {
		log.Fatalf("Failed to create echo runner: %v", err)
	}

	weatherRunner, err := adk.NewAgentRunner("hello_time_agent", weatherAgent, sessionService)
	if err != nil {
		log.Fatalf("Failed to create weather runner: %v", err)
	}

	runnerRegistry := adk.NewRunnerRegistry(macroRunner, orchestratorRunner, echoRunner, weatherRunner)

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

	routes.RegisterAgentEndpoints(api, "/agents", runnerRegistry)
	routes.RegisterDebugEndpoints(api, "/debug", runnerRegistry, lazyDB)
	routes.RegisterUserEndpoints(api, "/users", lazyDB)
	routes.RegisterMealEndpoints(api, "/meals", orchestratorRunner, lazyDB)
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
