# TODO-016: Inconsistent Error Handling (agents.go and debug.go)

## Status: Pending
## Priority: Medium

### Summary
1. `internal/routes/agents.go` uses `fmt.Errorf` directly in the return signature, while all other route handlers in the project use `huma.Error*` structured types. This results in inconsistent API error responses.
2. `internal/routes/debug.go` captures database errors but returns them as successful 200 responses with the error message in the body field.

### Affected Code
**File:** `internal/routes/agents.go`
**Lines:** 27, 32, 37, 57, 62, 67

```go
// Line 27
return nil, fmt.Errorf("echo_agent not found")
```

**File:** `internal/routes/debug.go`
**Lines:** 48-49, 52-56

```go
// Line 48
if err != nil {
    resp.Body.SessionIds = []string{fmt.Sprintf("DB error: %v", err)}
    return resp, nil // Returns 200 OK with the error message in the body field.
}
```

### Risk
Inconsistent API surface. Frontend clients cannot handle errors correctly if they are structured as 200 OK responses or as simple string errors.

### Proposed Fix
Standardize all route handlers to use `huma.Error*` functions for reporting failures. Ensure error responses use the correct HTTP status codes (e.g., 404, 500, 503).
