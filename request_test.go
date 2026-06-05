package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRequestBuildsMethodPathQueryHeadersAndBody(t *testing.T) {
	client := NewClient(
		"secret_test",
		WithBaseURL("https://notion.test/v1/"),
		WithVersion("2026-03-11"),
	)

	req, err := client.newRequest(
		context.Background(),
		http.MethodPost,
		"/search",
		struct {
			StartCursor string `url:"start_cursor,omitempty"`
			PageSize    int    `url:"page_size,omitempty"`
			Empty       string `url:"empty,omitempty"`
		}{
			StartCursor: "cursor_a",
			PageSize:    25,
		},
		map[string]any{"query": "Tasks"},
	)
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, "https://notion.test/v1/search?page_size=25&start_cursor=cursor_a", req.URL.String())
	require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
	require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
	require.Equal(t, "Tasks", body["query"])
}

func TestNewRequestSupportsURLValuesQuery(t *testing.T) {
	client := NewClient("secret_test", WithBaseURL("https://notion.test/v1"))

	query := url.Values{}
	query.Set("filter", "page")

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", query, nil)
	require.NoError(t, err)

	require.Equal(t, "https://notion.test/v1/search?filter=page", req.URL.String())
	require.Empty(t, req.Header.Get("Content-Type"))
}

func TestDoDecodesJSONResponse(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodGet, req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"object":"page","id":"page_a"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodGet, "pages/page_a", nil, nil)
	require.NoError(t, err)

	var out struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	}
	require.NoError(t, client.do(req, &out))
	require.Equal(t, "page", out.Object)
	require.Equal(t, "page_a", out.ID)
}

func TestDoHandlesEmptyBody(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodDelete, "blocks/block_a", nil, nil)
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, client.do(req, &out))
	require.Nil(t, out)
}

func TestDoRunsHooksWithRequestAndResponse(t *testing.T) {
	var requestHookCalled bool
	var responseHookCalled bool

	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.True(t, requestHookCalled)
		require.Equal(t, "hooked", req.Header.Get("X-Test-Hook"))

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"X-Response-Hook": []string{"seen"}},
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		}, nil
	})}
	client := NewClient(
		"secret_test",
		WithHTTPClient(httpClient),
		WithRequestHook(func(req *http.Request) error {
			requestHookCalled = true
			require.Equal(t, "https://api.notion.com/v1/search", req.URL.String())
			req.Header.Set("X-Test-Hook", "hooked")
			return nil
		}),
		WithResponseHook(func(resp *http.Response) error {
			responseHookCalled = true
			require.Equal(t, "seen", resp.Header.Get("X-Response-Hook"))
			return nil
		}),
	)

	req, err := client.newRequest(context.Background(), http.MethodPost, "search", nil, map[string]any{})
	require.NoError(t, err)

	require.NoError(t, client.do(req, nil))
	require.True(t, requestHookCalled)
	require.True(t, responseHookCalled)
}

func TestDoRetriesRetryableStatus(t *testing.T) {
	attempts := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"code":"service_unavailable","message":"try again"}`)),
		}, nil
	})}
	client := NewClient(
		"secret_test",
		WithHTTPClient(httpClient),
		WithRetryConfig(RetryConfig{MaxRetries: 2, Delay: time.Millisecond, MaxDelay: time.Millisecond, Jitter: -1}),
	)
	client.config.retrySleep = func(time.Duration) {}

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", nil, nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)
	require.Equal(t, 3, attempts)
}

func TestDoDoesNotRetryNonRetryableStatus(t *testing.T) {
	attempts := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"code":"invalid_request","message":"bad request"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", nil, nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)
	require.Equal(t, 1, attempts)
}

func TestDoRetriesUntilSuccess(t *testing.T) {
	attempts := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`temporary`)),
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"object":"page","id":"page_a"}`)),
		}, nil
	})}
	client := NewClient(
		"secret_test",
		WithHTTPClient(httpClient),
		WithRetryConfig(RetryConfig{MaxRetries: 2, Delay: time.Millisecond, MaxDelay: time.Millisecond, Jitter: -1}),
	)
	client.config.retrySleep = func(time.Duration) {}

	req, err := client.newRequest(context.Background(), http.MethodPost, "search", nil, map[string]any{"query": "Tasks"})
	require.NoError(t, err)

	var out struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	}
	require.NoError(t, client.do(req, &out))
	require.Equal(t, 2, attempts)
	require.Equal(t, "page", out.Object)
	require.Equal(t, "page_a", out.ID)
}

func TestDoHonorsRetryAfter(t *testing.T) {
	attempts := 0
	var sleeps []time.Duration
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"2"}},
				Body:       io.NopCloser(strings.NewReader(`rate limited`)),
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))
	client.config.retrySleep = func(delay time.Duration) {
		sleeps = append(sleeps, delay)
	}

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", nil, nil)
	require.NoError(t, err)

	require.NoError(t, client.do(req, nil))
	require.Equal(t, 2, attempts)
	require.Equal(t, []time.Duration{2 * time.Second}, sleeps)
}

func TestDoDoesNotRetryNonReplayableBody(t *testing.T) {
	attempts := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`temporary`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))
	client.config.retrySleep = func(time.Duration) {
		t.Fatal("retry sleep should not run")
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"https://api.notion.com/v1/search",
		io.NopCloser(strings.NewReader(`{"query":"Tasks"}`)),
	)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)
	require.Equal(t, 1, attempts)
}
