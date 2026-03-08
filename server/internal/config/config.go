package config

import (
	"os"
	"strconv"
	"time"
)

const (
	AppName   = "mugo"
	ModelName = "gemini-3-flash-preview"
)

// GetADKServerURL returns ADK server URL from environment variable.
// Defaults to "http://localhost:8080/api" if not set.
func GetADKServerURL() string {
	url := os.Getenv("ADK_SERVER_URL")
	if url == "" {
		return "http://localhost:8080/api"
	}
	return url
}

// GetTranscriptionServerURL returns whisper transcription server URL from environment variable.
// Defaults to "http://localhost:9000" if not set.
func GetTranscriptionServerURL() string {
	url := os.Getenv("TRANSCRIPTION_SERVER_URL")
	if url == "" {
		return "http://localhost:9000"
	}
	return url
}

// GetDatabaseURL returns the database URL from environment variable.
func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

// GetDatabaseMinConns returns the minimum number of connections in the pool.
// Defaults to 5 if not set.
func GetDatabaseMinConns() int {
	return getIntEnv("DB_MIN_CONNS", 5)
}

// GetDatabaseMaxConns returns the maximum number of connections in the pool.
// Defaults to 25 if not set.
func GetDatabaseMaxConns() int {
	return getIntEnv("DB_MAX_CONNS", 25)
}

// GetDatabaseMaxConnLifetime returns the maximum lifetime of a connection.
// Defaults to 1 hour if not set.
func GetDatabaseMaxConnLifetime() time.Duration {
	return getDurationEnv("DB_MAX_CONN_LIFETIME", 1*time.Hour)
}

// GetDatabaseMaxConnIdleTime returns the maximum idle time of a connection.
// Defaults to 30 minutes if not set.
func GetDatabaseMaxConnIdleTime() time.Duration {
	return getDurationEnv("DB_MAX_CONN_IDLE_TIME", 30*time.Minute)
}

// GetDatabaseHealthCheckPeriod returns the health check period for connections.
// Defaults to 1 minute if not set.
func GetDatabaseHealthCheckPeriod() time.Duration {
	return getDurationEnv("DB_HEALTH_CHECK_PERIOD", 1*time.Minute)
}

// GetDatabaseConnectTimeout returns the connection timeout.
// Defaults to 5 seconds if not set.
func GetDatabaseConnectTimeout() time.Duration {
	return getDurationEnv("DB_CONNECT_TIMEOUT", 5*time.Second)
}

func getIntEnv(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return duration
}
