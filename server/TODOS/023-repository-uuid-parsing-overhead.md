# 023 - Repository UUID Parsing Overhead

## Issue
The repository layer (e.g., `meal_repository.go`, `user_repository.go`) currently accepts `string` IDs for methods like `ListByUser`, `GetByID`, and `Delete`. Each call repeatedly invokes `pgutil.ParseUUID(id)` to convert the string to a `pgtype.UUID`.

## Potential Justifications
- **Defensive Programming**: The repository ensures its own integrity regardless of the caller (API, CLI, or background job).
- **Decoupling from External Types**: Avoids passing database-specific types like `pgtype.UUID` through higher-level application logic.

## Required Fix
1.  **Refactor Repository Interfaces**: Update all repository method signatures to accept `pgtype.UUID` (or a strongly typed `ID` alias) instead of raw `string`.
2.  **Move Validation to Edge**: Perform UUID parsing at the API boundary (e.g., in `huma` middleware or route handlers) where validation naturally belongs.
3.  **Update Call Sites**: Update all callers in `internal/routes/` to parse the UUID once and pass the typed value to the repository.

## Verification
- **CRITICAL**: Ensure that invalid UUIDs provided by the client are still caught early and return a `400 Bad Request`, not a `500 Internal Server Error` from the DB layer.
- Audit the mapping logic in `mapToMealLog` and similar functions to ensure consistency with how IDs are returned to the client.
- Run integration tests to ensure cross-repository data lookups (e.g., fetching a meal's user) are not broken by the type change.
