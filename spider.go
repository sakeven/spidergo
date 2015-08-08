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
	"github.com/sakeven/spidergo/lib/pool"
	"github.com/sakeven/spidergo/lib/raw"
	"github.com/sakeven/spidergo/lib/request"
	"github.com/sakeven/spidergo/lib/scheduler"
)

type Spider struct {
	_downloader  downloader.Downloader
	_analyser    analyser.Analyser
	_scheduler   scheduler.Scheduler
	pipelines    []pipe.Piper
	reqs         []*request.Request
	reqChan      chan *request.Request
	rawChan      chan *raw.Raw
	pageChan     chan *page.Page
	downloadPool *downloader.Pool
	pagePool     *pool.Pool
	analysePool  *pool.Pool

	delay      uint
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
	s.delay = 1
	s.reqChan = make(chan *request.Request, 8)
	s.rawChan = make(chan *raw.Raw, 8)
	s.pageChan = make(chan *page.Page, 8)
	return s
}

func (s *Spider) AddRequest(req *http.Request) *Spider {
	_req := request.New(req, 0)
	s.reqs = append(s.reqs, _req)

	return s
}

func (s *Spider) RegisterDownload(_downloads []downloader.Downloader) *Spider {
	//s._downloader = _download

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
		s._downloader = downloader.New("addd")
	}

	s._downloader.SetCallBack(s.rawChan)
	s.download()
	s.page()
	s.analyse()

	if s._scheduler == nil {
		s._scheduler = scheduler.New()
	}

	for _, req := range s.reqs {
		s._scheduler.SetMaxDepth(s.depth)
		s._scheduler.Add(req)
	}
}

// func (s *Spider) download() {
//     s.downloadPool = pool.NewPool(s.threadNum)
//     go func() {
//         for req := range s.reqChan {
//             s.downloadPool.Get()
//             go func() {
//                 defer s.downloadPool.Release()
//                 log.Println("download")
//                 s._downloader.Download(req)
//             }()
//         }
//     }()
// }

func (s *Spider) download() {

	var downloaders []downloader.Downloader
	for i := uint(0); i < s.threadNum; i++ {
		_downloader := downloader.New(fmt.Sprintf("downloader %d", i))
		_downloader.SetCallBack(s.rawChan)

		downloaders = append(downloaders, _downloader)
	}

	s.downloadPool = downloader.NewPool(downloaders)
	go func() {
		for req := range s.reqChan {
			_downloader := s.downloadPool.Get()
			go func() {
				defer s.downloadPool.Release(_downloader)
				log.Println("download")
				_downloader.Download(req)
			}()
		}
	}()
}

func (s *Spider) page() {
	s.pagePool = pool.NewPool(s.threadNum)
	go func() {
		for raw := range s.rawChan {
			s.pagePool.Get()
			go func() {
				defer s.pagePool.Release(0)

				page := page.NewPage(raw.Req, raw.Resp, s.oriCharset)
				log.Println("page")
				if page == nil {
					return
				}

				s.pageChan <- page
			}()
		}
	}()
}

func (s *Spider) analyse() {
	s.analysePool = pool.NewPool(s.threadNum)
	go func() {
		for page := range s.pageChan {
			s.analysePool.Get()
			go func() {
				defer s.analysePool.Release(0)
				log.Println("analyse")
				s._analyser.Analyse(page)
				log.Println("reqsss", len(page.NewReqs))
				for _, r := range page.NewReqs {
					s._scheduler.Add(request.New(r, page.Req.Depth+1))
				}
				log.Println("end", s._scheduler.Remain())
			}()
		}

	}()
}

func (s *Spider) Stop() {
	// TODO stop one by one
	close(s.reqChan)
	close(s.rawChan)
	close(s.pageChan)
}

// Run begin run spider.
func (s *Spider) Run() {

	s.register()

	start := time.Now()
	cnt := 0

	for s._scheduler.Remain() > 0 ||
		s.downloadPool.Used() > 0 ||
		s.pagePool.Used() > 0 ||
		s.analysePool.Used() > 0 ||
		len(s.rawChan) > 0 ||
		len(s.pageChan) > 0 ||
		len(s.reqChan) > 0 {

		req := s._scheduler.Get()
		if req == nil {
			time.Sleep(time.Second)
			runtime.Gosched()
			continue
		}
		cnt++
		log.Println("sched")
		s.reqChan <- req

		time.Sleep(time.Duration(s.delay) * time.Microsecond)
	}

	end := time.Now()

	log.Printf("strat at %s, end at %s, total %s\n", start.String(), end.String(), end.Sub(start).String())
	log.Printf("total urls %d/%d\n", cnt, s._scheduler.Total())

}
