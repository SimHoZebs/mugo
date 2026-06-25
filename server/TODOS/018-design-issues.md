# TODO-018: Design Issues (Context in Struct, Database Duplication)

## Status: Pending
## Priority: Low

### Summary
Several design issues exist in the codebase that violate Go best practices or introduce unnecessary code duplication.
1. `LazyDatabase` stores a `context.Context` in its struct field, which can lead to stale contexts or memory leaks.
2. `Database` and `TxDatabase` structs in `internal/db/database.go` duplicate all repository fields identically.
3. `http.DefaultClient` is used in `transcription.go` without a timeout.

### Affected Code
**File:** `internal/db/lazy.go`
**Lines:** 16-17, 22

```go
// Line 16
ctx context.Context
// Line 22
ld.ctx = ctx
```

**File:** `internal/db/database.go`
**Lines:** 11-17, 54-60

```go
// Line 11
type Database struct {
    UserRepository         UserRepository
    LoggingSessionRepository LoggingSessionRepository
    MealLogRepository     MealLogRepository
    NutritionRepository    NutritionRepository
    pool                  *pgxpool.Pool
}
// Line 54
type TxDatabase struct {
    UserRepository         UserRepository
    LoggingSessionRepository LoggingSessionRepository
    MealLogRepository     MealLogRepository
    NutritionRepository    NutritionRepository
    tx                    pgx.Tx
}
```

**File:** `internal/routes/transcription.go`
**Lines:** 78

```go
// Line 78
resp, err := http.DefaultClient.Do(req)
```

### Risk
Long-running requests that hang indefinitely without a timeout. Potential memory issues or stale state from storing contexts in structs. Increased code maintenance burden from duplicating repository fields.

### Proposed Fix
1. Pass `context.Context` as an argument to functions that need it (e.g., `GetDatabase(ctx)`).
2. Refactor `Database` and `TxDatabase` to share a common repository interface or struct.
3. Use a custom `http.Client` with a defined timeout for all network requests.
