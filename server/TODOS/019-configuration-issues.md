# TODO-019: Configuration Issues (raw os.Getenv, hardcoded paths)

## Status: Pending
## Priority: Low

### Summary
Several configuration issues exist that make the application harder to manage or deploy.
1. `internal/agents/model.go` uses `os.Getenv("GOOGLE_API_KEY")` directly, bypassing the `config` package.
2. `internal/tools/test.go` has a hardcoded relative path for reading a file.

### Affected Code
**File:** `internal/agents/model.go`
**Lines:** 17

```go
// Line 17
&genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")}
```

**File:** `internal/tools/test.go`
**Lines:** 22

```go
// Line 22
content, err := os.ReadFile("./tools/test.md")
```

### Risk
Unset environment variables or missing files will cause the application to behave incorrectly or crash with non-obvious error messages. Hardcoded relative paths are brittle and depend on the current working directory.

### Proposed Fix
Centralize all environment variables in the `config` package. Pass absolute file paths to functions that read files from the disk.
