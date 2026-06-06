package notion

import (
	"context"
	"net/http"
)

// CommentsService handles comment endpoints.
type CommentsService struct {
	client *Client
}

// ListCommentsRequest is the query for listing comments on a page or block.
type ListCommentsRequest struct {
	BlockID     string `json:"block_id,omitempty" url:"block_id,omitempty"`
	StartCursor string `json:"start_cursor,omitempty" url:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty" url:"page_size,omitempty"`
}

// CreateCommentRequest is the body for creating a comment.
type CreateCommentRequest struct {
	Parent       *Parent             `json:"parent,omitempty"`
	DiscussionID string              `json:"discussion_id,omitempty"`
	RichText     []RichTextObject    `json:"rich_text,omitempty"`
	Markdown     string              `json:"markdown,omitempty"`
	Attachments  []CommentAttachment `json:"attachments,omitempty"`
	DisplayName  any                 `json:"display_name,omitempty"`
}

// List retrieves one page of unresolved comments for a page or block.
func (s *CommentsService) List(ctx context.Context, request ListCommentsRequest) (*PaginatedResponse[Comment], error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "comments", request, nil)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[Comment]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a comment on a page, block, or existing discussion thread.
func (s *CommentsService) Create(ctx context.Context, request CreateCommentRequest) (*Comment, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "comments", nil, request)
	if err != nil {
		return nil, err
	}

	var comment Comment
	if err := s.client.do(req, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}
