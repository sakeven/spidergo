package pool

// import (
//     "log"
// )

type Pool struct {
	c     chan interface{}
	used  uint
	total uint
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

func New(voids []interface{}) *Pool {
	p := new(Pool)
	p.total = uint(len(voids))
	p.c = make(chan interface{}, p.total)
	p.used = 0

	for _, void := range voids {
		p.c <- void
	}

	return p
}
func (p *Pool) Get() interface{} {
	c := <-p.c
	p.used++
	// log.Println("get")
	return c
}

func (p *Pool) Release(void interface{}) {
	p.c <- void
	p.used--
	// log.Println("release")

}

func (p *Pool) Total() uint {
	return p.total
}

func (p *Pool) Used() uint {
	return p.used
}
