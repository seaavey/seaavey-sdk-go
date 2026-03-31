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

// Get calls a downloader endpoint by platform using the target URL as a query parameter.
func (s *DownloaderService) Get(ctx context.Context, platform, targetURL string) (*DownloadResponse, error) {
	platform = strings.TrimSpace(platform)
	targetURL = strings.TrimSpace(targetURL)

	switch {
	case platform == "":
		return nil, ErrMissingPlatform
	case targetURL == "":
		return nil, ErrMissingTargetURL
	}

	query := url.Values{}
	query.Set("url", targetURL)

	req, err := s.client.newRequest(ctx, "GET", "/downloader/"+url.PathEscape(platform), query, nil)
	if err != nil {
		return nil, err
	}

	var envelope downloadEnvelope
	if err := s.client.do(req, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Data, nil
}

// TikTok calls the TikTok downloader endpoint.
func (s *DownloaderService) TikTok(ctx context.Context, targetURL string) (*DownloadResponse, error) {
	return s.Get(ctx, "tiktok", targetURL)
}
