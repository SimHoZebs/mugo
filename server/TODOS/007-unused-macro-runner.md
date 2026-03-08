# Task: Remove Unused macro_estimator Runner

## Status: Pending
## Priority: Medium

### Summary
In `cmd/api/main.go:75`, a standalone runner is created for the `macro_estimator` agent. However, `macro_estimator` is a **sub-agent** of `meal_orchestrator` — the ADK framework handles sub-agent delegation internally. No route ever retrieves or uses the `macro_estimator` runner.

### Current Code (cmd/api/main.go:75)
```go
addRunner("macro_estimator", macroAgent)      // unused runner
addRunner("meal_orchestrator", orchestratorAgent)  // this is the one routes use
```

After the registry is created (line 80), only these runners are retrieved (lines 104-106):
```go
mealRunner, _ := runnersRegistry.Get("meal_orchestrator")
echoRunner, _ := runnersRegistry.Get("echo_agent")
weatherRunner, _ := runnersRegistry.Get("hello_time_agent")
```

`macro_estimator` is never retrieved.

### Impact
- Wastes an entry in the runner registry
- The `macroAgent` is still correctly used as a sub-agent within `orchestratorAgent` (passed via `agents.MealOrchestrator(macroAgent)` at line 46). The standalone runner is redundant.

### Proposed Fix
Remove line 75:
```go
// Remove this line:
addRunner("macro_estimator", macroAgent)
```

The `macroAgent` variable is still needed (line 46 passes it to `MealOrchestrator`), so only the runner creation is removed.

### Steps
- [ ] Remove `addRunner("macro_estimator", macroAgent)` from `cmd/api/main.go`
- [ ] Verify the server starts and meal creation still works (macro estimation handled by the orchestrator's sub-agent delegation)
