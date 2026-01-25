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
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize ADK client
	adkServerURL := config.GetADKServerURL()
	adkClient := adk.NewClient(adkServerURL)
	log.Printf("ADK client initialized with URL: %s", adkServerURL)

	lazyDB := db.NewLazyDatabase(ctx)
	defer lazyDB.Close()
	log.Println("Lazy database initialized - will connect on first use")

	r := chi.NewMux()
	api := humachi.New(r, huma.DefaultConfig("Mugo API", "0.1.0"))

	// Register GET /greeting/{name} handler.
	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	// Register agent endpoints (test endpoints only, no database)
	routes.RegisterAgentEndpoints(api, "/agents", adkClient)
	routes.RegisterDebugEndpoints(api, "/debug", adkClient, lazyDB)

	// Register user and meal endpoints (always registered, will connect on first use)
	routes.RegisterUserEndpoints(api, "/users", lazyDB)
	routes.RegisterMealEndpoints(api, "/meals", adkClient, lazyDB)
	routes.RegisterAnalyticsEndpoints(api, "/analytics", lazyDB)
	routes.RegisterConversationEndpoints(api, "/conversations", lazyDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	// Create HTTP server with explicit configuration for graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine so it doesn't block
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, starting graceful shutdown...", sig)

		// Give outstanding requests 5 seconds to complete
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
			srv.Close()
		}

		log.Println("Server stopped gracefully")
	}
}
