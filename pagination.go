package notion

import (
	"context"
	"errors"
)

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

// PageFunc fetches one cursor-based response page.
type PageFunc[T any] func(context.Context, *Pagination) (*PaginatedResponse[T], error)

// EachFunc handles one paginated result.
type EachFunc[T any] func(T) error

// ForEachPaginated iterates through every result from a cursor-based endpoint.
func ForEachPaginated[T any](ctx context.Context, pagination *Pagination, pageFn PageFunc[T], each EachFunc[T]) error {
	if pageFn == nil {
		return errors.New("notion: nil paginated page function")
	}
	if each == nil {
		return errors.New("notion: nil paginated callback")
	}

	next := clonePagination(pagination)
	for {
		page, err := pageFn(ctx, next)
		if err != nil {
			return err
		}
		if page == nil {
			return errors.New("notion: paginated page function returned nil response")
		}

		for _, result := range page.Results {
			if err := each(result); err != nil {
				return err
			}
		}

		if !page.HasMore {
			return nil
		}
		if page.NextCursor == nil || *page.NextCursor == "" {
			return errors.New("notion: paginated response has_more=true without next_cursor")
		}

		if next == nil {
			next = &Pagination{}
		}
		next.StartCursor = *page.NextCursor
	}
}

func clonePagination(pagination *Pagination) *Pagination {
	if pagination == nil {
		return nil
	}

	next := *pagination
	return &next
}
