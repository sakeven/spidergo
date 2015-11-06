package result

type Result struct {
	Items map[string]string
}

func New() *Result {
	return &Result{
		Items: make(map[string]string),
	}
}
