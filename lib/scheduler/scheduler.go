package scheduler

import (
	"github.com/sakeven/spidergo/lib/request"
)

type Scheduler interface {
	Add(req *request.Request)
	Get() *request.Request
	Remain() int
}
