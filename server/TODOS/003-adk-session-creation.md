# Task: Add ADK Session Creation Before runner.Run()

## Status: Completed
## Priority: Critical

### Summary
The ADK `runner.Run()` calls `sessionService.Get()` as its first operation and fails with "record not found" if the session does not exist. Our code never calls `sessionService.Create()` before invoking the runner, so every new conversation fails.

### Root Cause
The comment in `internal/adk/runner.go:85-86` incorrectly states that the runner handles session initialization. Reading the ADK source (`runner/runner.go:94-107`), it does not — it calls `Get()` and yields the error on failure.

Per the [ADK session lifecycle docs](https://google.github.io/adk-docs/sessions/session/), the application is responsible for calling `sessionService.Create()` before the runner is invoked.

### Affected Files
- `internal/adk/runner.go` — the `agentRunner.Run()` method (line 76)

### Proposed Fix
Add a get-or-create pattern inside `agentRunner.Run()` before calling `r.runner.Run()`:

```go
// Ensure ADK session exists
_, err := r.sessionService.Get(ctx, &session.GetRequest{
    AppName:   r.appName,
    UserID:    userID,
    SessionID: sessionID,
})
if err != nil {
    _, err = r.sessionService.Create(ctx, &session.CreateRequest{
        AppName:   r.appName,
        UserID:    userID,
        SessionID: sessionID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to initialize ADK session: %w", err)
    }
}
```

### Why Here (Not in Route Handlers)
Putting the get-or-create logic in `agentRunner.Run()` encapsulates the concern in one place. All callers (meal creation, meal updates, echo agent, weather agent) benefit automatically without leaking ADK internals into route handlers.

### Steps
- [x] Add get-or-create logic to `agentRunner.Run()` in `internal/adk/runner.go`
- [x] Remove the misleading comment on lines 85-88
- [ ] Test with a `POST /meals` request using a new user/session that doesn't exist in the ADK `sessions` table
