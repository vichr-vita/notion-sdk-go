package notion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestRoundTripFunc(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://api.notion.com/v1/search", nil)
	require.NoError(t, err)

	called := false
	rt := roundTripFunc(func(got *http.Request) (*http.Response, error) {
		called = true
		require.Same(t, req, got)

		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, called)
}
