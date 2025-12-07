package adk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/simhozebs/mugo/internal/httputil"

	"google.golang.org/adk/server/restapi/models"
	"google.golang.org/genai"
)

// Client is an HTTP client for communicating with the ADK REST API server.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// RunResult contains the result of running an agent.
type RunResult struct {
	Events    []models.Event // Full events for future manipulation
	FinalText string         // Extracted final text response
}

// NewClient creates a new ADK client with the given base URL.
// The baseURL should include the /api path, e.g., "http://localhost:8080/api"
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// sessionURL builds the session endpoint URL.
func (c *Client) sessionURL(appName, userID, sessionID string) string {
	return fmt.Sprintf("%s/apps/%s/users/%s/sessions/%s", c.baseURL, appName, userID, sessionID)
}

// doRequest executes an HTTP request and returns the response.
func (c *Client) doRequest(ctx context.Context, method, url string, body any) (*http.Response, error) {
	return httputil.DoRequest(ctx, c.httpClient, method, url, body)
}

// ListApps returns a list of available agent app names.
func (c *Client) ListApps(ctx context.Context) ([]string, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, c.baseURL+"/list-apps", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := httputil.CheckStatus(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var apps []string
	if err := httputil.DecodeJSON(resp, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

// CreateSession creates a new session for the given app, user, and session ID.
func (c *Client) CreateSession(ctx context.Context, appName, userID, sessionID string, state map[string]any) (*models.Session, error) {
	reqBody := models.CreateSessionRequest{State: state}

	resp, err := c.doRequest(ctx, http.MethodPost, c.sessionURL(appName, userID, sessionID), reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := httputil.CheckStatus(resp, http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}

	var session models.Session
	if err := httputil.DecodeJSON(resp, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSession retrieves an existing session.
// Returns nil, nil if the session is not found.
func (c *Client) GetSession(ctx context.Context, appName, userID, sessionID string) (*models.Session, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, c.sessionURL(appName, userID, sessionID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err := httputil.CheckStatus(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var session models.Session
	if err := httputil.DecodeJSON(resp, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSession deletes an existing session.
func (c *Client) DeleteSession(ctx context.Context, appName, userID, sessionID string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, c.sessionURL(appName, userID, sessionID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return httputil.CheckStatus(resp, http.StatusOK, http.StatusNoContent)
}

// Run executes an agent with the given request.
// It does NOT auto-create sessions. Use RunWithAutoSession for that.
func (c *Client) Run(ctx context.Context, runReq models.RunAgentRequest) (*RunResult, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, c.baseURL+"/run", runReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read body first since we need it for both error and success cases
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var events []models.Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &RunResult{
		Events:    events,
		FinalText: extractFinalText(events),
	}, nil
}

// RunWithAutoSession executes an agent, automatically creating a session if it doesn't exist.
// This is the recommended method for most use cases.
func (c *Client) RunWithAutoSession(ctx context.Context, runReq models.RunAgentRequest) (*RunResult, error) {
	result, err := c.Run(ctx, runReq)
	if err == nil {
		return result, nil
	}

	// Check if the error is due to session not found
	errStr := err.Error()
	if !strings.Contains(errStr, "session") || !strings.Contains(errStr, "not found") {
		fmt.Printf("Run failed with non-session error: %v\n", err)
		return nil, err
	}

	// Create the session and retry
	if _, createErr := c.CreateSession(ctx, runReq.AppName, runReq.UserId, runReq.SessionId, nil); createErr != nil {
		return nil, fmt.Errorf("failed to create session: %w (original error: %v)", createErr, err)
	}

	return c.Run(ctx, runReq)
}

// extractFinalText extracts the final text response from a list of events.
// It looks for the last event with model content that has text parts.
func extractFinalText(events []models.Event) string {
	var lastText string

	for _, event := range events {
		if event.Content == nil || event.Content.Role != string(genai.RoleModel) {
			continue
		}
		for _, part := range event.Content.Parts {
			if part.Text != "" {
				lastText = part.Text
			}
		}
	}

	return lastText
}
