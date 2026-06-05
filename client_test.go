package notion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClientDefaultConfig(t *testing.T) {
	client := NewClient("secret_test")

	require.Equal(t, "secret_test", client.config.token)
	require.Equal(t, "https://api.notion.com/v1", client.config.baseURL)
	require.Equal(t, "2026-03-11", client.config.version)
	require.Same(t, http.DefaultClient, client.config.httpClient)

	require.NotNil(t, client.Pages)
	require.Same(t, client, client.Pages.client)
	require.NotNil(t, client.Blocks)
	require.Same(t, client, client.Blocks.client)
	require.NotNil(t, client.DataSources)
	require.Same(t, client, client.DataSources.client)
	require.NotNil(t, client.Databases)
	require.Same(t, client, client.Databases.client)
	require.NotNil(t, client.Search)
	require.Same(t, client, client.Search.client)
	require.NotNil(t, client.Users)
	require.Same(t, client, client.Users.client)
	require.NotNil(t, client.Comments)
	require.Same(t, client, client.Comments.client)
	require.NotNil(t, client.FileUploads)
	require.Same(t, client, client.FileUploads.client)
	require.NotNil(t, client.CustomEmojis)
	require.Same(t, client, client.CustomEmojis.client)
}

func TestNewClientOptionOverrides(t *testing.T) {
	httpClient := &http.Client{}

	client := NewClient(
		"secret_test",
		WithHTTPClient(httpClient),
		WithBaseURL("https://notion.test/v1"),
		WithVersion("2025-09-03"),
	)

	require.Equal(t, "secret_test", client.config.token)
	require.Equal(t, "https://notion.test/v1", client.config.baseURL)
	require.Equal(t, "2025-09-03", client.config.version)
	require.Same(t, httpClient, client.config.httpClient)
}

func TestNewClientIgnoresEmptyOptionValues(t *testing.T) {
	client := NewClient(
		"secret_test",
		WithHTTPClient(nil),
		WithBaseURL(""),
		WithVersion(""),
	)

	require.Equal(t, "https://api.notion.com/v1", client.config.baseURL)
	require.Equal(t, "2026-03-11", client.config.version)
	require.Same(t, http.DefaultClient, client.config.httpClient)
}
