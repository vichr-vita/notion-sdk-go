package notion

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsersListBuildsPaginationQuery(t *testing.T) {
	client := newUsersTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/users?page_size=10&start_cursor=cursor_a", req.URL.String())
		require.Equal(t, "Bearer secret_test", req.Header.Get("Authorization"))
		require.Equal(t, "2026-03-11", req.Header.Get("Notion-Version"))
		require.Equal(t, "application/json", req.Header.Get("Accept"))
		require.Empty(t, req.Header.Get("Content-Type"))

		return `{
			"object": "list",
			"type": "user",
			"results": [
				{
					"object": "user",
					"id": "user_a",
					"name": "Ada",
					"avatar_url": null,
					"type": "person",
					"person": {"email": "ada@example.com"}
				}
			],
			"has_more": false,
			"next_cursor": null
		}`
	})

	resp, err := client.Users.List(context.Background(), &Pagination{
		StartCursor: "cursor_a",
		PageSize:    10,
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "user", resp.Type)
	require.False(t, resp.HasMore)
	require.Nil(t, resp.NextCursor)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "user_a", resp.Results[0].ID)
	require.Equal(t, "Ada", *resp.Results[0].Name)
	require.Equal(t, "ada@example.com", resp.Results[0].Person.Email)
}

func TestUsersGetBuildsPath(t *testing.T) {
	client := newUsersTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/users/user_a", req.URL.String())

		return `{
			"object": "user",
			"id": "user_a",
			"name": null,
			"avatar_url": null,
			"type": "person",
			"person": {}
		}`
	})

	user, err := client.Users.Get(context.Background(), "user_a")
	require.NoError(t, err)
	require.Equal(t, "user_a", user.ID)
	require.Nil(t, user.Name)
	require.NotNil(t, user.Person)
}

func TestUsersMeBuildsPath(t *testing.T) {
	client := newUsersTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/users/me", req.URL.String())

		return `{
			"object": "user",
			"id": "bot_a",
			"name": "Integration",
			"avatar_url": null,
			"type": "bot",
			"bot": {
				"owner": {"type": "workspace", "workspace": true},
				"workspace_name": "Workspace"
			}
		}`
	})

	user, err := client.Users.Me(context.Background())
	require.NoError(t, err)
	require.Equal(t, "bot_a", user.ID)
	require.Equal(t, "bot", user.Type)
	require.NotNil(t, user.Bot)
	require.JSONEq(t, `{"type":"workspace","workspace":true}`, string(user.Bot.Owner))
	require.Equal(t, "Workspace", *user.Bot.WorkspaceName)
}

func TestUsersReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     http.Header{"X-Request-Id": []string{"req_user"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":403,"code":"restricted_resource","message":"missing user capability"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	resp, err := client.Users.List(context.Background(), nil)
	require.Nil(t, resp)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	require.Equal(t, "restricted_resource", apiErr.Code)
	require.Equal(t, "missing user capability", apiErr.Message)
	require.Equal(t, "req_user", apiErr.RequestID)
}

func newUsersTestClient(t *testing.T, assert func(*http.Request) string) *Client {
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
