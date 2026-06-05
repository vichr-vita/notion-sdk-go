package notion

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DatabasesService handles legacy database compatibility endpoints.
//
// DataSources is the preferred service for querying database rows in the
// current Notion API.
type DatabasesService struct {
	client *Client
}

// CreateDatabaseRequest is the body for creating a legacy database container.
type CreateDatabaseRequest struct {
	Parent            Parent           `json:"parent"`
	Title             []RichTextObject `json:"title,omitempty"`
	Description       []RichTextObject `json:"description,omitempty"`
	IsInline          *bool            `json:"is_inline,omitempty"`
	InitialDataSource any              `json:"initial_data_source,omitempty"`
	Icon              *Icon            `json:"icon,omitempty"`
	Cover             *File            `json:"cover,omitempty"`
}

// UpdateDatabaseRequest is the body for updating a legacy database container.
type UpdateDatabaseRequest struct {
	Parent      *Parent          `json:"parent,omitempty"`
	Title       []RichTextObject `json:"title,omitempty"`
	Description []RichTextObject `json:"description,omitempty"`
	IsInline    *bool            `json:"is_inline,omitempty"`
	Icon        *Icon            `json:"icon,omitempty"`
	Cover       *File            `json:"cover,omitempty"`
	InTrash     *bool            `json:"in_trash,omitempty"`
	IsLocked    *bool            `json:"is_locked,omitempty"`
}

// Get retrieves a legacy database by ID.
func (s *DatabasesService) Get(ctx context.Context, databaseID string) (*Database, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, databasePath(databaseID), nil, nil)
	if err != nil {
		return nil, err
	}

	var database Database
	if err := s.client.do(req, &database); err != nil {
		return nil, err
	}
	return &database, nil
}

// Create creates a legacy database container.
func (s *DatabasesService) Create(ctx context.Context, request CreateDatabaseRequest) (*Database, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "databases", nil, request)
	if err != nil {
		return nil, err
	}

	var database Database
	if err := s.client.do(req, &database); err != nil {
		return nil, err
	}
	return &database, nil
}

// Update updates a legacy database by ID.
func (s *DatabasesService) Update(ctx context.Context, databaseID string, request UpdateDatabaseRequest) (*Database, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, databasePath(databaseID), nil, request)
	if err != nil {
		return nil, err
	}

	var database Database
	if err := s.client.do(req, &database); err != nil {
		return nil, err
	}
	return &database, nil
}

func databasePath(databaseID string) string {
	return fmt.Sprintf("databases/%s", url.PathEscape(databaseID))
}
