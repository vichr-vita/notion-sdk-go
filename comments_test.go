package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommentsListBuildsQuery(t *testing.T) {
	client := newCommentsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/comments?block_id=block_a&page_size=25&start_cursor=cursor_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{
			"object": "list",
			"type": "comment",
			"results": [
				{
					"object": "comment",
					"id": "comment_a",
					"parent": {"type": "block_id", "block_id": "block_a"},
					"discussion_id": "discussion_a",
					"created_time": "2026-06-05T12:00:00Z",
					"last_edited_time": null,
					"created_by": {"object": "user", "id": "user_a", "name": null, "avatar_url": null},
					"rich_text": [{"type": "text", "plain_text": "hello", "text": {"content": "hello"}}],
					"display_name": {"type": "user", "resolved_name": "Ada"},
					"attachments": [{"file": {"type": "external", "external": {"url": "https://example.com/file.txt"}}}]
				}
			],
			"has_more": false,
			"next_cursor": null
		}`
	})

	resp, err := client.Comments.List(context.Background(), ListCommentsRequest{
		BlockID:     "block_a",
		StartCursor: "cursor_a",
		PageSize:    25,
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "comment", resp.Type)
	require.False(t, resp.HasMore)
	require.Nil(t, resp.NextCursor)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "comment_a", resp.Results[0].ID)
	require.Equal(t, "block_a", resp.Results[0].Parent.BlockID)
	require.Equal(t, "discussion_a", resp.Results[0].DiscussionID)
	require.Equal(t, "hello", resp.Results[0].RichText[0].PlainText)
	require.Equal(t, "Ada", resp.Results[0].DisplayName.ResolvedName)
	require.Equal(t, "https://example.com/file.txt", resp.Results[0].Attachments[0].File.External.URL)
}

func TestCommentsCreateEncodesPageParent(t *testing.T) {
	client := newCommentsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/comments", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{
			"type":    "page_id",
			"page_id": "page_a",
		}, body["parent"])
		require.Equal(t, []any{
			map[string]any{
				"type": "text",
				"text": map[string]any{"content": "hello"},
			},
		}, body["rich_text"])
		require.NotContains(t, body, "discussion_id")
		require.NotContains(t, body, "markdown")

		return `{
			"object": "comment",
			"id": "comment_a",
			"parent": {"type": "page_id", "page_id": "page_a"},
			"discussion_id": "discussion_a",
			"created_time": "2026-06-05T12:00:00Z",
			"last_edited_time": null,
			"rich_text": [{"type": "text", "plain_text": "hello", "text": {"content": "hello"}}]
		}`
	})

	comment, err := client.Comments.Create(context.Background(), CreateCommentRequest{
		Parent:   ptr(PageParent("page_a")),
		RichText: []RichTextObject{{Type: "text", Text: &TextContent{Content: "hello"}}},
	})
	require.NoError(t, err)
	require.Equal(t, "comment_a", comment.ID)
	require.Equal(t, "page_a", comment.Parent.PageID)
	require.Equal(t, "discussion_a", comment.DiscussionID)
}

func TestCommentsCreateEncodesBlockParentMarkdown(t *testing.T) {
	client := newCommentsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, map[string]any{
			"type":     "block_id",
			"block_id": "block_a",
		}, body["parent"])
		require.Equal(t, "hello **there**", body["markdown"])
		require.NotContains(t, body, "rich_text")

		return `{"object":"comment","id":"comment_a"}`
	})

	comment, err := client.Comments.Create(context.Background(), CreateCommentRequest{
		Parent:   ptr(BlockParent("block_a")),
		Markdown: "hello **there**",
	})
	require.NoError(t, err)
	require.Equal(t, "comment_a", comment.ID)
}

func TestCommentsCreateEncodesDiscussionTarget(t *testing.T) {
	client := newCommentsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "discussion_a", body["discussion_id"])
		require.Equal(t, "hello", body["markdown"])
		require.NotContains(t, body, "parent")

		return `{"object":"comment","id":"comment_a","discussion_id":"discussion_a"}`
	})

	comment, err := client.Comments.Create(context.Background(), CreateCommentRequest{
		DiscussionID: "discussion_a",
		Markdown:     "hello",
	})
	require.NoError(t, err)
	require.Equal(t, "comment_a", comment.ID)
	require.Equal(t, "discussion_a", comment.DiscussionID)
}

func TestCommentsCreateEncodesAttachmentsAndDisplayName(t *testing.T) {
	client := newCommentsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, []any{
			map[string]any{
				"file": map[string]any{
					"type":     "external",
					"external": map[string]any{"url": "https://example.com/file.txt"},
				},
			},
		}, body["attachments"])
		require.Equal(t, map[string]any{
			"type": "custom",
			"name": "Build Bot",
		}, body["display_name"])

		return `{"object":"comment","id":"comment_a"}`
	})

	comment, err := client.Comments.Create(context.Background(), CreateCommentRequest{
		Parent:   ptr(PageParent("page_a")),
		Markdown: "hello",
		Attachments: []CommentAttachment{{
			File: File{
				Type:     "external",
				External: &ExternalFile{URL: "https://example.com/file.txt"},
			},
		}},
		DisplayName: map[string]any{
			"type": "custom",
			"name": "Build Bot",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "comment_a", comment.ID)
}

func TestCommentsReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     http.Header{"X-Request-Id": []string{"req_comment"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":403,"code":"restricted_resource","message":"missing comment capability"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	resp, err := client.Comments.List(context.Background(), ListCommentsRequest{BlockID: "block_a"})
	require.Nil(t, resp)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	require.Equal(t, "restricted_resource", apiErr.Code)
	require.Equal(t, "missing comment capability", apiErr.Message)
	require.Equal(t, "req_comment", apiErr.RequestID)
}

func newCommentsTestClient(t *testing.T, check func(*http.Request) string) *Client {
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

func ptr[T any](value T) *T {
	return &value
}
