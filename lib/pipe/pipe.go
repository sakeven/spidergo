package pipe

import "github.com/sakeven/spidergo/lib/result"

type DefaultPipeline struct {
}

func (p *DefaultPipeline) Write(res *result.Result, taskname string) error {
	return nil

}
