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

func TestDatabasesGetBuildsRequest(t *testing.T) {
	client := newDatabasesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/databases/database_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{"object":"database","id":"database_a","parent":{"type":"page_id","page_id":"page_a"},"properties":{"Name":{"type":"title"}},"archived":false,"in_trash":false}`
	})

	database, err := client.Databases.Get(context.Background(), "database_a")
	require.NoError(t, err)
	require.Equal(t, "database_a", database.ID)
	require.JSONEq(t, `{"type":"title"}`, string(database.Properties["Name"]))
}

func TestDatabasesCreateBuildsRequest(t *testing.T) {
	client := newDatabasesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/databases", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{"type": "page_id", "page_id": "page_a"}, body["parent"])
		require.Equal(t, []any{
			map[string]any{
				"type": "text",
				"text": map[string]any{"content": "Tasks"},
			},
		}, body["title"])
		require.Equal(t, true, body["is_inline"])
		require.Equal(t, map[string]any{
			"properties": map[string]any{
				"Name": map[string]any{"title": map[string]any{}},
			},
		}, body["initial_data_source"])

		return `{"object":"database","id":"database_b","parent":{"type":"page_id","page_id":"page_a"},"properties":{},"archived":false,"in_trash":false}`
	})

	isInline := true
	database, err := client.Databases.Create(context.Background(), CreateDatabaseRequest{
		Parent:   PageParent("page_a"),
		Title:    []RichTextObject{{Type: "text", Text: &TextContent{Content: "Tasks"}}},
		IsInline: &isInline,
		InitialDataSource: map[string]any{
			"properties": map[string]any{
				"Name": map[string]any{"title": map[string]any{}},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "database_b", database.ID)
}

func TestDatabasesUpdateBuildsRequest(t *testing.T) {
	client := newDatabasesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "https://notion.test/v1/databases/database_a", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{"type": "page_id", "page_id": "page_b"}, body["parent"])
		require.Equal(t, false, body["in_trash"])
		require.Equal(t, true, body["is_locked"])
		require.Equal(t, []any{
			map[string]any{
				"type": "text",
				"text": map[string]any{"content": "Updated"},
			},
		}, body["title"])

		return `{"object":"database","id":"database_a","parent":{"type":"page_id","page_id":"page_b"},"properties":{},"archived":false,"in_trash":false}`
	})

	inTrash := false
	isLocked := true
	parent := PageParent("page_b")
	database, err := client.Databases.Update(context.Background(), "database_a", UpdateDatabaseRequest{
		Parent:   &parent,
		Title:    []RichTextObject{{Type: "text", Text: &TextContent{Content: "Updated"}}},
		InTrash:  &inTrash,
		IsLocked: &isLocked,
	})
	require.NoError(t, err)
	require.Equal(t, "database_a", database.ID)
}

func TestDatabasesReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"X-Request-Id": []string{"req_database"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":404,"code":"object_not_found","message":"missing database"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	database, err := client.Databases.Get(context.Background(), "database_a")
	require.Nil(t, database)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	require.Equal(t, "object_not_found", apiErr.Code)
	require.Equal(t, "missing database", apiErr.Message)
	require.Equal(t, "req_database", apiErr.RequestID)
}

func newDatabasesTestClient(t *testing.T, check func(*http.Request) string) *Client {
	t.Helper()

	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body := check(req)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	return NewClient(
		"secret_test",
		WithBaseURL("https://notion.test/v1"),
		WithHTTPClient(httpClient),
	)
}
