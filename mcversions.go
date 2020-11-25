package mcversions

import (
	"encoding/json"
	"errors"
)

// APIResponse is used to parse the data from the api.
type APIResponse struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []APIVersionData `json:"versions"`
}

// APIVersionData is used to store version objects.
type APIVersionData struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

// NewMCVersions creates an MCVersions instance.
func NewMCVersions() (*MCVersions, error) {
	data := &MCVersions{}
	var err error
	data.data, err = data.grab()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MCVersions is the main struct for API requests.
type MCVersions struct {
	data *APIResponse
}

func (mcv *MCVersions) grab() (*APIResponse, error) {
	var err error
	data := &APIResponse{}

	bytes, err := Request("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}

	return data, err
}

// Get is used to get a version by id.
func (mcv *MCVersions) Get(id string) (*MCVersionDownloads, error) {
	mcvd := &MCVersionDownloads{}
	versionurl := ""
	for i := 0; i < len(mcv.data.Versions); i++ {
		if mcv.data.Versions[i].ID == id {
			versionurl = mcv.data.Versions[i].URL
		}
	}
	if versionurl == "" {
		return nil, errors.New("Missing version")
	}
	var err error
	mcvd.data, err = mcvd.grab(versionurl)
	if err != nil {
		return nil, err
	}
	return mcvd, nil
}

// GetLatestRelease is used to get the id for the latest release.
func (mcv *MCVersions) GetLatestRelease() string {
	return mcv.data.Latest.Release
}

// GetLatestSnapshot is used to get the id for the latest snapshot.
func (mcv *MCVersions) GetLatestSnapshot() string {
	return mcv.data.Latest.Snapshot
}

// List is used to get a list of all valid version ids.
func (mcv *MCVersions) List() []string {
	ids := make([]string, len(mcv.data.Versions))
	for i := 0; i < len(mcv.data.Versions); i++ {
		ids[i] = mcv.data.Versions[i].ID
	}
	return ids
}
