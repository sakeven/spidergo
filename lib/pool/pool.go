package pool

// import (
//     "log"
// )

type Pool struct {
	c     chan interface{}
	used  uint
	total uint
}

func New(n uint, void interface{}) *Pool {
	p := new(Pool)
	p.c = make(chan interface{}, n)
	p.total = n
	p.used = 0

	for i := uint(1); i <= n; i++ {
		p.c <- void
	}

	return p
}

func NewPool(n uint) *Pool {
	p := new(Pool)
	p.c = make(chan interface{}, n)
	p.total = n
	p.used = 0

	for i := uint(1); i <= n; i++ {
		p.c <- 0
	}

	return p
}
func (p *Pool) Get() interface{} {
	c := <-p.c
	p.used++
	// log.Println("get")
	return c
}

func (p *Pool) Release() {
	p.c <- 1
	p.used--
	// log.Println("release")

}

func (p *Pool) Total() uint {
	return p.total
}

func (p *Pool) Used() uint {
	return p.used
}
