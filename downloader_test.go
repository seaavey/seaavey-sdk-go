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
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"12 ms","data":{"author":"seaavey","title":"Example Video","video":"https://cdn.example/video.mp4","audio":"https://cdn.example/audio.mp3"}}`))
	})

	resp, err := client.Downloader.TikTok(context.Background(), "https://www.tiktok.com/@example/video/123")
	if err != nil {
		t.Fatalf("TikTok() error = %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Status)
	}

	if !resp.Success {
		t.Fatal("expected success to be true")
	}

	if resp.Latency != "12 ms" {
		t.Fatalf("expected latency to be parsed, got %q", resp.Latency)
	}

	if resp.Data.Video != "https://cdn.example/video.mp4" {
		t.Fatalf("expected data.video to be parsed, got %q", resp.Data.Video)
	}

	if resp.Data.Audio != "https://cdn.example/audio.mp3" {
		t.Fatalf("expected data.audio to be parsed, got %q", resp.Data.Audio)
	}

	if resp.Data.Title != "Example Video" {
		t.Fatalf("expected data.title to be parsed, got %q", resp.Data.Title)
	}

	if resp.Data.Author != "seaavey" {
		t.Fatalf("expected data.author to be parsed, got %q", resp.Data.Author)
	}
}

func TestDownloaderTikTokSupportsLegacyVidioField(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"9 ms","data":{"author":"seaavey","title":"Legacy Video Key","vidio":"https://cdn.example/video.mp4","audio":"https://cdn.example/audio.mp3"}}`))
	})

	resp, err := client.Downloader.TikTok(context.Background(), "https://www.tiktok.com/@example/video/123")
	if err != nil {
		t.Fatalf("TikTok() error = %v", err)
	}

	if resp.Data.Video != "https://cdn.example/video.mp4" {
		t.Fatalf("expected legacy vidio key to map into data.video, got %q", resp.Data.Video)
	}
}

func TestDownloaderTikTokImageResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"1235 ms","data":{"author":"Ayuni","title":"Photo Post","images":["https://cdn.example/image-1.jpeg","https://cdn.example/image-2.jpeg"],"audio":"https://cdn.example/audio.mp3"}}`))
	})

	resp, err := client.Downloader.TikTok(context.Background(), "https://www.tiktok.com/@example/photo/123")
	if err != nil {
		t.Fatalf("TikTok() error = %v", err)
	}

	if len(resp.Data.Images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(resp.Data.Images))
	}

	if resp.Data.Images[0] != "https://cdn.example/image-1.jpeg" {
		t.Fatalf("unexpected first image URL, got %q", resp.Data.Images[0])
	}

	if resp.Data.Audio != "https://cdn.example/audio.mp3" {
		t.Fatalf("expected audio to be parsed, got %q", resp.Data.Audio)
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

func TestDownloaderSoundCloudRequestAndResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
		}

		if r.URL.Path != "/downloader/soundcloud" {
			t.Fatalf("expected path /downloader/soundcloud, got %s", r.URL.Path)
		}

		if got := r.URL.Query().Get("url"); got != "https://soundcloud.com/example/track" {
			t.Fatalf("expected query url to be set, got %q", got)
		}

		if got := r.Header.Get("X-API-KEY"); got != "test-key" {
			t.Fatalf("expected X-API-KEY header, got %q", got)
		}

		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("expected Accept header, got %q", got)
		}

		if got := r.Header.Get("Content-Type"); got != "" {
			t.Fatalf("expected no Content-Type header for GET without body, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"15046 ms","data":{"title":"https://soundcloud.com/example/track","thumbnail":"https://logo.clearbit.com/soundcloud.com?size=256","download":"https://ricky4.savenow.to/pacific/?abc","alternatives":[{"type":"nip.io","url":"https://179-43-173-246.nip.io/pacific/?abc","has_ssl":true},{"type":"traefik.me","url":"http://179-43-173-246.traefik.me/pacific/?abc","has_ssl":false}]}}`))
	})

	resp, err := client.Downloader.SoundCloud(context.Background(), "https://soundcloud.com/example/track")
	if err != nil {
		t.Fatalf("SoundCloud() error = %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Status)
	}

	if !resp.Success {
		t.Fatal("expected success to be true")
	}

	if resp.Data.Title != "https://soundcloud.com/example/track" {
		t.Fatalf("expected title to be parsed, got %q", resp.Data.Title)
	}

	if resp.Data.Thumbnail != "https://logo.clearbit.com/soundcloud.com?size=256" {
		t.Fatalf("expected thumbnail to be parsed, got %q", resp.Data.Thumbnail)
	}

	if resp.Data.Download != "https://ricky4.savenow.to/pacific/?abc" {
		t.Fatalf("expected download to be parsed, got %q", resp.Data.Download)
	}

	if len(resp.Data.Alternatives) != 2 {
		t.Fatalf("expected 2 alternatives, got %d", len(resp.Data.Alternatives))
	}

	if resp.Data.Alternatives[0].Type != "nip.io" {
		t.Fatalf("expected first alternative type to be parsed, got %q", resp.Data.Alternatives[0].Type)
	}

	if !resp.Data.Alternatives[0].HasSSL {
		t.Fatal("expected first alternative has_ssl=true")
	}
}

func TestDownloaderSoundCloudAlternativesCanBeEmpty(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"14 ms","data":{"title":"https://soundcloud.com/example/track","thumbnail":"https://logo.clearbit.com/soundcloud.com?size=256","download":"https://ricky4.savenow.to/pacific/?abc","alternatives":[]}}`))
	})

	resp, err := client.Downloader.SoundCloud(context.Background(), "https://soundcloud.com/example/track")
	if err != nil {
		t.Fatalf("SoundCloud() error = %v", err)
	}

	if len(resp.Data.Alternatives) != 0 {
		t.Fatalf("expected alternatives to be empty, got %d", len(resp.Data.Alternatives))
	}
}

func TestDownloaderSoundCloudAlternativesCanBeMissing(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":200,"success":true,"latency":"14 ms","data":{"title":"https://soundcloud.com/example/track","thumbnail":"https://logo.clearbit.com/soundcloud.com?size=256","download":"https://ricky4.savenow.to/pacific/?abc"}}`))
	})

	resp, err := client.Downloader.SoundCloud(context.Background(), "https://soundcloud.com/example/track")
	if err != nil {
		t.Fatalf("SoundCloud() error = %v", err)
	}

	if len(resp.Data.Alternatives) != 0 {
		t.Fatalf("expected missing alternatives to decode as empty set, got %d", len(resp.Data.Alternatives))
	}
}

func TestDownloaderSoundCloudAPIError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":"invalid_url","message":"url is required"}`))
	})

	_, err := client.Downloader.SoundCloud(context.Background(), "https://soundcloud.com/example/track")
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
}
