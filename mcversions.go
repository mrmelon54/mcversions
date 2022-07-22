package mcversions

import (
	"encoding/json"
	"fmt"
	"github.com/mrmelon54/mcversions/structure"
	"github.com/wessie/appdirs"
	"io"
	"os"
	"path"
	"time"
)

const launcherMetaEndpoint = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

// NewMCVersions creates an MCVersions instance.
func NewMCVersions() (*MCVersions, error) {
	data := &MCVersions{}
	data.app = appdirs.New("mcversions", "MrMelon54", "")
	data.latest = new(MCVersionLatest)
	return data, nil
}

// MCVersions is the main struct for API requests.
type MCVersions struct {
	app    *appdirs.App
	data   *structure.PistonMetaVersionManifest
	latest *MCVersionLatest
}

type MCVersionLatest struct {
	release  *structure.PistonMetaVersionData
	snapshot *structure.PistonMetaVersionData
}

func (mcv *MCVersions) checkMemCache() error {
	if mcv.data == nil {
		return ErrCacheMissing
	}
	if time.Now().After(mcv.data.Expires) {
		return ErrCacheExpired
	}
	return nil
}

func (mcv *MCVersions) openCacheFile(write bool) (io.ReadWriteCloser, error) {
	// Generate path
	cacheDir := mcv.app.UserCache()
	vFile := path.Join(cacheDir, "versions.json")

	// Make missing directories
	err := os.MkdirAll(path.Dir(vFile), os.ModePerm)
	if err != nil {
		return nil, err
	}

	var body io.ReadWriteCloser
	if write {
		// Create the file
		body, err = os.Create(vFile)
	} else {
		// Open the file
		body, err = os.Open(vFile)
	}
	return body, err
}

// Grab automates the load and fetch calls into a single method for use in other programs
func (mcv *MCVersions) Grab() error {
	// Try load first
	err := mcv.Load()
	if err == nil {
		return err
	}

	// Then try fetch
	err2 := mcv.Fetch()
	if err2 == nil {
		return err2
	}

	return fmt.Errorf("failed to get Minecraft versions manifest:\n - Load from cache: %s\n - Fetch from Mojang: %s", err, err2)
}

func (mcv *MCVersions) Load() error {
	// Find and open the cache file
	body, err := mcv.openCacheFile(false)
	if err != nil {
		return ErrCacheMissing
	}

	defer func(body io.Closer) {
		_ = body.Close()
	}(body)

	// Decode the data
	data := structure.PistonMetaVersionManifest{}
	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		return ErrCacheMissing
	}

	// Check is cache is outdated
	if time.Now().After(data.Expires) {
		return ErrCacheExpired
	}
	mcv.data = &data
	return nil
}

func (mcv *MCVersions) Fetch() error {
	// Make request for new version manifest
	body, err := Request(launcherMetaEndpoint)
	if err != nil {
		return err
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	// Decode the data
	data := structure.PistonMetaVersionManifest{}

	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		return err
	}

	// Apply the new cache end time
	data.Expires = time.Now().Add(5 * time.Minute)

	// Open cache file to save
	f, err := mcv.openCacheFile(true)
	if err != nil {
		return err
	}

	defer func(f io.Closer) {
		_ = f.Close()
	}(f)

	// Create and save cache file
	err = json.NewEncoder(f).Encode(&data)
	if err != nil {
		return err
	}

	mcv.data = &data
	return nil
}

// GetVersion is used to get a version by id.
func (mcv *MCVersions) GetVersion(id *structure.PistonMetaId) (*structure.PistonMetaVersionData, error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}

	var version *structure.PistonMetaVersionData
	for i := 0; i < len(mcv.data.Versions); i++ {
		if mcv.data.Versions[i].ID.Equal(id) {
			version = mcv.data.Versions[i]
		}
	}
	return version, nil
}

// ListVersions is used to get a list of all valid version ids.
func (mcv *MCVersions) ListVersions() ([]*structure.PistonMetaVersionData, error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}
	return mcv.data.Versions, nil
}

func (mcv *MCVersions) LatestRelease() (_ *structure.PistonMetaVersionData, err error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}
	if mcv.latest.release == nil {
		mcv.latest.release, err = mcv.GetVersion(mcv.data.Latest.Release)
	}
	return mcv.latest.release, err
}

func (mcv *MCVersions) LatestSnapshot() (_ *structure.PistonMetaVersionData, err error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}
	if mcv.latest.snapshot == nil {
		mcv.latest.snapshot, err = mcv.GetVersion(mcv.data.Latest.Snapshot)
	}
	return mcv.latest.snapshot, err
}

func (mcv *MCVersions) GetVersionPackage(id *structure.PistonMetaId) (*structure.PistonMetaPackage, error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}
	version, err := mcv.GetVersion(id)
	data := &structure.PistonMetaPackage{}

	body, err := Request(version.URL)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(body).Decode(&data)
	return data, err
}
