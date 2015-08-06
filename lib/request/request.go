package request

import (
	"net/http"
)

type Request struct {
	ID        string
	Depth     uint
	Reference uint
	Req       *http.Request
}

func New(req *http.Request, depth uint) *Request {
	return &Request{
		Depth: depth,
		Req:   req,
	}
}
