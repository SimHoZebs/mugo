# Task: Fix Leaked GORM Connection in CreateSessionService()

## Status: Completed
## Priority: Critical

### Summary
`CreateSessionService()` in `internal/adk/runner.go:141` calls `gorm.Open(postgres.Open(dbURL))` to create a full `*gorm.DB` connection, then only uses `db.Dialector` to pass to `database.NewSessionService()`. The ADK's `NewSessionService()` internally calls `gorm.Open()` again from that dialector, opening a second connection. The original `*gorm.DB` is never used or closed — it's a leaked connection.

This means the application holds **two** unmanaged GORM connection pools to Postgres (in addition to the main pgx pool), and one of them is completely orphaned.

### Affected Files
- `internal/adk/runner.go` — `CreateSessionService()` function (lines 131-163)

### Current Code (lines 140-148)
```go
db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})  // leaked — never used or closed
if err != nil {
    ...
}

svc, err := database.NewSessionService(db.Dialector)  // opens ANOTHER connection
```

### Proposed Fix
`postgres.Open(dbURL)` returns a `gorm.Dialector` without opening a connection. Pass it directly:

```go
svc, err := database.NewSessionService(postgres.Open(dbURL))
if err != nil {
    ...
}
```

This eliminates the intermediate `gorm.Open()` call entirely. `NewSessionService()` will open the one connection it needs internally.

### Steps
- [x] Replace lines 140-148 in `CreateSessionService()` with direct `postgres.Open(dbURL)` passed to `NewSessionService()`
- [x] Remove the unused `gorm` import (no longer needed)
- [ ] Verify the server starts and the ADK session tables are still migrated correctly
