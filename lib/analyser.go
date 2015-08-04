package lib

type Analyser interface {
    Analyse(page *Page) *Result
}
