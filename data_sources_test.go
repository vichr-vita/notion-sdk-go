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

func TestDataSourcesGetBuildsRequest(t *testing.T) {
	client := newDataSourcesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/data_sources/data_source_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{"object":"data_source","id":"data_source_a","parent":{"type":"database_id","database_id":"database_a"},"properties":{"Name":{"type":"title"}},"archived":false,"in_trash":false}`
	})

	dataSource, err := client.DataSources.Get(context.Background(), "data_source_a")
	require.NoError(t, err)
	require.Equal(t, "data_source_a", dataSource.ID)
	require.JSONEq(t, `{"type":"title"}`, string(dataSource.Properties["Name"]))
}

func TestDataSourcesQueryBuildsRequest(t *testing.T) {
	client := newDataSourcesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/data_sources/data_source_a/query?filter_properties%5B%5D=Name&filter_properties%5B%5D=Status", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "cursor_a", body["start_cursor"])
		require.Equal(t, float64(25), body["page_size"])
		require.Equal(t, "page", body["result_type"])
		require.Equal(t, map[string]any{
			"property": "Status",
			"select":   map[string]any{"equals": "Done"},
		}, body["filter"])
		require.Equal(t, []any{
			map[string]any{"property": "Created", "direction": "descending"},
		}, body["sorts"])
		require.NotContains(t, body, "filter_properties")

		return `{
			"object":"list",
			"type":"page_or_data_source",
			"page_or_data_source":{},
			"results":[{"object":"page","id":"page_a","properties":{"Name":{"type":"title"}}}],
			"has_more":true,
			"next_cursor":"cursor_b"
		}`
	})

	resp, err := client.DataSources.Query(context.Background(), "data_source_a", QueryDataSourceRequest{
		Filter: map[string]any{
			"property": "Status",
			"select":   map[string]any{"equals": "Done"},
		},
		Sorts: []any{
			map[string]any{"property": "Created", "direction": "descending"},
		},
		StartCursor:      "cursor_a",
		PageSize:         25,
		ResultType:       "page",
		FilterProperties: []string{"Name", "Status"},
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "page_or_data_source", resp.Type)
	require.True(t, resp.HasMore)
	require.NotNil(t, resp.NextCursor)
	require.Equal(t, "cursor_b", *resp.NextCursor)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "page_a", resp.Results[0].ID)
	require.JSONEq(t, `{"object":"page","id":"page_a","properties":{"Name":{"type":"title"}}}`, string(resp.Results[0].Raw))
}

func TestDataSourcesCreateBuildsRequest(t *testing.T) {
	client := newDataSourcesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/data_sources", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{"database_id": "database_a"}, body["parent"])
		require.Equal(t, []any{
			map[string]any{
				"type": "text",
				"text": map[string]any{"content": "Tasks"},
			},
		}, body["title"])
		require.Equal(t, map[string]any{"title": map[string]any{}}, body["properties"].(map[string]any)["Name"])

		return `{"object":"data_source","id":"data_source_b","parent":{"type":"database_id","database_id":"database_a"},"properties":{},"archived":false,"in_trash":false}`
	})

	dataSource, err := client.DataSources.Create(context.Background(), CreateDataSourceRequest{
		Parent: Parent{DatabaseID: "database_a"},
		Title:  []RichTextObject{{Type: "text", Text: &TextContent{Content: "Tasks"}}},
		Properties: map[string]any{
			"Name": map[string]any{"title": map[string]any{}},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "data_source_b", dataSource.ID)
}

func TestDataSourcesUpdateBuildsRequest(t *testing.T) {
	client := newDataSourcesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "https://notion.test/v1/data_sources/data_source_a", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, false, body["in_trash"])
		require.Equal(t, map[string]any{"database_id": "database_b"}, body["parent"])
		require.Equal(t, []any{
			map[string]any{
				"type": "text",
				"text": map[string]any{"content": "Updated"},
			},
		}, body["title"])
		require.Nil(t, body["properties"].(map[string]any)["Old"])

		return `{"object":"data_source","id":"data_source_a","parent":{"type":"database_id","database_id":"database_b"},"properties":{},"archived":false,"in_trash":false}`
	})

	inTrash := false
	dataSource, err := client.DataSources.Update(context.Background(), "data_source_a", UpdateDataSourceRequest{
		Parent:  &Parent{DatabaseID: "database_b"},
		Title:   []RichTextObject{{Type: "text", Text: &TextContent{Content: "Updated"}}},
		InTrash: &inTrash,
		Properties: map[string]any{
			"Old": nil,
		},
	})
	require.NoError(t, err)
	require.Equal(t, "data_source_a", dataSource.ID)
}

func TestDataSourcesReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"X-Request-Id": []string{"req_data_source"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":404,"code":"object_not_found","message":"missing data source"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	dataSource, err := client.DataSources.Get(context.Background(), "data_source_a")
	require.Nil(t, dataSource)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	require.Equal(t, "object_not_found", apiErr.Code)
	require.Equal(t, "missing data source", apiErr.Message)
	require.Equal(t, "req_data_source", apiErr.RequestID)
}

func newDataSourcesTestClient(t *testing.T, check func(*http.Request) string) *Client {
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
