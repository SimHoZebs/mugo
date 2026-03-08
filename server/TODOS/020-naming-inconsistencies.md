# TODO-020: Naming Inconsistencies (UserId, camelCase files, yeet)

## Status: Pending
## Priority: Low

### Summary
Several naming inconsistencies exist in the codebase that deviate from standard Go conventions or the project's own established patterns.
1. `UserId`/`SessionId` fields used instead of `UserID`/`SessionID`.
2. Several files in `internal/agents/` use camelCase instead of snake_case.
3. A function named `yeet` in `internal/tools/test.go` is unprofessional and non-descriptive.

### Affected Code
**File:** `internal/routes/debug.go`
**Lines:** 14-15

```go
// Line 14
UserId string `json:"user_id"`
// Line 15
SessionId string `json:"session_id"`
```

**Files:** `internal/agents/`
- `mealOrchestrator.go`
- `macroEstimator.go`

**File:** `internal/tools/test.go`
**Line:** 20

```go
// Line 20
func yeet(ctx tool.Context, args TestArgs) (TestResponse, error) {
```

### Risk
Non-standard naming can make the codebase harder to maintain and less intuitive for new contributors. Inconsistent file naming breaks convention.

### Proposed Fix
Rename `UserId` to `UserID` and `SessionId` to `SessionID`. Standardize all Go file names to snake_case. Rename the `yeet` function to something more descriptive like `readTestFile`.
