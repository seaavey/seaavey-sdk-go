package seaavey

// DownloadResponse contains the common downloader response fields.
type DownloadResponse struct {
	Author   string `json:"author,omitempty"`
	Title    string `json:"title,omitempty"`
	VideoURL string `json:"vidio,omitempty"`
	AudioURL string `json:"audio,omitempty"`
}

type downloadEnvelope struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Latency string           `json:"latency,omitempty"`
	Data    DownloadResponse `json:"data"`
}
