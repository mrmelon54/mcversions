package structure

// APIDownloadResponse is used to parse the data from the api.
type APIDownloadResponse struct {
	Downloads   *APIDownloads `json:"downloads"`
	ID          string        `json:"id"`
	ReleaseTime string        `json:"releaseTime"`
	Type        string        `json:"type"`
}

// APIDownloads is used to store the client and server download information
type APIDownloads struct {
	Client *APIDownloadData `json:"client"`
	Server *APIDownloadData `json:"server"`
}

// APIDownloadData is used to store download objects.
type APIDownloadData struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}
