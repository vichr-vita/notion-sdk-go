package notion

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoParsesStructuredNotionError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"X-Request-Id": []string{"req_a"},
			},
			Body: io.NopCloser(strings.NewReader(`{"object":"error","status":400,"code":"validation_error","message":"Name is required"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodPost, "pages", nil, map[string]any{})
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	require.Equal(t, "validation_error", apiErr.Code)
	require.Equal(t, "Name is required", apiErr.Message)
	require.Equal(t, "req_a", apiErr.RequestID)
	require.Equal(t, `{"object":"error","status":400,"code":"validation_error","message":"Name is required"}`, apiErr.Body)
	require.Contains(t, apiErr.Error(), "validation_error")
}

func TestDoHandlesMalformedErrorBodyFallback(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("not json")),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", nil, nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	require.Empty(t, apiErr.Code)
	require.Empty(t, apiErr.Message)
	require.Equal(t, "not json", apiErr.Body)
	require.Contains(t, apiErr.Error(), `body="not json"`)
}

func TestDoPreservesNetworkErrors(t *testing.T) {
	netErr := &net.DNSError{Err: "lookup failed", Name: "api.notion.test"}
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, netErr
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodGet, "search", nil, nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.ErrorIs(t, err, netErr)

	var apiErr *APIError
	require.False(t, errors.As(err, &apiErr))
}

func TestDoPreservesJSONDecodeErrors(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"id":`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	req, err := client.newRequest(context.Background(), http.MethodGet, "pages/page_a", nil, nil)
	require.NoError(t, err)

	var out map[string]any
	err = client.do(req, &out)
	require.Error(t, err)

	var apiErr *APIError
	require.False(t, errors.As(err, &apiErr))
}
