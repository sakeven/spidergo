package pipe

import "github.com/sakeven/spidergo/lib/result"

type Piper interface {
	Write(res *result.Result, taskname string) error
}
