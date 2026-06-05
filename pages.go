package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// PagesService handles page endpoints.
type PagesService struct {
	client *Client
}

// CreatePageRequest is the body for creating a page.
type CreatePageRequest struct {
	Parent     Parent         `json:"parent"`
	Properties map[string]any `json:"properties,omitempty"`
	Children   []BlockRequest `json:"children,omitempty"`
	Icon       *Icon          `json:"icon,omitempty"`
	Cover      *File          `json:"cover,omitempty"`
}

// UpdatePageRequest is the body for updating page metadata and properties.
type UpdatePageRequest struct {
	Properties map[string]any `json:"properties,omitempty"`
	Icon       *Icon          `json:"icon,omitempty"`
	Cover      *File          `json:"cover,omitempty"`
	InTrash    *bool          `json:"in_trash,omitempty"`
	Archived   *bool          `json:"archived,omitempty"`
	IsArchived *bool          `json:"is_archived,omitempty"`
}

// PagePropertyItem is a flexible page property item response.
type PagePropertyItem struct {
	Object string          `json:"object,omitempty"`
	ID     string          `json:"id,omitempty"`
	Type   string          `json:"type,omitempty"`
	Raw    json.RawMessage `json:"-"`
}

// PageMarkdown is the response from retrieving page content as markdown.
type PageMarkdown struct {
	Object          string          `json:"object"`
	ID              string          `json:"id"`
	Markdown        string          `json:"markdown"`
	Truncated       bool            `json:"truncated"`
	UnknownBlockIDs []string        `json:"unknown_block_ids,omitempty"`
	Raw             json.RawMessage `json:"-"`
}

// Get retrieves a page by ID.
func (s *PagesService) Get(ctx context.Context, pageID string) (*Page, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, pagePath(pageID), nil, nil)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := s.client.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Create creates a page.
func (s *PagesService) Create(ctx context.Context, request CreatePageRequest) (*Page, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "pages", nil, request)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := s.client.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Update updates a page by ID.
func (s *PagesService) Update(ctx context.Context, pageID string, request UpdatePageRequest) (*Page, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, pagePath(pageID), nil, request)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := s.client.do(req, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// GetPropertyItem retrieves a page property item by page and property ID.
func (s *PagesService) GetPropertyItem(ctx context.Context, pageID, propertyID string, pagination *Pagination) (*PagePropertyItem, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, pagePropertyPath(pageID, propertyID), pagination, nil)
	if err != nil {
		return nil, err
	}

	var item PagePropertyItem
	if err := s.client.do(req, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// GetMarkdown retrieves page content as markdown.
func (s *PagesService) GetMarkdown(ctx context.Context, pageID string) (*PageMarkdown, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, pageMarkdownPath(pageID), nil, nil)
	if err != nil {
		return nil, err
	}

	var markdown PageMarkdown
	if err := s.client.do(req, &markdown); err != nil {
		return nil, err
	}
	return &markdown, nil
}

func pagePath(pageID string) string {
	return fmt.Sprintf("pages/%s", url.PathEscape(pageID))
}

func pagePropertyPath(pageID, propertyID string) string {
	return fmt.Sprintf("pages/%s/properties/%s", url.PathEscape(pageID), url.PathEscape(propertyID))
}

func pageMarkdownPath(pageID string) string {
	return fmt.Sprintf("pages/%s/markdown", url.PathEscape(pageID))
}

func (p *PagePropertyItem) UnmarshalJSON(data []byte) error {
	type pagePropertyItemAlias PagePropertyItem
	var v pagePropertyItemAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*p = PagePropertyItem(v)
	return nil
}

func (p *PageMarkdown) UnmarshalJSON(data []byte) error {
	type pageMarkdownAlias PageMarkdown
	var v pageMarkdownAlias
	if err := unmarshalWithRaw(data, &v, &v.Raw); err != nil {
		return err
	}
	*p = PageMarkdown(v)
	return nil
}
