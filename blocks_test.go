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

func TestBlocksGetBuildsRequest(t *testing.T) {
	client := newBlocksTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/blocks/block_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{"object":"block","id":"block_a","type":"paragraph","has_children":false,"archived":false,"in_trash":false}`
	})

	block, err := client.Blocks.Get(context.Background(), "block_a")
	require.NoError(t, err)
	require.Equal(t, "block_a", block.ID)
}

func TestBlocksUpdateBuildsRequest(t *testing.T) {
	client := newBlocksTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "https://notion.test/v1/blocks/block_a", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, false, body["in_trash"])
		require.Equal(t, map[string]any{
			"rich_text": []any{
				map[string]any{
					"type": "text",
					"text": map[string]any{"content": "updated"},
				},
			},
		}, body["paragraph"])

		return `{"object":"block","id":"block_a","type":"paragraph","paragraph":{"rich_text":[]},"has_children":false,"archived":false,"in_trash":false}`
	})

	block, err := client.Blocks.Update(context.Background(), "block_a", UpdateBlockRequest{
		"in_trash": false,
		"paragraph": map[string]any{
			"rich_text": []RichTextObject{
				{
					Type: "text",
					Text: &TextContent{Content: "updated"},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "block_a", block.ID)
}

func TestBlocksDeleteBuildsRequest(t *testing.T) {
	client := newBlocksTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodDelete, req.Method)
		require.Equal(t, "https://notion.test/v1/blocks/block_a", req.URL.String())
		require.Nil(t, req.Body)

		return `{"object":"block","id":"block_a","type":"paragraph","has_children":false,"archived":true,"in_trash":true}`
	})

	block, err := client.Blocks.Delete(context.Background(), "block_a")
	require.NoError(t, err)
	require.Equal(t, "block_a", block.ID)
	require.True(t, block.Archived)
	require.True(t, block.InTrash)
}

func TestBlocksListChildrenBuildsRequest(t *testing.T) {
	client := newBlocksTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/blocks/block_a/children?page_size=50&start_cursor=cursor_a", req.URL.String())

		return `{
			"object":"list",
			"type":"block",
			"results":[{"object":"block","id":"child_a","type":"paragraph","paragraph":{"rich_text":[]},"has_children":false,"archived":false,"in_trash":false}],
			"has_more":false,
			"next_cursor":null
		}`
	})

	resp, err := client.Blocks.ListChildren(context.Background(), "block_a", &Pagination{
		StartCursor: "cursor_a",
		PageSize:    50,
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "block", resp.Type)
	require.False(t, resp.HasMore)
	require.Nil(t, resp.NextCursor)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "child_a", resp.Results[0].ID)
	require.JSONEq(t, `{"rich_text":[]}`, string(resp.Results[0].Content))
}

func TestBlocksAppendChildrenBuildsRequest(t *testing.T) {
	client := newBlocksTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "https://notion.test/v1/blocks/block_a/children", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, []any{
			map[string]any{
				"type": "paragraph",
				"paragraph": map[string]any{
					"rich_text": []any{
						map[string]any{
							"type": "text",
							"text": map[string]any{"content": "hello"},
						},
					},
				},
			},
		}, body["children"])
		require.Equal(t, map[string]any{
			"type":        "after_block",
			"after_block": map[string]any{"id": "child_a"},
		}, body["position"])

		return `{
			"object":"list",
			"type":"block",
			"results":[{"object":"block","id":"child_b","type":"paragraph","has_children":false,"archived":false,"in_trash":false}],
			"has_more":false,
			"next_cursor":null
		}`
	})

	resp, err := client.Blocks.AppendChildren(context.Background(), "block_a", AppendBlockChildrenRequest{
		Children: []BlockRequest{
			{
				"type": "paragraph",
				"paragraph": map[string]any{
					"rich_text": []RichTextObject{
						{
							Type: "text",
							Text: &TextContent{Content: "hello"},
						},
					},
				},
			},
		},
		Position: &BlockPosition{
			Type:       "after_block",
			AfterBlock: &AfterBlockPosition{ID: "child_a"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "child_b", resp.Results[0].ID)
}

func TestBlocksForEachChildPaginates(t *testing.T) {
	cursorB := "cursor_b"
	calls := 0
	client := newBlocksTestClient(t, func(req *http.Request) string {
		calls++
		require.Equal(t, http.MethodGet, req.Method)
		switch calls {
		case 1:
			require.Equal(t, "https://notion.test/v1/blocks/block_a/children?page_size=1", req.URL.String())
			return `{
				"object":"list",
				"type":"block",
				"results":[{"object":"block","id":"child_a","type":"paragraph","has_children":false,"archived":false,"in_trash":false}],
				"has_more":true,
				"next_cursor":"cursor_b"
			}`
		case 2:
			require.Equal(t, "https://notion.test/v1/blocks/block_a/children?page_size=1&start_cursor=cursor_b", req.URL.String())
			return `{
				"object":"list",
				"type":"block",
				"results":[{"object":"block","id":"child_b","type":"paragraph","has_children":false,"archived":false,"in_trash":false}],
				"has_more":false,
				"next_cursor":null
			}`
		default:
			t.Fatalf("unexpected request %d", calls)
			return ""
		}
	})

	var seen []string
	err := client.Blocks.ForEachChild(context.Background(), "block_a", &Pagination{PageSize: 1}, func(block Block) error {
		seen = append(seen, block.ID)
		if block.ID == cursorB {
			return errors.New("unreachable")
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"child_a", "child_b"}, seen)
	require.Equal(t, 2, calls)
}

func TestBlocksReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"X-Request-Id": []string{"req_block"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":404,"code":"object_not_found","message":"missing block"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	block, err := client.Blocks.Get(context.Background(), "block_a")
	require.Nil(t, block)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	require.Equal(t, "object_not_found", apiErr.Code)
	require.Equal(t, "missing block", apiErr.Message)
	require.Equal(t, "req_block", apiErr.RequestID)
}

func newBlocksTestClient(t *testing.T, check func(*http.Request) string) *Client {
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
