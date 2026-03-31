# Downloader API

Read this file when updating `models.go`, `downloader.go`, `downloader_test.go`, or `cmd/livecheck/main.go`.

## Observed Response Shape

Observed from:

```text
GET /downloader/tiktok?url=<tiktok_url>
```

Example shape:

```json
{
  "status": 200,
  "success": true,
  "latency": "752 ms",
  "data": {
    "author": "raven.",
    "title": "Perfect - ed sheeran🎶 #perfect #lyrics #songs #foryou #fyp ",
    "vidio": "https://v19.tiktokcdn-us.com/...",
    "audio": "https://v16-ies-music.tiktokcdn-us.com/..."
  }
}
```

## Notes

- The useful payload is nested under `data`.
- The API uses `vidio` as the JSON key. Keep the Go field as `VideoURL` with `json:"vidio"`.
- The current SDK returns `DownloadResponse` only, not the envelope metadata.
- If the live response shape changes, update the code and this reference in the same change.
