package model

type Version struct {
	CurrentVersion string   `json:"current_version" dc:"Current version"`
	LatestVersion  string   `json:"latest_version" dc:"The current latest version"`
	DownloadURL    string   `json:"download_url,omitempty" dc:"Download link"`
	ForceUpdate    bool     `json:"force_update,omitempty" dc:"Whether or not to force an update"`
	ReleaseNotes   []string `json:"release_notes,omitempty" dc:"Changelog"`
}
