package mcversions

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Request is used for downloading json data.
func Request(url string) ([]byte, error) {
	var body []byte

	mcversionsClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return body, err
	}

	req.Header.Set("User-Agent", "mcversions")

	res, err := mcversionsClient.Do(req)
	if err != nil {
		return body, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}
