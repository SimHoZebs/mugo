# TODO-012: Ignored json.Unmarshal Errors in meal_repository.go

## Status: Completed
## Priority: High (Critical)

### Summary
In `internal/db/repository/meal_repository.go`, within the `mapToMealLog()` function, the errors returned by `json.Unmarshal` are discarded (`_`). If the JSON in the database is malformed or corrupted, the resulting struct fields will be empty zero-values without any indication of failure.

### Affected Code
**File:** `internal/db/repository/meal_repository.go`
**Lines:** 206, 211, 216

```go
// Line 206
_ = json.Unmarshal(m.Macros, &macros)
// Line 211
_ = json.Unmarshal(m.Assumptions, &assumptions)
// Line 216
_ = json.Unmarshal(m.RawResponse, &rawResponse)
```

### Risk
Silent data corruption. Users may see empty meals or missing macro data if the JSON storage becomes invalid, with no way for the server to log or report the issue.

### Proposed Fix
Capture the error and return it so the repository can properly report the failure to the caller.

```go
if err := json.Unmarshal(m.Macros, &macros); err != nil {
    return nil, fmt.Errorf("failed to unmarshal macros: %w", err)
}
```
