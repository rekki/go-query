package norm

import (
	"strings"

	ps "github.com/blevesearch/go-porterstemmer"
)

type PorterStemmer struct {
}

func NewPorterStemmer() *PorterStemmer {
	return &PorterStemmer{}
}

func (p *PorterStemmer) Apply(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))

	from := 0
	to := 0
	runned := []rune(s)

	for i, c := range runned {
		if c == ' ' {
			if to > from {
				for _, r := range ps.StemWithoutLowerCasing(runned[from : to-1]) {
					sb.WriteRune(r)
				}
			}
			from = i
			to = i
		} else {
			to++
		}
	}

	if from != to {
		for _, r := range ps.StemWithoutLowerCasing(runned[from:to]) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
