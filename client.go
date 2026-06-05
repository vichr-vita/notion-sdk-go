package notion

import "net/http"

// Client is the root Notion API client.
type Client struct {
	httpClient *http.Client
}
