package mcversions

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"errors"
)

var (
	defaultMcv   *MCVersions
	canRepeatArr = []error{
		ErrCacheMissing,
		ErrCacheExpired,
	}
)

func checkDefaultMcVersions() (*MCVersions, error) {
	var err error
	if defaultMcv == nil {
		defaultMcv, err = NewMCVersions()
		if err != nil {
			return nil, err
		}
	}
	return defaultMcv, nil
}

func canRepeatAfterGrab(err error) bool {
	for _, i := range canRepeatArr {
		if errors.Is(err, i) {
			return true
		}
	}
	return false
}

func runAndCheckMem(cb func() error) error {
	err := cb()
	if err != nil {
		if canRepeatAfterGrab(err) {
			return cb()
		}
		return err
	}
	return nil
}

// Version is a utility function to get version download information using the specific ID
func Version(id string) (*structure.APIVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data *structure.APIVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.GetVersion(id)
		return err
	})
	return data, err
}

// LatestRelease is a utility function to get the download information for the latest release
func LatestRelease() (*structure.APIVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}
	return Version(mcv.data.Latest.Release)
}

// LatestSnapshot is a utility function to get the download information for the latest snapshot
func LatestSnapshot() (*structure.APIVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}
	return Version(mcv.data.Latest.Snapshot)
}

// ListVersions is a utility function to get the download information for all versions
func ListVersions() ([]*structure.APIVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data []*structure.APIVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.ListVersions()
		return err
	})
	return data, err
}
