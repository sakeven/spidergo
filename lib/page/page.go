package page

import (
	"encoding/json"
	// "log"
	// "bufio"
	"bytes"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/sakeven/spidergo/lib/request"
)

type Page struct {
	Req         *request.Request
	Cookies     []*http.Cookie
	StatusCode  int
	ContentType string
	OriCharset  string
	Err         string
	Failed      bool

	Raw []byte

	NewReqs []*http.Request
}

func (p *Page) GetReq() *request.Request {
	return p.Req
}

func (p *Page) GetStatusCode() int {
	return p.StatusCode
}

func (p *Page) GetCookies() []*http.Cookie {
	return p.Cookies
}

func (p *Page) AddReq(req *http.Request) {
	p.NewReqs = append(p.NewReqs, req)
}

func (p *Page) ParseHtml() (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewReader(p.Raw))
}

func (p *Page) ParseJson() map[string]string {
	jsonMap := make(map[string]string)
	json.Unmarshal(p.Raw, jsonMap)
	return jsonMap
}

func (p *Page) ParseText() string {
	return string(p.Raw)
}
