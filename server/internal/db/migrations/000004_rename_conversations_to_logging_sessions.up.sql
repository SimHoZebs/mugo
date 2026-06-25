-- Rename conversations table to logging_sessions
ALTER TABLE conversations RENAME TO logging_sessions;

-- Rename indexes
ALTER INDEX IF EXISTS idx_conversations_user_id RENAME TO idx_logging_sessions_user_id;
ALTER INDEX IF EXISTS idx_conversations_session_id RENAME TO idx_logging_sessions_session_id;

-- Rename conversation_id column in meal_logs to logging_session_id
ALTER TABLE meal_logs RENAME COLUMN conversation_id TO logging_session_id;

-- Rename index for meal_logs
ALTER INDEX IF EXISTS idx_meal_logs_conversation_id RENAME TO idx_meal_logs_logging_session_id;
