# TODO-020: Naming Inconsistencies

## Status: Completed
## Priority: Low

### Summary
Several naming inconsistencies existed in the codebase that deviated from standard Go conventions or the project's own established patterns.
1. `UserId`/`SessionId` fields used standard Go initialism casing after cleanup.
2. Several files in `internal/agents/` used camelCase instead of snake_case.
3. The `internal/tools/test.go` test-tool function now has a descriptive name.

### Affected Code
**File:** `internal/routes/debug.go`
**Lines:** 14-15

```go
// Line 14
UserID string `json:"user_id"`
// Line 15
SessionID string `json:"session_id"`
```

**Files:** `internal/agents/`
- `meal_orchestrator.go`
- `macro_estimator.go`

**File:** `internal/tools/test.go`
**Line:** 20

```go
// Line 20
func readTestFile(ctx tool.Context, args TestArgs) (TestResponse, error) {
```

### Risk
Non-standard naming can make the codebase harder to maintain and less intuitive for new contributors. Inconsistent file naming breaks convention.

### Proposed Fix
Renamed `UserId` to `UserID` and `SessionId` to `SessionID`. Standardized Go file names to snake_case. Renamed `yeet` to `readTestFile`.
