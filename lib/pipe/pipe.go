package pipe

import (
	"fmt"
	"io"
	"sync"

	"github.com/sakeven/spidergo/lib/result"
)

type DefaultPipeline struct {
	locker sync.Locker
	Writer io.Writer
}

func makeSentence(taskname string, k, v string) string {
	return fmt.Sprintf("[%s-%s]: %s\n", taskname, k, v)
}

func (p *DefaultPipeline) Write(res *result.Result, taskname string) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	for k, v := range res.Items {
		p.Writer.Write([]byte(makeSentence(taskname, k, v)))
	}

	return nil
}
