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
		err2 := defaultMcv.Grab()
		if err2 != nil {
			return err
		}
		if canRepeatAfterGrab(err) {
			return cb()
		}
		return err
	}
	return nil
}

// Version is a utility function to get version download information using the specific ID
func Version(id *structure.PistonMetaId) (*structure.PistonMetaVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data *structure.PistonMetaVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.GetVersion(id)
		return err
	})
	return data, err
}

// LatestRelease is a utility function to get the download information for the latest release
func LatestRelease() (*structure.PistonMetaVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data *structure.PistonMetaVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.LatestRelease()
		return err
	})
	return data, err
}

// LatestSnapshot is a utility function to get the download information for the latest snapshot
func LatestSnapshot() (*structure.PistonMetaVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data *structure.PistonMetaVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.LatestSnapshot()
		return err
	})
	return data, err
}

// ListVersions is a utility function to get the download information for all versions
func ListVersions() ([]*structure.PistonMetaVersionData, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data []*structure.PistonMetaVersionData
	err = runAndCheckMem(func() error {
		data, err = mcv.ListVersions()
		return err
	})
	return data, err
}

func VersionPackage(id *structure.PistonMetaId) (*structure.PistonMetaPackage, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var out *structure.PistonMetaPackage
	err = runAndCheckMem(func() error {
		out, err = mcv.GetVersionPackage(id)
		return err
	})
	return out, err
}
