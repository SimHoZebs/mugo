# Task: Simplify RunnerRegistry to Direct Variables

## Status: Pending
## Priority: Low

### Summary
`RunnerRegistry` (`internal/adk/runner.go:38-57`) is a `map[string]AgentRunner` with a `Get(name)` method. It is created in `cmd/api/main.go:80`, queried 3 times on lines 104-106, and never used again. The `GetSessionService()` method on the registry (line 122) is dead code — nothing calls it.

### Current Flow (cmd/api/main.go)
```go
// Build
runnersRegistry := adk.NewRunnerRegistry(runners...)

// Immediately destructure — registry never used again
mealRunner, _ := runnersRegistry.Get("meal_orchestrator")
echoRunner, _ := runnersRegistry.Get("echo_agent")
weatherRunner, _ := runnersRegistry.Get("hello_time_agent")
```

### Problems
- The registry adds indirection without benefit — it's created and consumed in the same function scope
- `GetSessionService()` iterates runners to extract a session service but is never called
- The `_` on `Get()` silently ignores missing runners, making failures hard to trace

### Proposed Fix
Replace the registry with direct runner construction:

```go
mealRunner := createRunner("meal_orchestrator", orchestratorAgent, sessionService)
echoRunner := createRunner("echo_agent", echoAgent, sessionService)
weatherRunner := createRunner("hello_time_agent", weatherAgent, sessionService)
```

Where `createRunner` is a simple helper that returns nil on failure (like the current `addRunner` closure).

Then remove:
- `RunnerRegistry` struct
- `NewRunnerRegistry()` constructor
- `RunnerRegistry.Get()` method
- `RunnerRegistry.GetSessionService()` method (dead code)

### Steps
- [ ] Replace registry pattern in `cmd/api/main.go` with direct runner variables
- [ ] Remove `RunnerRegistry`, `NewRunnerRegistry`, `Get`, `GetSessionService` from `internal/adk/runner.go`
- [ ] Verify the server starts correctly

### Note
If a generic `/agents/{name}/run` endpoint is added later, a registry pattern would be justified. But that endpoint does not exist today — YAGNI.
