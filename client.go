package seaavey

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the root Seaavey API client.
type Client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client

	Downloader *DownloaderService
}

// NewClient creates a new Seaavey API client.
func NewClient(apiKey string) *Client {
	baseURL, _ := url.Parse("https://api.seaavey.com")

	c := &Client{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	c.Downloader = &DownloaderService{client: c}

	return c
}

// SetBaseURL overrides the API base URL.
func (c *Client) SetBaseURL(rawURL string) error {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBaseURL, err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: base URL must include scheme and host", ErrInvalidBaseURL)
	}

	c.baseURL = parsed
	return nil
}

// SetHTTPClient overrides the underlying HTTP client.
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		return
	}

	c.httpClient = httpClient
}
