package mcversions

import (
	"io"
	"net/http"
	"time"
)

// Request is used for downloading json data.
func Request(url string) (io.ReadCloser, error) {
	mcVersionsClient := http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mcversions/1.0.0 (+https://github.com/MrMelon54/mcversions)")

	res, err := mcVersionsClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
