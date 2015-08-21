# spidergo

:beetle: A high performance spider(crawler) written in go.

##Feature

* Concurrent
* Distributed
* Support analyse html page and download binary file

##Installation

```bash
go get github.com/sakeven/spidergo
```

##Example

There is a example in `github/sakeven/spidergo/example`. Blow is a file from example.

```go
package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sakeven/spidergo"
	"github.com/sakeven/spidergo/lib/analyser"
	"github.com/sakeven/spidergo/lib/page"
	"github.com/sakeven/spidergo/lib/result"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	runtime.GOMAXPROCS(runtime.NumCPU())
	req, _ := http.NewRequest("GET", "http://acm.hdu.edu.cn/listproblem.php?vol=1", nil)

	spidergo.New([]analyser.Analyser{NewAnalyser(), NewAnalyser()}).
		SetThreadNum(4).
		AddRequest(req).
		SetDelay(uint(100)).
		SetDepth(uint(4)).
		Run()
}

type Analyser struct {
}

func (a *Analyser) Analyse(pg *page.Page) *result.Result {
	if pg.Err != "" {
		log.Println(pg.Err)
		return nil
	}

	// log.Println(pg.Req.Req.URL.String())

	if pg.ContentType == "image/jpeg" {
		path := strings.Split(pg.Req.Req.URL.String(), "/")
		f, err := os.Create("out/" + path[len(path)-1])
		if err != nil {
			log.Println(err)
			return nil
		}
		defer f.Close()
		f.Write(pg.Raw)
		return nil
	}

	if pg.ContentType != "text/html" {
		return nil
	}
	pg.Doc.Find("a").Each(func(i int, se *goquery.Selection) {
		href, _ := se.Attr("href")

		if strings.HasPrefix(href, "list") {
			href = "http://acm.hdu.edu.cn/" + href
			req, err := http.NewRequest("GET", href, nil)
			if err != nil {
				log.Println(err)
				return
			}
			pg.AddReq(req)
		}

	})

	pg.Doc.Find("img").Each(func(i int, se *goquery.Selection) {
		href, _ := se.Attr("src")
		href = "http://acm.hdu.edu.cn/" + href
		href = page.FixUri(href)
		req, err := http.NewRequest("GET", href, nil)
		if err != nil {
			log.Println("req", err)
			return
		}
		pg.AddReq(req)
	})

	text := pg.Doc.Find("script").Text()
	titlePat := `p\((.*?)\);`
	titleRx := regexp.MustCompile(titlePat)
	match := titleRx.FindAllString(text, -1)
	for _, m := range match {
		pro := strings.Split(m, ",")
		if len(pro) >= 2 {
			href := "http://acm.hdu.edu.cn/showproblem.php?pid=" + pro[1]
			req, err := http.NewRequest("GET", href, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			pg.AddReq(req)
		}
	}
	return nil
}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

```

##License
Under [MIT](https://github.com/sakeven/spidergo/LICENSE)
