package notion

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.notion.com/v1"
	defaultVersion = "2026-03-11"
)

const (
	defaultMaxRetries    = 2
	defaultRetryDelay    = 500 * time.Millisecond
	defaultRetryMaxDelay = 4 * time.Second
	defaultRetryJitter   = 0.2
)

// Option configures a Client.
type Option func(*clientConfig)

// RequestHook runs after a request is built and before it is sent.
type RequestHook func(*http.Request) error

// ResponseHook runs after a response is received and before it is decoded.
type ResponseHook func(*http.Response) error

// RetryConfig controls retry behavior for transient API responses.
type RetryConfig struct {
	MaxRetries int
	Delay      time.Duration
	MaxDelay   time.Duration
	Jitter     float64
}

type clientConfig struct {
	token        string
	baseURL      string
	version      string
	httpClient   *http.Client
	requestHook  RequestHook
	responseHook ResponseHook
	retry        RetryConfig
	retrySleep   func(time.Duration)
	retryJitter  func() float64
}

// Client is the root Notion API client.
type Client struct {
	config clientConfig

	Pages        *PagesService
	Blocks       *BlocksService
	DataSources  *DataSourcesService
	Databases    *DatabasesService
	Search       *SearchService
	Users        *UsersService
	Comments     *CommentsService
	FileUploads  *FileUploadsService
	CustomEmojis *CustomEmojisService
}

// NewClient creates a Notion API client.
func NewClient(token string, opts ...Option) *Client {
	cfg := clientConfig{
		token:      token,
		baseURL:    defaultBaseURL,
		version:    defaultVersion,
		httpClient: http.DefaultClient,
		retry: RetryConfig{
			MaxRetries: defaultMaxRetries,
			Delay:      defaultRetryDelay,
			MaxDelay:   defaultRetryMaxDelay,
			Jitter:     defaultRetryJitter,
		},
		retryJitter: rand.Float64,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	c := &Client{config: cfg}
	c.Pages = &PagesService{client: c}
	c.Blocks = &BlocksService{client: c}
	c.DataSources = &DataSourcesService{client: c}
	c.Databases = &DatabasesService{client: c}
	c.Search = &SearchService{client: c}
	c.Users = &UsersService{client: c}
	c.Comments = &CommentsService{client: c}
	c.FileUploads = &FileUploadsService{client: c}
	c.CustomEmojis = &CustomEmojisService{client: c}

	return c
}

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(cfg *clientConfig) {
		if httpClient != nil {
			cfg.httpClient = httpClient
		}
	}
}

// WithBaseURL sets the base Notion API URL.
func WithBaseURL(baseURL string) Option {
	return func(cfg *clientConfig) {
		if baseURL != "" {
			cfg.baseURL = baseURL
		}
	}
}

// WithVersion sets the Notion API version header value.
func WithVersion(version string) Option {
	return func(cfg *clientConfig) {
		if version != "" {
			cfg.version = version
		}
	}
}

// WithRequestHook sets a hook that runs before API requests are sent.
func WithRequestHook(hook RequestHook) Option {
	return func(cfg *clientConfig) {
		cfg.requestHook = hook
	}
}

// WithResponseHook sets a hook that runs before API responses are decoded.
func WithResponseHook(hook ResponseHook) Option {
	return func(cfg *clientConfig) {
		cfg.responseHook = hook
	}
}

// WithRetryConfig sets retry behavior for transient API responses.
func WithRetryConfig(retry RetryConfig) Option {
	return func(cfg *clientConfig) {
		if retry.Delay == 0 {
			retry.Delay = defaultRetryDelay
		}
		if retry.MaxDelay == 0 {
			retry.MaxDelay = defaultRetryMaxDelay
		}
		cfg.retry = retry
	}
}
