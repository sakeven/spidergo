package raw

import (
	"net/http"

	"github.com/sakeven/spidergo/lib/request"
)

type Raw struct {
	Req  *request.Request
	Resp *http.Response
}
