package scheduler

import (
	"github.com/sakeven/spidergo/lib/request"
)

type Scheduler interface {
	Add(req *request.Request)
	Get() *request.Request
	SetMaxDepth(depth uint)
	Remain() int
	Total() int
}
