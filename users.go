package notion

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// UsersService handles user endpoints.
type UsersService struct {
	client *Client
}

// List retrieves one page of users.
func (s *UsersService) List(ctx context.Context, pagination *Pagination) (*PaginatedResponse[User], error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "users", pagination, nil)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[User]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a user by ID.
func (s *UsersService) Get(ctx context.Context, userID string) (*User, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, userPath(userID), nil, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := s.client.do(req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Me retrieves the bot user associated with the token.
func (s *UsersService) Me(ctx context.Context) (*User, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "users/me", nil, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := s.client.do(req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func userPath(userID string) string {
	return fmt.Sprintf("users/%s", url.PathEscape(userID))
}
