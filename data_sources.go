package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DataSourcesService handles data source endpoints.
type DataSourcesService struct {
	client *Client
}

// CreateDataSourceRequest is the body for creating a data source.
type CreateDataSourceRequest struct {
	Parent      Parent           `json:"parent"`
	Title       []RichTextObject `json:"title,omitempty"`
	Properties  map[string]any   `json:"properties,omitempty"`
	Description []RichTextObject `json:"description,omitempty"`
	Icon        *Icon            `json:"icon,omitempty"`
}

// UpdateDataSourceRequest is the body for updating a data source.
type UpdateDataSourceRequest struct {
	Parent     *Parent          `json:"parent,omitempty"`
	Title      []RichTextObject `json:"title,omitempty"`
	Properties map[string]any   `json:"properties,omitempty"`
	Icon       *Icon            `json:"icon,omitempty"`
	InTrash    *bool            `json:"in_trash,omitempty"`
}

// QueryDataSourceRequest is the body and query parameters for querying a data source.
type QueryDataSourceRequest struct {
	Sorts            []any    `json:"sorts,omitempty"`
	Filter           any      `json:"filter,omitempty"`
	StartCursor      string   `json:"start_cursor,omitempty"`
	PageSize         int      `json:"page_size,omitempty"`
	InTrash          *bool    `json:"in_trash,omitempty"`
	ResultType       string   `json:"result_type,omitempty"`
	FilterProperties []string `json:"-"`
}

// DataSourceQueryResponse is the paginated response from querying a data source.
type DataSourceQueryResponse struct {
	Object           string                  `json:"object"`
	Type             string                  `json:"type,omitempty"`
	PageOrDataSource json.RawMessage         `json:"page_or_data_source,omitempty"`
	Results          []DataSourceQueryResult `json:"results"`
	HasMore          bool                    `json:"has_more"`
	NextCursor       *string                 `json:"next_cursor"`
	RequestStatus    json.RawMessage         `json:"request_status,omitempty"`
}

// DataSourceQueryResult is a raw page or data source result.
type DataSourceQueryResult struct {
	Object string          `json:"object"`
	ID     string          `json:"id"`
	Raw    json.RawMessage `json:"-"`
}

// Get retrieves a data source by ID.
func (s *DataSourcesService) Get(ctx context.Context, dataSourceID string) (*DataSource, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, dataSourcePath(dataSourceID), nil, nil)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := s.client.do(req, &dataSource); err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// Query retrieves one page of data source children.
func (s *DataSourcesService) Query(ctx context.Context, dataSourceID string, request QueryDataSourceRequest) (*DataSourceQueryResponse, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, dataSourceQueryPath(dataSourceID), dataSourceQueryValues(request), request)
	if err != nil {
		return nil, err
	}

	var resp DataSourceQueryResponse
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a data source.
func (s *DataSourcesService) Create(ctx context.Context, request CreateDataSourceRequest) (*DataSource, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "data_sources", nil, request)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := s.client.do(req, &dataSource); err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// Update updates a data source by ID.
func (s *DataSourcesService) Update(ctx context.Context, dataSourceID string, request UpdateDataSourceRequest) (*DataSource, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, dataSourcePath(dataSourceID), nil, request)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := s.client.do(req, &dataSource); err != nil {
		return nil, err
	}
	return &dataSource, nil
}

func dataSourcePath(dataSourceID string) string {
	return fmt.Sprintf("data_sources/%s", url.PathEscape(dataSourceID))
}

func dataSourceQueryPath(dataSourceID string) string {
	return fmt.Sprintf("data_sources/%s/query", url.PathEscape(dataSourceID))
}

func dataSourceQueryValues(request QueryDataSourceRequest) url.Values {
	values := url.Values{}
	for _, property := range request.FilterProperties {
		if property != "" {
			values.Add("filter_properties[]", property)
		}
	}
	return values
}

func (r *DataSourceQueryResult) UnmarshalJSON(data []byte) error {
	type dataSourceQueryResultAlias DataSourceQueryResult
	var v dataSourceQueryResultAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*r = DataSourceQueryResult(v)
	return nil
}
