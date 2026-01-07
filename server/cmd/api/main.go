package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/db"
	"github.com/simhozebs/mugo/internal/routes"
	"log"
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
	go func() {
		<-sigChan
		cancel()
	}()

	// Initialize ADK client
	adkServerURL := config.GetADKServerURL()
	adkClient := adk.NewClient(adkServerURL)
	log.Printf("ADK client initialized with URL: %s", adkServerURL)

	// Initialize database
	var database *db.Database
	database, err := db.NewDatabase(ctx)
	if err != nil {
		if config.GetFailFastOnDBError() {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Server will run without database persistence")
	} else {
		defer database.Close()
		log.Println("Database connected successfully")
	}

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

	// Conversation endpoint (for testing with echo agent)
	huma.Post(api, "/conversation", func(ctx context.Context, input *struct {
		Body routes.ConversationRequest `body:""`
	}) (*routes.ConversationResponse, error) {
		return routes.ConversationHandler(ctx, adkClient, &struct {
			Body routes.ConversationRequest `body:""`
		}{Body: input.Body})
	})

	// Register agent endpoints with database
	routes.RegisterAgentEndpoints(api, "/agents", adkClient, database)
	routes.RegisterDebugEndpoints(api, "/debug", adkClient, database)

	// Register user and meal endpoints
	if database != nil {
		routes.RegisterUserEndpoints(api, "/users", database)
		routes.RegisterMealEndpoints(api, "/meals", database)
		routes.RegisterAnalyticsEndpoints(api, "/analytics", database)
		routes.RegisterConversationEndpoints(api, "/conversations", database)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
