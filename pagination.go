package notion

// Pagination controls cursor-based list endpoints.
type Pagination struct {
	StartCursor string `json:"start_cursor,omitempty" url:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty" url:"page_size,omitempty"`
}
