package seaavey

import "encoding/json"

// ResponseMeta contains common response envelope fields.
type ResponseMeta struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Latency string `json:"latency,omitempty"`
}

// DownloadResponse is a generic downloader response for platforms using mixed fields.
type DownloadResponse struct {
	ResponseMeta
	Data DownloadData `json:"data"`
}

// DownloadData contains generic downloader media fields.
type DownloadData struct {
	Author       string                `json:"author,omitempty"`
	Title        string                `json:"title,omitempty"`
	Video        string                `json:"video,omitempty"`
	Images       []string              `json:"images,omitempty"`
	Audio        string                `json:"audio,omitempty"`
	Thumbnail    string                `json:"thumbnail,omitempty"`
	Download     string                `json:"download,omitempty"`
	Alternatives []DownloadAlternative `json:"alternatives,omitempty"`
}

// DownloadAlternative contains an alternative download URL.
type DownloadAlternative struct {
	Type   string `json:"type,omitempty"`
	URL    string `json:"url,omitempty"`
	HasSSL bool   `json:"has_ssl"`
}

// UnmarshalJSON accepts both "video" and legacy "vidio" keys.
func (d *DownloadData) UnmarshalJSON(data []byte) error {
	type rawDownloadData struct {
		Author       string                `json:"author,omitempty"`
		Title        string                `json:"title,omitempty"`
		Video        string                `json:"video,omitempty"`
		Vidio        string                `json:"vidio,omitempty"`
		Images       []string              `json:"images,omitempty"`
		Audio        string                `json:"audio,omitempty"`
		Thumbnail    string                `json:"thumbnail,omitempty"`
		Download     string                `json:"download,omitempty"`
		Alternatives []DownloadAlternative `json:"alternatives,omitempty"`
	}

	var raw rawDownloadData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.Author = raw.Author
	d.Title = raw.Title
	d.Video = raw.Video
	if d.Video == "" {
		d.Video = raw.Vidio
	}
	d.Images = raw.Images
	d.Audio = raw.Audio
	d.Thumbnail = raw.Thumbnail
	d.Download = raw.Download
	d.Alternatives = raw.Alternatives

	return nil
}

// TikTokResponse is the typed response from /downloader/tiktok.
type TikTokResponse struct {
	ResponseMeta
	Data TikTokData `json:"data"`
}

// TikTokData contains media fields for TikTok downloader output.
type TikTokData struct {
	Author string   `json:"author,omitempty"`
	Title  string   `json:"title,omitempty"`
	Video  string   `json:"video,omitempty"`
	Images []string `json:"images,omitempty"`
	Audio  string   `json:"audio,omitempty"`
}

// UnmarshalJSON accepts both "video" and legacy "vidio" keys.
func (d *TikTokData) UnmarshalJSON(data []byte) error {
	type rawTikTokData struct {
		Author string   `json:"author,omitempty"`
		Title  string   `json:"title,omitempty"`
		Video  string   `json:"video,omitempty"`
		Vidio  string   `json:"vidio,omitempty"`
		Images []string `json:"images,omitempty"`
		Audio  string   `json:"audio,omitempty"`
	}

	var raw rawTikTokData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.Author = raw.Author
	d.Title = raw.Title
	d.Video = raw.Video
	if d.Video == "" {
		d.Video = raw.Vidio
	}
	d.Images = raw.Images
	d.Audio = raw.Audio

	return nil
}

// SoundCloudResponse is the typed response from /downloader/soundcloud.
type SoundCloudResponse struct {
	ResponseMeta
	Data SoundCloudData `json:"data"`
}

// SoundCloudData contains media fields for SoundCloud downloader output.
type SoundCloudData struct {
	Title        string                `json:"title,omitempty"`
	Thumbnail    string                `json:"thumbnail,omitempty"`
	Download     string                `json:"download,omitempty"`
	Alternatives []DownloadAlternative `json:"alternatives,omitempty"`
}

// FacebookResponse is the typed response from /downloader/facebook.
type FacebookResponse struct {
	ResponseMeta
	Data FacebookData `json:"data"`
}

// FacebookData contains media fields for Facebook downloader output.
type FacebookData struct {
	Title     string             `json:"title,omitempty"`
	Downloads []FacebookDownload `json:"downloads,omitempty"`
}

// FacebookDownload contains download information with quality and URL.
type FacebookDownload struct {
	Quality string `json:"quality,omitempty"`
	URL     string `json:"url,omitempty"`
}
