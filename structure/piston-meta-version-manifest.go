package structure

import "time"

// PistonMetaVersionManifest is used to parse the data from the api.
type PistonMetaVersionManifest struct {
	Expires  time.Time                `json:"expires"`
	Latest   *PistonMetaLatestVersion `json:"latest"`
	Versions []*PistonMetaVersionData `json:"versions"`
}

// PistonMetaLatestVersion used to store the latest release and snapshot versions
type PistonMetaLatestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

// PistonMetaVersionData is used to store version objects.
type PistonMetaVersionData struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	Time            string `json:"time"`
	ReleaseTime     string `json:"releaseTime"`
	Sha1            string `json:"sha1"`
	ComplianceLevel int    `json:"complianceLevel"`
}
