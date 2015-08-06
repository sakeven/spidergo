package page

import (
	"encoding/json"
	"io/ioutil"
	"log"
	// "log"
	"bytes"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

type Page struct {
	Req         *http.Request
	Cookies     []*http.Cookie
	StatusCode  int
	ContentType string
	OriCharset  string
	Err         string

	Raw     []byte
	Doc     *goquery.Document
	JsonMap map[string]string
	Body    string

	NewReqs []*http.Request
}

func NewPage(res *http.Response, charset string) *Page {

	page := new(Page)
	page.NewReqs = make([]*http.Request, 0)

	page.ContentType = res.Header.Get("Content-type")
	page.Cookies = res.Cookies()
	page.StatusCode = res.StatusCode
	page.OriCharset = charset

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	page.Raw = b

	contentType := page.ContentType
	switch {
	case contain(contentType, "text/html"):
		page.ConvertCharset()
		page.ParseHtml()
	case contain(contentType, "application/json"):
		page.ParseJson()
	case contain(contentType, "text/plain"):
		page.ParseText()
	default:
	}

	return page
}

func (p *Page) AddReq(req *http.Request) {
	p.NewReqs = append(p.NewReqs, req)
}

func (p *Page) getOriCharset() string {
	var idx = 0
	if idx = strings.Index(p.ContentType, "charset="); idx < 0 {
		return p.OriCharset
	}
	return p.ContentType[idx:]
}

//TODO charset
func (p *Page) ConvertCharset() {
	charset := p.getOriCharset()
	if charset != "utf-8" {

		raw := make([]byte, len(p.Raw)*2)
		_, _, err := iconv.Convert(p.Raw, raw, charset, "utf-8")
		if err != nil {
			log.Println(err)
			return
		}
		p.Raw = raw
	}
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
