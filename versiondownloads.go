package mcversions

import "encoding/json"

// APIDownloadResponse is used to parse the data from the api.
type APIDownloadResponse struct {
	Downloads   APIDownloads `json:"downloads"`
	ID          string       `json:"id"`
	ReleaseTime string       `json:"releaseTime"`
	Type        string       `json:"type"`
}

// APIDownloads is used to store the client and server download information
type APIDownloads struct {
	Client APIDownloadData `json:"client"`
	Server APIDownloadData `json:"server"`
}

// APIDownloadData is used to store download objects.
type APIDownloadData struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// MCVersionDownloads is the main struct for download API requests.
type MCVersionDownloads struct {
	data *APIDownloadResponse
}

func (vd *MCVersionDownloads) grab(url string) (*APIDownloadResponse, error) {
	var err error
	data := &APIDownloadResponse{}

	body, err := Request(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}

// GetID is used to get the id of the version.
func (vd *MCVersionDownloads) GetID() string {
	return vd.data.ID
}

// GetReleaseTime is used to get the release time of the version.
func (vd *MCVersionDownloads) GetReleaseTime() string {
	return vd.data.ReleaseTime
}

// GetType is used to get the release type of the version.
func (vd *MCVersionDownloads) GetType() string {
	return vd.data.Type
}

// GetClient is used to get the APIDownloadData struct for the client jar.
func (vd *MCVersionDownloads) GetClient() APIDownloadData {
	return vd.data.Downloads.Client
}

// GetServer is used to get the APIDownloadData struct for the server jar.
func (vd *MCVersionDownloads) GetServer() APIDownloadData {
	return vd.data.Downloads.Server
}
