package downloader

import (
	"net/http"
)

type DefaultDownloader struct {
}

func (d *DefaultDownloader) Download(req *http.Request) (*http.Response, error) {

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func New() *DefaultDownloader {
	return &DefaultDownloader{}
}
