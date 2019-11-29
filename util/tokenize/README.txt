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

func Tokenize(s string, tokenizers ...Tokenizer) []string
type Custom struct{ ... }
    func NewCustom(f func([]string) []string) *Custom
type LeftEdge struct{ ... }
    func NewLeftEdge(n int) *LeftEdge
type Tokenizer interface{ ... }
type Unique struct{}
    func NewUnique() *Unique
type Whitespace struct{}
    func NewWhitespace() *Whitespace
