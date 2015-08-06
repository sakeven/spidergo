package downloader

import (
	"net/http"
)

type Downloader interface {
	Download(req *http.Request) (*http.Response, error)
}
