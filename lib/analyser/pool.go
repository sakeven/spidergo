package analyser

import "github.com/sakeven/spidergo/lib/pool"

type Pool struct {
	pool *pool.Pool
}

func NewPool(analysers []Analyser) *Pool {
	p := &Pool{}

	var voids []interface{}
	for _, ays := range analysers {
		voids = append(voids, ays)
	}

	p.pool = pool.New(voids)

	return p
}

func (p *Pool) Get() Analyser {
	return p.pool.Get().(Analyser)
}

func (p *Pool) Release(analyser Analyser) {
	p.pool.Release(analyser)
}

func (p *Pool) Used() uint {
	return p.pool.Used()
}

func (p *Pool) Total() uint {
	return p.pool.Total()
}
