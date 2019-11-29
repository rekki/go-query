package tokenize

import (
	"testing"
)

type TestCase struct {
	in  string
	out []string
	t   []Tokenizer
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func testMany(t *testing.T, cases []TestCase) {
	for _, c := range cases {
		tokenized := Tokenize(c.in, c.t...)
		if !Equal(tokenized, c.out) {
			t.Fatalf("in: <%s>, out: <%v>, expected: <%v>", c.in, tokenized, c.out)
		}
	}
}

func TestUnique(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello hello world",
			out: []string{"hello", "world"},
			t:   []Tokenizer{NewWhitespace(), NewUnique()},
		},
	}
	testMany(t, cases)
}

func TestLegtEdge(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello",
			out: []string{"he", "hel", "hell", "hello"},
			t:   []Tokenizer{NewLeftEdge(2)},
		},
		TestCase{
			in:  "hello",
			out: []string{"hello"},
			t:   []Tokenizer{NewLeftEdge(20)},
		},
		TestCase{
			in:  "hello",
			out: []string{"h", "he", "hel", "hell", "hello"},
			t:   []Tokenizer{NewLeftEdge(1)},
		},
	}
	testMany(t, cases)
}

func TestWhitespace(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello",
			out: []string{"hello"},
			t:   []Tokenizer{NewWhitespace()},
		},
		TestCase{
			in:  "",
			out: []string{},
			t:   []Tokenizer{NewWhitespace()},
		},
		TestCase{
			in:  "     ",
			out: []string{},
			t:   []Tokenizer{NewWhitespace()},
		},
		TestCase{
			in:  "     a     b",
			out: []string{"a", "b"},
			t:   []Tokenizer{NewWhitespace()},
		},
		TestCase{
			in: ` a
b
c	g
d  f
`,
			out: []string{"a", "b", "c", "g", "d", "f"},
			t:   []Tokenizer{NewWhitespace()},
		},
	}
	testMany(t, cases)
}

func TestComplex(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello world hellz",
			out: []string{"h", "he", "hel", "hello", "w", "wo", "wor", "world", "hellz"},
			t: []Tokenizer{NewWhitespace(), NewLeftEdge(1), NewUnique(), NewCustom(func(c []string) []string {
				out := []string{}
				for _, v := range c {
					if len(v) != 4 {
						out = append(out, v)
					}
				}
				return out
			})},
		},
	}
	testMany(t, cases)
}
