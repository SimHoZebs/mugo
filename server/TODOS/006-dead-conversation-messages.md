# Task: Remove Dead conversation_messages Table and Queries

## Status: Pending
## Priority: Medium

### Summary
The `conversation_messages` table is defined in the schema and has sqlc-generated queries, but **no application code ever writes to or reads from it**. The ADK `events` table (managed by GORM auto-migration) already stores the full conversation history with richer structure (tool calls, grounding metadata, usage stats, state deltas).

### Evidence of Dead Code

**Schema definition** — `internal/db/migrations/000001_initial_schema.up.sql:48-58`:
- `CREATE TABLE conversation_messages` with columns: id, conversation_id, role, content, metadata, created_at

**sqlc queries** — `internal/db/queries/messages.sql`:
- `CreateMessage`, `GetMessage`, `GetMessagesByConversation`, `GetRecentMessages`, `DeleteMessagesByConversation`

**Generated Go code** — `internal/db/dbgenerated/messages.sql.go`:
- Full Go functions for all 5 queries above

**No repository layer** — There is no `MessageRepository` in `internal/db/repository/`. No route handler or service imports or calls any message-related function.

### Why It Exists
Likely created early in development before the ADK was integrated. The ADK's `events` table supersedes it entirely — events capture not just user/model messages but also tool calls, function results, state mutations, and metadata.

### Proposed Fix
- [ ] Remove `internal/db/queries/messages.sql`
- [ ] Re-run `make sqlc` to regenerate without message queries (deletes `messages.sql.go`)
- [ ] Optionally create a new migration (`000003_drop_conversation_messages.up.sql`) to drop the table
- [ ] Alternatively, leave the table in the schema and just remove the queries — the table does no harm sitting empty

### Technical Notes
- The `conversation_messages` table has a foreign key to `conversations(id)`. Dropping it has no cascade impact since nothing references `conversation_messages`.
- If data exists in any environment (unlikely since nothing writes to it), a migration would need a safety check.
