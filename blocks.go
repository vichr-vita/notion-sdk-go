package notion

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// BlocksService handles block endpoints.
type BlocksService struct {
	client *Client
}

// UpdateBlockRequest is the body for updating a block.
type UpdateBlockRequest map[string]any

// BlockRequest is a flexible block body used when creating child blocks.
type BlockRequest map[string]any

// AppendBlockChildrenRequest is the body for appending children to a block.
type AppendBlockChildrenRequest struct {
	Children []BlockRequest `json:"children"`
	Position *BlockPosition `json:"position,omitempty"`
}

// BlockPosition controls where appended children are inserted.
type BlockPosition struct {
	Type       string              `json:"type"`
	AfterBlock *AfterBlockPosition `json:"after_block,omitempty"`
}

// AfterBlockPosition identifies the block after which new children are inserted.
type AfterBlockPosition struct {
	ID string `json:"id"`
}

// Get retrieves a block by ID.
func (s *BlocksService) Get(ctx context.Context, blockID string) (*Block, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, blockPath(blockID), nil, nil)
	if err != nil {
		return nil, err
	}

	var block Block
	if err := s.client.do(req, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// Update updates a block by ID.
func (s *BlocksService) Update(ctx context.Context, blockID string, request UpdateBlockRequest) (*Block, error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, blockPath(blockID), nil, request)
	if err != nil {
		return nil, err
	}

	var block Block
	if err := s.client.do(req, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// Delete archives a block by ID.
func (s *BlocksService) Delete(ctx context.Context, blockID string) (*Block, error) {
	req, err := s.client.newRequest(ctx, http.MethodDelete, blockPath(blockID), nil, nil)
	if err != nil {
		return nil, err
	}

	var block Block
	if err := s.client.do(req, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// ListChildren retrieves one page of child blocks.
func (s *BlocksService) ListChildren(ctx context.Context, blockID string, pagination *Pagination) (*PaginatedResponse[Block], error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, blockChildrenPath(blockID), pagination, nil)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[Block]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppendChildren appends child blocks to a block.
func (s *BlocksService) AppendChildren(ctx context.Context, blockID string, request AppendBlockChildrenRequest) (*PaginatedResponse[Block], error) {
	req, err := s.client.newRequest(ctx, http.MethodPatch, blockChildrenPath(blockID), nil, request)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[Block]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ForEachChild iterates through every child block.
func (s *BlocksService) ForEachChild(ctx context.Context, blockID string, pagination *Pagination, each func(Block) error) error {
	return ForEachPaginated(ctx, pagination, func(ctx context.Context, pagination *Pagination) (*PaginatedResponse[Block], error) {
		return s.ListChildren(ctx, blockID, pagination)
	}, each)
}

func blockPath(blockID string) string {
	return fmt.Sprintf("blocks/%s", url.PathEscape(blockID))
}

func blockChildrenPath(blockID string) string {
	return fmt.Sprintf("blocks/%s/children", url.PathEscape(blockID))
}
