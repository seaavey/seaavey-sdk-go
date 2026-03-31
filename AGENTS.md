# AGENTS.md

## Project Overview

This repository contains `seaavey-sdk-go`, a Go SDK for interacting with the Seaavey API at `https://api.seaavey.com`.

Primary goals:

- Provide a stable, idiomatic Go client for Seaavey API consumers.
- Keep the public API small, predictable, and backward compatible.
- Centralize HTTP behavior so authentication, headers, encoding, and error handling are consistent.
- Organize endpoints into services exposed from a root `Client`.

This file is the operating contract for AI agents and contributors. Follow it strictly.

## Core Rules

- Follow existing repository patterns before introducing new abstractions.
- Preserve backward compatibility unless a breaking change is explicitly requested.
- Keep logic centralized; do not duplicate request, header, or error handling code.
- Every HTTP request must include `X-API-KEY: <API_KEY>`.
- Never use `Authorization: Bearer`.
- Store the API key on the `Client` struct.
- Inject required headers from a single centralized request builder.
- All public SDK methods must accept `context.Context`.

## Project Structure

Prefer this layout unless the repository already defines an equivalent pattern:

```text
.
├── client.go          # Client definition and constructor
├── request.go         # Centralized request creation and execution helpers
├── errors.go          # API and transport error types
├── downloader.go      # Downloader service methods for platform routes
├── models.go          # Shared request/response models
├── go.mod
└── README.md
```

Structure rules:

- Keep the root package focused on the public SDK surface.
- Group endpoints by service, not by individual route files.
- Put shared request/response models near the package they belong to, or in shared model files when reused.
- Keep internal helper functions unexported unless they are part of the intended public API.

## Build And Test Commands

Run these before finishing any change:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Use these when validating coverage or a specific package:

```bash
go test -cover ./...
go test ./... -run TestName
```

## Client Standard

The SDK must expose a reusable `Client` that stores the API key and an `*http.Client`.

Requirements:

- `Client` stores `apiKey string`.
- `Client` stores a reusable `http.Client`.
- `Client` stores the base URL.
- `Client` exposes the `Downloader` service.
- `NewClient(apiKey string)` is the required constructor.
- Avoid constructing ad hoc `http.Client` values inside endpoint methods.

Example:

```go
package seaavey

import (
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.seaavey.com"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Downloader *DownloaderService
}

func NewClient(apiKey string) *Client {
	c := &Client{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	c.Downloader = &DownloaderService{client: c}

	return c
}
```

## Request Rules

All HTTP request creation must go through one centralized helper such as `newRequest`.

Requirements:

- All public service methods must accept `context.Context`.
- Use `http.NewRequestWithContext`.
- Automatically attach:
  - `X-API-KEY: <API_KEY>`
  - `Accept: application/json`
  - `Content-Type: application/json` when a request body exists
- Build query parameters centrally from the request helper.
- For downloader endpoints, pass the target media URL as the `url` query parameter.
- Do not set auth headers in service methods.
- Do not manually duplicate JSON encoding logic across methods.
- Ensure relative paths are resolved against the client base URL.
- Do not manually concatenate raw query strings inside service methods.

Example:

```go
package seaavey

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	u := base.ResolveReference(rel)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var buf *bytes.Buffer
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	var reader *bytes.Buffer
	if buf != nil {
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
```

## Service Structure

Endpoints must be grouped into service types attached to `Client`.

Rules:

- Expose only services that actually exist in the repository.
- In the current SDK, use `Downloader` as the service attached to `Client`.
- Each service holds a reference to the root client.
- Keep service methods thin; shared behavior belongs in the client helpers.
- Map downloader methods directly to supported platform routes.
- Prefer method names such as `TikTok`, `Instagram`, or `YouTube` when they mirror the API endpoint.
- For routes like `/downloader/tiktok?url=<target>`, accept the target URL as a method parameter, not as a path segment.

Example:

```go
package seaavey

import (
	"context"
	"net/url"
)

type DownloaderService struct {
	client *Client
}

type DownloadResponse struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (s *DownloaderService) TikTok(ctx context.Context, targetURL string) (*DownloadResponse, error) {
	query := url.Values{}
	query.Set("url", targetURL)

	req, err := s.client.newRequest(ctx, "GET", "/downloader/tiktok", query, nil)
	if err != nil {
		return nil, err
	}

	var out DownloadResponse
	if err := s.client.do(req, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
```

Required usage style:

```go
ctx := context.Background()
client := NewClient("your-api-key")

resp, err := client.Downloader.TikTok(ctx, "https://www.tiktok.com/@example/video/123")
if err != nil {
	// handle error
}

_ = resp
```

## API Design Conventions

- Use `NewClient(apiKey string)` as the primary entry point.
- Public methods must be context-first: `TikTok(ctx context.Context, targetURL string)`.
- Return concrete response models and `error`.
- Use pointer results for single-resource responses when zero values are ambiguous.
- Avoid leaking raw `http.Response` from public methods unless the SDK explicitly supports low-level APIs.
- Keep naming idiomatic Go: short, descriptive, exported only when necessary.
- Prefer typed request/response structs over `map[string]any`.
- Use JSON struct tags consistently.
- Keep pagination, filtering, and optional fields explicit through option structs when needed.
- When an endpoint expects query parameters, model them explicitly and send them through the centralized request builder.

## Error Handling Rules

- Centralize response execution and decoding in a helper such as `do`.
- Treat all non-2xx responses as errors.
- Return typed API errors when the server responds with an error payload.
- Wrap transport and decoding errors with operation context.
- Preserve `context.Canceled` and `context.DeadlineExceeded` semantics; do not hide them.
- Always close response bodies.
- Do not panic for request, transport, or parsing failures.
- Do not ignore JSON decoding errors.

Recommended pattern:

```go
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}
```

Implementation expectations:

- `do` should execute the request with the reusable `http.Client`.
- `do` should decode successful JSON responses into the destination value.
- `do` should parse error payloads when possible and fall back to status-based errors when not.

## Coding Guidelines

- Write idiomatic Go.
- Keep functions small and single-purpose.
- Prefer composition over unnecessary inheritance-style abstractions.
- Avoid package-level mutable state.
- Use the standard library unless a dependency provides clear value.
- Keep exported API surface intentional and documented.
- Use receiver names consistently, typically `c` for client and `s` for service.
- Avoid reflection and generics unless they materially reduce duplication without harming clarity.
- Run `go fmt` on every change.
- Add doc comments for exported types and functions when they are part of the public SDK.

## Testing

Testing requirements are mandatory for new request and response behavior.

Rules:

- Use `httptest` for HTTP-level tests.
- Test request construction, including:
  - method
  - path
  - query parameters such as `url`
  - `X-API-KEY` header
  - `Accept` header
  - `Content-Type` header when applicable
- Test response parsing for successful responses.
- Test error parsing for non-2xx responses.
- Test context propagation when relevant.
- Prefer table-driven tests where they improve clarity.

Minimum coverage for new endpoints:

- One success-path test.
- One server error test.
- One request-building assertion test.

## Contribution Workflow

Follow this sequence for every change:

1. Inspect the current code to identify existing client, service, and test patterns.
2. Reuse existing helpers before adding new ones.
3. Implement the smallest change that fits the current public API design.
4. Add or update `httptest`-based tests.
5. Run formatting, vetting, and tests.
6. Update documentation when the public API changes.

Review checklist:

- No breaking changes to exported names, signatures, or behavior unless explicitly requested.
- No duplicated request/header/error logic.
- All new public methods accept `context.Context`.
- All requests include `X-API-KEY`.
- No `Authorization: Bearer` usage introduced.

## DO And DO NOT

### DO

- Follow existing patterns in the repository.
- Keep authentication header injection centralized.
- Reuse the shared `http.Client`.
- Keep service methods small and readable.
- Pass downloader target links through the `url` query parameter when the API expects it.
- Add tests for request building and response parsing.
- Return typed, actionable errors.
- Preserve backward compatibility.

### DO NOT

- Do not use `Authorization: Bearer`.
- Do not omit `X-API-KEY`.
- Do not set required headers in multiple places.
- Do not construct a new `http.Client` per request.
- Do not duplicate JSON encoding or response parsing logic.
- Do not hardcode query strings in endpoint methods.
- Do not introduce breaking API changes without explicit approval.
- Do not bypass `context.Context`.
- Do not add dependencies for simple standard-library tasks.

## Downloader Endpoint Pattern

Downloader routes should follow the API path directly and pass the downloadable target as a query parameter.

Pattern:

```text
GET /downloader/{platform}?url=<target_url>
```

Example:

```text
GET /downloader/tiktok?url=https://www.tiktok.com/@example/video/123
```

SDK mapping example:

```go
resp, err := client.Downloader.TikTok(ctx, "https://www.tiktok.com/@example/video/123")
```

## Non-Negotiable Authentication Policy

Every outbound request must include:

```http
X-API-KEY: <API_KEY>
```

This repository does not use bearer token authentication.

Forbidden:

```http
Authorization: Bearer <token>
```

Authentication must be implemented only through the client-held API key and the centralized request builder.
