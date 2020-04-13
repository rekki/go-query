package analyzer

import (
	norm "github.com/rekki/go-query-normalize"
	"github.com/rekki/go-query/util/tokenize"
)

type Analyzer struct {
	n      []norm.Normalizer
	search []tokenize.Tokenizer
	index  []tokenize.Tokenizer
}

func NewAnalyzer(normalizer []norm.Normalizer, search []tokenize.Tokenizer, index []tokenize.Tokenizer) *Analyzer {
	return &Analyzer{n: normalizer, search: search, index: index}
}

func (a *Analyzer) AnalyzeIndex(s string) []string {
	return tokenize.Tokenize(norm.Normalize(s, a.n...), a.index...)
}

func (a *Analyzer) AnalyzeSearch(s string) []string {
	return tokenize.Tokenize(norm.Normalize(s, a.n...), a.search...)
}
