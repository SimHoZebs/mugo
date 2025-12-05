package config

import "os"

const (
	AppName   = "mugo"
	ModelName = "gemini-2.5-flash"
)

// AgentMapping maps API route names to ADK agent app names.
// Key: route name (e.g., "nutrition")
// Value: ADK app name (e.g., "macro_estimator")
var AgentMapping = map[string]string{
	"nutrition": "macro_estimator",
	"weather":   "hello_time_agent",
	"echo":      "echo_agent",
}

// GetADKServerURL returns the ADK server URL from environment variable.
// Defaults to "http://localhost:8080/api" if not set.
func GetADKServerURL() string {
	url := os.Getenv("ADK_SERVER_URL")
	if url == "" {
		return "http://localhost:8080/api"
	}
	return url
}
