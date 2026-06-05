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

func TestPagesGetBuildsRequest(t *testing.T) {
	client := newPagesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/pages/page_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{"object":"page","id":"page_a","parent":{"type":"page_id","page_id":"parent_a"},"properties":{"Name":{"type":"title"}},"archived":false,"in_trash":false}`
	})

	page, err := client.Pages.Get(context.Background(), "page_a")
	require.NoError(t, err)
	require.Equal(t, "page_a", page.ID)
	require.JSONEq(t, `{"type":"title"}`, string(page.Properties["Name"]))
}

func TestPagesCreateBuildsRequest(t *testing.T) {
	client := newPagesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/pages", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{
			"type":           "data_source_id",
			"data_source_id": "data_source_a",
		}, body["parent"])
		require.Equal(t, map[string]any{
			"title": []any{
				map[string]any{
					"type": "text",
					"text": map[string]any{"content": "Task"},
				},
			},
		}, body["properties"].(map[string]any)["Name"])
		require.Equal(t, []any{
			map[string]any{
				"type": "paragraph",
				"paragraph": map[string]any{
					"rich_text": []any{
						map[string]any{
							"type": "text",
							"text": map[string]any{"content": "Body"},
						},
					},
				},
			},
		}, body["children"])

		return `{"object":"page","id":"page_b","parent":{"type":"data_source_id","data_source_id":"data_source_a"},"properties":{},"archived":false,"in_trash":false}`
	})

	page, err := client.Pages.Create(context.Background(), CreatePageRequest{
		Parent: DataSourceParent("data_source_a"),
		Properties: map[string]any{
			"Name": Title("Task"),
		},
		Children: []BlockRequest{
			{
				"type": "paragraph",
				"paragraph": map[string]any{
					"rich_text": []RichTextObject{
						{
							Type: "text",
							Text: &TextContent{Content: "Body"},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "page_b", page.ID)
}

func TestPagesUpdateBuildsRequest(t *testing.T) {
	client := newPagesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "https://notion.test/v1/pages/page_a", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, false, body["in_trash"])
		require.Equal(t, true, body["is_archived"])
		require.Equal(t, map[string]any{
			"rich_text": []any{
				map[string]any{
					"type": "text",
					"text": map[string]any{"content": "Updated"},
				},
			},
		}, body["properties"].(map[string]any)["Notes"])

		return `{"object":"page","id":"page_a","parent":{"type":"page_id","page_id":"parent_a"},"properties":{},"archived":false,"in_trash":false}`
	})

	inTrash := false
	isArchived := true
	page, err := client.Pages.Update(context.Background(), "page_a", UpdatePageRequest{
		InTrash:    &inTrash,
		IsArchived: &isArchived,
		Properties: map[string]any{
			"Notes": RichText("Updated"),
		},
	})
	require.NoError(t, err)
	require.Equal(t, "page_a", page.ID)
}

func TestPagesGetPropertyItemBuildsRequest(t *testing.T) {
	client := newPagesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/pages/page_a/properties/prop_a?page_size=25&start_cursor=cursor_a", req.URL.String())

		return `{"object":"property_item","id":"prop_a","type":"title","title":{"plain_text":"Task"}}`
	})

	item, err := client.Pages.GetPropertyItem(context.Background(), "page_a", "prop_a", &Pagination{
		StartCursor: "cursor_a",
		PageSize:    25,
	})
	require.NoError(t, err)
	require.Equal(t, "property_item", item.Object)
	require.Equal(t, "title", item.Type)
	require.JSONEq(t, `{"object":"property_item","id":"prop_a","type":"title","title":{"plain_text":"Task"}}`, string(item.Raw))
}

func TestPagesGetMarkdownBuildsRequest(t *testing.T) {
	client := newPagesTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/pages/page_a/markdown", req.URL.String())
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{"object":"page_markdown","id":"page_a","markdown":"# Title\nBody","truncated":true,"unknown_block_ids":["block_a"]}`
	})

	markdown, err := client.Pages.GetMarkdown(context.Background(), "page_a")
	require.NoError(t, err)
	require.Equal(t, "page_a", markdown.ID)
	require.Equal(t, "# Title\nBody", markdown.Markdown)
	require.True(t, markdown.Truncated)
	require.Equal(t, []string{"block_a"}, markdown.UnknownBlockIDs)
	require.JSONEq(t, `{"object":"page_markdown","id":"page_a","markdown":"# Title\nBody","truncated":true,"unknown_block_ids":["block_a"]}`, string(markdown.Raw))
}

func TestPagesReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"X-Request-Id": []string{"req_page"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":404,"code":"object_not_found","message":"missing page"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	page, err := client.Pages.Get(context.Background(), "page_a")
	require.Nil(t, page)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	require.Equal(t, "object_not_found", apiErr.Code)
	require.Equal(t, "missing page", apiErr.Message)
	require.Equal(t, "req_page", apiErr.RequestID)
}

func newPagesTestClient(t *testing.T, check func(*http.Request) string) *Client {
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
