package seaavey

import (
	"context"
	"net/url"
	"strings"
)

// DownloaderService handles downloader endpoints.
type DownloaderService struct {
	client *Client
}

func (s *DownloaderService) getInto(ctx context.Context, platform, targetURL string, dst any) error {
	platform = strings.TrimSpace(platform)
	targetURL = strings.TrimSpace(targetURL)

	switch {
	case platform == "":
		return ErrMissingPlatform
	case targetURL == "":
		return ErrMissingTargetURL
	}

	query := url.Values{}
	query.Set("url", targetURL)

	req, err := s.client.newRequest(ctx, "GET", "/downloader/"+url.PathEscape(platform), query, nil)
	if err != nil {
		return err
	}

	if err := s.client.do(req, dst); err != nil {
		return err
	}

	return nil
}

// Get calls a downloader endpoint by platform using the target URL as a query parameter.
// It returns a generic response shape for mixed or unknown platform payloads.
func (s *DownloaderService) Get(ctx context.Context, platform, targetURL string) (*DownloadResponse, error) {
	var out DownloadResponse
	if err := s.getInto(ctx, platform, targetURL, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// TikTok calls the TikTok downloader endpoint.
func (s *DownloaderService) TikTok(ctx context.Context, targetURL string) (*TikTokResponse, error) {
	var out TikTokResponse
	if err := s.getInto(ctx, "tiktok", targetURL, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// SoundCloud calls the SoundCloud downloader endpoint.
func (s *DownloaderService) SoundCloud(ctx context.Context, targetURL string) (*SoundCloudResponse, error) {
	var out SoundCloudResponse
	if err := s.getInto(ctx, "soundcloud", targetURL, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
