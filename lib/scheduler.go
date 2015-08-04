package lib

import (
    "net/http"
)

type Scheduler interface {
    Add(req *http.Request)
    Get() *http.Request
    Count() int
}
