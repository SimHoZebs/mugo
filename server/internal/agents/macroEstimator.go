package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/simhozebs/mugo/internal/config"
	"github.com/simhozebs/mugo/internal/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

// MacroEstimator creates the nutrition estimation agent.
func MacroEstimator() (agent.Agent, error) {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx,
		config.ModelName,
		&genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"macros": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"calories": {Type: genai.TypeNumber, Description: "kilocalories"},
					"protein":  {Type: genai.TypeNumber, Description: "protein grams"},
					"carbs":    {Type: genai.TypeNumber, Description: "carbohydrate grams"},
					"fat":      {Type: genai.TypeNumber, Description: "fat grams"},
				},
				Required: []string{"calories", "protein", "carbs", "fat"},
			},
			"assumptions": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"id":            {Type: genai.TypeString, Description: "assumption id"},
						"text":          {Type: genai.TypeString, Description: "assumption text"},
						"category":      {Type: genai.TypeString},
						"field":         {Type: genai.TypeString},
						"assumed_value": {Type: genai.TypeNumber},
						"confidence":    {Type: genai.TypeString, Description: "low|medium|high"},
						"rationale":     {Type: genai.TypeString},
					},
					Required: []string{"assumed_value"},
				},
			},
		},
		Required: []string{"macros", "assumptions"},
	}

	// afterModel callback: strict unmarshal into NutritionPayload, assign IDs, error if schema mismatch
	onAfterModelAssignIDs := llmagent.AfterModelCallback(func(ctx agent.CallbackContext, resp *adkmodel.LLMResponse, respErr error) (*adkmodel.LLMResponse, error) {
		if respErr != nil {
			return nil, respErr
		}
		if resp == nil || resp.Content == nil || len(resp.Content.Parts) == 0 {
			return resp, nil
		}
		text := resp.Content.Parts[0].Text
		if text == "" {
			return resp, nil
		}

		var payload models.NutritionPayload
		if err := json.Unmarshal([]byte(text), &payload); err != nil {
			// Strict mode: fail the agent invocation so schema violations are visible
			return nil, fmt.Errorf("nutrition agent: response did not match expected schema: %w", err)
		}

		// Assign sequential IDs if missing and default unit to 'g' if empty
		for i := range payload.Assumptions {
			if payload.Assumptions[i].ID == "" {
				payload.Assumptions[i].ID = fmt.Sprintf("A%d", i+1)
			}
			if payload.Assumptions[i].Unit == "" {
				payload.Assumptions[i].Unit = "g"
			}
		}

		newBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("nutrition agent: failed to marshal normalized payload: %w", err)
		}
		resp.Content.Parts[0].Text = string(newBytes)
		return resp, nil
	})

	return llmagent.New(llmagent.Config{
		Name:        "macro_estimator",
		Model:       model,
		Description: "Estimates nutritional value (macros) and lists assumptions based on food description.",
		Instruction: `You are a nutritional estimation assistant.
Your goal is to estimate the macronutrients for the food described by the user.
You MUST provide:
1. A short, descriptive name for the meal (e.g., "Grilled Chicken Caesar Salad", "Homemade Beef Tacos")
2. The estimated macronutrients (calories, protein, carbs, fat)
3. A list of assumptions you made to reach these estimates
`,
		OutputSchema:        schema,
		AfterModelCallbacks: []llmagent.AfterModelCallback{onAfterModelAssignIDs},
	})
}
