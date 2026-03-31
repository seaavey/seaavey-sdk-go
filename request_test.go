package seaavey

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNewRequestSetsHeadersQueryAndBody(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key")

	req, err := client.newRequest(
		context.Background(),
		http.MethodPost,
		"/downloader/tiktok",
		url.Values{"url": []string{"https://www.tiktok.com/@example/video/123"}},
		map[string]string{"quality": "hd"},
	)
	if err != nil {
		t.Fatalf("newRequest() error = %v", err)
	}

	if req.Header.Get("X-API-KEY") != "test-key" {
		t.Fatalf("expected X-API-KEY header, got %q", req.Header.Get("X-API-KEY"))
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Fatalf("expected Accept header, got %q", req.Header.Get("Accept"))
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type header, got %q", req.Header.Get("Content-Type"))
	}

	if req.URL.Path != "/downloader/tiktok" {
		t.Fatalf("expected path /downloader/tiktok, got %s", req.URL.Path)
	}

	if req.URL.Query().Get("url") != "https://www.tiktok.com/@example/video/123" {
		t.Fatalf("expected query url to be set, got %q", req.URL.Query().Get("url"))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("reading request body: %v", err)
	}

	if !strings.Contains(string(body), `"quality":"hd"`) {
		t.Fatalf("expected JSON body to contain quality field, got %s", string(body))
	}
}

func TestNewRequestRequiresAPIKey(t *testing.T) {
	t.Parallel()

	client := NewClient("")

	_, err := client.newRequest(context.Background(), http.MethodGet, "/downloader/tiktok", nil, nil)
	if !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("expected ErrMissingAPIKey, got %v", err)
	}
}
