package notion

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomEmojisListBuildsPaginationQuery(t *testing.T) {
	client := newCustomEmojisTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/custom_emojis?page_size=25&start_cursor=cursor_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{
			"object": "list",
			"type": "custom_emoji",
			"results": [
				{
					"id": "emoji_a",
					"name": "ship",
					"url": "https://notion.test/emoji.png"
				}
			],
			"has_more": false,
			"next_cursor": null
		}`
	})

	resp, err := client.CustomEmojis.List(context.Background(), &Pagination{
		StartCursor: "cursor_a",
		PageSize:    25,
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "custom_emoji", resp.Type)
	require.False(t, resp.HasMore)
	require.Nil(t, resp.NextCursor)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "emoji_a", resp.Results[0].ID)
	require.Equal(t, "ship", resp.Results[0].Name)
	require.Equal(t, "https://notion.test/emoji.png", resp.Results[0].URL)
	require.JSONEq(t, `{"id":"emoji_a","name":"ship","url":"https://notion.test/emoji.png"}`, string(resp.Results[0].Raw))
}

func TestCustomEmojisReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     http.Header{"X-Request-Id": []string{"req_emoji"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":403,"code":"restricted_resource","message":"missing emoji capability"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	resp, err := client.CustomEmojis.List(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	require.Equal(t, "restricted_resource", apiErr.Code)
	require.Equal(t, "missing emoji capability", apiErr.Message)
	require.Equal(t, "req_emoji", apiErr.RequestID)
}

func newCustomEmojisTestClient(t *testing.T, assert func(*http.Request) string) *Client {
	t.Helper()

	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body := assert(req)
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
