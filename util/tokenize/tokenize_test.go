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

func TestCharNgram(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "rome",
			out: []string{"ro", "om", "me"},
			t:   []Tokenizer{NewCharNgram(2)},
		},
		TestCase{
			in:  "rome",
			out: []string{"$ro", "om", "me$"},
			t:   []Tokenizer{NewCharNgram(2), NewSurround("$")},
		},
		TestCase{
			in:  "rome",
			out: []string{"rom", "ome"},
			t:   []Tokenizer{NewCharNgram(3)},
		},
		TestCase{
			in:  "ro",
			out: []string{"ro"},
			t:   []Tokenizer{NewCharNgram(3)},
		},
		TestCase{
			in:  "",
			out: []string{""},
			t:   []Tokenizer{NewCharNgram(3)},
		},
		TestCase{
			in:  "rome",
			out: []string{"r", "o", "m", "e"},
			t:   []Tokenizer{NewCharNgram(1)},
		},
		TestCase{
			in:  "rome",
			out: []string{"rome"},
			t:   []Tokenizer{NewCharNgram(4)},
		},
	}
	testMany(t, cases)
}

func TestShingles(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "",
			out: []string{""},
			t:   []Tokenizer{NewShingles(3)},
		},
		TestCase{
			in:  "new york",
			out: []string{"newyork"},
			t:   []Tokenizer{NewWhitespace(), NewShingles(2)},
		},
		TestCase{
			in:  "new york",
			out: []string{"new", "york"},
			t:   []Tokenizer{NewWhitespace(), NewShingles(3)},
		},
		TestCase{
			in:  "new york",
			out: []string{"new", "york"},
			t:   []Tokenizer{NewWhitespace(), NewShingles(1)},
		},
		TestCase{
			in:  "new york city",
			out: []string{"newyork", "yorkcity"},
			t:   []Tokenizer{NewWhitespace(), NewShingles(2)},
		},
		TestCase{
			in: "new york city",
			out: []string{
				"newyorkcity",
			},
			t: []Tokenizer{NewWhitespace(), NewShingles(3)},
		},
		TestCase{
			in: "new york city killa",
			out: []string{
				"newyorkcity", "yorkcitykilla",
			},
			t: []Tokenizer{NewWhitespace(), NewShingles(3)},
		},
		TestCase{
			in: "new york city killa gorilla",
			out: []string{
				"newyorkcity", "yorkcitykilla", "citykillagorilla",
			},
			t: []Tokenizer{NewWhitespace(), NewShingles(3)},
		},
	}
	testMany(t, cases)
}

func TestSurround(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello abc world",
			out: []string{"$hello", "abc", "world$"},
			t:   []Tokenizer{NewWhitespace(), NewSurround("$"), NewUnique()},
		},
		TestCase{
			in:  "",
			out: []string{},
			t:   []Tokenizer{NewWhitespace(), NewSurround("$"), NewUnique()},
		},
		TestCase{
			in:  "a",
			out: []string{"$a$"},
			t:   []Tokenizer{NewWhitespace(), NewSurround("$"), NewUnique()},
		},
	}
	testMany(t, cases)
}

func TestSoundex(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello hallo abc world warld",
			out: []string{"H400", "H400", "A120", "W643", "W643"},
			t:   []Tokenizer{NewWhitespace(), NewSoundex()},
		},
		TestCase{
			in:  "",
			out: []string{},
			t:   []Tokenizer{NewWhitespace(), NewSoundex()},
		},
	}

	testMany(t, cases)
}

func TestNoop(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello hallo abc world warld",
			out: []string{"hello hallo abc world warld"},
			t:   []Tokenizer{NewNoop()},
		},
	}

	testMany(t, cases)
}

func TestEmpty(t *testing.T) {
	cases := []TestCase{
		TestCase{
			in:  "hello hallo abc world warld",
			out: []string{},
			t:   []Tokenizer{},
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
