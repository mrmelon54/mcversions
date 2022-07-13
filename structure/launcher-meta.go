package structure

import "time"

// APIResponse is used to parse the data from the api.
type APIResponse struct {
	Expires  time.Time         `json:"expires"`
	Latest   *APILatestVersion `json:"latest"`
	Versions []*APIVersionData `json:"versions"`
}

// APILatestVersion used to store the latest release and snapshot versions
type APILatestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

// APIVersionData is used to store version objects.
type APIVersionData struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}
