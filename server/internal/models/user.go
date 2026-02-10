package models

type User struct {
	ID        string `json:"id" doc:"Unique user ID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string `json:"username" doc:"Unique username" example:"johndoe"`
	CreatedAt string `json:"created_at" doc:"When the user was created" example:"2025-01-01T10:00:00Z"`
	UpdatedAt string `json:"updated_at" doc:"When the user was last updated" example:"2025-01-01T10:00:00Z"`
}
