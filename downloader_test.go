package seaavey

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	client := NewClient("test-key")
	if err := client.SetBaseURL("https://api.seaavey.test"); err != nil {
		t.Fatalf("SetBaseURL() error = %v", err)
	}

	client.SetHTTPClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			return rec.Result(), nil
		}),
	})

	return client
}

func TestDownloaderTikTokRequestAndResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/downloader/tiktok" {
			t.Fatalf("expected path /downloader/tiktok, got %s", r.URL.Path)
		}

		if got := r.URL.Query().Get("url"); got != "https://www.tiktok.com/@example/video/123" {
			t.Fatalf("expected query url to be set, got %q", got)
		}

		if got := r.Header.Get("X-API-KEY"); got != "test-key" {
			t.Fatalf("expected X-API-KEY header, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected Accept header, got %q", got)
		}

		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("expected Authorization header to be empty, got %q", got)
		}

		if got := r.Header.Get("Content-Type"); got != "" {
			t.Fatalf("expected no Content-Type header for GET without body, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"12 ms","data":{"author":"seaavey","title":"Example Video","vidio":"https://cdn.example/video.mp4","audio":"https://cdn.example/audio.mp3"}}`))
	})

	resp, err := client.Downloader.TikTok(context.Background(), "https://www.tiktok.com/@example/video/123")
	if err != nil {
		t.Fatalf("TikTok() error = %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.VideoURL != "https://cdn.example/video.mp4" {
		t.Fatalf("expected VideoURL to be parsed, got %q", resp.VideoURL)
	}

	if resp.AudioURL != "https://cdn.example/audio.mp3" {
		t.Fatalf("expected AudioURL to be parsed, got %q", resp.AudioURL)
	}

	if resp.Title != "Example Video" {
		t.Fatalf("expected Title to be parsed, got %q", resp.Title)
	}

	if resp.Author != "seaavey" {
		t.Fatalf("expected Author to be parsed, got %q", resp.Author)
	}
}

func TestDownloaderTikTokAPIError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":"invalid_url","message":"url is required"}`))
	})

	_, err := client.Downloader.TikTok(context.Background(), "https://www.tiktok.com/@example/video/123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}

	if apiErr.Code != "invalid_url" {
		t.Fatalf("expected code invalid_url, got %q", apiErr.Code)
	}

	if apiErr.Message != "url is required" {
		t.Fatalf("expected message to be parsed, got %q", apiErr.Message)
	}
}

func TestDownloaderTikTokValidatesTargetURL(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key")

	_, err := client.Downloader.TikTok(context.Background(), "  ")
	if !errors.Is(err, ErrMissingTargetURL) {
		t.Fatalf("expected ErrMissingTargetURL, got %v", err)
	}
}
