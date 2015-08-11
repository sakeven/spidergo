package page

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	// "log"
	// "bufio"
	"bytes"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sakeven/spidergo/lib/request"
	"golang.org/x/net/html/charset"
)

type Page struct {
	Req         *request.Request
	Cookies     []*http.Cookie
	StatusCode  int
	ContentType string
	OriCharset  string
	Err         string
	Failed      bool

	Raw     []byte
	Doc     *goquery.Document
	JsonMap map[string]string
	Body    string

	NewReqs []*http.Request
}

func New(req *request.Request, res *http.Response) *Page {
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

	body, err := page.Iconv(res.Body)
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

	contentType := page.ContentType
	switch {
	case contain(contentType, "text/html"):
		page.ParseHtml()
	case contain(contentType, "application/json"):
		page.ParseJson()
	case contain(contentType, "text/plain"):
		page.ParseText()
	default:
	}

	return page
}

func (p *Page) Iconv(reader io.Reader) (io.Reader, error) {
	contentType := p.ContentType
	switch {
	case contain(contentType, "text"):
		return charset.NewReader(reader, contentType)
	}

	return reader, nil

}

func (p *Page) AddReq(req *http.Request) {
	p.NewReqs = append(p.NewReqs, req)
}

func (p *Page) ParseHtml() {
	var err error
	p.Doc, err = goquery.NewDocumentFromReader(bytes.NewReader(p.Raw))
	if err != nil {
		p.Err = err.Error()
	}
}

func (p *Page) ParseJson() {
	json.Unmarshal(p.Raw, &p.JsonMap)
}

func (p *Page) ParseText() {
	p.Body = string(p.Raw)
}

func contain(src string, dst string) bool {
	return strings.Index(src, dst) >= 0
}
