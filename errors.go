package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIError represents a non-2xx response from the Notion API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RequestID  string
	Body       string
}

func (e *APIError) Error() string {
	if e.Code != "" || e.Message != "" {
		return fmt.Sprintf("notion API error: status=%d code=%q message=%q", e.StatusCode, e.Code, e.Message)
	}

	return fmt.Sprintf("notion API error: status=%d body=%q", e.StatusCode, e.Body)
}

func newAPIError(resp *http.Response, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Request-Id"),
		Body:       bodySnippet(body),
	}

	var payload struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		apiErr.Code = payload.Code
		apiErr.Message = payload.Message
	}

	return apiErr
}

func bodySnippet(body []byte) string {
	const max = 1024

	snippet := strings.TrimSpace(string(body))
	if len(snippet) <= max {
		return snippet
	}

	return snippet[:max]
}
