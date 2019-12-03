// Simlpe tokenizer chain
//
// Example:
//  package main
//  import t "github.com/jackdoe/go-query/util/tokenize"
//  func main() {
//  	tokenizer := []t.Tokenizer{t.NewWhitespace(), t.NewLeftEdge(1), t.NewUnique()}
//  	tokens := t.Tokenize("hello world", tokenizer...)
//
//  	fmt.Printf("%v",tokens) // [h he hel hell hello w wo wor worl world]
//  }
//
package tokenize

import (
	"strings"

	"github.com/jackdoe/go-query/util/common"
)

type Tokenizer interface {
	Apply([]string) []string
}

func Tokenize(s string, tokenizers ...Tokenizer) []string {
	out := []string{}
	if len(tokenizers) == 0 {
		return out
	}
	out = tokenizers[0].Apply([]string{s})
	for i := 1; i < len(tokenizers); i++ {
		out = tokenizers[i].Apply(out)
	}

	return out
}

type LeftEdge struct {
	n int
}

func NewLeftEdge(n int) *LeftEdge {
	return &LeftEdge{n: n - 1}
}
func (e *LeftEdge) Apply(current []string) []string {
	out := []string{}
	for _, s := range current {
		if len(s) < e.n {
			out = append(out, s)
		} else {
			for i := e.n; i < len(s); i++ {
				out = append(out, s[:i+1])
			}
		}
	}
	return out
}

type Whitespace struct{}

func NewWhitespace() *Whitespace {
	return &Whitespace{}
}
func (w *Whitespace) Apply(current []string) []string {
	out := []string{}
	var sb strings.Builder
	for _, s := range current {

		hasChar := true
		for _, c := range s {
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				if sb.Len() > 0 && hasChar {
					out = append(out, sb.String())
					sb.Reset()
					hasChar = false
				}
			} else {
				hasChar = true
				sb.WriteRune(c)
			}
		}
		if sb.Len() > 0 {
			out = append(out, sb.String())
			sb.Reset()
		}
	}
	return out
}

type Custom struct {
	f func([]string) []string
}

func NewCustom(f func([]string) []string) *Custom {
	return &Custom{f: f}
}

func (l *Custom) Apply(s []string) []string {
	return l.f(s)
}

type Unique struct {
}

func NewUnique() *Unique {
	return &Unique{}
}
func (w *Unique) Apply(current []string) []string {
	return common.Unique(current)
}

type CharNgram struct {
	size int
}

func NewCharNgram(size int) *CharNgram {
	return &CharNgram{size: size}
}

func (w *CharNgram) Apply(current []string) []string {
	out := []string{}
	for _, s := range current {
		if len(s) < w.size {
			out = append(out, s)
		} else {
			for i := 0; i <= len(s)-w.size; i++ {
				out = append(out, s[i:i+w.size])
			}
		}
	}
	return out
}

// NewSurround("$").Apply([]string{"h","he","hel"}) -> []string{"$h","he","hel$"}
type Surround struct {
	s string
}

func NewSurround(s string) *Surround {
	return &Surround{s: s}
}

func (w *Surround) Apply(current []string) []string {
	if len(current) == 0 {
		return current
	}
	out := make([]string, len(current))
	copy(out, current)
	out[0] = w.s + out[0]
	out[len(out)-1] = out[len(out)-1] + w.s
	return out
}
