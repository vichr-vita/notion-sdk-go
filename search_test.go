package notion

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchQueryBuildsRequest(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/search", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "meeting notes", body["query"])
		require.Equal(t, float64(25), body["page_size"])
		require.Equal(t, "cursor_a", body["start_cursor"])
		require.Equal(t, map[string]any{
			"property": "object",
			"value":    "page",
		}, body["filter"])
		require.Equal(t, map[string]any{
			"timestamp": "last_edited_time",
			"direction": "descending",
		}, body["sort"])

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"object":"list","results":[],"has_more":false,"next_cursor":null}`)),
		}, nil
	})}
	client := NewClient(
		"secret_test",
		WithBaseURL("https://notion.test/v1"),
		WithHTTPClient(httpClient),
	)

	resp, err := client.Search.Query(context.Background(), SearchRequest{
		Query:       "meeting notes",
		StartCursor: "cursor_a",
		PageSize:    25,
		Filter:      &SearchFilter{Property: "object", Value: "page"},
		Sort:        &SearchSort{Timestamp: "last_edited_time", Direction: "descending"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "list", resp.Object)
}

func TestSearchQueryDecodesResponse(t *testing.T) {
	cursor := "cursor_b"
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"object": "list",
				"type": "page_or_data_source",
				"page_or_data_source": {},
				"results": [
					{"object": "page", "id": "page_a", "properties": {"Name": {"id": "title"}}},
					{"object": "data_source", "id": "data_source_a", "title": []}
				],
				"has_more": true,
				"next_cursor": "cursor_b",
				"request_status": {}
			}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	resp, err := client.Search.Query(context.Background(), SearchRequest{Query: "Tasks"})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "page_or_data_source", resp.Type)
	require.True(t, resp.HasMore)
	require.Equal(t, &cursor, resp.NextCursor)
	require.Len(t, resp.Results, 2)
	require.Equal(t, "page", resp.Results[0].Object)
	require.Equal(t, "page_a", resp.Results[0].ID)
	require.JSONEq(t, `{"object":"page","id":"page_a","properties":{"Name":{"id":"title"}}}`, string(resp.Results[0].Raw))
	require.JSONEq(t, `{}`, string(resp.PageOrDataSource))
	require.JSONEq(t, `{}`, string(resp.RequestStatus))
}

func TestSearchQueryReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"X-Request-Id": []string{"req_search"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":401,"code":"unauthorized","message":"bad token"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	resp, err := client.Search.Query(context.Background(), SearchRequest{Query: "Tasks"})
	require.Nil(t, resp)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	require.Equal(t, "unauthorized", apiErr.Code)
	require.Equal(t, "bad token", apiErr.Message)
	require.Equal(t, "req_search", apiErr.RequestID)
}
