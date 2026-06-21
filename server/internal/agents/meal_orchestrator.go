package agents

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/simhozebs/mugo/internal/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	adkmodel "google.golang.org/adk/model"
)

func NormalizeMealsBatchResponse(text string) (string, *models.MealsBatchPayload, error) {
	text = StripMarkdownFences(text)

	if !strings.HasPrefix(text, "{") {
		return text, nil, nil
	}

	var batch models.MealsBatchPayload
	if err := json.Unmarshal([]byte(text), &batch); err != nil {
		return text, nil, fmt.Errorf("meal orchestrator: final response did not match expected schema: %w\nContent: %s", err, text)
	}

	return text, &batch, nil
}

func MealOrchestrator(model adkmodel.LLM, macroEstimator agent.Agent) (agent.Agent, error) {

	onAfterModelNormalize := llmagent.AfterModelCallback(func(ctx agent.CallbackContext, resp *adkmodel.LLMResponse, respErr error) (*adkmodel.LLMResponse, error) {
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

		cleanText, batch, err := NormalizeMealsBatchResponse(text)
		if err != nil {
			return nil, err
		}
		if batch == nil {
			return resp, nil
		}

		resp.Content.Parts[0].Text = cleanText
		return resp, nil
	})

	return llmagent.New(llmagent.Config{
		Name:        "meal_orchestrator",
		Model:       model,
		Description: "Orchestrates meal log creation by identifying explicitly distinct eating occasions and delegating each to the macro_estimator sub-agent.",
		SubAgents:   []agent.Agent{macroEstimator},
		Instruction: `You are a meal log orchestration assistant.
The user's message begins with today's date in YYYY-MM-DD format.
Your job is to:
1. Identify each explicitly distinct eating occasion mentioned in the user's message.
   - Create multiple meal logs only when the user clearly describes separate eating occasions (for example, breakfast and dinner, yesterday's lunch and today's snack).
   - If multiple foods belong to the same eating occasion, keep them together as one meal log for now.
2. Resolve any relative date references (e.g. "today", "yesterday") using the date provided.
3. For each identified eating occasion, transfer to the macro_estimator sub-agent providing:
   - The meal description
   - Its resolved date in YYYY-MM-DD format
   Format the transfer message as: "Date: YYYY-MM-DD. Meal: <description>"
4. After ALL meals have been processed by macro_estimator, combine the results and reply with ONLY a JSON object in this exact format:
{"meals": [<macro_estimator result 1>, <macro_estimator result 2>, ...]}
Do not include any other text in your final reply.
`,
		AfterModelCallbacks: []llmagent.AfterModelCallback{onAfterModelNormalize},
	})
}
