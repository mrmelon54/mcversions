package mcversions

import "errors"

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

func Get(id string) (*MCVersionDownloads, error) {
	mcv, err := checkDefaultMcVersions()
	if err != nil {
		return nil, err
	}

	var data *MCVersionDownloads
	err = runAndCheckMem(func() error {
		data, err = mcv.Get(id)
		return err
	})
	return data, err
}

func LatestRelease() (*MCVersionDownloads, error) {
	// TODO
	return nil, errors.New("hi")
}
