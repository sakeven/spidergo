package scheduler

import (
	// "log"

	"github.com/sakeven/spidergo/lib/request"
)

type DefaultScheduler struct {
	Reqs    map[string]*request.Request
	Handled map[string]*request.Request
}

func (s *DefaultScheduler) Add(req *request.Request) {
	if _, ok := s.Reqs[req.Req.URL.String()]; ok {
		return
	}

	if _, ok := s.Handled[req.Req.URL.String()]; ok {
		return
	}

	s.Reqs[req.Req.URL.String()] = req

}

func (s *DefaultScheduler) Get() *request.Request {
	for url, req := range s.Reqs {
		s.Handled[url] = req
		delete(s.Reqs, url)
		return req
	}

	return nil
}

func (s *DefaultScheduler) Remain() int {

	return len(s.Reqs)
}

func New() *DefaultScheduler {
	s := new(DefaultScheduler)
	s.Reqs = make(map[string]*request.Request)
	s.Handled = make(map[string]*request.Request)

	return s
}
