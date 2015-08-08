package scheduler

import (
	// "log"

	"sync"

	"github.com/sakeven/spidergo/lib/request"
)

type DefaultScheduler struct {
	Reqs     map[string]*request.Request
	Handled  map[string]*request.Request
	maxDepth uint
	lock     sync.Locker
}

func New() *DefaultScheduler {
	s := new(DefaultScheduler)
	s.Reqs = make(map[string]*request.Request)
	s.Handled = make(map[string]*request.Request)
	s.lock = &sync.Mutex{}

	return s
}

func (s *DefaultScheduler) SetMaxDepth(depth uint) {
	s.maxDepth = depth
}

// Add adds a new request
func (s *DefaultScheduler) Add(req *request.Request) {
	defer s.lock.Unlock()
	s.lock.Lock()

	// out of sched handle depth
	if req.Depth > s.maxDepth {
		return
	}

	if _, ok := s.Reqs[req.Req.URL.String()]; ok {
		return
	}

	if _, ok := s.Handled[req.Req.URL.String()]; ok {
		return
	}

	s.Reqs[req.Req.URL.String()] = req

}

func (s *DefaultScheduler) Get() *request.Request {
	defer s.lock.Unlock()
	s.lock.Lock()

	for url, req := range s.Reqs {
		s.Handled[url] = req
		delete(s.Reqs, url)
		return req
	}

	return nil
}

func (s *DefaultScheduler) Remain() int {
	defer s.lock.Unlock()
	s.lock.Lock()

	return len(s.Reqs)
}

func (s *DefaultScheduler) Total() int {
	defer s.lock.Unlock()

	s.lock.Lock()

	return len(s.Reqs) + len(s.Handled)
}
