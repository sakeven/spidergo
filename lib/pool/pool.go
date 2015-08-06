package pool

// import (
//     "log"
// )

type Pool struct {
	c     chan uint
	used  uint
	total uint
}

func NewPool(n uint) *Pool {
	p := new(Pool)
	p.c = make(chan uint, n)
	p.total = n
	p.used = 0

	for i := uint(1); i <= n; i++ {
		p.c <- 1
	}

	return p
}
func (p *Pool) Get() {
	<-p.c
	p.used++
	// log.Println("get")
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
