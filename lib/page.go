package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	// "log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Req         *http.Request
	Cookies     []*http.Cookie
	StatusCode  int
	ContentType string
	Err         string

	Raw     []byte
	Doc     *goquery.Document
	JsonMap map[string]string
	Body    string

	NewReqs []*http.Request
}

func NewPage(res *http.Response) *Page {

	page := new(Page)
	page.NewReqs = make([]*http.Request, 0)

	page.ContentType = res.Header.Get("Content-type")
	page.Cookies = res.Cookies()
	page.StatusCode = res.StatusCode

	//TODO charset

	contentType := page.ContentType
	switch {
	case contain(contentType, "text/html"):
		page.ParseHtml(res)
	case contain(contentType, "application/json"):
		page.ParseJson(res)
	case contain(contentType, "text/plain"):
		page.ParseText(res)
	default:
	}

	return page
}

func (p *Page) AddReq(req *http.Request) {
	// log.Println(*req)
	p.NewReqs = append(p.NewReqs, req)
	// log.Println(p.NewReqs)
}

func copyResponse(dst *http.Response, src *http.Response) {

	*dst = *src
	var bodyBytes []byte

	if src.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(src.Body)
	}

	// Restore the io.ReadCloser to its original state
	src.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	dst.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

}

func (p *Page) ParseHtml(res *http.Response) {
	var err error
	p.Doc, err = goquery.NewDocumentFromResponse(res)
	if err != nil {
		p.Err = err.Error()
	}
}

func (p *Page) ParseJson(res *http.Response) {
	json.NewDecoder(res.Body).Decode(&p.JsonMap)
}

func (p *Page) ParseText(res *http.Response) {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	p.Body = string(b)
}

func contain(src string, dst string) bool {
	return strings.Index(src, dst) >= 0
}
