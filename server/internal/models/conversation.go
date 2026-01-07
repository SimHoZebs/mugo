package models

type Conversation struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	SessionID string  `json:"session_id"`
	Title     *string `json:"title,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
