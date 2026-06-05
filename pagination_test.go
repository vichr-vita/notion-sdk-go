package notion

import (
	"context"
	"encoding/json"
	"errors"
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

func TestForEachPaginatedStopsOnHasMoreFalse(t *testing.T) {
	cursorB := "cursor_b"
	calls := 0
	var seen []string

	err := ForEachPaginated(context.Background(), &Pagination{PageSize: 2}, func(_ context.Context, pagination *Pagination) (*PaginatedResponse[string], error) {
		calls++

		switch calls {
		case 1:
			require.Equal(t, 2, pagination.PageSize)
			require.Empty(t, pagination.StartCursor)
			return &PaginatedResponse[string]{
				Results:    []string{"a", "b"},
				HasMore:    true,
				NextCursor: &cursorB,
			}, nil
		case 2:
			require.Equal(t, 2, pagination.PageSize)
			require.Equal(t, "cursor_b", pagination.StartCursor)
			return &PaginatedResponse[string]{
				Results: []string{"c"},
			}, nil
		default:
			t.Fatalf("unexpected page call %d", calls)
			return nil, nil
		}
	}, func(value string) error {
		seen = append(seen, value)
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, 2, calls)
	require.Equal(t, []string{"a", "b", "c"}, seen)
}

func TestForEachPaginatedReturnsCallbackError(t *testing.T) {
	wantErr := errors.New("stop")
	calls := 0

	err := ForEachPaginated(context.Background(), nil, func(_ context.Context, _ *Pagination) (*PaginatedResponse[string], error) {
		calls++
		return &PaginatedResponse[string]{
			Results: []string{"a", "b"},
		}, nil
	}, func(value string) error {
		if value == "b" {
			return wantErr
		}
		return nil
	})

	require.ErrorIs(t, err, wantErr)
	require.Equal(t, 1, calls)
}
