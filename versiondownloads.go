package mcversions

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"encoding/json"
)

// MCPistonMeta is the main struct for download API requests.
type MCPistonMeta struct {
	data *structure.APIDownloadResponse
}

func NewPistonMeta(url string) (*MCPistonMeta, error) {
	var err error
	data := &structure.APIDownloadResponse{}

	body, err := Request(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(body).Decode(&data)
	return &MCPistonMeta{data: data}, err
}

// GetID is used to get the id of the version.
func (pm *MCPistonMeta) GetID() string {
	return pm.data.ID
}

// GetReleaseTime is used to get the release time of the version.
func (pm *MCPistonMeta) GetReleaseTime() string {
	return pm.data.ReleaseTime
}

// GetType is used to get the release type of the version.
func (pm *MCPistonMeta) GetType() string {
	return pm.data.Type
}

// GetClient is used to get the APIDownloadData struct for the client jar.
func (pm *MCPistonMeta) GetClient() *structure.APIDownloadData {
	return pm.data.Downloads.Client
}

// GetServer is used to get the APIDownloadData struct for the server jar.
func (pm *MCPistonMeta) GetServer() *structure.APIDownloadData {
	return pm.data.Downloads.Server
}
