# Task: Make LazySessionService Truly Lazy or Fail Fast

## Status: Pending
## Priority: Medium

### Summary
`LazySessionService` in `internal/adk/session.go` is misnamed. Despite the "Lazy" prefix, it captures an error once at startup and returns that same error on every subsequent call forever. There is no retry logic. Compare with `LazyDatabase` (`internal/db/lazy.go`) which genuinely retries connection on each call via `getOrConnect()`.

### Current Behavior
```go
type LazySessionService struct {
    service session.Service
    err     error              // captured once, returned forever
}
```

If the database is temporarily unavailable at startup, the session service is permanently broken for the lifetime of the process. The app starts silently, and every agent call fails with `"session service unavailable: <original error>"`.

### Options

**Option A — Fail Fast (simpler)**
- Change `CreateSessionService()` to return `(session.Service, error)`
- In `main.go`, if session service creation fails, log clearly and skip agent registration entirely
- Removes the `LazySessionService` type
- Pros: Honest failure, simpler code, obvious at startup
- Cons: Server won't start agent features if DB is temporarily slow

**Option B — True Lazy (matches LazyDatabase pattern)**
- Store a factory function instead of an error, retry on each call:
```go
type LazySessionService struct {
    mu      sync.RWMutex
    service session.Service
    initFn  func() (session.Service, error)
}
```
- Pros: Resilient to startup timing issues, consistent with `LazyDatabase` pattern
- Cons: More complex, GORM migration retry on each call could be expensive

### Recommendation
Option A is preferred unless there are deployment scenarios where the DB might be unavailable at startup but available shortly after (e.g., container orchestration ordering). In that case, Option B with a cached-after-success pattern makes sense.

### Steps
- [ ] Decide between Option A and Option B
- [ ] Implement the chosen approach
- [ ] Update `cmd/api/main.go` if the `CreateSessionService()` signature changes
- [ ] Remove or rename `LazySessionService` accordingly
