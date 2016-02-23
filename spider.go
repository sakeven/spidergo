package spidergo

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/sakeven/spidergo/lib/analyser"
	"github.com/sakeven/spidergo/lib/downloader"
	"github.com/sakeven/spidergo/lib/page"
	"github.com/sakeven/spidergo/lib/pipe"
	"github.com/sakeven/spidergo/lib/raw"
	"github.com/sakeven/spidergo/lib/request"
	"github.com/sakeven/spidergo/lib/scheduler"
)

type Spider struct {
	_scheduler   scheduler.Scheduler
	pipelines    []pipe.Piper
	reqs         []*request.Request
	reqChan      chan *request.Request
	rawChan      chan *raw.Raw
	pageChan     chan *page.Page
	downloadPool *downloader.Pool
	pagePool     *page.Pool
	analyserPool *analyser.Pool

	delay     uint
	threadNum uint
	depth     uint

	OnWatch bool
}

type Result struct {
}

func New(_analysers []analyser.Analyser) *Spider {
	s := new(Spider)
	s.threadNum = 1
	s.delay = 1
	s.reqChan = make(chan *request.Request, 8)
	s.pageChan = make(chan *page.Page, 8)
	s.rawChan = make(chan *raw.Raw, 8)

	s.analyserPool = analyser.NewPool(_analysers)
	return s
}

func (s *Spider) AddRequest(req *http.Request) *Spider {
	_req := request.New(req, 0)
	s.reqs = append(s.reqs, _req)

	return s
}

func (s *Spider) RegisterDownload(downloaders []downloader.Downloader) *Spider {
	s.rawChan = make(chan *raw.Raw, len(downloaders))

	for _, _downloader := range downloaders {
		_downloader.SetCallBack(s.rawChan)
	}

	s.downloadPool = downloader.NewPool(downloaders)
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

// SetDelay set delay time after fetched a url.
func (s *Spider) SetDelay(delay uint) *Spider {
	s.delay = delay

	return s
}

// SetDepth set how deep we dig in.
func (s *Spider) SetDepth(depth uint) *Spider {
	s.depth = depth

	return s
}

// register register all components which wasn't registered.
func (s *Spider) register() {

	if s.pagePool == nil {
		var pageProcessors []page.PageProcessor
		for i := uint(0); i < s.threadNum; i++ {
			pageProcessor := page.NewPageProcessor()
			pageProcessors = append(pageProcessors, pageProcessor)
		}
		s.pagePool = page.NewPool(pageProcessors)
	}

	if s.downloadPool == nil {
		var downloaders []downloader.Downloader
		for i := uint(0); i < s.threadNum; i++ {
			_downloader := downloader.New(fmt.Sprintf("down %d", i))
			_downloader.SetCallBack(s.rawChan)
			downloaders = append(downloaders, _downloader)
		}
		s.downloadPool = downloader.NewPool(downloaders)
	}

	s.download()
	s.page()
	s.analyse()

	if s._scheduler == nil {
		s._scheduler = scheduler.New()
	}

	s._scheduler.SetMaxDepth(s.depth)
	for _, req := range s.reqs {
		s._scheduler.Add(req)
	}

	if s.OnWatch {
		s.Watch()
	}

	time.Sleep(1 * time.Second)
}

func (s *Spider) download() {

	go func() {
		for req := range s.reqChan {
			_downloader := s.downloadPool.Get()
			go func(req *request.Request) {
				defer s.downloadPool.Release(_downloader)
				_downloader.Download(req)
			}(req)
		}
	}()
}

func (s *Spider) page() {
	go func() {
		for r := range s.rawChan {
			if r == nil {
				continue
			}
			pageProcessor := s.pagePool.Get()
			go func(r *raw.Raw) {
				defer s.pagePool.Release(pageProcessor)

				page := pageProcessor.Process(r.Req, r.Resp)
				if page == nil {
					return
				}

				s.pageChan <- page
			}(r)
		}
	}()
}

func (s *Spider) analyse() {
	go func() {
		for p := range s.pageChan {
			_analyser := s.analyserPool.Get()
			go func(p *page.Page) {
				defer s.analyserPool.Release(_analyser)
				res := _analyser.Analyse(p)
				for _, pipeline := range s.pipelines {
					pipeline.Write(res, "")
				}

				for _, r := range p.NewReqs {
					s._scheduler.Add(request.New(r, p.Req.Depth+1))
				}
			}(p)
		}

	}()
}

func (s *Spider) Stop() {
	// TODO stop one by one
	close(s.reqChan)
	close(s.rawChan)
	close(s.pageChan)
}

func (s *Spider) Watch() {

	go func() {
		for true {
			if s.CanStop() {
				return
			}
			log.Println(s._scheduler.Remain(),
				s._scheduler.Total(),
				s.downloadPool.Used(),
				s.pagePool.Used(),
				s.analyserPool.Used(),
				len(s.rawChan),
				len(s.pageChan),
				len(s.reqChan))
			time.Sleep(time.Second)

		}
	}()
}

func (s *Spider) CanStop() bool {

	if s._scheduler.Remain() > 0 ||
		s.downloadPool.Used() > 0 ||
		s.pagePool.Used() > 0 ||
		s.analyserPool.Used() > 0 ||
		len(s.rawChan) > 0 ||
		len(s.pageChan) > 0 ||
		len(s.reqChan) > 0 {
		return false

	}
	return true
}

// Run begin run spider.
func (s *Spider) Run() {
	defer s.Stop()

	s.register()

	cnt := 0
	stopCnt := 1
	start := time.Now()

	for true {
		if s.CanStop() {
			if stopCnt == 8 {
				break
			}
			stopCnt *= 2
			runtime.Gosched()
			log.Printf("can stop %d\n", stopCnt)
			time.Sleep(time.Duration(stopCnt) * time.Second)
			runtime.Gosched()
			continue
		}
		stopCnt = 1
		req := s._scheduler.Get()
		if req == nil {
			time.Sleep(time.Second)
			runtime.Gosched()
			continue
		}
		cnt++
		s.reqChan <- req

		time.Sleep(time.Duration(s.delay) * time.Microsecond)
	}

	end := time.Now()

	log.Printf("start at %s, end at %s, total %s\n", start.String(), end.String(), end.Sub(start).String())
	log.Printf("total urls %d/%d\n", cnt, s._scheduler.Total())

}
