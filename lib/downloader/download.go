package downloader

import (
	"log"
	"net/http"
	"time"

	"github.com/sakeven/spidergo/lib/pool"
	"github.com/sakeven/spidergo/lib/raw"
	"github.com/sakeven/spidergo/lib/request"
)

type DefaultDownloader struct {
	Name string
	c    chan<- *raw.Raw
}

func (d *DefaultDownloader) Download(req *request.Request) {

	client := &http.Client{Timeout: 30 * time.Second}

	res, err := client.Do(req.Req)
	if err != nil {
		log.Println(err)
		return
	}

	d.c <- &raw.Raw{Req: req, Resp: res}
}

func (d *DefaultDownloader) SetCallBack(c chan<- *raw.Raw) {
	d.c = c
}

func New(name string) *DefaultDownloader {
	return &DefaultDownloader{Name: name}
}

type Pool struct {
	pool *pool.Pool
}

func NewPool(downloaders []Downloader) *Pool {
	p := &Pool{}
	var d []interface{}
	for _, downloader := range downloaders {
		d = append(d, downloader)
	}

	p.pool = pool.New(d)
	return p
}

func (p *Pool) Get() Downloader {
	return p.pool.Get().(Downloader)
}

func (p *Pool) Release(downloader Downloader) {
	p.pool.Release(downloader)
}

func (p *Pool) Total() uint {
	return p.pool.Total()
}

func (p *Pool) Used() uint {
	return p.pool.Used()
}
