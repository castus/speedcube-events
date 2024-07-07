package dataFetch

import (
	"bytes"
	"io"
	"os"
)

type FileFetcher struct{}

func (k FileFetcher) Fetch(URL string) (r io.Reader, ok bool) {
	file, err := os.Open("kalendarz-imprez.html")
	if err != nil {
		log.Error("Couldn't open file", "error", err)
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Error("Unable to read file: ", "error", err)
	}

	return bytes.NewReader(data), true
}
