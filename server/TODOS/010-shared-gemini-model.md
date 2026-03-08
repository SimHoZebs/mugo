# Task: Share a Single Gemini Model Instance Across Agents

## Status: Pending
## Priority: Low

### Summary
Each agent factory (`MacroEstimator()`, `MealOrchestrator()`, `Weather()`) independently calls `NewGeminiModel()`, which calls `gemini.NewModel()` each time. This creates 3 separate model client instances, each potentially holding its own HTTP client and connection pool to the Gemini API.

### Current Code

`internal/agents/model.go`:
```go
func NewGeminiModel() (model.LLM, error) {
    return gemini.NewModel(
        context.Background(),
        config.ModelName,
        &genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")},
    )
}
```

Called independently in:
- `internal/agents/macroEstimator.go:35` — `model, err := NewGeminiModel()`
- `internal/agents/mealOrchestrator.go:30` — `model, err := NewGeminiModel()`
- `internal/agents/weather.go:13` — `model, err := NewGeminiModel()`

### Proposed Fix
Create the model once in `main.go` and pass it to each agent constructor:

```go
// main.go
model, err := agents.NewGeminiModel()
if err != nil {
    log.Fatalf("Failed to create Gemini model: %v", err)
}

macroAgent, err := agents.MacroEstimatorWithModel(model)
orchestratorAgent, err := agents.MealOrchestratorWithModel(model, macroAgent)
weatherAgent, err := agents.WeatherWithModel(model)
```

Update each agent constructor to accept a `model.LLM` parameter instead of creating its own.

### Prerequisites
- Verify that `model.LLM` returned by `gemini.NewModel()` is safe for concurrent use. Given that the ADK framework runs agents concurrently, this is almost certainly the case, but should be confirmed.

### Steps
- [ ] Add `WithModel` variants to each agent constructor (or change existing signatures)
- [ ] Create model once in `main.go` and pass to all agents
- [ ] Standardize constructor naming (`New` prefix or not — currently inconsistent)
- [ ] Verify concurrent safety of shared model instance
