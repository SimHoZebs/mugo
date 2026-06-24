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

func NormalizeNutritionResponse(text string) (string, *models.NutritionPayload, error) {
	text = StripMarkdownFences(text)

	var payload models.NutritionPayload
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return text, nil, fmt.Errorf("nutrition agent: response did not match expected schema: %w\nContent: %s", err, text)
	}

	return text, &payload, nil
}

func CreateMacroEstimator(model adkmodel.LLM) (agent.Agent, error) {

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
						"category":      {Type: genai.TypeString, Description: "Category of assumption (e.g. portion, ingredient)"},
						"field":         {Type: genai.TypeString, Description: "The field being assumed (e.g. weight, quantity)"},
						"assumed_value": {Type: genai.TypeNumber, Description: "The numerical value used for the assumption"},
						"confidence":    {Type: genai.TypeString, Description: "low|medium|high"},
						"rationale":     {Type: genai.TypeString, Description: "Reasoning behind the assumption"},
					},
					Required: []string{"assumed_value", "rationale"},
				},
			},
			"meal_type": {
				Type:        genai.TypeString,
				Description: "The type of meal (breakfast, brunch, lunch, dinner, snack, or unknown)",
			},
		},
		Required: []string{"name", "date", "macros", "assumptions", "meal_type"},
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

		cleanText, _, err := NormalizeNutritionResponse(text)
		if err != nil {
			return nil, err
		}

		resp.Content.Parts[0].Text = cleanText
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
5. The meal type (breakfast, brunch, lunch, dinner, snack, or unknown)
`,
		OutputSchema:        schema,
		AfterModelCallbacks: []llmagent.AfterModelCallback{onAfterModelAssignIDs},
		OutputKey:           "meal_result",
	})
}
