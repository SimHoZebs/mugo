package models

type LoggingSession struct {
	ID        string  `json:"id" doc:"Unique logging session ID" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string  `json:"user_id" doc:"User ID" example:"user-123"`
	SessionID string  `json:"session_id" doc:"ADK session ID" example:"session-456"`
	Title     *string `json:"title,omitempty" doc:"Optional title for the logging session" example:"Lunch at Joe's"`
	CreatedAt string  `json:"created_at" doc:"When the logging session started" example:"2025-01-07T12:00:00Z"`
	UpdatedAt string  `json:"updated_at" doc:"When the logging session was last active" example:"2025-01-07T12:05:00Z"`
}
