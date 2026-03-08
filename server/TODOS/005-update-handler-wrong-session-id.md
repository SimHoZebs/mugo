# Task: Fix Update Handler Passing Wrong ID as ADK Session ID

## Status: Completed
## Priority: Critical

### Summary
The meal update handler in `internal/routes/meals.go:150-153` passes `meal.ConversationID` (the `conversations.id` UUID primary key) as the ADK session ID. But the ADK session was registered under `conversations.session_id` (a different string generated at conversation creation time). This means meal corrections cannot find the original ADK session and the agent loses all conversation context.

### Affected Files
- `internal/routes/meals.go` — update handler (lines 144-165)

### Current Code (lines 150-165)
```go
sessionID := ""
if meal.ConversationID != nil {
    sessionID = *meal.ConversationID  // BUG: this is conversations.id (UUID PK), not conversations.session_id
}
// ...
result, err := mealRunner.Run(ctx, meal.UserID, sessionID, input.Body.Correction)
```

### The Two Identifiers
| Field | Value | Example |
|---|---|---|
| `conversations.id` (UUID PK) | Internal database ID, FK target for `meal_logs` | `a1b2c3d4-...` |
| `conversations.session_id` | String used as the ADK session identifier | `e5f6g7h8-...` (different UUID) |

`meal.ConversationID` maps to `conversations.id`. The ADK needs `conversations.session_id`.

### Proposed Fix
Look up the conversation record to get the real session ID:

```go
adkSessionID := ""
if meal.ConversationID != nil {
    conv, err := txDB.ConversationRepository.GetByID(ctx, *meal.ConversationID)
    if err != nil {
        return fmt.Errorf("failed to get conversation for meal: %w", err)
    }
    adkSessionID = conv.SessionID
}
// ...
result, err := mealRunner.Run(ctx, meal.UserID, adkSessionID, input.Body.Correction)
```

`ConversationRepository` is already available on `TxDatabase` (see `internal/db/database.go:56,76`).

### Steps
- [x] Replace direct `meal.ConversationID` usage with a conversation lookup via `txDB.ConversationRepository.GetByID()`
- [x] Rename the variable from `sessionID` to `adkSessionID` for clarity
- [ ] Test with a meal update request on a meal that has a valid conversation
