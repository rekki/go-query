// Simlpe tokenizer chain
//
// Example:
//  package main
//  import t "github.com/rekki/go-query/util/tokenize"
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
)

type Token struct {
	Text     string
	Position int
	LineNo   int
}

func (t Token) Clone(s string) Token {
	return Token{Text: s, LineNo: t.LineNo, Position: t.Position}
}

type Tokenizer interface {
	Apply([]Token) []Token
}

func Tokenize(s string, tokenizers ...Tokenizer) []string {
	return tokensToString(TokenizeT(s, tokenizers...))
}

func TokenizeT(s string, tokenizers ...Tokenizer) []Token {
	out := []Token{}
	if len(tokenizers) == 0 {
		return out
	}
	out = tokenizers[0].Apply([]Token{{Text: s}})
	for i := 1; i < len(tokenizers); i++ {
		out = tokenizers[i].Apply(out)
	}

	return out
}

func tokensToString(in []Token) []string {
	out := make([]string, len(in))
	for i, t := range in {
		out[i] = t.Text
	}
	return out
}

type LeftEdge struct {
	n int
}

func NewLeftEdge(n int) *LeftEdge {
	return &LeftEdge{n: n - 1}
}
func (e *LeftEdge) Apply(current []Token) []Token {
	out := []Token{}
	for _, s := range current {
		if len(s.Text) < e.n {
			out = append(out, s)
		} else {
			for i := e.n; i < len(s.Text); i++ {
				out = append(out, s.Clone(s.Text[:i+1]))
			}
		}
	}
	return out
}

type Whitespace struct{}

func NewWhitespace() *Whitespace {
	return &Whitespace{}
}
func (w *Whitespace) Apply(current []Token) []Token {
	out := []Token{}
	var sb strings.Builder
	lineNo := 0
	position := 0
	for _, s := range current {
		hasChar := true
		for _, c := range s.Text {
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				if sb.Len() > 0 && hasChar {
					out = append(out, Token{Text: sb.String(), Position: position, LineNo: lineNo})
					sb.Reset()
					hasChar = false

					position++
				}
			} else {
				hasChar = true
				sb.WriteRune(c)
			}

			if c == '\n' {
				lineNo++
			}
		}
		if sb.Len() > 0 {
			out = append(out, Token{Text: sb.String(), Position: position, LineNo: lineNo})
			sb.Reset()
			position++
		}
	}
	return out
}

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}
func (w *Noop) Apply(current []Token) []Token {
	return current
}

type Custom struct {
	f func([]Token) []Token
}

func NewCustom(f func([]Token) []Token) *Custom {
	return &Custom{f: f}
}

func (l *Custom) Apply(s []Token) []Token {
	return l.f(s)
}

type Unique struct {
}

func NewUnique() *Unique {
	return &Unique{}
}
func (w *Unique) Apply(current []Token) []Token {
	seen := map[string]bool{}
	out := []Token{}
	for _, v := range current {
		_, ok := seen[v.Text]
		if !ok {
			out = append(out, v)
			seen[v.Text] = true
		}
	}
	return out
}

type CharNgram struct {
	size int
}

func NewCharNgram(size int) *CharNgram {
	return &CharNgram{size: size}
}

func (w *CharNgram) Apply(current []Token) []Token {
	out := []Token{}
	for _, s := range current {
		if len(s.Text) < w.size {
			out = append(out, s)
		} else {
			for i := 0; i <= len(s.Text)-w.size; i++ {
				out = append(out, s.Clone(s.Text[i:i+w.size]))
			}
		}
	}
	return out
}

// Shingles tokenizer (n-gram for words)
type Shingles struct {
	size int
}

// NewShingles creates new Shingles struct
func NewShingles(size int) *Shingles {
	return &Shingles{size: size}
}

// Apply applies semi shingles tokenizer
// it creates permutations "new","york","city" -> "new","newyork","york","yorkcity"
// it is very handy becuase when people search sometimes they just dont put space
func (shingles *Shingles) Apply(current []Token) []Token {
	out := make([]Token, 0, len(current)*shingles.size)
	length := len(current)
	if shingles.size > length || shingles.size < 2 {
		return current
	}

	// new york city
	// new newyork york yorkcity
	for idx, s := range current {
		out = append(out, s)
		end := idx + shingles.size

		if end <= length {
			out = append(out, joinTokens(current[idx:end]))
		}

	}
	return out
}

func joinTokens(in []Token) Token {
	if len(in) == 0 {
		return Token{}
	}
	var sb strings.Builder

	for _, t := range in {
		sb.WriteString(t.Text)
	}

	return in[0].Clone(sb.String())
}

// NewSurround("$").Apply([]string{"h","he","hel"}) -> []string{"$h","he","hel$"}
type Surround struct {
	s string
}

func NewSurround(s string) *Surround {
	return &Surround{s: s}
}

func (w *Surround) Apply(current []Token) []Token {
	if len(current) == 0 {
		return current
	}
	out := make([]Token, len(current))
	copy(out, current)

	first := out[0]
	first = first.Clone(w.s + first.Text)
	out[0] = first

	last := out[len(out)-1]
	last = last.Clone(last.Text + w.s)
	out[len(out)-1] = last

	return out
}
