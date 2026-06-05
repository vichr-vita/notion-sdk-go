package notion

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginationEncodesQueryFields(t *testing.T) {
	client := NewClient("secret_test", WithBaseURL("https://notion.test/v1"))

	req, err := client.newRequest(context.Background(), http.MethodGet, "blocks/block_a/children", &Pagination{
		StartCursor: "cursor_a",
		PageSize:    100,
	}, nil)
	require.NoError(t, err)

	require.Equal(t, "https://notion.test/v1/blocks/block_a/children?page_size=100&start_cursor=cursor_a", req.URL.String())
}

func TestPaginationOmitsZeroQueryFields(t *testing.T) {
	client := NewClient("secret_test", WithBaseURL("https://notion.test/v1"))

	req, err := client.newRequest(context.Background(), http.MethodGet, "blocks/block_a/children", &Pagination{}, nil)
	require.NoError(t, err)

	require.Equal(t, "https://notion.test/v1/blocks/block_a/children", req.URL.String())
}

func TestPaginatedResponseDecodesShell(t *testing.T) {
	data := []byte(`{
		"object": "list",
		"type": "block",
		"results": [{"id": "block_a"}],
		"has_more": true,
		"next_cursor": "cursor_b"
	}`)

	var resp PaginatedResponse[struct {
		ID string `json:"id"`
	}]
	require.NoError(t, json.Unmarshal(data, &resp))

	require.Equal(t, "list", resp.Object)
	require.Equal(t, "block", resp.Type)
	require.True(t, resp.HasMore)
	require.NotNil(t, resp.NextCursor)
	require.Equal(t, "cursor_b", *resp.NextCursor)
	require.Equal(t, []struct {
		ID string `json:"id"`
	}{{ID: "block_a"}}, resp.Results)
}

func TestPaginatedResponseDecodesNullCursor(t *testing.T) {
	data := []byte(`{
		"object": "list",
		"results": [],
		"has_more": false,
		"next_cursor": null
	}`)

	var resp PaginatedResponse[json.RawMessage]
	require.NoError(t, json.Unmarshal(data, &resp))

	require.False(t, resp.HasMore)
	require.Nil(t, resp.NextCursor)
	require.Empty(t, resp.Results)
}
