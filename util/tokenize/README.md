# tokenize
--
    import "github.com/jackdoe/go-query/util/tokenize"

Simlpe tokenizer chain

Example:

    package main
    import t "github.com/jackdoe/go-query/util/tokenize"
    func main() {
    	tokenizer := []t.Tokenizer{t.NewWhitespace(), t.NewLeftEdge(1), t.NewUnique()}
    	tokens := t.Tokenize("hello world", tokenizer...)

    	fmt.Printf("%v",tokens) // [h he hel hell hello w wo wor worl world]
    }

## Usage

#### func  Tokenize

```go
func Tokenize(s string, tokenizers ...Tokenizer) []string
```

#### type Custom

```go
type Custom struct {
}
```


#### func  NewCustom

```go
func NewCustom(f func([]string) []string) *Custom
```

#### func (*Custom) Apply

```go
func (l *Custom) Apply(s []string) []string
```

#### type LeftEdge

```go
type LeftEdge struct {
}
```


#### func  NewLeftEdge

```go
func NewLeftEdge(n int) *LeftEdge
```

#### func (*LeftEdge) Apply

```go
func (e *LeftEdge) Apply(current []string) []string
```

#### type Tokenizer

```go
type Tokenizer interface {
	Apply([]string) []string
}
```


#### type Unique

```go
type Unique struct {
}
```


#### func  NewUnique

```go
func NewUnique() *Unique
```

#### func (*Unique) Apply

```go
func (w *Unique) Apply(current []string) []string
```

#### type Whitespace

```go
type Whitespace struct{}
```


#### func  NewWhitespace

```go
func NewWhitespace() *Whitespace
```

#### func (*Whitespace) Apply

```go
func (w *Whitespace) Apply(current []string) []string
```
