package mcversions

import (
	"errors"
)

var (
	// ErrCacheMissing can usually be fixed by calling MCVersions.Grab()
	ErrCacheMissing = errors.New("cache missing")

	// ErrCacheExpired can usually be fixed by calling MCVersions.Grab()
	ErrCacheExpired = errors.New("cache expired")
)
