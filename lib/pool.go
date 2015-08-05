package lib

// import (
//     "log"
// )

type Pool struct {
	c chan uint
	n int
}

func NewPool(n uint) *Pool {
	p := new(Pool)
	p.c = make(chan uint, n)

	for i := uint(1); i <= n; i++ {
		p.c <- 1
	}

	return p
}
func (p *Pool) Get() {
	<-p.c
	p.n++
	// log.Println("get")
}

func (p *Pool) Release() {
	p.c <- 1
	p.n--
	// log.Println("release")

}

func (p *Pool) Count() int {
	return p.n
}
