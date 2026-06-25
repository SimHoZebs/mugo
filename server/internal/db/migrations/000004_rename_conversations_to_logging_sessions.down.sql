-- Revert logging_sessions back to conversations
ALTER TABLE logging_sessions RENAME TO conversations;

-- Revert indexes
ALTER INDEX IF EXISTS idx_logging_sessions_user_id RENAME TO idx_conversations_user_id;
ALTER INDEX IF EXISTS idx_logging_sessions_session_id RENAME TO idx_conversations_session_id;

-- Revert logging_session_id column in meal_logs back to conversation_id
ALTER TABLE meal_logs RENAME COLUMN logging_session_id TO conversation_id;

-- Revert index for meal_logs
ALTER INDEX IF EXISTS idx_meal_logs_logging_session_id RENAME TO idx_meal_logs_conversation_id;
