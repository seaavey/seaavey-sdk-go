---
name: seaavey-go-sdk
description: Implement and update the Seaavey Go SDK in this repository. Use when working on downloader endpoints, client/request helpers, response models, API error handling, live smoke tests, or httptest coverage for seaavey-sdk-go.
---

# Seaavey Go SDK

## Core Rules

- Keep the root package as `seaavey`.
- Keep `NewClient(apiKey string)` as the public constructor.
- Keep the API key on `Client`; do not move auth state into call sites.
- Reuse one `http.Client` per `Client`.
- Build every HTTP request through `newRequest`.
- Always send `X-API-KEY`; never use `Authorization: Bearer`.
- Set `Accept: application/json` on every request.
- Set `Content-Type: application/json` only when a body exists.
- Keep all public SDK methods context-first.
- Avoid breaking exported names or signatures unless the user explicitly asks.

## Repository Map

- [client.go](../../client.go): `Client`, `NewClient`, base URL handling, injected `Downloader` service.
- [request.go](../../request.go): centralized request creation, header injection, `do`, and `APIError` decoding.
- [downloader.go](../../downloader.go): `DownloaderService`, generic platform fetch, and `TikTok` wrapper.
- [models.go](../../models.go): downloader response model and private response envelope.
- [downloader_test.go](../../downloader_test.go): request/response tests with `httptest`.
- [request_test.go](../../request_test.go): request builder tests for headers, query params, and API key enforcement.
- [cmd/livecheck/main.go](../../cmd/livecheck/main.go): real API smoke test entrypoint using env vars.

## Endpoint Pattern

- Model downloader routes as `GET /downloader/{platform}?url=<target_url>`.
- Use `Downloader.Get(ctx, platform, targetURL)` for generic platform support.
- Add thin wrappers like `Downloader.TikTok(ctx, targetURL)` only when they improve call-site clarity.
- Build query params with `url.Values`; do not concatenate raw query strings by hand.
- Keep service methods thin; shared logic belongs in `newRequest` or `do`.

## Response Shape

Read [references/downloader-api.md](./references/downloader-api.md) when changing models, tests, or live-check behavior.

Current expectations:

- The API response is wrapped in a top-level envelope with `status`, `success`, `latency`, and `data`.
- The SDK currently returns only `data` to callers.
- The API field for the video link is spelled `vidio`, not `video`.
- Map `data.vidio` to `DownloadResponse.VideoURL`.
- Map `data.audio` to `DownloadResponse.AudioURL`.

## Workflow

1. Read `client.go`, `request.go`, `downloader.go`, and `models.go` before changing behavior.
2. Reuse the centralized request path before adding any new helper.
3. If the API shape changes, update `models.go`, the downloader decoding path, the tests, and `references/downloader-api.md` together.
4. Keep auth/header logic in one place only.
5. Add or update `httptest` coverage for request building and response parsing.
6. Use `cmd/livecheck/main.go` for a real API smoke test when the user provides an API key and target URL.

## Validation

In this workspace, Go's default build cache may point to a read-only location. Prefer `GOCACHE=/tmp/...` when running verification commands.

Run:

```bash
go fmt ./...
GOCACHE=/tmp/go-build-cache go test ./...
GOCACHE=/tmp/go-build-cache-race go test -race ./...
GOCACHE=/tmp/go-build-cache-vet go vet ./...
```

For a live API smoke test:

```bash
SEAAVEY_API_KEY=... SEAAVEY_TARGET_URL=... GOCACHE=/tmp/go-build-cache-live go run ./cmd/livecheck
```

## Do Not

- Do not send bearer tokens.
- Do not add headers in service methods.
- Do not create a new `http.Client` per request.
- Do not return raw envelope metadata unless the user asks for that API change.
- Do not rename `vidio` at the JSON tag level; keep the API mapping accurate and expose a sane Go field name instead.
