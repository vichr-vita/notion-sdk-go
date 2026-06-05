package notion

import (
	"context"
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
