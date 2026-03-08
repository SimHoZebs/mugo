package agents

import (
	"encoding/json"
	"fmt"

	"github.com/simhozebs/mugo/internal/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	adkmodel "google.golang.org/adk/model"
	"google.golang.org/genai"
)

func NormalizeNutritionResponse(text string) (*models.NutritionPayload, error) {
	text = StripMarkdownFences(text, "macro_estimator")

	var payload models.NutritionPayload
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return nil, fmt.Errorf("nutrition agent: response did not match expected schema: %w\nContent: %s", err, text)
	}

	for i := range payload.Assumptions {
		if payload.Assumptions[i].ID == "" {
			payload.Assumptions[i].ID = fmt.Sprintf("A%d", i+1)
		}
		if payload.Assumptions[i].Unit == "" {
			payload.Assumptions[i].Unit = "g"
		}
	}

	return &payload, nil
}

func MacroEstimator() (agent.Agent, error) {
	model, err := NewGeminiModel()
	if err != nil {
		return nil, err
	}

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type:        genai.TypeString,
				Description: "A short, descriptive name for the meal",
			},
			"date": {
				Type:        genai.TypeString,
				Description: "Date of the meal in YYYY-MM-DD format",
			},
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
			"meal_type": {
				Type:        genai.TypeString,
				Description: "The type of meal (breakfast, lunch, dinner, or snack)",
			},
		},
		Required: []string{"name", "date", "macros", "assumptions"},
	}

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

		payload, err := NormalizeNutritionResponse(text)
		if err != nil {
			return nil, err
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
		Description: "Estimates nutritional macros and assumptions for a single described meal. Returns structured JSON for one meal.",
		Instruction: `You are a nutritional estimation assistant for a single meal.
You will receive a meal description that includes the meal's date in YYYY-MM-DD format.
You MUST provide:
1. A short, descriptive name for the meal (e.g., "Grilled Chicken Caesar Salad")
2. The date of the meal in YYYY-MM-DD format (as provided in the input)
3. The estimated macronutrients (calories, protein, carbs, fat)
4. A list of assumptions you made to reach these estimates
5. The meal type (breakfast, lunch, dinner, or snack)
`,
		OutputSchema:        schema,
		AfterModelCallbacks: []llmagent.AfterModelCallback{onAfterModelAssignIDs},
		OutputKey:           "meal_result",
	})
}
