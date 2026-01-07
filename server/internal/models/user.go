package models

type User struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}
