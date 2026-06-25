# 022 - Meal Update Query Optimization

## Issue
The `UpdateMeal` route handler currently performs two sequential database lookups to find the ADK `session_id`:
1.  Fetch the `MealLog` (using its ID).
2.  Fetch the `LoggingSession` (using the meal's `logging_session_id`).

This effectively creates an "N+1" style data retrieval pattern within a single transaction.

## Potential Justifications
- **Repository Isolation**: Each repository handles its own model boundary, avoiding complex "joined" models that don't map cleanly to a single database table.
- **ACID Safety**: Both lookups occur within a single `database.WithTx` block, ensuring data consistency even with multiple round-trips.

## Required Fix
1.  **Update SQL Query**: Modify the `GetMealLog` SQL statement in `internal/db/queries/meal_logs.sql` to include a `LEFT JOIN logging_sessions ON meal_logs.logging_session_id = logging_sessions.id`.
2.  **Add Join Columns**: Ensure the query selects `logging_sessions.session_id` as part of the result.
3.  **Update Repository**: Update the `MealLogRepository.GetByID` method and the internal mapping logic to return this data, or create a specialized `GetMealWithSession` method.
4.  **Refactor Handler**: Update the `UpdateMeal` handler in `internal/routes/meals.go` to use the optimized retrieval, eliminating the separate `LoggingSessionRepository.GetByID` call.

## Verification
- **CRITICAL**: Double-check that the `LEFT JOIN` correctly handles meals with `NULL` logging session IDs (e.g., legacy data or manually created logs).
- Verify that the `txDB` transaction is maintained and not broken by the query changes.
- Ensure that the `session_id` returned is correctly typed and matches the ADK's requirements.
