package spidergo

import (
	// "log"

	"net/http"
	// "runtime"

	"github.com/sakeven/spidergo/lib/analyser"
	"github.com/sakeven/spidergo/lib/downloader"
	"github.com/sakeven/spidergo/lib/page"
	"github.com/sakeven/spidergo/lib/pipe"
	"github.com/sakeven/spidergo/lib/pool"
	"github.com/sakeven/spidergo/lib/request"
	"github.com/sakeven/spidergo/lib/scheduler"
)

type Spider struct {
	_downloader downloader.Downloader
	_analyser   analyser.Analyser
	_scheduler  scheduler.Scheduler
	pipelines   []pipe.Piper
	reqs        []*request.Request

	threadNum  uint
	oriCharset string
	depth      uint
}

type Result struct {
}

func New(_analyser analyser.Analyser) *Spider {
	s := new(Spider)
	s.threadNum = 1
	s._analyser = _analyser
	return s
}

func (s *Spider) AddRequest(req *http.Request) *Spider {
	_req := request.New(req, 0)
	s.reqs = append(s.reqs, _req)

	return s
}

func (s *Spider) RegisterDownload(_download downloader.Downloader) *Spider {
	s._downloader = _download

	return s
}

func (s *Spider) RegisterScheduler(_scheduler scheduler.Scheduler) *Spider {
	s._scheduler = _scheduler
	return s
}

func (s *Spider) AddPipeline(pipeline pipe.Piper) *Spider {
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

func (s *Spider) SetThreadNum(n uint) *Spider {
	s.threadNum = n

	return s
}

// SetOriCharset set pages' original charset.
// Sometimes we can't get charset info from HTTP header Content-type.
func (s *Spider) SetOriCharset(charset string) *Spider {
	s.oriCharset = charset

	return s
}

// SetDepth set how deep we dig in.
func (s *Spider) SetDepth(depth uint) *Spider {
	s.depth = depth

	return s
}

// register register all components which wasn't registered.
func (s *Spider) register() {
	if s._downloader == nil {
		s._downloader = downloader.New()
	}

	if s._scheduler == nil {
		s._scheduler = scheduler.New()
	}

	for _, req := range s.reqs {
		s._scheduler.Add(req)
	}
}

// Run begin run spider.
func (s *Spider) Run() {

	s.register()

	pool := pool.NewPool(s.threadNum)

	for pool.Used() > 0 || s._scheduler.Remain() > 0 {
		req := s._scheduler.Get()
		if req == nil {
			continue
		}
		pool.Get()

		go func() {
			defer pool.Release()

			res, _ := s._downloader.Download(req.Req)

			page := page.NewPage(req.Req, res, s.oriCharset)
			s._analyser.Analyse(page)
			for _, r := range page.NewReqs {
				s._scheduler.Add(request.New(r, req.Depth+1))
			}
		}()
	}

}
