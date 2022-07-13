package mcversions

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"encoding/json"
	"fmt"
	"github.com/wessie/appdirs"
	"io"
	"os"
	"path"
	"time"
)

const launcherMetaEndpoint = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

// NewMCVersions creates an MCVersions instance.
func NewMCVersions() (*MCVersions, error) {
	data := &MCVersions{}
	data.app = appdirs.New("mcversions", "MrMelon54", "")
	return data, nil
}

// MCVersions is the main struct for API requests.
type MCVersions struct {
	app  *appdirs.App
	data *structure.APIResponse
}

func (mcv *MCVersions) checkMemCache() error {
	if mcv.data == nil {
		return ErrCacheMissing
	}
	if mcv.data.Expires.After(time.Now()) {
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
		return err
	}

	defer func(body io.Closer) {
		_ = body.Close()
	}(body)

	// Decode the data
	data := structure.APIResponse{}
	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		return err
	}

	// Check is cache is outdated
	if data.Expires.After(time.Now()) {
		return fmt.Errorf("cache end time '%s' is after the current time", data.Expires.Format(time.UnixDate))
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
	data := structure.APIResponse{}
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
func (mcv *MCVersions) GetVersion(id string) (*structure.APIVersionData, error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}

	var version *structure.APIVersionData
	for i := 0; i < len(mcv.data.Versions); i++ {
		if mcv.data.Versions[i].ID == id {
			version = mcv.data.Versions[i]
		}
	}
	return version, nil
}

// ListVersions is used to get a list of all valid version ids.
func (mcv *MCVersions) ListVersions() ([]*structure.APIVersionData, error) {
	if err := mcv.checkMemCache(); err != nil {
		return nil, err
	}
	return mcv.data.Versions, nil
}
