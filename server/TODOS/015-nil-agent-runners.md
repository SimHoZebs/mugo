# TODO-015: Nil Agent Runners Passed to Route Handlers

## Status: Completed
## Priority: High (Critical)

### Summary
In `cmd/api/main.go`, if initializing the shared Gemini model or any agent runner fails, the error is logged but the application continues to start. The resulting `nil` runners are then passed to the route handlers, leading to potential panics or 503 errors at request time rather than a failure at startup.

### Affected Code
**File:** `cmd/api/main.go`
**Lines:** 43-46, 51-58, 84-86

```go
// Line 43
sharedModel, err := agents.NewGeminiModel()
if err != nil {
    log.Printf("CRITICAL: Shared Gemini Model could not be initialized: %v", err)
}
// sharedModel is nil here...
// Line 51
macroAgent, err := agents.MacroEstimator(sharedModel)
// macroAgent is nil here...
// Line 110
routes.RegisterAgentEndpoints(api, "/agents", echoRunner, weatherRunner)
// echoRunner / weatherRunner are potentially nil...
```

### Risk
Requests that attempt to use these nil runners will fail with errors (if checked) or panics (if not). A server starting in a degraded state is harder to monitor and debug than a server that fails to start when critical components are missing.

### Proposed Fix
Decide whether the application should fail to start if the Gemini model or session service is unavailable. If it should start degraded, ensure the route handlers gracefully handle nil runners using proper Huma error responses.

```go
if sharedModel == nil {
    log.Fatalf("Fatal: Gemini model initialization failed.")
}
```
