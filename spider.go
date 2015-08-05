package spidergo

import (
	// "log"
	"net/http"
	// "runtime"

	// "github.com/PuerkitoBio/goquery"
	"github.com/sakeven/spidergo/lib"
)

type Spider struct {
	downloader lib.Downloader
	analyser   lib.Analyser
	scheduler  lib.Scheduler
	threadNum  uint
}

type Result struct {
}

func New(analyser lib.Analyser) *Spider {
	s := new(Spider)
	s.threadNum = 1
	s.analyser = analyser
	return s
}

func (s *Spider) AddRequest(req *http.Request) *Spider {
	s.scheduler.Add(req)

	return s
}

func (s *Spider) RegisterDownload(download lib.Downloader) *Spider {
	s.downloader = download

	return s
}

func (s *Spider) RegisterScheduler(scheduler lib.Scheduler) *Spider {
	s.scheduler = scheduler
	return s
}

func (s *Spider) SetThreadNum(n uint) *Spider {
	s.threadNum = n

	return s
}

func (s *Spider) Run() {
	pool := lib.NewPool(s.threadNum)

	for pool.Count() > 0 || s.scheduler.Count() > 0 {
		req := s.scheduler.Get()
		if req == nil {
			continue
		}
		pool.Get()

		go func() {
			defer pool.Release()

			res, _ := s.downloader.Download(req)

			page := lib.NewPage(res)
			s.analyser.Analyse(page)
			for _, req := range page.NewReqs {
				s.scheduler.Add(req)
			}
		}()
	}

	// for true {
	// runtime.Gosched()
	//}

}
