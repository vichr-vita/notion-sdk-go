package notion

import (
	"context"
	"encoding/json"
	"net/http"
)

// SearchService handles search endpoints.
type SearchService struct {
	client *Client
}

// SearchRequest is the body for the search endpoint.
type SearchRequest struct {
	Query       string        `json:"query,omitempty"`
	Sort        *SearchSort   `json:"sort,omitempty"`
	Filter      *SearchFilter `json:"filter,omitempty"`
	StartCursor string        `json:"start_cursor,omitempty"`
	PageSize    int           `json:"page_size,omitempty"`
}

// SearchSort controls search result ordering.
type SearchSort struct {
	Timestamp string `json:"timestamp,omitempty"`
	Direction string `json:"direction,omitempty"`
}

// SearchFilter limits search results by object type.
type SearchFilter struct {
	Property string `json:"property,omitempty"`
	Value    string `json:"value,omitempty"`
}

// SearchResult is a raw page or data source result from search.
type SearchResult struct {
	Object string          `json:"object"`
	ID     string          `json:"id"`
	Raw    json.RawMessage `json:"-"`
}

// SearchResponse is the paginated response from search.
type SearchResponse struct {
	Object           string          `json:"object"`
	Type             string          `json:"type,omitempty"`
	PageOrDataSource json.RawMessage `json:"page_or_data_source,omitempty"`
	Results          []SearchResult  `json:"results"`
	HasMore          bool            `json:"has_more"`
	NextCursor       *string         `json:"next_cursor"`
	RequestStatus    json.RawMessage `json:"request_status,omitempty"`
}

// Query searches pages and data sources by title.
func (s *SearchService) Query(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "search", nil, request)
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *SearchResult) UnmarshalJSON(data []byte) error {
	type searchResultAlias SearchResult
	var v searchResultAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*r = SearchResult(v)
	return nil
}
