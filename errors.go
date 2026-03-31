package seaavey

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingAPIKey is returned when the client has no API key configured.
	ErrMissingAPIKey = errors.New("seaavey: missing API key")
	// ErrInvalidBaseURL is returned when the client base URL is invalid.
	ErrInvalidBaseURL = errors.New("seaavey: invalid base URL")
	// ErrMissingPlatform is returned when a downloader platform is empty.
	ErrMissingPlatform = errors.New("seaavey: missing downloader platform")
	// ErrMissingTargetURL is returned when a downloader target URL is empty.
	ErrMissingTargetURL = errors.New("seaavey: missing target url")
)

// APIError represents a non-2xx API response.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Body       []byte
}

// Error returns the formatted API error.
func (e *APIError) Error() string {
	switch {
	case e == nil:
		return ""
	case e.Code != "" && e.Message != "":
		return fmt.Sprintf("seaavey API error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("seaavey API error: status=%d message=%s", e.StatusCode, e.Message)
	default:
		return fmt.Sprintf("seaavey API error: status=%d", e.StatusCode)
	}
}
