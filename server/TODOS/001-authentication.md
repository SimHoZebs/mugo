# Task: Implement Authentication (AuthN)

## Status: Pending
## Priority: High

### Summary
The Mugo API currently has no authentication mechanism. All endpoints, including user listing and debug tools, are publicly accessible. This task involves adding a middleware layer to verify user identity.

### Proposed Steps
- [ ] Research and select an authentication strategy (e.g., JWT, Session Tokens, or Simple API Keys).
- [ ] Create a Chi middleware in `internal/routes/middleware.go` to validate credentials.
- [ ] Update `cmd/api/main.go` to apply the middleware to sensitive route groups.
- [ ] Integrate authentication into Huma's OpenAPI/Swagger documentation using `SecurityScheme`.
- [ ] Update `internal/routes/users.go` to require authentication for listing users.

### Technical Notes
- Huma v2 supports security schemes directly in the `huma.Config`.
- Context should be used to pass the authenticated `UserID` to the route handlers.
