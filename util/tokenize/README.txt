package tokenize // import "github.com/jackdoe/go-query/util/tokenize"

Simlpe tokenizer chain

Example:

    package main
    import t "github.com/jackdoe/go-query/util/tokenize"
    func main() {
    	tokenizer := []t.Tokenizer{t.NewWhitespace(), t.NewLeftEdge(1), t.NewUnique()}
    	tokens := t.Tokenize("hello world", tokenizer...)

    	fmt.Printf("%v",tokens) // [h he hel hell hello w wo wor worl world]
    }

FUNCTIONS

func Tokenize(s string, tokenizers ...Tokenizer) []string

TYPES

type Custom struct {
	// Has unexported fields.
}

func NewCustom(f func([]string) []string) *Custom

func (l *Custom) Apply(s []string) []string

type LeftEdge struct {
	// Has unexported fields.
}

func NewLeftEdge(n int) *LeftEdge

func (e *LeftEdge) Apply(current []string) []string

type Tokenizer interface {
	Apply([]string) []string
}

type Unique struct {
}

func NewUnique() *Unique

func (w *Unique) Apply(current []string) []string

type Whitespace struct{}

func NewWhitespace() *Whitespace

func (w *Whitespace) Apply(current []string) []string

