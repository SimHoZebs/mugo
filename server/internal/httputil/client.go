package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DoRequest executes an HTTP request and returns the response.
// It handles request creation, JSON marshaling of body, and basic error wrapping.
func DoRequest(ctx context.Context, client *http.Client, method, url string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// CheckStatus verifies the response status code is one of the accepted statuses.
// If not, it reads the body and returns an error with the status and body content.
func CheckStatus(resp *http.Response, acceptedStatuses ...int) error {
	for _, status := range acceptedStatuses {
		if resp.StatusCode == status {
			return nil
		}
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
}

// DecodeJSON decodes the response body into the result.
func DecodeJSON(resp *http.Response, result any) error {
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}
