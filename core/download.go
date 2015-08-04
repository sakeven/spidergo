package core

import (
    "net/http"
)

type Downloader struct {
}

func (d *Downloader) Download(req *http.Request) (*http.Response, error) {

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func NewDownloader() *Downloader {
    return &Downloader{}
}
