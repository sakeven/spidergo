package core

import (
	// "log"
	"net/http"
)

type Scheduler struct {
	Reqs    map[string]*http.Request
	Handled map[string]*http.Request
}

func (s *Scheduler) Add(req *http.Request) {
	if _, ok := s.Reqs[req.URL.String()]; ok {
		return
	}

	if _, ok := s.Handled[req.URL.String()]; ok {
		return
	}

	s.Reqs[req.URL.String()] = req

	// log.Println(s.Count())

}

func (s *Scheduler) Get() *http.Request {
	for url, req := range s.Reqs {
		s.Handled[url] = req
		delete(s.Reqs, url)
		return req
	}

	// log.Println("here")

	return nil
}

func (s *Scheduler) Count() int {

	return len(s.Reqs)
}

func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.Reqs = make(map[string]*http.Request)
	s.Handled = make(map[string]*http.Request)

	return s
}
