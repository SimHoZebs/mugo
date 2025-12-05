package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/simhozebs/mugo/internal/adk"
	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/routes"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	// Initialize ADK client
	adkServerURL := config.GetADKServerURL()
	adkClient := adk.NewClient(adkServerURL)
	fmt.Printf("ADK client initialized with URL: %s\n", adkServerURL)

	r := chi.NewMux()
	api := humachi.New(r, huma.DefaultConfig("Conversation API", "0.1.0"))

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

	routes.RegisterAgentEndpoints(api, "/agents", adkClient)
	routes.RegisterDebugEndpoints(api, "/debug", adkClient)

	fmt.Println("Server starting at http://localhost:8888")
	http.ListenAndServe(":8888", r)
}
