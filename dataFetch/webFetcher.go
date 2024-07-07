package dataFetch

import (
	"bytes"
	"io"
	"net/http"
)

type WebFetcher struct{}

func (k WebFetcher) Fetch(URL string) (r io.Reader, ok bool) {
	res, err := http.Get(URL)
	if err != nil {
		log.Error("Couldn't fetch page to scrap", "error", err, "url", URL)
		return nil, false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Status code error", "status code", res.StatusCode, "status", res.Status)
		return nil, false
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Unable to read file: ", "error", err)
	}

	return bytes.NewReader(data), true
}
