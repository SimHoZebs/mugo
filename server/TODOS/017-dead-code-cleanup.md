# TODO-017: Dead Code Cleanup (httputil, hello.go, etc.)

## Status: Pending
## Priority: Medium

### Summary
Multiple files and functions exist in the project that are never imported, called, or registered. These add to the maintenance burden and cognitive load without providing value.

### Affected Code
**File:** `internal/httputil/client.go`
**Lines:** entire file (59 lines)
**Issue:** No other file in the project imports this package.

**File:** `internal/routes/hello.go`
**Lines:** entire file (9 lines)
**Issue:** `Hello()` is never registered in any router.

**File:** `internal/adk/runner.go`
**Lines:** 121-131
**Issue:** `GetSession()` is defined but not part of any interface and never called.

**File:** `internal/db/migrations.go`
**Lines:** 61-69
**Issue:** `tableExists()` is unused (comment says "preserved for potential future use").

**File:** `internal/db/pgutil/pgutil.go`
**Lines:** 41-46
**Issue:** `TextPtr()` is defined but never called.

**File:** `internal/routes/meals.go`
**Lines:** 55-59
**Issue:** `ListMealsByDateRangeRequest` struct is defined but never used (handler uses inline anonymous struct).

### Proposed Fix
Delete the unused code or integrate it if it's genuinely needed.
