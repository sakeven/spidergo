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
	"github.com/sakeven/spidergo/lib/page"
	"github.com/sakeven/spidergo/lib/result"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	req, _ := http.NewRequest("GET", "http://acm.hdu.edu.cn/listproblem.php?vol=1", nil)

	spidergo.New(NewAnalyser()).
		SetThreadNum(4).
		AddRequest(req).
		SetOriCharset("gb2312").
		Run()
}

type Analyser struct {
}

func (a *Analyser) Analyse(page *page.Page) *result.Result {
	if page.Err != "" {
		log.Println(page.Err)
		return nil
	}

	if page.ContentType == "image/jpeg" {
		log.Println(page.Req.URL.String())
		path := strings.Split(page.Req.URL.String(), "/")
		f, err := os.Create("out/" + path[len(path)-1])
		if err != nil {
			log.Println(err)
			return nil
		}
		defer f.Close()
		f.Write(page.Raw)
		return nil
	}

	if page.ContentType != "text/html" {
		return nil
	}
	page.Doc.Find("a").Each(func(i int, se *goquery.Selection) {
		href, _ := se.Attr("href")

		if strings.HasPrefix(href, "list") {
			href = "http://acm.hdu.edu.cn/" + href
			req, err := http.NewRequest("GET", href, nil)
			if err != nil {
				log.Println(err)
				return
			}
			page.AddReq(req)
		}

	})

	page.Doc.Find("img").Each(func(i int, se *goquery.Selection) {
		href, _ := se.Attr("src")
		href = "http://acm.hdu.edu.cn/" + href
		req, err := http.NewRequest("GET", href, nil)
		if err != nil {
			log.Println(err)
			return
		}
		page.AddReq(req)
	})

	text := page.Doc.Find("script").Text()
	titlePat := `p\((.*?)\);`
	titleRx := regexp.MustCompile(titlePat)
	match := titleRx.FindAllString(text, -1)
	for _, m := range match {
		log.Println(m)
		pro := strings.Split(m, ",")
		if len(pro) >= 2 {
			href := "http://acm.hdu.edu.cn/showproblem.php?pid=" + pro[1]
			req, err := http.NewRequest("GET", href, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			page.AddReq(req)
		}
	}
	return nil
}

func NewAnalyser() *Analyser {
	return &Analyser{}
}
