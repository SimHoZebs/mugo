# 021 - Agent JSON Redundancy & Double Unmarshal

## Issue
The system currently performs redundant JSON operations during agent execution and route handling.
- `AfterModelCallback` in `internal/agents/` unmarshals LLM output to normalize it, then marshals it back to a string.
- `internal/routes/meals.go` then unmarshals that same string *again* into a Go struct.

## Potential Justifications
- **Interface Genericity**: Keeping `AgentRunner.Run` returning a string allows it to remain agnostic of the specific agent's response schema, avoiding complex Go generics at the ADK wrapper level.
- **Protocol Stability**: ADK's internal event logging and session history may require `resp.Content.Parts[0].Text` to be a valid JSON string for persistence.

## Required Fix
1.  **Refactor AgentRunner**: Modify the `AgentRunner` interface and implementation to return a structured result (or the already-unmarshaled payload) alongside the raw text.
2.  **Eliminate redundant Marshal**: Stop calling `json.Marshal` inside `AfterModelCallbacks` if the data is just being passed to the next layer that will unmarshal it anyway.
3.  **Consolidate Normalization**: Ensure `NormalizeMealsBatchResponse` and `NormalizeNutritionResponse` are only called once per lifecycle.

## Verification
- **CRITICAL**: Double-check that removing `json.Marshal` from callbacks doesn't break ADK's internal event logging or session persistence, which may rely on `resp.Content.Parts[0].Text` being a valid JSON string.
- Verify that `CreateMeal` and `UpdateMeal` handlers in `meals.go` still correctly process the full batch of results.
- Run `make tests` to ensure agent orchestration logic remains intact.
