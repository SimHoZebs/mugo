package agents

import (
	"context"

	"github.com/simhozebs/mugo/internal/config"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func NewGeminiModel() (model.LLM, error) {
	return gemini.NewModel(
		context.Background(),
		config.ModelName,
		&genai.ClientConfig{APIKey: config.GetGoogleAPIKey()},
	)
}
