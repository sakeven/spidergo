package downloader

import (
	"net/http"
	"time"

	"github.com/sakeven/spidergo/lib/pool"
	"github.com/sakeven/spidergo/lib/raw"
	"github.com/sakeven/spidergo/lib/request"
)

type DefaultDownloader struct {
	ID string
	c  chan<- *raw.Raw
}

func (d *DefaultDownloader) Download(req *request.Request) {
	client := &http.Client{Timeout: 30 * time.Second}

	res, err := client.Do(req.Req)
	if err != nil {
		return
	}

	d.c <- &raw.Raw{req, res}
}

func (d *DefaultDownloader) SetCallBack(c chan<- *raw.Raw) {
	d.c = c
}

func New() *DefaultDownloader {
	return &DefaultDownloader{}
}

type Pool struct {
	*pool.Pool
}

func NewPool(total uint) *Pool {
	p := &Pool{}
	p.Pool = pool.New(total, New())
	return p
}

func (p *Pool) Get() Downloader {
	return p.Pool.Get().(Downloader)
}
