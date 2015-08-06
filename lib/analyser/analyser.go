package analyser

import (
	"github.com/sakeven/spidergo/lib/page"
	"github.com/sakeven/spidergo/lib/result"
)

type Analyser interface {
	Analyse(page *page.Page) *result.Result
}
