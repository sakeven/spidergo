package page

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sakeven/spidergo/lib/request"
	"golang.org/x/net/html/charset"
)

type DefaultPageProcessor struct {
}

func NewPageProcessor() PageProcessor {
	return &DefaultPageProcessor{}

}

func (p *DefaultPageProcessor) Process(req *request.Request, res *http.Response) *Page {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
		res.Body.Close()
	}()
	page := new(Page)
	page.NewReqs = make([]*http.Request, 0)

	page.ContentType = res.Header.Get("Content-type")
	page.Cookies = res.Cookies()
	page.StatusCode = res.StatusCode
	page.Req = req

	body, err := p.Iconv(res.Body, page.ContentType)
	if err != nil {
		log.Println(err)
		return nil
	}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
		return nil
	}

	page.Raw = b

	return page
}

func (p *DefaultPageProcessor) Iconv(reader io.Reader, contentType string) (io.Reader, error) {
	switch {
	case contain(contentType, "text"):
		return charset.NewReader(reader, contentType)
	}

	return reader, nil
}

func contain(src string, dst string) bool {
	return strings.Index(src, dst) >= 0
}
