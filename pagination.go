package notion

// Pagination controls cursor-based list endpoints.
type Pagination struct {
	StartCursor string `json:"start_cursor,omitempty" url:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty" url:"page_size,omitempty"`
}

// PaginatedResponse is the common response envelope for Notion list endpoints.
type PaginatedResponse[T any] struct {
	Object     string  `json:"object"`
	Type       string  `json:"type,omitempty"`
	Results    []T     `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}
