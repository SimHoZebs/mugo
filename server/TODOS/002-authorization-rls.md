# Task: Implement Row Level Security (RLS) and Authorization

## Status: Pending
## Priority: Medium

### Summary
Currently, any authenticated user (or in the current state, anyone) can access all data. This task aims to restrict data access so users only see their own records (meal logs, conversations, etc.).

### Proposed Steps
- [ ] Enable Row Level Security (RLS) on database tables (`meal_logs`, `conversations`, `daily_nutrition_summaries`, etc.).
- [ ] Create a database migration for PostgreSQL `CREATE POLICY` statements.
- [ ] Update repository queries to handle user-based filtering.
- [ ] Implement a mechanism for the application to set the `app.current_user_id` context for each database connection/transaction.

### Technical Notes
- RLS provides a strong safety net by enforcing data access at the database level.
- This would require a database migration and a middleware or handler to set the appropriate Postgres session variables.
