package downloader

import (
	"github.com/sakeven/spidergo/lib/raw"
	"github.com/sakeven/spidergo/lib/request"
)

type Downloader interface {
	SetCallBack(c chan<- *raw.Raw)
	Download(req *request.Request)
}
