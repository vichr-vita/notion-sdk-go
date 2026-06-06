package notion

import (
	"context"
	"net/http"
)

// CustomEmojisService handles custom emoji endpoints.
type CustomEmojisService struct {
	client *Client
}

// List retrieves one page of workspace custom emojis.
func (s *CustomEmojisService) List(ctx context.Context, pagination *Pagination) (*PaginatedResponse[CustomEmoji], error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "custom_emojis", pagination, nil)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[CustomEmoji]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
