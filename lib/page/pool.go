package page

import (
	"net/http"

	"github.com/sakeven/spidergo/lib/pool"
	"github.com/sakeven/spidergo/lib/request"
)

type Pool struct {
	pool *pool.Pool
}

type PageProcessor interface {
	Process(req *request.Request, resp *http.Response) *Page
}

type DefaultPageProcessor struct {
}

func NewPageProcessor() PageProcessor {
	return &DefaultPageProcessor{}

}

func (d DefaultPageProcessor) Process(req *request.Request, resp *http.Response) *Page {
	return New(req, resp)
}

func NewPool(processors []PageProcessor) *Pool {
	p := &Pool{}
	var voids []interface{}
	for _, pro := range processors {
		voids = append(voids, pro)
	}

	p.pool = pool.New(voids)

	return p

}

func (p *Pool) Get() PageProcessor {
	return p.pool.Get().(PageProcessor)
}

func (p *Pool) Release(processor PageProcessor) {
	p.pool.Release(processor)
}

func (p *Pool) Used() uint {
	return p.pool.Used()
}

func (p *Pool) Total() uint {
	return p.pool.Total()
}
