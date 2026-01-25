package api

import "github.com/simhozebs/mugo/internal/models"

// EchoRequest is the request body for the echo endpoint (testing).
type EchoRequest struct {
	Body struct {
		UserID    string `json:"user_id" example:"user_12345" doc:"User ID of the requester"`
		SessionID string `json:"session_id" example:"session_12345" doc:"Session ID for the conversation"`
		Message   string `json:"message" example:"Hello, world!" doc:"Message to echo back"`
	}
}

// EchoResponse is the response body for the echo endpoint (testing).
type EchoResponse struct {
	Body struct {
		Echo string `json:"echo" example:"Hello, world!" doc:"Echoed message"`
	}
}

// NutritionRequest is the request body for the nutrition endpoint.
type NutritionRequest struct {
	Body struct {
		UserID    string `json:"user_id" example:"user_12345" doc:"User ID of the requester"`
		SessionID string `json:"session_id" example:"session_12345" doc:"Session ID for the conversation"`
		Text      string `json:"text" example:"I ate a chicken sandwich" doc:"Description of food eaten"`
	}
}

// NutritionResponse is the response body for the nutrition endpoint.
type NutritionResponse struct {
	Body struct {
		Analysis  models.NutritionPayload `json:"analysis" doc:"Nutritional analysis and assumptions"`
		SessionID string                  `json:"session_id" example:"session_67890" doc:"Session ID for continued conversation"`
	}
}

// WeatherRequest is the request body for the weather endpoint.
type WeatherRequest struct {
	Body struct {
		UserID    string `json:"user_id" example:"user_12345" doc:"User ID of the requester"`
		SessionID string `json:"session_id" example:"session_12345" doc:"Session ID for the conversation"`
		City      string `json:"city" example:"San Francisco" doc:"City to get weather for"`
	}
}

// WeatherResponse is the response body for the weather endpoint.
type WeatherResponse struct {
	Body struct {
		Forecast string `json:"forecast" example:"Sunny with a high of 75°F" doc:"Weather forecast for the specified city"`
	}
}
